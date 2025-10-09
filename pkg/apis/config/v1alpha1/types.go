// SPDX-FileCopyrightText: SAP SE or an SAP affiliate company and Gardener contributors
//
// SPDX-License-Identifier: Apache-2.0

package v1alpha1

import (
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object

// GardenShootTrustConfiguratorConfiguration defines the configuration for the Gardener garden-shoot-trust-configurator.
type GardenShootTrustConfiguratorConfiguration struct {
	metav1.TypeMeta `json:",inline"`

	// LogLevel is the level/severity for the logs. Must be one of [info,debug,error].
	LogLevel string `json:"logLevel"`
	// LogFormat is the output format for the logs. Must be one of [text,json].
	LogFormat string `json:"logFormat"`
	// Controllers defines the configuration of the controllers.
	Controllers ControllerConfiguration `json:"controllers"`
}

// ControllerConfiguration defines the configuration of the controllers.
type ControllerConfiguration struct {
	// GarbageCollector is the configuration for the garbage-collector controller.
	GarbageCollector GarbageCollectorControllerConfig `json:"garbageCollector"`
}

// GarbageCollectorControllerConfig is the configuration for the garbage-collector controller.
type GarbageCollectorControllerConfig struct {
	// SyncPeriod is the duration how often the controller performs its reconciliation.
	// +optional
	SyncPeriod *metav1.Duration `json:"syncPeriod,omitempty"`
	// MinimumObjectLifetime is the minimum age an object must have before it is considered for garbage collection.
	// +optional
	MinimumObjectLifetime *metav1.Duration `json:"minimumObjectLifetime,omitempty"`
}
