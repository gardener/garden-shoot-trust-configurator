// SPDX-FileCopyrightText: SAP SE or an SAP affiliate company and Gardener contributors
//
// SPDX-License-Identifier: Apache-2.0

package v1alpha1

import (
	"time"

	"github.com/gardener/gardener/pkg/logger"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
	componentbaseconfigv1alpha1 "k8s.io/component-base/config/v1alpha1"
)

func addDefaultingFuncs(scheme *runtime.Scheme) error {
	return RegisterDefaults(scheme)
}

// SetDefaults_GardenShootTrustConfiguratorConfiguration sets defaults for the configuration of the garden shoot trust configurator.
func SetDefaults_GardenShootTrustConfiguratorConfiguration(obj *GardenShootTrustConfiguratorConfiguration) {
	if obj.LogLevel == "" {
		obj.LogLevel = logger.InfoLevel
	}
	if obj.LogFormat == "" {
		obj.LogFormat = logger.FormatJSON
	}
	if obj.LeaderElection == nil {
		obj.LeaderElection = &componentbaseconfigv1alpha1.LeaderElectionConfiguration{}
	}
}

// SetDefaults_GarbageCollectorControllerConfig sets defaults for the GarbageCollectorControllerConfig object.
func SetDefaults_GarbageCollectorControllerConfig(obj *GarbageCollectorControllerConfig) {
	if obj.SyncPeriod == nil {
		obj.SyncPeriod = &metav1.Duration{Duration: time.Hour}
	}
	if obj.MinimumObjectLifetime == nil {
		obj.MinimumObjectLifetime = &metav1.Duration{Duration: 10 * time.Minute}
	}
}

// SetDefaults_ShootControllerConfig sets defaults for the ShootControllerConfig object.
func SetDefaults_ShootControllerConfig(obj *ShootControllerConfig) {
	if obj.SyncPeriod == nil {
		obj.SyncPeriod = &metav1.Duration{Duration: time.Hour}
	}
	if obj.OIDCConfig == nil {
		obj.OIDCConfig = &OIDCConfig{}
	}
}

// SetDefaults_OIDCConfig sets defaults for the OIDCConfig object.
func SetDefaults_OIDCConfig(obj *OIDCConfig) {
	if len(obj.Audiences) == 0 {
		obj.Audiences = []string{DefaultAudience}
	}
	if obj.MaxTokenExpiration == nil {
		obj.MaxTokenExpiration = &metav1.Duration{Duration: DefaultMaxTokenExpiration}
	}
}

// SetDefaults_ServerConfiguration sets defaults for the ServerConfiguration object.
func SetDefaults_ServerConfiguration(obj *ServerConfiguration) {
	if obj.HealthProbes == nil {
		obj.HealthProbes = &Server{}
	}
}

// SetDefaults_Server sets defaults for the Server object.
func SetDefaults_Server(obj *Server) {
	if obj.Port == 0 {
		obj.Port = 8081
	}
}

// SetDefaults_HTTPSServer sets defaults for the HTTPSServer object.
func SetDefaults_HTTPSServer(obj *HTTPSServer) {
	if obj.Port == 0 {
		obj.Port = 10443
	}
}

// SetDefaults_TLS sets defaults for the TLS object.
func SetDefaults_TLS(obj *TLS) {
	if obj.ServerCertDir == "" {
		obj.ServerCertDir = DefaultVolumeMountPathCertificates
	}
}

// SetDefaults_LeaderElectionConfiguration sets defaults for the LeaderElectionConfiguration object.
func SetDefaults_LeaderElectionConfiguration(obj *componentbaseconfigv1alpha1.LeaderElectionConfiguration) {
	if obj.ResourceLock == "" {
		obj.ResourceLock = "leases"
	}

	componentbaseconfigv1alpha1.RecommendedDefaultLeaderElectionConfiguration(obj)

	if obj.ResourceNamespace == "" {
		obj.ResourceNamespace = DefaultLockObjectNamespace
	}
	if obj.ResourceName == "" {
		obj.ResourceName = DefaultLockObjectName
	}
}
