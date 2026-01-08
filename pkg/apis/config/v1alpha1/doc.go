// SPDX-FileCopyrightText: SAP SE or an SAP affiliate company and Gardener contributors
//
// SPDX-License-Identifier: Apache-2.0

// +k8s:deepcopy-gen=package
// +k8s:openapi-gen=true
// +k8s:defaulter-gen=TypeMeta

//go:generate gen-crd-api-reference-docs -api-dir . -config ../../../../hack/api-reference/config.json -template-dir $GARDENER_HACK_DIR/api-reference/template -out-file ../../../../docs/api-reference/config.md

// Package v1alpha1 contains the shoot trust confogurator configuration.
// +groupName=config.trust-configurator.gardener.cloud
package v1alpha1 // import "github.com/gardener/garden-shoot-trust-configurator/pkg/apis/config/v1alpha1"
