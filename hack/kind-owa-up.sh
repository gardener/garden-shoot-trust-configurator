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

repo_root="$(readlink -f $(dirname ${0})/..)"

OIDC_WEBHOOK_AUTH_NAME="oidc-webhook-authenticator"
OIDC_WEBHOOK_AUTH_REPO_NAME=github.com/gardener/$OIDC_WEBHOOK_AUTH_NAME

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
cert_dir="$repo_root/dev/$OIDC_WEBHOOK_AUTH_NAME/cfssl"

ca_key="$cert_dir/ca.key"
ca_crt="$cert_dir/ca.crt"
tls_key="$cert_dir/tls.key"
tls_csr="$cert_dir/tls.csr"
tls_crt="$cert_dir/tls.crt"

if [[ ! -s "$ca_crt" || ! -s "$ca_key" ]]; then
    echo "No CA found. Generating new CA key and certificate."

    openssl genrsa -out "$ca_key" 3072

    openssl req -x509 -new -nodes \
        -key "$ca_key" \
        -sha256 \
        -days 3650 \
        -out "$ca_crt" \
        -subj "/CN=webhook-ca" \
        -addext "basicConstraints=CA:TRUE" \
        -addext "keyUsage=keyCertSign,cRLSign" \
        -addext "subjectKeyIdentifier=hash"
fi

if [[ -s "$ca_key" && -s "$ca_crt" ]]; then
    if openssl x509 -checkend 86400 -in "$ca_crt" >/dev/null 2>&1; then
        echo "CA certificate is valid and will be reused."
    else
        echo "CA certificate has expired. Generating a new one."
        openssl req -x509 -new -nodes \
            -key "$ca_key" \
            -sha256 \
            -days 3650 \
            -out "$ca_crt" \
            -subj "/CN=webhook-ca" \
            -addext "basicConstraints=CA:TRUE" \
            -addext "keyUsage=keyCertSign,cRLSign" \
            -addext "subjectKeyIdentifier=hash"
    fi
fi

SANs="DNS:localhost,DNS:oidc-webhook-authenticator,DNS:oidc-webhook-authenticator.garden,DNS:oidc-webhook-authenticator.garden.svc,DNS:oidc-webhook-authenticator.garden.svc.cluster.local,IP:127.0.0.1"

if [[ -s "$tls_key" && -s "$tls_crt" ]]; then
  if openssl x509 -checkend 86400 -in "$tls_crt" >/dev/null 2>&1; then
    echo "Development TLS certificate is valid and will be reused."
    should_generate_cert=false
  else
    echo "Development TLS certificate has expired. Regenerating."
    should_generate_cert=true
  fi
else
  echo "No TLS cert found. Generating new one."
  should_generate_cert=true
fi

if [[ "$should_generate_cert" == true ]]; then
  openssl genrsa -out "$tls_key" 3072

  openssl req -new -key "$tls_key" -out "$tls_csr" \
    -subj "/CN=oidc-webhook-authenticator.garden.svc.cluster.local" \
    -addext "subjectAltName=$SANs"

  openssl x509 -req \
    -in "$tls_csr" \
    -CA "$ca_crt" -CAkey "$ca_key" \
    -out "$tls_crt" -days 365 -sha256 \
    -extfile <(printf "subjectAltName=%s" "$SANs")

  rm -f "$tls_csr"
  echo "Development TLS certificate generated successfully."
fi
# Finish generating certificates

charts_dir="$repo_root/dev/$OIDC_WEBHOOK_AUTH_NAME/charts/$OIDC_WEBHOOK_AUTH_NAME"
values_file="$charts_dir/values_$(date +%s).yaml"
cp "$charts_dir/values.yaml" "$values_file"

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

# Patch garden  local to point to OWA
kubectl patch garden local \
  --kubeconfig $KIND_KUBECONFIG \
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

echo "Cleaning up."
rm -rf $values_file

echo "Done."
