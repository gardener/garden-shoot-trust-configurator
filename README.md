# garden-shoot-trust-configurator

[![REUSE status](https://api.reuse.software/badge/github.com/gardener/garden-shoot-trust-configurator)](https://api.reuse.software/info/github.com/gardener/garden-shoot-trust-configurator)


### Get your reuse badge

// TODO(theoddora): Get your reuse badge
To get your project reuse compliant you should register it [here](https://api.reuse.software/register) using your SAP email address. After confirming your email, an initial reuse check is done by the reuse API.


Enable shoot clusters with [`Managed Service Account Issuer`](https://gardener.cloud/docs/gardener/security/shoot_serviceaccounts/#Managed-Service-Account-Issuer) to be registered as trusted clusters in the Garden cluster. This reduces the need for manual service account token management and allows more secure, direct communication between shoots and the Garden cluster. This project is part of the [Gardener](https://gardener.cloud/) ecosystem for managing Kubernetes clusters.

## Development
As a prerequisite you need to have a Garden cluster up and running. Follow the [Gardener's local setup guide](https://github.com/gardener/gardener/blob/master/docs/deployment/getting_started_locally.md#alternative-way-to-set-up-garden-and-seed-leveraging-gardener-operator) which explains how to set up Gardener.

For local development, make sure to install oidc-webhook-authenticator. [Details are outlined here](docs/getting-started-locally.md)

Once the Garden cluster is up and running, export the `kubeconfig` pointing to the cluster as an environment variable.

```bash
gardener_repo_path=$(pwd)/../gardener # change this if needed

export KUBECONFIG=$gardener_repo_path/dev-setup/kubeconfigs/virtual-garden/kubeconfig
```

```bash
make start
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
