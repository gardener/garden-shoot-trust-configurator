// SPDX-FileCopyrightText: SAP SE or an SAP affiliate company and Gardener contributors
//
// SPDX-License-Identifier: Apache-2.0

package reconciler

import (
	"context"
	"time"

	"github.com/go-logr/logr"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"
	logf "sigs.k8s.io/controller-runtime/pkg/log"
)

// Reconciler reconciles shoot trust configurator information.
type Reconciler struct {
	Client client.Client
	Log    logr.Logger

	ResyncPeriod time.Duration
}

// Reconcile handles reconciliation requests for Shoots marked to be trusted in the Garden cluster.
func (r *Reconciler) Reconcile(ctx context.Context, req ctrl.Request) (ctrl.Result, error) {
	log := logf.FromContext(ctx)

	log.Info("Shoot reconcile finished")

	return ctrl.Result{}, nil
}
