// SPDX-FileCopyrightText: SAP SE or an SAP affiliate company and Gardener contributors
//
// SPDX-License-Identifier: Apache-2.0

package reconciler

import (
	"strconv"
	"time"

	gardencorev1beta1 "github.com/gardener/gardener/pkg/apis/core/v1beta1"
	v1beta1constants "github.com/gardener/gardener/pkg/apis/core/v1beta1/constants"
	"golang.org/x/time/rate"
	"k8s.io/client-go/util/workqueue"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/builder"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/controller"
	"sigs.k8s.io/controller-runtime/pkg/event"
	"sigs.k8s.io/controller-runtime/pkg/predicate"
	"sigs.k8s.io/controller-runtime/pkg/reconcile"
)

const (
	// ControllerName is the name of the shoot trust configurator.
	ControllerName = "shoot-trust-configurator"

	// AnnotationTrustedShoot is the annotation that marks a Shoot to be trusted in the Garden cluster.
	AnnotationTrustedShoot = "authentication.gardener.cloud/trusted"
)

// SetupWithManager specifies how the controller is built
// to watch Shoots with the "authentication.gardener.cloud/trusted" annotation set to "true"
func (r *Reconciler) SetupWithManager(mgr ctrl.Manager) error {
	if r.Client == nil {
		r.Client = mgr.GetClient()
	}

	return builder.ControllerManagedBy(mgr).
		Named(ControllerName).
		For(&gardencorev1beta1.Shoot{}, builder.WithPredicates(r.ShootPredicate())).
		WithOptions(controller.Options{
			MaxConcurrentReconciles: 50,
			RateLimiter: workqueue.NewTypedMaxOfRateLimiter(
				workqueue.NewTypedItemExponentialFailureRateLimiter[reconcile.Request](5*time.Second, 2*time.Minute),
				&workqueue.TypedBucketRateLimiter[reconcile.Request]{Limiter: rate.NewLimiter(rate.Limit(10), 100)},
			),
		}).
		Complete(r)
}

// ShootPredicate returns a predicate to filter Shoot events
func (r *Reconciler) ShootPredicate() predicate.Predicate {
	return predicate.Funcs{
		CreateFunc:  func(e event.CreateEvent) bool { return r.IsRelevantShoot(e.Object) },
		UpdateFunc:  func(e event.UpdateEvent) bool { return r.IsRelevantShootUpdate(e.ObjectOld, e.ObjectNew) },
		DeleteFunc:  func(e event.DeleteEvent) bool { return r.IsRelevantShoot(e.Object) },
		GenericFunc: func(_ event.GenericEvent) bool { return false },
	}
}

// IsRelevantShoot checks whether the given object is a Shoot with the "authentication.gardener.cloud/trusted" annotation set to "true"
func (r *Reconciler) IsRelevantShoot(obj client.Object) bool {
	shoot, ok := obj.(*gardencorev1beta1.Shoot)
	if !ok {
		return false
	}
	if shoot.Annotations[v1beta1constants.AnnotationAuthenticationIssuer] != v1beta1constants.AnnotationAuthenticationIssuerManaged {
		return false
	}
	// Specifies whether the Shoot should be registered as a trusted cluster in the Garden cluster.
	if trusted, _ := strconv.ParseBool(shoot.Annotations[AnnotationTrustedShoot]); !trusted {
		return false
	}
	return true
}

// IsRelevantShootUpdate checks whether the old or new object is a relevant Shoot
func (r *Reconciler) IsRelevantShootUpdate(oldObj, newObj client.Object) bool {
	return r.IsRelevantShoot(newObj) || r.IsRelevantShoot(oldObj)
}
