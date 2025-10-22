// SPDX-FileCopyrightText: SAP SE or an SAP affiliate company and Gardener contributors
//
// SPDX-License-Identifier: Apache-2.0

package v1alpha1_test

import (
	"time"

	"github.com/gardener/gardener/pkg/logger"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	. "github.com/onsi/gomega/gstruct"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"

	. "github.com/gardener/garden-shoot-trust-configurator/pkg/apis/config/v1alpha1"
)

var _ = Describe("Defaults", func() {
	Describe("#SetDefaults_GardenShootTrustConfiguratorConfiguration", func() {
		var obj *GardenShootTrustConfiguratorConfiguration

		BeforeEach(func() {
			obj = &GardenShootTrustConfiguratorConfiguration{}
		})

		Context("LogLevel", func() {
			It("should default log level", func() {
				SetDefaults_GardenShootTrustConfiguratorConfiguration(obj)

				Expect(obj.LogLevel).To(Equal(logger.InfoLevel))
			})

			It("should not overwrite already set value for log level", func() {
				obj.LogLevel = "warning"

				SetDefaults_GardenShootTrustConfiguratorConfiguration(obj)

				Expect(obj.LogLevel).To(Equal("warning"))
			})
		})

		Context("LogFormat", func() {
			It("should default log format", func() {
				SetDefaults_GardenShootTrustConfiguratorConfiguration(obj)

				Expect(obj.LogFormat).To(Equal(logger.FormatJSON))
			})

			It("should not overwrite already set value for log format", func() {
				obj.LogFormat = "md"

				SetDefaults_GardenShootTrustConfiguratorConfiguration(obj)

				Expect(obj.LogFormat).To(Equal("md"))
			})
		})
	})

	Describe("#SetDefaults_GarbageCollectorControllerConfig", func() {
		var obj *GarbageCollectorControllerConfig

		BeforeEach(func() {
			obj = &GarbageCollectorControllerConfig{}
		})

		Context("SyncPeriod", func() {
			It("should default sync period", func() {
				SetDefaults_GarbageCollectorControllerConfig(obj)

				Expect(obj.SyncPeriod).To(PointTo(Equal(metav1.Duration{Duration: time.Hour})))
			})

			It("should not overwrite already set value for sync period", func() {
				obj.SyncPeriod = &metav1.Duration{Duration: time.Minute}

				SetDefaults_GarbageCollectorControllerConfig(obj)

				Expect(obj.SyncPeriod).To(PointTo(Equal(metav1.Duration{Duration: time.Minute})))
			})
		})

		Context("MinimumObjectLifetime", func() {
			It("should default minimum object lifetime", func() {
				SetDefaults_GarbageCollectorControllerConfig(obj)

				Expect(obj.MinimumObjectLifetime).To(PointTo(Equal(metav1.Duration{Duration: 10 * time.Minute})))
			})

			It("should not overwrite already set value for minimum object lifetime", func() {
				obj.MinimumObjectLifetime = &metav1.Duration{Duration: 5 * time.Minute}

				SetDefaults_GarbageCollectorControllerConfig(obj)

				Expect(obj.MinimumObjectLifetime).To(PointTo(Equal(metav1.Duration{Duration: 5 * time.Minute})))
			})
		})
	})

	Describe("#SetDefaults_ShootControllerConfig", func() {
		var obj *ShootControllerConfig

		BeforeEach(func() {
			obj = &ShootControllerConfig{}
		})

		Context("SyncPeriod", func() {
			It("should default sync period", func() {
				SetDefaults_ShootControllerConfig(obj)

				Expect(obj.SyncPeriod).To(PointTo(Equal(metav1.Duration{Duration: time.Hour})))
			})

			It("should not overwrite already set value for sync period", func() {
				obj.SyncPeriod = &metav1.Duration{Duration: time.Minute}

				SetDefaults_ShootControllerConfig(obj)

				Expect(obj.SyncPeriod).To(PointTo(Equal(metav1.Duration{Duration: time.Minute})))
			})
		})

		Context("OIDCConfig", func() {
			It("should initialize OIDC config when nil", func() {
				SetDefaults_ShootControllerConfig(obj)

				Expect(obj.OIDCConfig).NotTo(BeNil())
			})

			It("should not overwrite already set OIDC config", func() {
				existingConfig := &OIDCConfig{
					Audiences:          []string{"custom-audience"},
					MaxTokenExpiration: &metav1.Duration{Duration: 1 * time.Hour},
				}
				obj.OIDCConfig = existingConfig

				SetDefaults_ShootControllerConfig(obj)

				Expect(obj.OIDCConfig).To(Equal(existingConfig))
			})
		})
	})

	Describe("#SetDefaults_OIDCConfig", func() {
		var obj *OIDCConfig

		BeforeEach(func() {
			obj = &OIDCConfig{}
		})

		Context("Audiences", func() {
			It("should default audiences", func() {
				SetDefaults_OIDCConfig(obj)

				Expect(obj.Audiences).To(Equal([]string{"garden"}))
			})

			It("should not overwrite already set value for audiences", func() {
				obj.Audiences = []string{"custom-audience"}

				SetDefaults_OIDCConfig(obj)

				Expect(obj.Audiences).To(Equal([]string{"custom-audience"}))
			})
		})

		Context("MaxTokenExpiration", func() {
			It("should default max token expiration", func() {
				SetDefaults_OIDCConfig(obj)

				Expect(obj.MaxTokenExpiration).To(PointTo(Equal(metav1.Duration{Duration: 2 * time.Hour})))
			})

			It("should not overwrite already set value for max token expiration", func() {
				obj.MaxTokenExpiration = &metav1.Duration{Duration: 1 * time.Hour}

				SetDefaults_OIDCConfig(obj)

				Expect(obj.MaxTokenExpiration).To(PointTo(Equal(metav1.Duration{Duration: 1 * time.Hour})))
			})
		})
	})
})
