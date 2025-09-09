// SPDX-FileCopyrightText: SAP SE or an SAP affiliate company and Gardener contributors
//
// SPDX-License-Identifier: Apache-2.0

package v1alpha1_test

import (
	"github.com/gardener/gardener/pkg/logger"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"

	. "github.com/gardener/garden-shoot-trust-configurator/pkg/apis/config/v1alpha1"
)

var _ = Describe("Defaults", func() {
	var obj *GardenShootTrustConfiguratorConfiguration

	BeforeEach(func() {
		obj = &GardenShootTrustConfiguratorConfiguration{}
	})

	Describe("GardenShootTrustConfiguratorConfiguration defaulting", func() {
		It("should default GardenShootTrustConfiguratorConfiguration correctly", func() {
			SetDefaults_GardenShootTrustConfiguratorConfiguration(obj)

			Expect(obj.LogLevel).To(Equal(logger.InfoLevel))
			Expect(obj.LogFormat).To(Equal(logger.FormatJSON))
		})

		It("should not default fields that are set", func() {
			obj = &GardenShootTrustConfiguratorConfiguration{
				LogLevel:  "warning",
				LogFormat: "md",
			}
			SetDefaults_GardenShootTrustConfiguratorConfiguration(obj)

			Expect(obj.LogLevel).To(Equal("warning"))
			Expect(obj.LogFormat).To(Equal("md"))
		})
	})

})

var _ = Describe("Constants", func() {
	It("should have the same values as the corresponding constants in the logger package", func() {
		Expect(LogLevelDebug).To(Equal(logger.DebugLevel))
		Expect(LogLevelInfo).To(Equal(logger.InfoLevel))
		Expect(LogLevelError).To(Equal(logger.ErrorLevel))
		Expect(LogFormatJSON).To(Equal(logger.FormatJSON))
		Expect(LogFormatText).To(Equal(logger.FormatText))
	})
})
