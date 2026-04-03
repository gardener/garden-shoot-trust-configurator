# garden-shoot-trust-configurator

[![REUSE status](https://api.reuse.software/badge/github.com/gardener/garden-shoot-trust-configurator)](https://api.reuse.software/info/github.com/gardener/garden-shoot-trust-configurator)
[![Build](https://github.com/gardener/garden-shoot-trust-configurator/actions/workflows/non-release.yaml/badge.svg)](https://github.com/gardener/garden-shoot-trust-configurator/actions/workflows/non-release.yaml)
[![Go Report Card](https://goreportcard.com/badge/github.com/gardener/garden-shoot-trust-configurator)](https://goreportcard.com/report/github.com/gardener/garden-shoot-trust-configurator)

Enable shoot clusters with [`Managed Service Account Issuer`](https://gardener.cloud/docs/gardener/security/shoot_serviceaccounts/#Managed-Service-Account-Issuer) to be registered as trusted clusters in the Garden cluster. This reduces the need for manual service account token management and allows more secure, direct communication between shoots and the Garden cluster. This project is part of the [Gardener](https://gardener.cloud/) ecosystem for managing Kubernetes clusters.

## Development
As a prerequisite you need to have a Garden cluster up and running. Follow the [Gardener's local setup guide](https://github.com/gardener/gardener/blob/master/docs/deployment/getting_started_locally.md#alternative-way-to-set-up-garden-and-seed-leveraging-gardener-operator) which explains how to set up Gardener.

Once the `Garden` cluster is up and running locally, we expect to have two KUBECONFIGs for:
- Virtual Garden cluster
- Runtime Garden cluster

```bash
GARDENER_REPO_ROOT=$(pwd)/../gardener # change this if needed

export KUBECONFIG_VIRTUAL=$GARDENER_REPO_ROOT/dev-setup/kubeconfigs/virtual-garden/kubeconfig
export KUBECONFIG_RUNTIME=$GARDENER_REPO_ROOT/dev-setup/kubeconfigs/runtime/kubeconfig
```

For local development, make sure to install the dependency `oidc-webhook-authenticator`, [more details are outlined here](docs/getting-started-locally.md).

Now start the `garden-shoot-trust-configurator`
```bash
make start
```

Alternatively you can deploy the trust-configurator in the local cluster with the following command.

```bash
make server-up
```

## Feedback and Support

Feedback and contributions are always welcome!

Please report bugs or suggestions as [GitHub issues](https://github.com/gardener/garden-shoot-trust-configurator/issues) or reach out on [Slack](https://gardener-cloud.slack.com/) (join the workspace [here](https://gardener.cloud/community)).

## Learn more

Please find further resources about our project here:

* [Our landing page gardener.cloud](https://gardener.cloud/)
* ["Gardener, the Kubernetes Botanist" blog on kubernetes.io](https://kubernetes.io/blog/2018/05/17/gardener/)
* ["Gardener Project Update" blog on kubernetes.io](https://kubernetes.io/blog/2019/12/02/gardener-project-update/)
* [Gardener Extensions Golang library](https://godoc.org/github.com/gardener/gardener/extensions/pkg)
* [GEP-1 (Gardener Enhancement Proposal) on extensibility](https://github.com/gardener/gardener/blob/master/docs/proposals/01-extensibility.md)
* [Extensibility API documentation](https://github.com/gardener/gardener/tree/master/docs/extensions)
