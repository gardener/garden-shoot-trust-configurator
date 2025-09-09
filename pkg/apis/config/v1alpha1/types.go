// SPDX-FileCopyrightText: SAP SE or an SAP affiliate company and Gardener contributors
//
// SPDX-License-Identifier: Apache-2.0

package v1alpha1

import (
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

const (
	// LogLevelDebug is the debug log level, i.e. the most verbose.
	LogLevelDebug = "debug"
	// LogLevelInfo is the default log level.
	LogLevelInfo = "info"
	// LogLevelError is a log level where only errors are logged.
	LogLevelError = "error"

	// LogFormatJSON is the output type that produces a JSON object per log line.
	LogFormatJSON = "json"
	// LogFormatText outputs the log as human-readable text.
	LogFormatText = "text"
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
	// ShootController is the configuration for the Shoot controller.
	ShootController ShootControllerConfig `json:"shoot"`
}

// ShootControllerConfig is the configuration for the Shoot controller.
type ShootControllerConfig struct {
	// ResyncPeriod is the duration how often the controller performs its reconciliation.
	// +optional
	ResyncPeriod *metav1.Duration `json:"resyncPeriod,omitempty"`
}
