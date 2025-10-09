// SPDX-FileCopyrightText: SAP SE or an SAP affiliate company and Gardener contributors
//
// SPDX-License-Identifier: Apache-2.0

package garbagecollector

import (
	"context"

	"sigs.k8s.io/controller-runtime/pkg/client"
	logf "sigs.k8s.io/controller-runtime/pkg/log"
	"sigs.k8s.io/controller-runtime/pkg/reconcile"

	configv1alpha1 "github.com/gardener/garden-shoot-trust-configurator/pkg/apis/config/v1alpha1"
)

// Reconciler performs garbage collection.
type Reconciler struct {
	Client client.Client
	Config configv1alpha1.GarbageCollectorControllerConfig
}

// Reconcile performs the main reconciliation logic.
func (r *Reconciler) Reconcile(ctx context.Context, _ reconcile.Request) (reconcile.Result, error) {
	log := logf.FromContext(ctx)

	log.Info("Starting garbage collection")
	defer log.Info("Garbage collection finished")

	return reconcile.Result{
		Requeue:      true,
		RequeueAfter: r.Config.SyncPeriod.Duration,
	}, nil
}
