#!/usr/bin/env bash
# SPDX-FileCopyrightText: SAP SE or an SAP affiliate company and Gardener contributors
#
# SPDX-License-Identifier: Apache-2.0

set -o errexit
set -o nounset
set -o pipefail

VIRTUAL_KUBECONFIG=""
KIND_KUBECONFIG=""

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

OIDC_WEBHOOK_AUTH_NAME="oidc-webhook-authenticator"
OIDC_WEBHOOK_AUTH_REPO=""gardener/$OIDC_WEBHOOK_AUTH_NAME
OIDC_WEBHOOK_AUTH_REPO_NAME=github.com/$OIDC_WEBHOOK_AUTH_REPO

# TODO(theoddora) Use renovate bot that will be updating the version automatically
# See: https://docs.renovatebot.com/modules/versioning/
owa_version=$(go list -m -f '{{.Version}}' $OIDC_WEBHOOK_AUTH_REPO_NAME)

echo "Using OIDC Webhook Authenticator version: $owa_version"
if [[ ! -d "$repo_root/dev/$OIDC_WEBHOOK_AUTH_NAME" || -z "$(ls -A "$repo_root/dev/$OIDC_WEBHOOK_AUTH_NAME" 2>/dev/null)" ]]; then
  mkdir -p "$repo_root/dev"
  git clone --quiet -c advice.detachedHead=false https://github.com/gardener/oidc-webhook-authenticator.git "$repo_root/dev/$OIDC_WEBHOOK_AUTH_NAME" --branch "$owa_version"
else
  cd "$repo_root/dev/$OIDC_WEBHOOK_AUTH_NAME"
  git fetch
fi

cd "$repo_root/dev/$OIDC_WEBHOOK_AUTH_NAME"

git checkout "$owa_version"

# Generate certificates
dev_owa_dir="$repo_root/dev/owa"
cert_dir="$dev_owa_dir/certs"

"$repo_root"/hack/generate-certs.sh \
  "$dev_owa_dir" \
  "oidc-webhook-authenticator.garden.svc.cluster.local" \
  "DNS:localhost,DNS:oidc-webhook-authenticator,DNS:oidc-webhook-authenticator.garden,DNS:oidc-webhook-authenticator.garden.svc,DNS:oidc-webhook-authenticator.garden.svc.cluster.local,IP:127.0.0.1"
# Finish generating certificates

charts_dir="$repo_root/dev/$OIDC_WEBHOOK_AUTH_NAME/charts/$OIDC_WEBHOOK_AUTH_NAME"
values_file="$dev_owa_dir/values.yaml"
cp "$charts_dir/values.yaml" "$values_file"

repo_image="europe-docker.pkg.dev/gardener-project/public/gardener/oidc-webhook-authenticator"
yq -i '.runtime.image.repository = "'"$repo_image"'"' "$values_file"
yq -i '.runtime.image.tag = "'"$owa_version"'"' "$values_file"

# Virtual cluster installation
echo "Patching Helm values: $values_file"
yq -i '
  .application.webhookConfig.caBundle = load_str("'"$cert_dir/ca.crt"'") 
  | (.application.webhookConfig.caBundle style="literal")
' "$values_file"

helm upgrade \
  --install \
  --wait \
  --history-max=4 \
  --values "$values_file" \
  --set application.enabled="true" \
  --set application.virtualGarden.enabled="true" \
  --set runtime.enabled="false" \
  --namespace garden \
  --kubeconfig "$VIRTUAL_KUBECONFIG" \
  oidc-webhook-authenticator \
  ./charts/oidc-webhook-authenticator 

echo "OIDC Webhook Authenticator installed successfully in the virtual cluster."

echo "Generating kubeconfig for the OIDC Webhook Authenticator to access the hosting cluster."

cluster_ca_data=$(kubectl config view --kubeconfig "$VIRTUAL_KUBECONFIG" --raw -o jsonpath='{.clusters[0].cluster.certificate-authority-data}')
cluster_server=$(kubectl config view --kubeconfig "$VIRTUAL_KUBECONFIG" --raw -o jsonpath='{.clusters[0].cluster.server}')
token=$(kubectl -n garden create token oidc-webhook-authenticator --duration 48h --kubeconfig "$VIRTUAL_KUBECONFIG")

tmpfile=$(mktemp)

cat <<EOF > "$tmpfile"
---
apiVersion: v1
kind: Config
current-context: cluster
contexts:
- context:
    cluster: cluster
    user: oidc-webhook-authenticator
  name: cluster
clusters:
- cluster:
    certificate-authority-data: $cluster_ca_data
    server: $cluster_server
  name: cluster
users:
- name: oidc-webhook-authenticator
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
  .runtime.webhookConfig.tls.crt = load_str("'"$cert_dir/tls.crt"'") 
  | (.runtime.webhookConfig.tls.crt style="literal")
' "$values_file"
yq -i '
  .runtime.webhookConfig.tls.key = load_str("'"$cert_dir/tls.key"'") 
  | (.runtime.webhookConfig.tls.key style="literal")
' "$values_file"

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

helm upgrade \
  --install \
  --wait \
  --history-max=4 \
  --values "$values_file" \
  --set application.enabled="false" \
  --set application.virtualGarden.enabled="false" \
  --set runtime.enabled="true" \
  --namespace garden \
  --kubeconfig "$KIND_KUBECONFIG" \
  oidc-webhook-authenticator \
  ./charts/oidc-webhook-authenticator 

# Patch garden local to point to OWA
kubectl patch garden local \
  --kubeconfig "$KIND_KUBECONFIG" \
  --type='merge' \
  -p '{
    "spec": {
      "virtualCluster": {
        "kubernetes": {
          "kubeAPIServer": {
            "authentication": {
              "webhook": {
                "kubeconfigSecretName": "oidc-webhook-authenticator-kubeconfig",
                "cacheTTL": "0s"
              }
            }
          }
        }
      }
    }
  }'

echo "OIDC Webhook Authenticator installed successfully in the hosting kind cluster."

echo "Done."
