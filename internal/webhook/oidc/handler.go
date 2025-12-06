// SPDX-FileCopyrightText: SAP SE or an SAP affiliate company and Gardener contributors
//
// SPDX-License-Identifier: Apache-2.0

package oidc

import (
	"context"
	"fmt"
	"net/http"

	authenticationv1alpha1 "github.com/gardener/oidc-webhook-authenticator/apis/authentication/v1alpha1"
	"github.com/go-logr/logr"
	admissionv1 "k8s.io/api/admission/v1"
	"sigs.k8s.io/controller-runtime/pkg/webhook/admission"

	"github.com/gardener/garden-shoot-trust-configurator/pkg/apis/constants"
)

// Handler is an admission webhook handler that restricts updates to certain fields
// of managed OpenIDConnect resources.
type Handler struct {
	Logger  logr.Logger
	Decoder admission.Decoder
}

// Handle handles an admission request for an OIDC resource and restricts updates to labels
// if the resource is managed by the Garden Shoot Trust Configurator.
func (h *Handler) Handle(_ context.Context, req admission.Request) admission.Response {
	h.Logger.Info("OIDC Handler invoked",
		"operation", req.Operation,
		"resource", req.Resource.Resource,
		"name", req.Name,
		"username", req.UserInfo.Username,
	)

	if req.Operation != admissionv1.Update {
		return admission.Allowed("")
	}

	oldObj := &authenticationv1alpha1.OpenIDConnect{}
	if err := h.Decoder.DecodeRaw(req.OldObject, oldObj); err != nil {
		return admission.Errored(http.StatusBadRequest, err)
	}

	// If the old object is not managed, allow the update
	if oldObj.Labels[constants.LabelManagedByKey] != constants.LabelManagedByValue {
		return admission.Allowed("")
	}

	newObj := &authenticationv1alpha1.OpenIDConnect{}
	if err := h.Decoder.DecodeRaw(req.Object, newObj); err != nil {
		return admission.Errored(http.StatusBadRequest, err)
	}

	if newObj.Labels[constants.LabelManagedByKey] != constants.LabelManagedByValue {
		return admission.Denied(fmt.Sprintf("removing or changing label %q for managed OpenIDConnect is not allowed", constants.LabelManagedByKey))
	}

	return admission.Allowed("")
}
