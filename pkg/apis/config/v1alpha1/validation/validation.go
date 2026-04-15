// SPDX-FileCopyrightText: SAP SE or an SAP affiliate company and Gardener contributors
//
// SPDX-License-Identifier: Apache-2.0

package validation

import (
	"time"

	"github.com/gardener/gardener/pkg/logger"
	validationutils "github.com/gardener/gardener/pkg/utils/validation"
	"k8s.io/apimachinery/pkg/util/sets"
	"k8s.io/apimachinery/pkg/util/validation/field"

	configv1alpha1 "github.com/gardener/garden-shoot-trust-configurator/pkg/apis/config/v1alpha1"
)

// ValidateGardenShootTrustConfiguratorConfiguration validates the given [*configv1alpha1.GardenShootTrustConfiguratorConfiguration].
func ValidateGardenShootTrustConfiguratorConfiguration(conf *configv1alpha1.GardenShootTrustConfiguratorConfiguration) field.ErrorList {
	allErrs := field.ErrorList{}

	if conf.LogLevel != "" {
		if !sets.New(logger.AllLogLevels...).Has(conf.LogLevel) {
			allErrs = append(allErrs, field.NotSupported(field.NewPath("logLevel"), conf.LogLevel, logger.AllLogLevels))
		}
	}

	if conf.LogFormat != "" {
		if !sets.New(logger.AllLogFormats...).Has(conf.LogFormat) {
			allErrs = append(allErrs, field.NotSupported(field.NewPath("logFormat"), conf.LogFormat, logger.AllLogFormats))
		}
	}

	allErrs = append(allErrs, validateControllers(&conf.Controllers, field.NewPath("controllers"))...)
	allErrs = append(allErrs, validationutils.ValidateLeaderElectionConfiguration(conf.LeaderElection, field.NewPath("leaderElection"))...)
	allErrs = append(allErrs, validateServerConfiguration(&conf.Server, field.NewPath("server"))...)

	return allErrs
}

// validateControllers validates the controllers configuration.
func validateControllers(controllers *configv1alpha1.ControllerConfiguration, fldPath *field.Path) field.ErrorList {
	allErrs := field.ErrorList{}

	allErrs = append(allErrs, validateShootControllerConfig(&controllers.Shoot, fldPath.Child("shoot"))...)
	allErrs = append(allErrs, validateGarbageCollectorControllerConfig(&controllers.GarbageCollector, fldPath.Child("garbageCollector"))...)

	return allErrs
}

// validateShootControllerConfig validates the shoot controller configuration.
func validateShootControllerConfig(config *configv1alpha1.ShootControllerConfig, fldPath *field.Path) field.ErrorList {
	allErrs := field.ErrorList{}

	if config.SyncPeriod != nil && config.SyncPeriod.Duration <= 0 {
		allErrs = append(allErrs, field.Invalid(fldPath.Child("syncPeriod"), config.SyncPeriod.Duration.String(), "must be positive"))
	}
	if config.OIDCConfig != nil {
		allErrs = append(allErrs, validateOIDCConfig(config.OIDCConfig, fldPath.Child("oidcConfig"))...)
	}

	return allErrs
}

// validateGarbageCollectorControllerConfig validates the garbage collector controller configuration.
func validateGarbageCollectorControllerConfig(config *configv1alpha1.GarbageCollectorControllerConfig, fldPath *field.Path) field.ErrorList {
	allErrs := field.ErrorList{}

	if config.SyncPeriod != nil && config.SyncPeriod.Duration <= 0 {
		allErrs = append(allErrs, field.Invalid(fldPath.Child("syncPeriod"), config.SyncPeriod.Duration.String(), "must be positive"))
	}
	if config.MinimumObjectLifetime != nil && config.MinimumObjectLifetime.Duration <= 0 {
		allErrs = append(allErrs, field.Invalid(fldPath.Child("minimumObjectLifetime"), config.MinimumObjectLifetime.Duration.String(), "must be positive"))
	}

	return allErrs
}

// validateOIDCConfig validates the OIDC configuration.
func validateOIDCConfig(config *configv1alpha1.OIDCConfig, fldPath *field.Path) field.ErrorList {
	allErrs := field.ErrorList{}

	if config.MaxTokenExpiration != nil {
		duration := config.MaxTokenExpiration.Duration
		if duration < 5*time.Minute {
			allErrs = append(allErrs, field.Forbidden(fldPath.Child("maxTokenExpiration"), "must be at least 5 minutes"))
		}
		if duration > 24*time.Hour {
			allErrs = append(allErrs, field.Forbidden(fldPath.Child("maxTokenExpiration"), "must not exceed 24 hours"))
		}
	}

	for i, audience := range config.Audiences {
		if audience == "" {
			allErrs = append(allErrs, field.Required(fldPath.Child("audiences").Index(i), "audience must not be empty"))
		}
	}

	return allErrs
}

// validateServerConfiguration validates the server configuration.
func validateServerConfiguration(config *configv1alpha1.ServerConfiguration, fldPath *field.Path) field.ErrorList {
	allErrs := field.ErrorList{}
	allErrs = append(allErrs, validatePortField(config.HealthProbes.Port, fldPath.Child("healthProbes", "port"))...)
	allErrs = append(allErrs, validatePortField(config.Webhooks.Port, fldPath.Child("webhooks", "port"))...)

	if config.Webhooks.TLS.ServerCertDir == "" {
		allErrs = append(allErrs, field.Required(fldPath.Child("webhooks", "tls", "serverCertDir"), "server certificate directory is required"))
	}

	return allErrs
}

// validatePortField validates that a port number is in the valid range [1, 65535].
func validatePortField(port int, fldPath *field.Path) field.ErrorList {
	if port == 0 {
		return field.ErrorList{field.Required(fldPath, "port is required")}
	}
	if port < 0 || port > 65535 {
		return field.ErrorList{field.Invalid(fldPath, port, "port must be between 1 and 65535")}
	}
	return nil
}
