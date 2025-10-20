// SPDX-FileCopyrightText: SAP SE or an SAP affiliate company and Gardener contributors
//
// SPDX-License-Identifier: Apache-2.0

package v1alpha1

import (
	"time"

	"github.com/gardener/gardener/pkg/logger"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
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
}

// SetDefaults_OIDCConfig sets defaults for the OIDCConfig object.
func SetDefaults_OIDCConfig(obj *OIDCConfig) {
	if len(obj.Audiences) == 0 {
		obj.Audiences = []string{DefaultOIDCClientID}
	}
}
