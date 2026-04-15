// SPDX-FileCopyrightText: SAP SE or an SAP affiliate company and Gardener contributors
//
// SPDX-License-Identifier: Apache-2.0

package validation_test

import (
	"time"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	. "github.com/onsi/gomega/gstruct"
	gomegatypes "github.com/onsi/gomega/types"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/util/validation/field"
	componentbaseconfigv1alpha1 "k8s.io/component-base/config/v1alpha1"
	"k8s.io/utils/ptr"

	"github.com/gardener/garden-shoot-trust-configurator/pkg/apis/config/v1alpha1"
	. "github.com/gardener/garden-shoot-trust-configurator/pkg/apis/config/v1alpha1/validation"
)

var _ = Describe("#ValidateGardenShootTrustConfiguratorConfiguration", func() {
	var conf *v1alpha1.GardenShootTrustConfiguratorConfiguration

	BeforeEach(func() {
		conf = &v1alpha1.GardenShootTrustConfiguratorConfiguration{
			LogLevel:  "info",
			LogFormat: "json",
			Controllers: v1alpha1.ControllerConfiguration{
				Shoot: v1alpha1.ShootControllerConfig{
					SyncPeriod: &metav1.Duration{Duration: time.Hour},
					OIDCConfig: &v1alpha1.OIDCConfig{
						Audiences:          []string{"garden"},
						MaxTokenExpiration: &metav1.Duration{Duration: 2 * time.Hour},
					},
				},
				GarbageCollector: v1alpha1.GarbageCollectorControllerConfig{
					SyncPeriod:            &metav1.Duration{Duration: time.Hour},
					MinimumObjectLifetime: &metav1.Duration{Duration: 10 * time.Minute},
				},
			},
			Server: v1alpha1.ServerConfiguration{
				Webhooks: v1alpha1.HTTPSServer{
					Server: v1alpha1.Server{
						Port: 10443,
					},
					TLS: v1alpha1.TLS{
						ServerCertDir: "/etc/garden-shoot-trust-configurator/webhooks/tls",
					},
				},
				HealthProbes: &v1alpha1.Server{
					Port: 8081,
				},
			},
			LeaderElection: &componentbaseconfigv1alpha1.LeaderElectionConfiguration{
				LeaderElect:       ptr.To(true),
				LeaseDuration:     metav1.Duration{Duration: 15 * time.Second},
				RenewDeadline:     metav1.Duration{Duration: 10 * time.Second},
				RetryPeriod:       metav1.Duration{Duration: 2 * time.Second},
				ResourceLock:      "configmapsleases",
				ResourceName:      "garden-shoot-trust-configurator-leader-election",
				ResourceNamespace: "garden",
			},
		}
	})

	It("should pass validation with valid configuration", func() {
		Expect(ValidateGardenShootTrustConfiguratorConfiguration(conf)).To(BeEmpty())
	})

	Describe("#LeaderElectionConfiguration", func() {
		BeforeEach(func() {
			v1alpha1.SetDefaults_LeaderElectionConfiguration(conf.LeaderElection)
		})

		It("should allow default leader election configuration with required fields", func() {
			Expect(ValidateGardenShootTrustConfiguratorConfiguration(conf)).To(BeEmpty())
		})

		It("should allow omitting leader election config", func() {
			conf.LeaderElection = nil

			Expect(ValidateGardenShootTrustConfiguratorConfiguration(conf)).To(BeEmpty())
		})

		It("should allow not enabling leader election", func() {
			conf.LeaderElection.LeaderElect = nil

			Expect(ValidateGardenShootTrustConfiguratorConfiguration(conf)).To(BeEmpty())
		})

		It("should allow disabling leader election", func() {
			conf.LeaderElection.LeaderElect = ptr.To(false)

			Expect(ValidateGardenShootTrustConfiguratorConfiguration(conf)).To(BeEmpty())
		})

		It("should reject leader election config with missing required fields", func() {
			conf.LeaderElection.ResourceNamespace = ""

			Expect(ValidateGardenShootTrustConfiguratorConfiguration(conf)).To(ConsistOf(
				PointTo(MatchFields(IgnoreExtras, Fields{
					"Type":  Equal(field.ErrorTypeRequired),
					"Field": Equal("leaderElection.resourceNamespace"),
				})),
			))
		})
	})

	Describe("#Controllers", func() {
		Describe("#ShootControllerConfig", func() {
			Context("syncPeriod", func() {
				It("should forbid a zero sync period", func() {
					conf.Controllers.Shoot.SyncPeriod = &metav1.Duration{Duration: 0}

					Expect(ValidateGardenShootTrustConfiguratorConfiguration(conf)).To(ConsistOf(PointTo(
						MatchFields(IgnoreExtras, Fields{
							"Type":  Equal(field.ErrorTypeInvalid),
							"Field": Equal("controllers.shoot.syncPeriod"),
						}),
					)))
				})

				It("should forbid a negative sync period", func() {
					conf.Controllers.Shoot.SyncPeriod = &metav1.Duration{Duration: -1 * time.Second}

					Expect(ValidateGardenShootTrustConfiguratorConfiguration(conf)).To(ConsistOf(PointTo(
						MatchFields(IgnoreExtras, Fields{
							"Type":  Equal(field.ErrorTypeInvalid),
							"Field": Equal("controllers.shoot.syncPeriod"),
						}),
					)))
				})
			})

			Describe("#OIDCConfig", func() {
				It("should pass validation when OIDCConfig is nil", func() {
					conf.Controllers.Shoot.OIDCConfig = nil

					Expect(ValidateGardenShootTrustConfiguratorConfiguration(conf)).To(BeEmpty())
				})

				DescribeTable("MaxTokenExpiration",
					func(maxTokenExpiration *metav1.Duration, matcher gomegatypes.GomegaMatcher) {
						conf.Controllers.Shoot.OIDCConfig.MaxTokenExpiration = maxTokenExpiration
						Expect(ValidateGardenShootTrustConfiguratorConfiguration(conf)).To(matcher)
					},
					Entry("should allow nil value", nil, BeEmpty()),
					Entry("should allow value between min and max", &metav1.Duration{Duration: 2 * time.Hour}, BeEmpty()),
					Entry("should allow exactly 5 minutes (minimum)", &metav1.Duration{Duration: 5 * time.Minute}, BeEmpty()),
					Entry("should allow exactly 24 hours (maximum)", &metav1.Duration{Duration: 24 * time.Hour}, BeEmpty()),
					Entry("should forbid value less than 5 minutes",
						&metav1.Duration{Duration: 2 * time.Minute},
						ConsistOf(PointTo(MatchFields(IgnoreExtras, Fields{
							"Type":   Equal(field.ErrorTypeForbidden),
							"Field":  Equal("controllers.shoot.oidcConfig.maxTokenExpiration"),
							"Detail": ContainSubstring("must be at least 5 minutes"),
						}))),
					),
					Entry("should forbid value greater than 24 hours",
						&metav1.Duration{Duration: 25 * time.Hour},
						ConsistOf(PointTo(MatchFields(IgnoreExtras, Fields{
							"Type":   Equal(field.ErrorTypeForbidden),
							"Field":  Equal("controllers.shoot.oidcConfig.maxTokenExpiration"),
							"Detail": ContainSubstring("must not exceed 24 hours"),
						}))),
					),
				)

				It("should forbid empty string in audiences", func() {
					conf.Controllers.Shoot.OIDCConfig.Audiences = []string{"garden", ""}

					Expect(ValidateGardenShootTrustConfiguratorConfiguration(conf)).To(ConsistOf(PointTo(
						MatchFields(IgnoreExtras, Fields{
							"Type":  Equal(field.ErrorTypeRequired),
							"Field": Equal("controllers.shoot.oidcConfig.audiences[1]"),
						}),
					)))
				})
			})
		})

		Describe("#GarbageCollectorControllerConfig", func() {
			It("should forbid a zero sync period", func() {
				conf.Controllers.GarbageCollector.SyncPeriod = &metav1.Duration{Duration: 0}

				Expect(ValidateGardenShootTrustConfiguratorConfiguration(conf)).To(ConsistOf(PointTo(
					MatchFields(IgnoreExtras, Fields{
						"Type":  Equal(field.ErrorTypeInvalid),
						"Field": Equal("controllers.garbageCollector.syncPeriod"),
					}),
				)))
			})

			It("should forbid a negative sync period", func() {
				conf.Controllers.GarbageCollector.SyncPeriod = &metav1.Duration{Duration: -1 * time.Second}

				Expect(ValidateGardenShootTrustConfiguratorConfiguration(conf)).To(ConsistOf(PointTo(
					MatchFields(IgnoreExtras, Fields{
						"Type":  Equal(field.ErrorTypeInvalid),
						"Field": Equal("controllers.garbageCollector.syncPeriod"),
					}),
				)))
			})

			It("should forbid a zero minimum object lifetime", func() {
				conf.Controllers.GarbageCollector.MinimumObjectLifetime = &metav1.Duration{Duration: 0}

				Expect(ValidateGardenShootTrustConfiguratorConfiguration(conf)).To(ConsistOf(PointTo(
					MatchFields(IgnoreExtras, Fields{
						"Type":  Equal(field.ErrorTypeInvalid),
						"Field": Equal("controllers.garbageCollector.minimumObjectLifetime"),
					}),
				)))
			})

			It("should forbid a negative minimum object lifetime", func() {
				conf.Controllers.GarbageCollector.MinimumObjectLifetime = &metav1.Duration{Duration: -1 * time.Second}

				Expect(ValidateGardenShootTrustConfiguratorConfiguration(conf)).To(ConsistOf(PointTo(
					MatchFields(IgnoreExtras, Fields{
						"Type":  Equal(field.ErrorTypeInvalid),
						"Field": Equal("controllers.garbageCollector.minimumObjectLifetime"),
					}),
				)))
			})
		})
	})

	Describe("#ServerConfiguration", func() {
		Context("health probes port", func() {
			It("should allow a valid port", func() {
				conf.Server.HealthProbes.Port = 65535

				Expect(ValidateGardenShootTrustConfiguratorConfiguration(conf)).To(BeEmpty())
			})

			It("should return an error when health port is misconfigured", func() {
				conf.Server.HealthProbes.Port = 0

				errs := ValidateGardenShootTrustConfiguratorConfiguration(conf)
				Expect(errs).To(ConsistOf(PointTo(MatchFields(IgnoreExtras, Fields{
					"Type":   Equal(field.ErrorTypeRequired),
					"Field":  Equal("server.healthProbes.port"),
					"Detail": ContainSubstring("port is required"),
				}))))
			})

			It("should return an error when port is invalid", func() {
				conf.Server.HealthProbes.Port = -1

				errs := ValidateGardenShootTrustConfiguratorConfiguration(conf)
				Expect(errs).To(ConsistOf(PointTo(MatchFields(IgnoreExtras, Fields{
					"Type":   Equal(field.ErrorTypeInvalid),
					"Field":  Equal("server.healthProbes.port"),
					"Detail": ContainSubstring("port must be between 1 and 65535"),
				}))))
			})

			It("should return an error when port exceeds 65535", func() {
				conf.Server.HealthProbes.Port = 70000

				errs := ValidateGardenShootTrustConfiguratorConfiguration(conf)
				Expect(errs).To(ConsistOf(PointTo(MatchFields(IgnoreExtras, Fields{
					"Type":   Equal(field.ErrorTypeInvalid),
					"Field":  Equal("server.healthProbes.port"),
					"Detail": ContainSubstring("port must be between 1 and 65535"),
				}))))
			})
		})

		Context("webhooks port", func() {
			It("should allow a valid webhooks port", func() {
				conf.Server.Webhooks.Port = 65535

				Expect(ValidateGardenShootTrustConfiguratorConfiguration(conf)).To(BeEmpty())
			})

			It("should return an error when webhooks port is misconfigured", func() {
				conf.Server.Webhooks.Port = 0

				errs := ValidateGardenShootTrustConfiguratorConfiguration(conf)
				Expect(errs).To(ConsistOf(PointTo(MatchFields(IgnoreExtras, Fields{
					"Type":   Equal(field.ErrorTypeRequired),
					"Field":  Equal("server.webhooks.port"),
					"Detail": ContainSubstring("port is required"),
				}))))
			})

			It("should return an error when webhooks port exceeds 65535", func() {
				conf.Server.Webhooks.Port = 70000

				errs := ValidateGardenShootTrustConfiguratorConfiguration(conf)
				Expect(errs).To(ConsistOf(PointTo(MatchFields(IgnoreExtras, Fields{
					"Type":   Equal(field.ErrorTypeInvalid),
					"Field":  Equal("server.webhooks.port"),
					"Detail": ContainSubstring("port must be between 1 and 65535"),
				}))))
			})

			It("should return an error when webhooks port is invalid", func() {
				conf.Server.Webhooks.Port = -1

				errs := ValidateGardenShootTrustConfiguratorConfiguration(conf)
				Expect(errs).To(ConsistOf(PointTo(MatchFields(IgnoreExtras, Fields{
					"Type":   Equal(field.ErrorTypeInvalid),
					"Field":  Equal("server.webhooks.port"),
					"Detail": ContainSubstring("port must be between 1 and 65535"),
				}))))
			})
		})

		Context("webhooks TLS", func() {
			It("should forbid empty Webhooks TLS ServerCertDir", func() {
				conf.Server.Webhooks.TLS.ServerCertDir = ""

				Expect(ValidateGardenShootTrustConfiguratorConfiguration(conf)).To(ConsistOf(PointTo(
					MatchFields(IgnoreExtras, Fields{
						"Type":  Equal(field.ErrorTypeRequired),
						"Field": Equal("server.webhooks.tls.serverCertDir"),
					}),
				)))
			})
		})
	})
})
