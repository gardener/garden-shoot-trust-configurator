#!/usr/bin/env bash
# SPDX-FileCopyrightText: SAP SE or an SAP affiliate company and Gardener contributors
#
# SPDX-License-Identifier: Apache-2.0

set -o errexit
set -o nounset
set -o pipefail

VIRTUAL_KUBECONFIG=""
KIND_KUBECONFIG=""

TRUST_CONFIGURATOR_NAME="garden-shoot-trust-configurator"

parse_flags() {
  while test $# -gt 0; do
    case "$1" in
    --virtual-kubeconfig)
      shift; VIRTUAL_KUBECONFIG="$1"
      ;;
    --kind-kubeconfig)
      shift; KIND_KUBECONFIG="$1"
      ;;
    esac

    shift
  done
}

parse_flags "$@"

if [[ -z "${VIRTUAL_KUBECONFIG:-}" || -z "${KIND_KUBECONFIG:-}" ]]; then
  echo "Usage: $0 --virtual-kubeconfig <virtual-garden-kubeconfig> --kind-kubeconfig <kind-gardener-kubeconfig>"
  exit 1
fi

repo_root="$(readlink -f "$(dirname "${0}")/..")"

charts_dir="$repo_root/charts/$TRUST_CONFIGURATOR_NAME"
temp_dir="$repo_root/dev/trust-configurator"
mkdir -p "$temp_dir"
values_file="$temp_dir/values.yaml"
cp "$charts_dir/values.yaml" "$values_file"

repo_image="europe-docker.pkg.dev/gardener-project/public/gardener/garden-shoot-trust-configurator"
yq -i '.runtime.image.repository = "'"$repo_image"'"' "$values_file"
# yq -i '.runtime.image.tag = "'"$trust_configurator_version"'"' "$values_file"

# Virtual cluster installation
echo "Patching Helm values: $values_file"

helm upgrade \
  --install \
  --wait \
  --history-max=4 \
  --values "$values_file" \
  --set application.enabled="true" \
  --set runtime.enabled="false" \
  --namespace garden \
  --kubeconfig "$VIRTUAL_KUBECONFIG" \
  $TRUST_CONFIGURATOR_NAME \
  ./charts/$TRUST_CONFIGURATOR_NAME 

echo "garden-shoot-trust-configurator installed successfully in the virtual cluster."

echo "Generating kubeconfig for the garden-shoot-trust-configurator to access the runtime cluster."

cluster_ca_data=$(kubectl config view --kubeconfig "$VIRTUAL_KUBECONFIG" --raw -o jsonpath='{.clusters[0].cluster.certificate-authority-data}')
cluster_server=$(kubectl config view --kubeconfig "$VIRTUAL_KUBECONFIG" --raw -o jsonpath='{.clusters[0].cluster.server}')
token=$(kubectl -n garden create token $TRUST_CONFIGURATOR_NAME --duration 48h --kubeconfig "$VIRTUAL_KUBECONFIG")

tmpfile=$(mktemp)

cat <<EOF > "$tmpfile"
---
apiVersion: v1
kind: Config
current-context: cluster
contexts:
- context:
    cluster: cluster
    user: garden-shoot-trust-configurator
  name: cluster
clusters:
- cluster:
    certificate-authority-data: $cluster_ca_data
    server: $cluster_server
  name: cluster
users:
- name: garden-shoot-trust-configurator
  user:
    token: $token
EOF

yq -i '
  .runtime.kubeconfig = load_str("'"$tmpfile"'")
  | (.runtime.kubeconfig style="literal")
' "$values_file"

rm -f "$tmpfile"

# Kind Gardener cluster installation
yq -i '
  .runtime.additionalLabels.deployment = {
    "high-availability-config.resources.gardener.cloud/type": "server",
    "networking.gardener.cloud/to-dns": "allowed",
    "networking.gardener.cloud/to-public-networks": "allowed",
    "networking.resources.gardener.cloud/to-virtual-garden-kube-apiserver-tcp-443": "allowed",
    "networking.resources.gardener.cloud/to-all-istio-ingresses-istio-ingressgateway-tcp-9443": "allowed"
  }
  |
  .runtime.additionalLabels.hpa = {
    "high-availability-config.resources.gardener.cloud/type": "server"
  }
  |
  .runtime.additionalAnnotations.service = {
    "networking.resources.gardener.cloud/from-all-webhook-targets-allowed-ports": "[{\"protocol\":\"TCP\",\"port\":10443}]",
    "networking.resources.gardener.cloud/from-all-garden-scrape-targets-allowed-ports": "[{\"protocol\":\"TCP\",\"port\":10443}]"
  }
' "$values_file"

echo "Installing garden-shoot-trust-configurator in the runtime cluster."

helm upgrade \
  --install \
  --wait \
  --history-max=4 \
  --values "$values_file" \
  --set application.enabled="false" \
  --set runtime.enabled="true" \
  --namespace garden \
  --kubeconfig "$KIND_KUBECONFIG" \
  $TRUST_CONFIGURATOR_NAME \
  ./charts/$TRUST_CONFIGURATOR_NAME 

echo "garden-shoot-trust-configurator installed successfully in the runtime cluster."

echo "Done."
