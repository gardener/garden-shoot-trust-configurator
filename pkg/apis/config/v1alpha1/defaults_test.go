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
	var obj *GardenShootTrustConfiguratorConfiguration

	BeforeEach(func() {
		obj = &GardenShootTrustConfiguratorConfiguration{}
	})

	Describe("#SetDefaults_GardenShootTrustConfiguratorConfiguration", func() {
		It("should default the log level and format", func() {
			SetDefaults_GardenShootTrustConfiguratorConfiguration(obj)

			Expect(obj.LogLevel).To(Equal(logger.InfoLevel))
			Expect(obj.LogFormat).To(Equal(logger.FormatJSON))
		})

		It("should not override existing values", func() {
			obj = &GardenShootTrustConfiguratorConfiguration{
				LogLevel:  "warning",
				LogFormat: "md",
			}
			SetDefaults_GardenShootTrustConfiguratorConfiguration(obj)

			Expect(obj.LogLevel).To(Equal("warning"))
			Expect(obj.LogFormat).To(Equal("md"))
		})
	})

	Describe("#SetDefaults_GarbageCollectorControllerConfig", func() {
		var obj *GarbageCollectorControllerConfig

		BeforeEach(func() {
			obj = &GarbageCollectorControllerConfig{}
		})

		It("should default the object", func() {
			SetDefaults_GarbageCollectorControllerConfig(obj)

			Expect(obj.SyncPeriod).To(PointTo(Equal(metav1.Duration{Duration: time.Hour})))
			Expect(obj.MinimumObjectLifetime).To(PointTo(Equal(metav1.Duration{Duration: 10 * time.Minute})))
		})

		It("should not overwrite existing values", func() {
			obj := &GarbageCollectorControllerConfig{
				SyncPeriod:            &metav1.Duration{Duration: time.Minute},
				MinimumObjectLifetime: &metav1.Duration{Duration: 5 * time.Minute},
			}

			SetDefaults_GarbageCollectorControllerConfig(obj)

			Expect(obj.SyncPeriod).To(PointTo(Equal(metav1.Duration{Duration: time.Minute})))
			Expect(obj.MinimumObjectLifetime).To(PointTo(Equal(metav1.Duration{Duration: 5 * time.Minute})))
		})
	})
})
