// SPDX-FileCopyrightText: SAP SE or an SAP affiliate company and Gardener contributors
//
// SPDX-License-Identifier: Apache-2.0

package oidc_test

import (
	"context"
	"net/http"

	"github.com/gardener/gardener/pkg/client/kubernetes"
	authenticationv1alpha1 "github.com/gardener/oidc-webhook-authenticator/apis/authentication/v1alpha1"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	admissionv1 "k8s.io/api/admission/v1"
	authenticationv1 "k8s.io/api/authentication/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/runtime/serializer/json"
	"sigs.k8s.io/controller-runtime/pkg/webhook/admission"

	"github.com/gardener/garden-shoot-trust-configurator/internal/webhook/oidc"
)

var _ = Describe("#Handler", func() {
	var (
		ctx context.Context

		handler admission.Handler
		request admission.Request
		encoder runtime.Encoder

		responseAllowed admission.Response
	)

	BeforeEach(func() {
		ctx = context.Background()

		scheme := runtime.NewScheme()
		Expect(kubernetes.AddGardenSchemeToScheme(scheme)).To(Succeed())
		Expect(authenticationv1.AddToScheme(scheme)).To(Succeed())

		handler = oidc.NewHandler(admission.NewDecoder(scheme))

		encoder = &json.Serializer{}
		request.UserInfo = authenticationv1.UserInfo{
			Username: "garden-shoot-trust-configurator",
		}
		request.Resource = metav1.GroupVersionResource{
			Resource: "openidconnects",
		}
		request.Operation = admissionv1.Update

		responseAllowed = admission.Response{
			AdmissionResponse: admissionv1.AdmissionResponse{
				Allowed: true,
				Result: &metav1.Status{
					Code: int32(http.StatusOK),
				},
			},
		}
	})

	Describe("#Handle", func() {
		It("should allow update that does not change any managed label", func() {
			objData, err := runtime.Encode(encoder, &authenticationv1alpha1.OpenIDConnect{
				ObjectMeta: metav1.ObjectMeta{
					Name: "example-oidc",
					Labels: map[string]string{
						"app.kubernetes.io/managed-by": "garden-shoot-trust-configurator",
						"env":                          "prod",
					},
				},
			})
			Expect(err).NotTo(HaveOccurred())
			request.OldObject.Raw = objData

			objData, err = runtime.Encode(encoder, &authenticationv1alpha1.OpenIDConnect{
				ObjectMeta: metav1.ObjectMeta{
					Name: "example-oidc",
					Labels: map[string]string{
						"app.kubernetes.io/managed-by": "garden-shoot-trust-configurator",
						"env":                          "prod",
					},
				},
			})
			Expect(err).NotTo(HaveOccurred())
			request.Object.Raw = objData

			Expect(handler.Handle(ctx, request)).To(Equal(responseAllowed))
		})

		It("should deny update that removes a managed label", func() {
			objData, err := runtime.Encode(encoder, &authenticationv1alpha1.OpenIDConnect{
				ObjectMeta: metav1.ObjectMeta{
					Name: "example-oidc",
					Labels: map[string]string{
						"app.kubernetes.io/managed-by": "garden-shoot-trust-configurator",
						"env":                          "prod",
					},
				},
			})
			Expect(err).NotTo(HaveOccurred())
			request.OldObject.Raw = objData

			objData, err = runtime.Encode(encoder, &authenticationv1alpha1.OpenIDConnect{
				ObjectMeta: metav1.ObjectMeta{
					Name: "example-oidc",
					Labels: map[string]string{
						"env": "prod",
					},
				},
			})
			Expect(err).NotTo(HaveOccurred())
			request.Object.Raw = objData

			response := handler.Handle(ctx, request)
			Expect(response.Allowed).To(BeFalse())
			Expect(response.Result.Message).To(ContainSubstring(`removing or changing label "app.kubernetes.io/managed-by" for managed OpenIDConnect is not allowed`))
		})

		It("should deny update that changes a managed label", func() {
			objData, err := runtime.Encode(encoder, &authenticationv1alpha1.OpenIDConnect{
				ObjectMeta: metav1.ObjectMeta{
					Name: "example-oidc",
					Labels: map[string]string{
						"app.kubernetes.io/managed-by": "garden-shoot-trust-configurator",
						"env":                          "prod",
					},
				},
			})
			Expect(err).NotTo(HaveOccurred())
			request.OldObject.Raw = objData

			objData, err = runtime.Encode(encoder, &authenticationv1alpha1.OpenIDConnect{
				ObjectMeta: metav1.ObjectMeta{
					Name: "example-oidc",
					Labels: map[string]string{
						"app.kubernetes.io/managed-by": "some-other-manager",
						"env":                          "prod",
					},
				},
			})
			Expect(err).NotTo(HaveOccurred())
			request.Object.Raw = objData

			response := handler.Handle(ctx, request)
			Expect(response.Allowed).To(BeFalse())
			Expect(response.Result.Message).To(ContainSubstring(`removing or changing label "app.kubernetes.io/managed-by" for managed OpenIDConnect is not allowed`))
		})

		It("should allow update of an unmanaged OIDC resource", func() {
			objData, err := runtime.Encode(encoder, &authenticationv1alpha1.OpenIDConnect{
				ObjectMeta: metav1.ObjectMeta{
					Name: "example-oidc",
					Labels: map[string]string{
						"env": "prod",
					},
				},
			})
			Expect(err).NotTo(HaveOccurred())
			request.OldObject.Raw = objData

			objData, err = runtime.Encode(encoder, &authenticationv1alpha1.OpenIDConnect{
				ObjectMeta: metav1.ObjectMeta{
					Name: "example-oidc",
					Labels: map[string]string{
						"env": "staging",
					},
				},
			})
			Expect(err).NotTo(HaveOccurred())
			request.Object.Raw = objData

			Expect(handler.Handle(ctx, request)).To(Equal(responseAllowed))
		})

		It("should allow non-update operations", func() {
			request.Operation = admissionv1.Create

			objData, err := runtime.Encode(encoder, &authenticationv1alpha1.OpenIDConnect{
				ObjectMeta: metav1.ObjectMeta{
					Name: "example-oidc",
					Labels: map[string]string{
						"app.kubernetes.io/managed-by": "garden-shoot-trust-configurator",
						"env":                          "prod",
					},
				},
			})
			Expect(err).NotTo(HaveOccurred())
			request.Object.Raw = objData

			Expect(handler.Handle(ctx, request)).To(Equal(responseAllowed))
		})

		It("should allow update of labels when OIDC is managed", func() {
			objData, err := runtime.Encode(encoder, &authenticationv1alpha1.OpenIDConnect{
				ObjectMeta: metav1.ObjectMeta{
					Name: "example-oidc",
					Labels: map[string]string{
						"app.kubernetes.io/managed-by": "garden-shoot-trust-configurator",
					},
				},
			})
			Expect(err).NotTo(HaveOccurred())
			request.OldObject.Raw = objData

			objData, err = runtime.Encode(encoder, &authenticationv1alpha1.OpenIDConnect{
				ObjectMeta: metav1.ObjectMeta{
					Name: "example-oidc",
					Labels: map[string]string{
						"app.kubernetes.io/managed-by": "garden-shoot-trust-configurator",
						"env":                          "staging",
					},
				},
			})
			Expect(err).NotTo(HaveOccurred())
			request.Object.Raw = objData

			Expect(handler.Handle(ctx, request)).To(Equal(responseAllowed))
		})

		It("should return an error if decoding fails", func() {
			request.OldObject.Raw = []byte("invalid-json")
			request.Object.Raw = []byte("invalid-json")

			response := handler.Handle(ctx, request)
			Expect(response.Allowed).To(BeFalse())
			Expect(response.Result.Code).To(Equal(int32(http.StatusBadRequest)))
		})
	})
})
