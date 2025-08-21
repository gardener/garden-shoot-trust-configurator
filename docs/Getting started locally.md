
# Local KinD Setup with gardener-operator and OWA installed
This document will walk you through deploying Gardener on your local machine, installing OWA in it and running the shoot-trust-configurator.
TODO(theoddora): Later it will contain the steps to install garden-shoot-trust-configurator.

### Prerequisites
```yaml
- Docker
- KinD
- Go (>=1.24)
- Helm
- yq
- kubectl
- openssl
```

### 1. Set up local Gardener operator (`gardener/gardener`)

Follow [Gardener's local setup guide](https://github.com/gardener/gardener/blob/master/docs/deployment/getting_started_locally.md#alternative-way-to-set-up-garden-and-seed-leveraging-gardener-operator).

Quick overview of the steps in `gardener/gardener`:
```shell
make kind-multi-zone-up
```

To create your local `Garden` and install `gardenlet` into the KinD cluster use this command:

```shell
make operator-seed-up
```

### 2. **Install oidc-webhook-authenticator** in `virtual garden cluster` and `kind cluster`

**Note:** The kubeconfig for the virtual garden is at gardener repo ->  `./dev-setup/kubeconfigs/virtual-garden/kubeconfig`.

```shell
export VIRTUAL_KUBECONFIG={gardener-path}/dev-setup/kubeconfigs/virtual-garden/kubeconfig
```

**Note:** Gardener repo:

```shell
export KIND_KUBECONFIG={gardener-path}/example/gardener-local/kind/multi-zone/kubeconfig
```

To install OWA resources in the locally set up Gardener cluster:
```shell
./hack/kind-owa-up.sh --virtual-kubeconfig $VIRTUAL_KUBECONFIG --kind-kubeconfig $KIND_KUBECONFIG
```

### Troubleshooting

- OWA setup uses short-lived service account tokens. If authentication fails, re-run the `kind-owa-up.sh` scripts.
- Ensure your kubeconfig files exist at the right location and are valid.