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
					TLS: v1alpha1.TLSServer{
						ServerCertDir: "/etc/garden-shoot-trust-configurator/tls",
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
		errorList := ValidateGardenShootTrustConfiguratorConfiguration(conf)
		Expect(errorList).To(BeEmpty())
	})

	It("should pass validation when OIDCConfig is nil", func() {
		conf.Controllers.Shoot.OIDCConfig = nil

		errorList := ValidateGardenShootTrustConfiguratorConfiguration(conf)
		Expect(errorList).To(BeEmpty())
	})

	It("should pass validation when LeaderElectionConfiguration is nil", func() {
		conf.LeaderElection = nil

		errorList := ValidateGardenShootTrustConfiguratorConfiguration(conf)
		Expect(errorList).To(BeEmpty())
	})

	DescribeTable("MaxTokenExpiration",
		func(maxTokenExpiration *metav1.Duration, matcher gomegatypes.GomegaMatcher) {
			conf.Controllers.Shoot.OIDCConfig.MaxTokenExpiration = maxTokenExpiration

			errs := ValidateGardenShootTrustConfiguratorConfiguration(conf)
			Expect(errs).To(matcher)
		},
		Entry("should allow value between min and max",
			&metav1.Duration{Duration: 2 * time.Hour},
			BeEmpty(),
		),
		Entry("should allow exactly 5 minutes (minimum)",
			&metav1.Duration{Duration: 5 * time.Minute},
			BeEmpty(),
		),
		Entry("should allow exactly 24 hours (maximum)",
			&metav1.Duration{Duration: 24 * time.Hour},
			BeEmpty(),
		),
		Entry("should allow nil value",
			nil,
			BeEmpty(),
		),
		Entry("should forbid value less than 5 minutes",
			&metav1.Duration{Duration: 2 * time.Minute},
			ConsistOf(PointTo(
				MatchFields(IgnoreExtras, Fields{
					"Type":   Equal(field.ErrorTypeForbidden),
					"Field":  Equal("controllers.shoot.oidcConfig.maxTokenExpiration"),
					"Detail": ContainSubstring("must be at least 5 minutes"),
				}),
			)),
		),
		Entry("should forbid value greater than 24 hours",
			&metav1.Duration{Duration: 25 * time.Hour},
			ConsistOf(PointTo(
				MatchFields(IgnoreExtras, Fields{
					"Type":   Equal(field.ErrorTypeForbidden),
					"Field":  Equal("controllers.shoot.oidcConfig.maxTokenExpiration"),
					"Detail": ContainSubstring("must not exceed 24 hours"),
				}),
			)),
		),
	)

	Describe("ServerConfiguration", func() {
		It("should forbid negative HealthProbes port", func() {
			conf.Server.HealthProbes.Port = -1

			errs := ValidateGardenShootTrustConfiguratorConfiguration(conf)
			Expect(errs).To(ConsistOf(PointTo(
				MatchFields(IgnoreExtras, Fields{
					"Type":  Equal(field.ErrorTypeInvalid),
					"Field": Equal("server.healthProbes.port"),
				}),
			)))
		})

		It("should forbid negative Webhooks port", func() {
			conf.Server.Webhooks.Port = -1

			errs := ValidateGardenShootTrustConfiguratorConfiguration(conf)
			Expect(errs).To(ConsistOf(PointTo(
				MatchFields(IgnoreExtras, Fields{
					"Type":  Equal(field.ErrorTypeInvalid),
					"Field": Equal("server.webhooks.port"),
				}),
			)))
		})

		It("should forbid empty Webhooks TLS ServerCertDir", func() {
			conf.Server.Webhooks.TLS.ServerCertDir = ""

			errs := ValidateGardenShootTrustConfiguratorConfiguration(conf)
			Expect(errs).To(ConsistOf(PointTo(
				MatchFields(IgnoreExtras, Fields{
					"Type":  Equal(field.ErrorTypeRequired),
					"Field": Equal("server.webhooks.tls.serverCertDir"),
				}),
			)))
		})
	})
})
