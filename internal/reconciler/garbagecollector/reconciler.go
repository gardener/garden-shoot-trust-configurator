// SPDX-FileCopyrightText: SAP SE or an SAP affiliate company and Gardener contributors
//
// SPDX-License-Identifier: Apache-2.0

package garbagecollector

import (
	"context"
	"errors"
	"strconv"
	"strings"

	gardencorev1beta1 "github.com/gardener/gardener/pkg/apis/core/v1beta1"
	"github.com/gardener/gardener/pkg/controllerutils"
	authenticationv1alpha1 "github.com/gardener/oidc-webhook-authenticator/apis/authentication/v1alpha1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/types"
	"k8s.io/apimachinery/pkg/util/sets"
	"k8s.io/utils/clock"
	"sigs.k8s.io/controller-runtime/pkg/client"
	logf "sigs.k8s.io/controller-runtime/pkg/log"
	"sigs.k8s.io/controller-runtime/pkg/reconcile"

	configv1alpha1 "github.com/gardener/garden-shoot-trust-configurator/pkg/apis/config/v1alpha1"
	constants "github.com/gardener/garden-shoot-trust-configurator/pkg/apis/constants"
)

// Reconciler performs garbage collection.
type Reconciler struct {
	Client client.Client
	Config configv1alpha1.GarbageCollectorControllerConfig
	Clock  clock.Clock
}

// Reconcile performs the main reconciliation logic.
func (r *Reconciler) Reconcile(ctx context.Context, _ reconcile.Request) (reconcile.Result, error) {
	log := logf.FromContext(ctx)

	ctx, cancel := controllerutils.GetMainReconciliationContext(ctx, r.Config.SyncPeriod.Duration)
	defer cancel()

	log.Info("Starting garbage collection")

	var (
		label                   = client.MatchingLabels{constants.LabelManagedByKey: constants.LabelManagedByValue}
		objectsToGarbageCollect = sets.New[string]()
	)

	objList := &metav1.PartialObjectMetadataList{}
	objList.SetGroupVersionKind(authenticationv1alpha1.GroupVersion.WithKind("OpenIDConnectList"))
	if err := r.Client.List(ctx, objList, label); err != nil {
		return reconcile.Result{}, err
	}

	for _, obj := range objList.Items {
		if obj.CreationTimestamp.Add(r.Config.MinimumObjectLifetime.Duration).UTC().After(r.Clock.Now().UTC()) {
			// Do not consider recently created objects for garbage collection.
			continue
		}
		objectsToGarbageCollect.Insert(obj.Name)
	}

	for oidcName := range objectsToGarbageCollect {
		oidc := authenticationv1alpha1.OpenIDConnect{}
		if err := r.Client.Get(ctx, client.ObjectKey{Name: oidcName}, &oidc); err != nil {
			if client.IgnoreNotFound(err) != nil {
				log.Error(err, "Error retrieving OIDC resource", "oidc", oidc.Name)
				continue
			}
		}

		shootNamespacedName, err := parseOIDCResourceName(&oidc)
		if err != nil {
			log.Error(err, "Skipping OIDC resource as it has an invalid name", "oidc", oidc.Name)
			continue
		}

		shoot := &gardencorev1beta1.Shoot{}
		err = r.Client.Get(ctx, shootNamespacedName, shoot)
		if err != nil {
			if client.IgnoreNotFound(err) != nil {
				log.Error(err, "Error retrieving shoot", "shoot", shootNamespacedName)
				continue
			}

			log.Info("Shoot not found, deleting OIDC resource", "shoot", shootNamespacedName, "oidc", oidc.Name)
			if err := r.Client.Delete(ctx, &oidc); err != nil {
				if client.IgnoreNotFound(err) != nil {
					log.Error(err, "Error deleting OIDC resource", "oidc", oidc.Name)
					continue
				}
			}
			log.Info("Deleted OIDC resource", "oidc", oidc.Name)
			continue
		}

		if trusted, _ := strconv.ParseBool(shoot.Annotations[constants.AnnotationTrustedShoot]); !trusted {
			log.Info("Shoot is not trusted anymore, deleting OIDC resource", "shoot", shootNamespacedName, "oidc", oidc.Name)
			if err := r.Client.Delete(ctx, &oidc); err != nil {
				if client.IgnoreNotFound(err) != nil {
					log.Error(err, "Error deleting OIDC resource", "oidc", oidc.Name)
					continue
				}
			}
			log.Info("Deleted OIDC resource", "oidc", oidc.Name)
		}
	}

	log.Info("Garbage collection finished")
	return reconcile.Result{RequeueAfter: r.Config.SyncPeriod.Duration}, nil
}

// parseOIDCResourceName parses the OIDC resource name and returns the shoot's namespace and name.
// The expected format is "<namespace>--<name>--<uid>".
func parseOIDCResourceName(oidc *authenticationv1alpha1.OpenIDConnect) (types.NamespacedName, error) {
	parts := strings.SplitN(oidc.Name, constants.Separator, 3)
	if len(parts) != 3 {
		return types.NamespacedName{}, errors.New("invalid OIDC resource name format")
	}
	return types.NamespacedName{
		Namespace: parts[0],
		Name:      parts[1],
	}, nil
}
