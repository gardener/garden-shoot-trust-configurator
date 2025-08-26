# Getting Started Locally

## Local KinD Setup with Gardener Operator and OIDC Webhook Authenticator Installed
This document will walk you through running a Gardener KinD cluster on your local machine, installing oidc-webhook-authenticator ([**OWA**](https://github.com/gardener/oidc-webhook-authenticator)) in it and running the garden-shoot-trust-configurator.
TODO(theoddora): Later it will contain the steps to install garden-shoot-trust-configurator.

### Prerequisites

- Docker
- KinD
- Go (>=1.24)
- Helm
- yq
- kubectl
- openssl


## 1. Set Up Local Gardener Operator (`gardener/gardener`)

Follow [Gardener's local setup guide](https://github.com/gardener/gardener/blob/master/docs/deployment/getting_started_locally.md#alternative-way-to-set-up-garden-and-seed-leveraging-gardener-operator).

Quick overview of the steps in `gardener/gardener`:
```bash
make kind-multi-zone-up operator-up
```

To create your local `Garden` and install `gardenlet` into the KinD cluster use this command:

```bash
make operator-seed-up
```

## 2. Install OIDC Webhook Authenticator

> [!NOTE]
> We will use kubeconfigs for the virtual garden and KinD clusters from the gardener repository.

```bash
gardener_repo_path=$(pwd)/../gardener # change this if needed

export VIRTUAL_KUBECONFIG=$gardener_repo_path/dev-setup/kubeconfigs/virtual-garden/kubeconfig

export KIND_KUBECONFIG=$gardener_repo_path/example/gardener-local/kind/multi-zone/kubeconfig
```

To install OWA resources in the locally set up garden cluster:
```bash
./hack/kind-owa-up.sh --virtual-kubeconfig $VIRTUAL_KUBECONFIG --kind-kubeconfig $KIND_KUBECONFIG
```

## Troubleshooting

- OWA setup uses short-lived service account tokens. If authentication fails, re-run the `kind-owa-up.sh` scripts.
- Ensure your kubeconfig files exist at the right location and are valid.