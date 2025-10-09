// SPDX-FileCopyrightText: SAP SE or an SAP affiliate company and Gardener contributors
//
// SPDX-License-Identifier: Apache-2.0

package app

import (
	"context"
	"fmt"
	"net"
	"time"

	"github.com/gardener/gardener/pkg/client/kubernetes"
	gardenerhealthz "github.com/gardener/gardener/pkg/healthz"
	"github.com/gardener/gardener/pkg/logger"
	authenticationv1alpha1 "github.com/gardener/oidc-webhook-authenticator/apis/authentication/v1alpha1"
	"github.com/go-logr/logr"
	"github.com/spf13/cobra"
	"github.com/spf13/pflag"
	"k8s.io/apimachinery/pkg/runtime"
	utilruntime "k8s.io/apimachinery/pkg/util/runtime"
	"k8s.io/component-base/version"
	"k8s.io/component-base/version/verflag"
	"k8s.io/utils/clock"
	"k8s.io/utils/ptr"
	ctrl "sigs.k8s.io/controller-runtime"
	controllerconfig "sigs.k8s.io/controller-runtime/pkg/config"
	"sigs.k8s.io/controller-runtime/pkg/healthz"
	logf "sigs.k8s.io/controller-runtime/pkg/log"
	metricsserver "sigs.k8s.io/controller-runtime/pkg/metrics/server"

	"github.com/gardener/garden-shoot-trust-configurator/internal/reconciler/garbagecollector"
	shootreconciler "github.com/gardener/garden-shoot-trust-configurator/internal/reconciler/shoot"
	configv1alpha1 "github.com/gardener/garden-shoot-trust-configurator/pkg/apis/config/v1alpha1"
)

// AppName is the name of the application.
const AppName = "garden-shoot-trust-configurator"

// NewCommand is the root command for Garden shoot trust configurator server.
func NewCommand() *cobra.Command {
	opt := newOptions()

	cmd := &cobra.Command{
		Use:   AppName,
		Short: "Launch the " + AppName,
		RunE: func(cmd *cobra.Command, _ []string) error {
			if err := opt.Complete(); err != nil {
				return err
			}

			if err := opt.Validate(); err != nil {
				return fmt.Errorf("cannot validate options: %w", err)
			}

			logLevel, logFormat := opt.LogConfig()
			log, err := logger.NewZapLogger(logLevel, logFormat)
			if err != nil {
				return fmt.Errorf("error instantiating zap logger: %w", err)
			}
			logf.SetLogger(log)

			log.Info("Starting application", "app", AppName, "version", version.Get())
			cmd.Flags().VisitAll(func(flag *pflag.Flag) {
				log.Info("Flag", "name", flag.Name, "value", flag.Value, "default", flag.DefValue)
			})

			return run(cmd.Context(), log, opt.config)
		},
		PreRunE: func(_ *cobra.Command, _ []string) error {
			verflag.PrintAndExitIfRequested()
			return nil
		},
	}

	flags := cmd.Flags()
	opt.addFlags(flags)
	verflag.AddFlags(flags)

	return cmd
}

func run(ctx context.Context, log logr.Logger, conf *configv1alpha1.GardenShootTrustConfiguratorConfiguration) error {
	cfg, err := ctrl.GetConfig()
	if err != nil {
		return err
	}

	scheme := runtime.NewScheme()
	utilruntime.Must(kubernetes.AddGardenSchemeToScheme(scheme))
	utilruntime.Must(authenticationv1alpha1.AddToScheme(scheme))

	mgr, err := ctrl.NewManager(cfg, ctrl.Options{
		Logger: log.WithName("manager"),
		Scheme: scheme,
		Metrics: metricsserver.Options{
			BindAddress: "0",
		},
		GracefulShutdownTimeout: ptr.To(5 * time.Second),
		// TODO(theoddora): Consider enabling the support for leader election + source/target clusters
		LeaderElection:         false,
		PprofBindAddress:       "",
		HealthProbeBindAddress: net.JoinHostPort("", "8081"),
		Controller: controllerconfig.Controller{
			RecoverPanic: ptr.To(true),
		},
	})
	if err != nil {
		return fmt.Errorf("unable to create manager: %w", err)
	}

	if err := mgr.AddHealthzCheck("ping", healthz.Ping); err != nil {
		return err
	}
	if err := mgr.AddHealthzCheck("informer-sync", gardenerhealthz.NewCacheSyncHealthzWithDeadline(mgr.GetLogger(), clock.RealClock{}, mgr.GetCache(), gardenerhealthz.DefaultCacheSyncDeadline)); err != nil {
		return err
	}
	if err := mgr.AddReadyzCheck("informer-sync", gardenerhealthz.NewCacheSyncHealthz(mgr.GetCache())); err != nil {
		return err
	}

	// Setup all Controllers
	if err := (&shootreconciler.Reconciler{}).SetupWithManager(mgr); err != nil {
		return fmt.Errorf("unable to create shoot reconcile controller: %w", err)
	}

	if err := (&garbagecollector.Reconciler{
		Config: conf.Controllers.GarbageCollector,
	}).SetupWithManager(mgr); err != nil {
		return fmt.Errorf("unable to create garbage collector controller: %w", err)
	}

	return mgr.Start(ctx)
}
