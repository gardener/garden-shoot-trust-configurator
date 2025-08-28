// SPDX-FileCopyrightText: SAP SE or an SAP affiliate company and Gardener contributors
//
// SPDX-License-Identifier: Apache-2.0

//go:build tools
// +build tools

// This package imports things required by build scripts, to force `go mod` to see them as dependencies
package tools

import (
	_ "github.com/gardener/oidc-webhook-authenticator/cmd/oidc-webhook-authenticator"
	_ "github.com/ironcore-dev/vgopath"
	_ "go.uber.org/mock/mockgen"
)
