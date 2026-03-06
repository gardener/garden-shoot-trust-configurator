#!/bin/bash

# SPDX-FileCopyrightText: 2024 SAP SE or an SAP affiliate company and Gardener contributors
#
# SPDX-License-Identifier: Apache-2.0
# 

set -o errexit
set -o pipefail

repo_root="$(readlink -f "$(dirname "${0}")"/..)"
dev_trust_config_dir="$repo_root/dev/trust-configurator"
cert_dir="$dev_trust_config_dir/certs"

"$repo_root"/hack/generate-certs.sh \
  "$dev_trust_config_dir" \
  "garden-shoot-trust-configurator.garden.svc" \
  "DNS:localhost,DNS:garden-shoot-trust-configurator,DNS:garden-shoot-trust-configurator.garden,DNS:garden-shoot-trust-configurator.garden.svc,DNS:garden-shoot-trust-configurator.garden.svc.cluster.local,IP:127.0.0.1"

kubectl apply -f <(cat <<EOF
---
apiVersion: v1
kind: Secret
metadata:
  name: garden-shoot-trust-configurator-tls
  namespace: garden
type: Opaque
data:
  tls.key: $(base64 -w 0 "$cert_dir"/tls.key)
  tls.crt: $(base64 -w 0 "$cert_dir"/tls.crt)
EOF
)
