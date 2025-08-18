# garden-shoot-trust-configurator

[![reuse compliant](https://reuse.software/badge/reuse-compliant.svg)](https://reuse.software/)

## How to use this repository template

TODO(theoddora):
- Set the repository description in the "About" section of your repository
- Describe the new component in additional sections in this `README.md`
- Ask the [Owner of the gardener github organisation](https://github.com/orgs/gardener/people?query=role%3Aowner)
  - to double-check the initial content of this repository
  - to create the maintainer team for this new repository
  - to make this repository public
  - protect at least the master branch requiring mandatory code review by the maintainers defined in CODEOWNERS
  - grant admin permission to the maintainers team of the new repository defined in CODEOWNERS

## Maintain copyright and license information
By default all source code files are under `Apache 2.0` and all markdown files are under `Creative Commons` license.

When creating new source code files the license and copyright information should be provided using corresponding SPDX headers.

Example for go source code files (replace `<year>` with the current year)
```
/*
 * SPDX-FileCopyrightText: <year> SAP SE or an SAP affiliate company and Gardener contributors
 *
 * SPDX-License-Identifier: Apache-2.0
 */
```

### Third-party source code

If you copy third-party code into this repository or fork a repository, you must keep the license and copyright information (usually defined in the header of the file).

In addition you should adapt the `.reuse/dep5` file and assign the correct copyright and license information.

### Get your reuse badge

TODO(theoddora):
To get your project reuse compliant you should register it [here](https://api.reuse.software/register) using your SAP email address. After confirming your email, an inital reuse check is done by the reuse API.

To add the badge to your project's `README.md` file, use the snipped provided by the reuse API.


## Usage

TODO(theoddora)

## Feedback and Support

Feedback and contributions are always welcome!

Please report bugs or suggestions as [GitHub issues](https://github.com/gardener/garden-shoot-trust-configurator/issues) or reach out on [Slack](https://gardener-cloud.slack.com/) (join the workspace [here](https://gardener.cloud/community)).
## Learn more

Please find further resources about out project here:

* [Our landing page gardener.cloud](https://gardener.cloud/)
* ["Gardener, the Kubernetes Botanist" blog on kubernetes.io](https://kubernetes.io/blog/2018/05/17/gardener/)
* ["Gardener Project Update" blog on kubernetes.io](https://kubernetes.io/blog/2019/12/02/gardener-project-update/)
* [Gardener Extensions Golang library](https://godoc.org/github.com/gardener/gardener/extensions/pkg)
* [GEP-1 (Gardener Enhancement Proposal) on extensibility](https://github.com/gardener/gardener/blob/master/docs/proposals/01-extensibility.md)
* [Extensibility API documentation](https://github.com/gardener/gardener/tree/master/docs/extensions)
