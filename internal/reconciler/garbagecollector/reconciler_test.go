// SPDX-FileCopyrightText: SAP SE or an SAP affiliate company and Gardener contributors
//
// SPDX-License-Identifier: Apache-2.0

package garbagecollector_test

import (
	"context"
	"time"

	gardencorev1beta1 "github.com/gardener/gardener/pkg/apis/core/v1beta1"
	"github.com/gardener/gardener/pkg/client/kubernetes"
	authenticationv1alpha1 "github.com/gardener/oidc-webhook-authenticator/apis/authentication/v1alpha1"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
	testclock "k8s.io/utils/clock/testing"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/client/fake"
	logf "sigs.k8s.io/controller-runtime/pkg/log"
	logzap "sigs.k8s.io/controller-runtime/pkg/log/zap"
	"sigs.k8s.io/controller-runtime/pkg/reconcile"

	garbagecollectorcontroller "github.com/gardener/garden-shoot-trust-configurator/internal/reconciler/garbagecollector"
	configv1alpha1 "github.com/gardener/garden-shoot-trust-configurator/pkg/apis/config/v1alpha1"
)

var _ = Describe("Controller", func() {
	var (
		ctx = logf.IntoContext(context.Background(), logzap.New(logzap.WriteTo(GinkgoWriter)))

		gc         *garbagecollectorcontroller.Reconciler
		fakeClient client.Client

		creationTimestamp = metav1.Date(2000, 5, 5, 5, 30, 0, 0, time.Local)
		fakeClock         = testclock.NewFakeClock(creationTimestamp.Add(time.Minute / 2))
	)

	BeforeEach(func() {
		scheme := runtime.NewScheme()
		Expect(kubernetes.AddGardenSchemeToScheme(scheme)).To(Succeed())
		Expect(authenticationv1alpha1.AddToScheme(scheme)).To(Succeed())

		fakeClient = fake.NewClientBuilder().WithScheme(scheme).Build()
		gc = &garbagecollectorcontroller.Reconciler{
			Client: fakeClient,
			Clock:  fakeClock,
			Config: configv1alpha1.GarbageCollectorControllerConfig{
				SyncPeriod:            &metav1.Duration{Duration: time.Hour},
				MinimumObjectLifetime: &metav1.Duration{Duration: time.Minute},
			},
		}
	})

	Describe("#GarbageCollect Reconcile Successful", func() {
		var (
			unlabeledOIDC *authenticationv1alpha1.OpenIDConnect

			labeledOIDC1 *authenticationv1alpha1.OpenIDConnect
			labeledOIDC2 *authenticationv1alpha1.OpenIDConnect
			labeledOIDC3 *authenticationv1alpha1.OpenIDConnect
			labeledOIDC4 *authenticationv1alpha1.OpenIDConnect
			labeledOIDC5 *authenticationv1alpha1.OpenIDConnect
			labeledOIDC6 *authenticationv1alpha1.OpenIDConnect
			labeledOIDC7 *authenticationv1alpha1.OpenIDConnect
			labeledOIDC8 *authenticationv1alpha1.OpenIDConnect
			labeledOIDC9 *authenticationv1alpha1.OpenIDConnect

			trustedShoot1    *gardencorev1beta1.Shoot
			trustedShoot2    *gardencorev1beta1.Shoot
			nonTrustedShoot3 *gardencorev1beta1.Shoot
		)

		BeforeEach(func() {
			unlabeledOIDC = &authenticationv1alpha1.OpenIDConnect{
				ObjectMeta: metav1.ObjectMeta{
					Name: "unlabeledOIDC",
				},
			}

			labeledOIDC1 = createLabeledOIDC("garden--shoot-1--UID")
			labeledOIDC2 = createLabeledOIDC("garden--shoot-2--UID")
			labeledOIDC3 = createLabeledOIDC("garden--shoot-3--UID")
			labeledOIDC4 = createLabeledOIDC("garden--shoot-4--UID")
			labeledOIDC5 = createLabeledOIDC("garden--shoot-5--UID")
			labeledOIDC6 = createLabeledOIDC("garden--shoot-6--UID")
			labeledOIDC7 = createLabeledOIDC("garden--shoot-7--UID")
			labeledOIDC8 = createLabeledOIDC("garden--shoot-8--UID")
			labeledOIDC9 = createLabeledOIDC("garden--shoot-9--UID")
			labeledOIDC9.CreationTimestamp = creationTimestamp

			trustedShoot1 = &gardencorev1beta1.Shoot{
				ObjectMeta: metav1.ObjectMeta{
					Name:      "shoot-1",
					Namespace: "garden",
					UID:       "UID",
					Annotations: map[string]string{
						"authentication.gardener.cloud/issuer":  "managed",
						"authentication.gardener.cloud/trusted": "true",
					},
				},
			}

			trustedShoot2 = &gardencorev1beta1.Shoot{
				ObjectMeta: metav1.ObjectMeta{
					Name:      "shoot-2",
					Namespace: "garden",
					UID:       "UID",
					Annotations: map[string]string{
						"authentication.gardener.cloud/issuer":  "managed",
						"authentication.gardener.cloud/trusted": "True",
					},
				},
			}

			nonTrustedShoot3 = &gardencorev1beta1.Shoot{
				ObjectMeta: metav1.ObjectMeta{
					Name:      "shoot-3",
					Namespace: "garden",
					UID:       "UID",
					Annotations: map[string]string{
						"authentication.gardener.cloud/issuer":  "managed",
						"authentication.gardener.cloud/trusted": "false",
					},
				},
			}
		})

		It("should do nothing because no OIDCs are found", func() {
			oidcList := &authenticationv1alpha1.OpenIDConnectList{}
			Expect(fakeClient.List(ctx, oidcList)).To(Succeed())
			Expect(oidcList.Items).To(BeEmpty())

			res, err := gc.Reconcile(ctx, reconcile.Request{})
			Expect(err).NotTo(HaveOccurred())
			Expect(res).To(Equal(reconcile.Result{RequeueAfter: time.Hour}))

			oidcList = &authenticationv1alpha1.OpenIDConnectList{}
			Expect(fakeClient.List(ctx, oidcList)).To(Succeed())
			Expect(oidcList.Items).To(BeEmpty())
		})

		It("should delete nothing because no labeled OIDCs are found", func() {
			Expect(fakeClient.Create(ctx, unlabeledOIDC)).To(Succeed())

			oidcList := &authenticationv1alpha1.OpenIDConnectList{}
			Expect(fakeClient.List(ctx, oidcList)).To(Succeed())
			Expect(oidcList.Items).To(ConsistOf(*unlabeledOIDC))

			res, err := gc.Reconcile(ctx, reconcile.Request{})
			Expect(err).NotTo(HaveOccurred())
			Expect(res).To(Equal(reconcile.Result{RequeueAfter: time.Hour}))

			oidcList = &authenticationv1alpha1.OpenIDConnectList{}
			Expect(fakeClient.List(ctx, oidcList)).To(Succeed())
			Expect(oidcList.Items).To(ConsistOf(*unlabeledOIDC))
		})

		It("should delete the unused resources", func() {
			Expect(fakeClient.Create(ctx, labeledOIDC1)).To(Succeed())
			Expect(fakeClient.Create(ctx, labeledOIDC2)).To(Succeed())
			Expect(fakeClient.Create(ctx, labeledOIDC3)).To(Succeed())
			Expect(fakeClient.Create(ctx, labeledOIDC4)).To(Succeed())
			Expect(fakeClient.Create(ctx, labeledOIDC5)).To(Succeed())
			Expect(fakeClient.Create(ctx, labeledOIDC6)).To(Succeed())
			Expect(fakeClient.Create(ctx, labeledOIDC7)).To(Succeed())
			Expect(fakeClient.Create(ctx, labeledOIDC8)).To(Succeed())
			Expect(fakeClient.Create(ctx, labeledOIDC9)).To(Succeed())

			oidcList := &authenticationv1alpha1.OpenIDConnectList{}
			Expect(fakeClient.List(ctx, oidcList)).To(Succeed())
			Expect(oidcList.Items).To(ConsistOf(
				*labeledOIDC1, *labeledOIDC2, *labeledOIDC3,
				*labeledOIDC4, *labeledOIDC6, *labeledOIDC7, *labeledOIDC8,
				*labeledOIDC9, *labeledOIDC5,
			))

			Expect(fakeClient.Create(ctx, trustedShoot1)).To(Succeed())
			Expect(fakeClient.Create(ctx, trustedShoot2)).To(Succeed())
			Expect(fakeClient.Create(ctx, nonTrustedShoot3)).To(Succeed())

			res, err := gc.Reconcile(ctx, reconcile.Request{})
			Expect(err).NotTo(HaveOccurred())
			Expect(res).To(Equal(reconcile.Result{RequeueAfter: time.Hour}))

			oidcList = &authenticationv1alpha1.OpenIDConnectList{}
			Expect(fakeClient.List(ctx, oidcList)).To(Succeed())
			Expect(oidcList.Items).To(ConsistOf(
				*labeledOIDC1, *labeledOIDC2, *labeledOIDC9,
			))
		})
	})

	Describe("#GarbageCollect Reconcile With Invalid Resources", func() {
		var invalidNameOIDC *authenticationv1alpha1.OpenIDConnect

		BeforeEach(func() {
			invalidNameOIDC = createLabeledOIDC("invalid--name")
		})

		It("should do nothing because OIDC name is invalid", func() {
			Expect(fakeClient.Create(ctx, invalidNameOIDC)).To(Succeed())

			oidcList := &authenticationv1alpha1.OpenIDConnectList{}
			Expect(fakeClient.List(ctx, oidcList)).To(Succeed())
			Expect(oidcList.Items).To(ConsistOf(*invalidNameOIDC))

			res, err := gc.Reconcile(ctx, reconcile.Request{})
			Expect(err).ToNot(HaveOccurred())
			Expect(res).To(Equal(reconcile.Result{RequeueAfter: time.Hour}))

			oidcList = &authenticationv1alpha1.OpenIDConnectList{}
			Expect(fakeClient.List(ctx, oidcList)).To(Succeed())
			Expect(oidcList.Items).To(ConsistOf(*invalidNameOIDC))
		})
	})
})

func createLabeledOIDC(name string) *authenticationv1alpha1.OpenIDConnect {
	return &authenticationv1alpha1.OpenIDConnect{
		ObjectMeta: metav1.ObjectMeta{
			Name: name,
			Labels: map[string]string{
				"app.kubernetes.io/managed-by": "garden-shoot-trust-configurator",
			},
		},
	}
}
