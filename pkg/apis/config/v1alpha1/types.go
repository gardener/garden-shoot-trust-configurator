// SPDX-FileCopyrightText: SAP SE or an SAP affiliate company and Gardener contributors
//
// SPDX-License-Identifier: Apache-2.0

package v1alpha1

import (
	"time"

	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	componentbaseconfigv1alpha1 "k8s.io/component-base/config/v1alpha1"
)

const (
	// DefaultAudience is the default audience used in the OIDC resources for trusted shoots.
	DefaultAudience = "garden"
	// DefaultMaxTokenExpiration is the default maximum token expiration duration (2 hours).
	DefaultMaxTokenExpiration = 2 * time.Hour
	// DefaultLockObjectNamespace is the default lock namespace for leader election.
	DefaultLockObjectNamespace = "kube-system"
	// DefaultLockObjectName is the default lock name for leader election.
	DefaultLockObjectName = "garden-shoot-trust-configurator-leader-election"
	// DefaultVolumeMountPathCertificates is the default directory for the webhook server TLS certificate and key.
	DefaultVolumeMountPathCertificates = "/etc/garden-shoot-trust-configurator/webhooks/tls"
)

// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object

// GardenShootTrustConfiguratorConfiguration defines the configuration for the Gardener garden-shoot-trust-configurator.
type GardenShootTrustConfiguratorConfiguration struct {
	metav1.TypeMeta `json:",inline"`

	// LeaderElection defines the configuration of leader election client.
	// +optional
	LeaderElection *componentbaseconfigv1alpha1.LeaderElectionConfiguration `json:"leaderElection,omitempty"`
	// LogLevel is the level/severity for the logs. Must be one of [info,debug,error].
	LogLevel string `json:"logLevel"`
	// LogFormat is the output format for the logs. Must be one of [text,json].
	LogFormat string `json:"logFormat"`
	// Controllers defines the configuration of the controllers.
	Controllers ControllerConfiguration `json:"controllers"`
	// Server defines the configuration of the HTTP server.
	Server ServerConfiguration `json:"server"`
}

// ControllerConfiguration defines the configuration of the controllers.
type ControllerConfiguration struct {
	// Shoot is the configuration for the shoot controller.
	Shoot ShootControllerConfig `json:"shoot"`
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

// ShootControllerConfig is the configuration for the shoot controller.
type ShootControllerConfig struct {
	// SyncPeriod is the duration how often the controller performs its reconciliation.
	// +optional
	SyncPeriod *metav1.Duration `json:"syncPeriod,omitempty"`
	// OIDCConfig is the configuration for the OIDC resources which are created for trusted shoots.
	// +optional
	OIDCConfig *OIDCConfig `json:"oidcConfig,omitempty"`
}

// OIDCConfig is the configuration for the OIDC resources created for trusted shoots.
type OIDCConfig struct {
	// Audiences is the list of audience identifiers used in the OIDC resources for trusted shoots.
	// Defaults to ["garden"].
	// +optional
	Audiences []string `json:"audiences,omitempty"`
	// MaxTokenExpiration sets a limit to the maximum validity duration of a token.
	// Tokens issued with validity greater than this value will not be verified.
	// Must be between 5 minutes and 24 hours. Defaults to 2 hours.
	// +optional
	MaxTokenExpiration *metav1.Duration `json:"maxTokenExpiration,omitempty"`
}

// ServerConfiguration contains details for the HTTP(S) servers.
type ServerConfiguration struct {
	// Webhooks is the configuration for the HTTPS webhook server.
	Webhooks HTTPSServer `json:"webhooks"`
	// HealthProbes is the configuration for serving the healthz and readyz endpoints.
	// +optional
	HealthProbes *Server `json:"healthProbes,omitempty"`
	// Metrics is the configuration for serving the metrics endpoint.
	// +optional
	Metrics *Server `json:"metrics,omitempty"`
}

// Server contains information for HTTP(S) server configuration.
type Server struct {
	// Port is the port on which to serve requests.
	Port int `json:"port"`
	// BindAddress is the IP address on which to listen for the specified port.
	BindAddress string `json:"bindAddress"`
}

// HTTPSServer is the configuration for the HTTPSServer server.
type HTTPSServer struct {
	// Server is the configuration for the bind address and the port.
	Server `json:",inline"`

	// TLS contains information about the TLS configuration for a HTTPS server.
	TLS TLS `json:"tls"`
}

// TLS contains information about the TLS configuration for a HTTPS server.
type TLS struct {
	// ServerCertDir is the path to a directory containing the server's TLS certificate and key (the files must be
	// named tls.crt and tls.key respectively).
	ServerCertDir string `json:"serverCertDir"`
}
