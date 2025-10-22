// SPDX-FileCopyrightText: SAP SE or an SAP affiliate company and Gardener contributors
//
// SPDX-License-Identifier: Apache-2.0

package validation

import (
	"time"

	"github.com/gardener/gardener/pkg/logger"
	"k8s.io/apimachinery/pkg/util/sets"
	"k8s.io/apimachinery/pkg/util/validation/field"

	"github.com/gardener/garden-shoot-trust-configurator/pkg/apis/config/v1alpha1"
)

// ValidateGardenShootTrustConfiguratorConfiguration validates the given `GardenShootTrustConfiguratorConfiguration`.
func ValidateGardenShootTrustConfiguratorConfiguration(conf *v1alpha1.GardenShootTrustConfiguratorConfiguration) field.ErrorList {
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

	return allErrs
}

// validateControllers validates the controllers configuration.
func validateControllers(controllers *v1alpha1.ControllerConfiguration, fldPath *field.Path) field.ErrorList {
	allErrs := field.ErrorList{}

	if controllers.Shoot.OIDCConfig != nil {
		allErrs = append(allErrs, validateOIDCConfig(controllers.Shoot.OIDCConfig, fldPath.Child("shoot", "oidcConfig"))...)
	}

	return allErrs
}

// validateOIDCConfig validates the OIDC configuration.
func validateOIDCConfig(config *v1alpha1.OIDCConfig, fldPath *field.Path) field.ErrorList {
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

	return allErrs
}
