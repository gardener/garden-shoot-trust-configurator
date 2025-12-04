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
	componentbaseconfigv1alpha1 "k8s.io/component-base/config/v1alpha1"
	"k8s.io/utils/ptr"

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

		Context("LeaderElection", func() {
			It("should initialize LeaderElection when nil", func() {
				SetDefaults_GardenShootTrustConfiguratorConfiguration(obj)

				Expect(obj.LeaderElection).NotTo(BeNil())
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

	Describe("#SetDefaults_ServerConfiguration", func() {
		var obj *ServerConfiguration

		BeforeEach(func() {
			obj = &ServerConfiguration{}
		})

		Context("HealthProbes", func() {
			It("should default HealthProbes when nil", func() {
				SetDefaults_ServerConfiguration(obj)

				Expect(obj.HealthProbes).NotTo(BeNil())
			})
		})
	})

	Describe("#SetDefaults_Server", func() {
		var obj *Server

		BeforeEach(func() {
			obj = &Server{}
		})

		Context("Port", func() {
			It("should default port", func() {
				SetDefaults_Server(obj)

				Expect(obj.Port).To(Equal(8081))
			})

			It("should not overwrite already set value for port", func() {
				obj.Port = 9090

				SetDefaults_Server(obj)

				Expect(obj.Port).To(Equal(9090))
			})
		})
	})

	Describe("#SetDefaults_HTTPSServer", func() {
		var obj *HTTPSServer

		BeforeEach(func() {
			obj = &HTTPSServer{}
		})

		Context("Port", func() {
			It("should default port", func() {
				SetDefaults_HTTPSServer(obj)

				Expect(obj.Port).To(Equal(10443))
			})

			It("should not overwrite already set value for port", func() {
				obj.Port = 9090

				SetDefaults_HTTPSServer(obj)

				Expect(obj.Port).To(Equal(9090))
			})
		})
	})

	Describe("#SetDefaults_TLSServer", func() {
		var obj *TLSServer

		BeforeEach(func() {
			obj = &TLSServer{}
		})

		Context("ServerCertDir", func() {
			It("should default server cert dir", func() {
				SetDefaults_TLSServer(obj)

				Expect(obj.ServerCertDir).To(Equal(DefaultVolumeMountPathCertificates))
			})

			It("should not overwrite already set value for server cert dir", func() {
				obj.ServerCertDir = "/custom/dir"

				SetDefaults_TLSServer(obj)

				Expect(obj.ServerCertDir).To(Equal("/custom/dir"))
			})
		})
	})

	Describe("#SetDefaults_LeaderElectionConfiguration", func() {
		var obj *componentbaseconfigv1alpha1.LeaderElectionConfiguration

		BeforeEach(func() {
			obj = &componentbaseconfigv1alpha1.LeaderElectionConfiguration{}
		})

		Context("should default to recommended leader election values", func() {
			It("should set default recommended leader election values", func() {
				SetDefaults_LeaderElectionConfiguration(obj)

				expectedLeaderElectionConfig := &componentbaseconfigv1alpha1.LeaderElectionConfiguration{
					LeaderElect:       ptr.To(true),
					LeaseDuration:     metav1.Duration{Duration: 15 * time.Second},
					RenewDeadline:     metav1.Duration{Duration: 10 * time.Second},
					RetryPeriod:       metav1.Duration{Duration: 2 * time.Second},
					ResourceLock:      "leases",
					ResourceName:      DefaultLockObjectName,
					ResourceNamespace: DefaultLockObjectNamespace,
				}
				Expect(obj).To(Equal(expectedLeaderElectionConfig))
			})

			It("should not overwrite already set values for leader election", func() {
				obj.LeaderElect = ptr.To(false)
				obj.LeaseDuration = metav1.Duration{Duration: 30 * time.Second}
				obj.RenewDeadline = metav1.Duration{Duration: 20 * time.Second}
				obj.RetryPeriod = metav1.Duration{Duration: 5 * time.Second}
				obj.ResourceLock = "lock"
				obj.ResourceName = "name"
				obj.ResourceNamespace = "namespace"

				SetDefaults_LeaderElectionConfiguration(obj)

				expectedLeaderElectionConfig := &componentbaseconfigv1alpha1.LeaderElectionConfiguration{
					LeaderElect:       ptr.To(false),
					LeaseDuration:     metav1.Duration{Duration: 30 * time.Second},
					RenewDeadline:     metav1.Duration{Duration: 20 * time.Second},
					RetryPeriod:       metav1.Duration{Duration: 5 * time.Second},
					ResourceLock:      "lock",
					ResourceName:      "name",
					ResourceNamespace: "namespace",
				}
				Expect(obj).To(Equal(expectedLeaderElectionConfig))
			})
		})
	})
})
