// SPDX-FileCopyrightText: SAP SE or an SAP affiliate company and Gardener contributors
//
// SPDX-License-Identifier: Apache-2.0

package reconciler

import (
	"context"
	"fmt"
	"strconv"

	gardencorev1beta1 "github.com/gardener/gardener/pkg/apis/core/v1beta1"
	v1beta1constants "github.com/gardener/gardener/pkg/apis/core/v1beta1/constants"
	"github.com/gardener/gardener/pkg/controllerutils"
	authenticationv1alpha1 "github.com/gardener/oidc-webhook-authenticator/apis/authentication/v1alpha1"
	"github.com/go-logr/logr"
	apierrors "k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/controller/controllerutil"
	logf "sigs.k8s.io/controller-runtime/pkg/log"
	"sigs.k8s.io/controller-runtime/pkg/reconcile"

	configv1alpha1 "github.com/gardener/garden-shoot-trust-configurator/pkg/apis/config/v1alpha1"
)

const (
	// labelManagedByKey is a constant for a key of a label on an OIDC resource describing who is managing it.
	labelManagedByKey = "app.kubernetes.io/managed-by"
	// labelManagedByValue is a constant for a value of a label on a OIDC describing the value 'garden-shoot-trust-configurator'.
	labelManagedByValue = "garden-shoot-trust-configurator"
)

// Reconciler reconciles shoot trust configurator information.
type Reconciler struct {
	Client client.Client
	Config configv1alpha1.ShootControllerConfig
}

// Reconcile handles reconciliation requests for Shoots marked to be trusted in the Garden cluster.
func (r *Reconciler) Reconcile(ctx context.Context, req ctrl.Request) (ctrl.Result, error) {
	log := logf.FromContext(ctx)

	ctx, cancel := controllerutils.GetMainReconciliationContext(ctx, controllerutils.DefaultReconciliationTimeout)
	defer cancel()

	shoot := &gardencorev1beta1.Shoot{}
	if err := r.Client.Get(ctx, req.NamespacedName, shoot); err != nil {
		if apierrors.IsNotFound(err) {
			log.Info("Object is gone, stop reconciling")
			// We don't have the shoot object here, so we cannot pass it to deleteOIDCResource to construct the OIDC resource name.
			// We have a garbage collection mechanism to clean up old OIDC resources that are not referenced by any shoot anymore.
			return reconcile.Result{}, nil
		}
		return reconcile.Result{}, fmt.Errorf("error retrieving shoot from store: %w", err)
	}

	if shoot.DeletionTimestamp != nil {
		log.Info("Shoot is being deleted, cleaning up OIDC resource")
		return r.handleDeletion(ctx, log, shoot)
	}

	if shoot.Annotations[v1beta1constants.AnnotationAuthenticationIssuer] != v1beta1constants.AnnotationAuthenticationIssuerManaged {
		log.Info("Shoot does not have expected annotation or their value is not 'managed'", "annotation", v1beta1constants.AnnotationAuthenticationIssuer, "value", shoot.Annotations[v1beta1constants.AnnotationAuthenticationIssuer])
		return r.handleDeletion(ctx, log, shoot)
	}

	if trusted, _ := strconv.ParseBool(shoot.Annotations[AnnotationTrustedShoot]); !trusted {
		log.Info("Shoot does not have expected annotation or their value is not 'true', clean up OIDC resource", "annotation", AnnotationTrustedShoot, "value", shoot.Annotations[AnnotationTrustedShoot])
		return r.handleDeletion(ctx, log, shoot)
	}

	if !controllerutil.ContainsFinalizer(shoot, FinalizerName) {
		log.Info("Adding finalizer")
		if err := controllerutils.AddFinalizers(ctx, r.Client, shoot, FinalizerName); err != nil {
			return ctrl.Result{}, fmt.Errorf("could not add finalizer to shoot: %w", err)
		}
	}

	var issuerURL string
	for _, adr := range shoot.Status.AdvertisedAddresses {
		if adr.Name == v1beta1constants.AdvertisedAddressServiceAccountIssuer {
			issuerURL = adr.URL
			break
		}
	}

	if issuerURL == "" {
		return ctrl.Result{}, fmt.Errorf("shoot does not have service-account-issuer in its status.advertisedAddresses: %s", shoot.Status.AdvertisedAddresses)
	}

	// TODO(theoddora): Add proper validation that a single issuer is not registered more than once
	// This should check if another OIDC resource with the same issuerURL already exists

	var (
		userNameClaim  = "sub"
		groupsClaim    = "groups"
		prefix         = buildPrefix(shoot)
		userNamePrefix = prefix
		groupsPrefix   = prefix
		clientID       = r.Config.OIDCConfig.ClientID
	)

	oidc := emptyOIDC(shoot)
	if _, err := controllerutils.GetAndCreateOrMergePatch(ctx, r.Client, oidc, func() error {
		oidc.Annotations = nil
		oidc.Labels = map[string]string{
			labelManagedByKey: labelManagedByValue,
		}
		oidc.Spec = authenticationv1alpha1.OIDCAuthenticationSpec{
			IssuerURL:      issuerURL,
			ClientID:       clientID,
			UsernameClaim:  &userNameClaim,
			UsernamePrefix: &userNamePrefix,
			GroupsClaim:    &groupsClaim,
			GroupsPrefix:   &groupsPrefix,
		}
		return nil
	}); err != nil {
		return ctrl.Result{}, err
	}

	log.Info("Successfully created or updated OIDC resource for shoot", "oidc", client.ObjectKeyFromObject(oidc))
	return ctrl.Result{RequeueAfter: r.Config.SyncPeriod.Duration}, nil
}

// handleDeletion handles the deletion of a shoot and its associated OIDC resource
func (r *Reconciler) handleDeletion(ctx context.Context, log logr.Logger, shoot *gardencorev1beta1.Shoot) (ctrl.Result, error) {
	// Clean up the OIDC resource
	if err := r.deleteOIDCResource(ctx, log, shoot); err != nil {
		return ctrl.Result{}, err
	}

	log.Info("Removing finalizer")
	if err := controllerutils.RemoveFinalizers(ctx, r.Client, shoot, FinalizerName); err != nil {
		return reconcile.Result{}, fmt.Errorf("failed to remove finalizer: %w", err)
	}

	return reconcile.Result{}, nil
}

func (r *Reconciler) deleteOIDCResource(ctx context.Context, log logr.Logger, shoot *gardencorev1beta1.Shoot) error {
	oidc := emptyOIDC(shoot)
	oidcObjectKey := client.ObjectKeyFromObject(oidc)
	err := r.Client.Get(ctx, oidcObjectKey, oidc)
	if err != nil {
		if apierrors.IsNotFound(err) {
			log.Info("OIDC resource not found, nothing to do", "oidc", oidcObjectKey)
			return nil
		}
		return fmt.Errorf("failed to get OIDC: %w", err)
	}

	if err := r.Client.Delete(ctx, oidc); err != nil {
		if apierrors.IsNotFound(err) {
			log.Info("OIDC resource not found, nothing to do", "oidc", oidcObjectKey)
			return nil
		}
		return fmt.Errorf("failed to delete OIDC: %w", err)
	}
	log.Info("Successfully deleted OIDC resource", "oidc", oidcObjectKey)
	return nil
}

func emptyOIDC(shoot *gardencorev1beta1.Shoot) *authenticationv1alpha1.OpenIDConnect {
	return &authenticationv1alpha1.OpenIDConnect{
		ObjectMeta: metav1.ObjectMeta{
			Name: getOIDCResourceName(shoot),
		},
	}
}

func getOIDCResourceName(shoot *gardencorev1beta1.Shoot) string {
	return fmt.Sprintf("%s--%s--%s", shoot.Namespace, shoot.Name, shoot.UID)
}

func buildPrefix(shoot *gardencorev1beta1.Shoot) string {
	return fmt.Sprintf("ns:%s:shoot:%s:%s:", shoot.Namespace, shoot.Name, string(shoot.UID))
}
