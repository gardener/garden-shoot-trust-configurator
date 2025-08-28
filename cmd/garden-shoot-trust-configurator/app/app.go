// SPDX-FileCopyrightText: SAP SE or an SAP affiliate company and Gardener contributors
//
// SPDX-License-Identifier: Apache-2.0

package app

import (
	"context"
	"flag"
	"fmt"
	"net"
	"time"

	"github.com/gardener/gardener/pkg/client/kubernetes"
	gardenerhealthz "github.com/gardener/gardener/pkg/healthz"
	"github.com/gardener/gardener/pkg/logger"
	"github.com/go-logr/logr"
	"github.com/spf13/cobra"
	"github.com/spf13/pflag"
	"k8s.io/component-base/version"
	"k8s.io/component-base/version/verflag"
	"k8s.io/utils/clock"
	"k8s.io/utils/ptr"
	ctrl "sigs.k8s.io/controller-runtime"
	controllerconfig "sigs.k8s.io/controller-runtime/pkg/config"
	"sigs.k8s.io/controller-runtime/pkg/healthz"
	logf "sigs.k8s.io/controller-runtime/pkg/log"
	metricsserver "sigs.k8s.io/controller-runtime/pkg/metrics/server"
)

// AppName is the name of the application.
const AppName = "garden-shoot-trust-configurator"

// NewCommand is the root command for Gardener discovery server.
func NewCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use: AppName,
		RunE: func(cmd *cobra.Command, _ []string) error {
			logLevel, logFormat := "info", "json"
			log, err := logger.NewZapLogger(logLevel, logFormat)
			if err != nil {
				return fmt.Errorf("error instantiating zap logger: %w", err)
			}
			logf.SetLogger(log)

			log.Info("Starting application", "app", AppName, "version", version.Get())
			cmd.Flags().VisitAll(func(flag *pflag.Flag) {
				log.Info("Flag", "name", flag.Name, "value", flag.Value, "default", flag.DefValue)
			})

			return run(cmd.Context(), log)
		},
		PreRunE: func(_ *cobra.Command, _ []string) error {
			verflag.PrintAndExitIfRequested()
			return nil
		},
	}

	fs := cmd.Flags()
	verflag.AddFlags(fs)
	fs.AddGoFlagSet(flag.CommandLine)

	return cmd
}

func run(ctx context.Context, log logr.Logger) error {
	cfg, err := ctrl.GetConfig()
	if err != nil {
		return err
	}

	mgr, err := ctrl.NewManager(cfg, ctrl.Options{
		Logger: log.WithName("manager"),
		Scheme: kubernetes.GardenScheme,
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

	return mgr.Start(ctx)
}
