// SPDX-FileCopyrightText: SAP SE or an SAP affiliate company and Gardener contributors
//
// SPDX-License-Identifier: Apache-2.0

package handler

import (
	"net/http"

	"github.com/gardener/oidc-webhook-authenticator/webhook/authentication"
)

func Handle() http.Handler {
	// Dummy return handler
	owa := authentication.Webhook{}
	return owa.Build()
}
