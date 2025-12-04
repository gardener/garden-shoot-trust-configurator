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

echo "Exporting Generic Token Kubeconfig secret name."

generic_kubeconfig_secret_name=$(kubectl --kubeconfig "$KIND_KUBECONFIG" -n garden get secret -o=custom-columns='name:.metadata.name' | grep generic)
if [[ $generic_kubeconfig_secret_name == "" ]]; then
  fail "Generic Token Kubeconfig not found"
fi

# Generate certificates
dev_trust_config_dir="$repo_root/dev/trust-configurator"
cert_dir="$dev_trust_config_dir/certs"

"$repo_root"/hack/generate-certs.sh \
  "$repo_root/dev/trust-configurator" \
  "garden-shoot-trust-configurator.garden.svc.cluster.local" \
  "DNS:localhost,DNS:garden-shoot-trust-configurator,DNS:garden-shoot-trust-configurator.garden,DNS:garden-shoot-trust-configurator.garden.svc,DNS:garden-shoot-trust-configurator.garden.svc.cluster.local,IP:127.0.0.1"
# Finish generating certificates

yq -i ' .global.webhooks.tls.caBundle = load_str("'"$cert_dir/ca.crt"'") | (.application.webhookConfig.caBundle style="literal") ' "$values_file" 
yq -i ' .global.webhooks.tls.crt = load_str("'"$cert_dir/tls.crt"'") | (.application.webhookConfig.crt style="literal") ' "$values_file" 
yq -i ' .global.webhooks.tls.key = load_str("'"$cert_dir/tls.key"'") | (.application.webhookConfig.key style="literal") ' "$values_file"

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

echo "Installing garden-shoot-trust-configurator in the runtime cluster."

helm upgrade \
  --install \
  --wait \
  --history-max=4 \
  --values "$values_file" \
  --set application.enabled="false" \
  --set runtime.enabled="true" \
  --set runtime.projectedKubeconfig.baseMountPath="/var/run/secrets/gardener.cloud/shoot/generic-kubeconfig" \
  --set runtime.projectedKubeconfig.genericKubeconfigSecretName="$generic_kubeconfig_secret_name" \
  --namespace garden \
  --kubeconfig "$KIND_KUBECONFIG" \
  $TRUST_CONFIGURATOR_NAME \
  ./charts/$TRUST_CONFIGURATOR_NAME 

echo "garden-shoot-trust-configurator installed successfully in the runtime cluster."

echo "Done."
