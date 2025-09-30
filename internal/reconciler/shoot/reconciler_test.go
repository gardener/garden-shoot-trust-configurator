// SPDX-FileCopyrightText: SAP SE or an SAP affiliate company and Gardener contributors
//
// SPDX-License-Identifier: Apache-2.0

package reconciler_test

import (
	"context"
	"fmt"
	"time"

	gardencorev1beta1 "github.com/gardener/gardener/pkg/apis/core/v1beta1"
	v1beta1constants "github.com/gardener/gardener/pkg/apis/core/v1beta1/constants"
	"github.com/gardener/gardener/pkg/client/kubernetes"
	. "github.com/gardener/gardener/pkg/utils/test/matchers"
	authenticationv1alpha1 "github.com/gardener/oidc-webhook-authenticator/apis/authentication/v1alpha1"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/types"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/client/fake"
	logf "sigs.k8s.io/controller-runtime/pkg/log"
	logzap "sigs.k8s.io/controller-runtime/pkg/log/zap"
	"sigs.k8s.io/controller-runtime/pkg/reconcile"

	shootreconciler "github.com/gardener/garden-shoot-trust-configurator/internal/reconciler/shoot"
)

var _ = Describe("Reconciler", func() {
	const (
		shootName      = "my-shoot"
		shootNamespace = "garden-abc"
		resyncPeriod   = time.Second
	)

	var (
		ctx = logf.IntoContext(context.Background(), logzap.New(logzap.WriteTo(GinkgoWriter)))

		reconciler *shootreconciler.Reconciler
		fakeClient client.Client

		shoot          *gardencorev1beta1.Shoot
		shootUID       = types.UID("39f6d713-99c6-424a-827b-6bc532329b77")
		shootObjectKey client.ObjectKey

		oidc          *authenticationv1alpha1.OpenIDConnect
		oidcObjectKey client.ObjectKey
	)

	BeforeEach(func() {
		scheme := runtime.NewScheme()
		err := kubernetes.AddGardenSchemeToScheme(scheme)
		Expect(err).ToNot(HaveOccurred())
		err = authenticationv1alpha1.AddToScheme(scheme)
		Expect(err).ToNot(HaveOccurred())

		fakeClient = fake.NewClientBuilder().WithScheme(scheme).Build()
		reconciler = &shootreconciler.Reconciler{
			Client:       fakeClient,
			ResyncPeriod: resyncPeriod,
		}

		shoot = &gardencorev1beta1.Shoot{
			ObjectMeta: metav1.ObjectMeta{
				Name:      shootName,
				Namespace: shootNamespace,
				UID:       shootUID,
				Annotations: map[string]string{
					"authentication.gardener.cloud/issuer":  "managed",
					"authentication.gardener.cloud/trusted": "true",
				},
			},
			Status: gardencorev1beta1.ShootStatus{
				AdvertisedAddresses: []gardencorev1beta1.ShootAdvertisedAddress{
					{
						Name: v1beta1constants.AdvertisedAddressServiceAccountIssuer,
						URL:  "https://shoot/issuer",
					},
				},
			},
		}
		shootObjectKey = client.ObjectKey{Namespace: shootNamespace, Name: shootName}

		oidc = &authenticationv1alpha1.OpenIDConnect{
			ObjectMeta: metav1.ObjectMeta{
				Name: fmt.Sprintf("%s--%s--%s", shoot.Namespace, shoot.Name, shoot.UID),
			},
		}
		oidcObjectKey = client.ObjectKey{Name: oidc.Name}
	})

	It("should create OIDC resource", func() {
		Expect(fakeClient.Create(ctx, shoot)).To(Succeed())
		Expect(fakeClient.Get(ctx, shootObjectKey, shoot)).To(Succeed())

		res, err := reconciler.Reconcile(ctx, reconcile.Request{NamespacedName: shootObjectKey})
		Expect(err).ToNot(HaveOccurred())
		Expect(res).To(Equal(ctrl.Result{}))

		Expect(fakeClient.Get(ctx, oidcObjectKey, oidc)).To(Succeed())
		Expect(oidc.Spec.IssuerURL).To(Equal("https://shoot/issuer"))
		Expect(oidc.Spec.ClientID).To(Equal("garden"))
		Expect(oidc.Spec.UsernameClaim).ToNot(BeNil())
		Expect(*oidc.Spec.UsernameClaim).To(Equal("sub"))
		Expect(oidc.Spec.UsernamePrefix).ToNot(BeNil())
		Expect(*oidc.Spec.UsernamePrefix).To(Equal(fmt.Sprintf("ns:%s:shoot:%s:%s:", shoot.Namespace, shoot.Name, string(shoot.UID))))
		Expect(oidc.Spec.GroupsClaim).ToNot(BeNil())
		Expect(*oidc.Spec.GroupsClaim).To(Equal("groups"))
		Expect(oidc.Spec.GroupsPrefix).ToNot(BeNil())
		Expect(*oidc.Spec.GroupsPrefix).To(Equal(fmt.Sprintf("ns:%s:shoot:%s:%s:", shoot.Namespace, shoot.Name, string(shoot.UID))))
	})

	It("should do nothing because shoot has no managed issuer annotation", func() {
		shoot.Annotations = map[string]string{}
		Expect(fakeClient.Create(ctx, shoot)).To(Succeed())
		Expect(fakeClient.Get(ctx, shootObjectKey, shoot)).To(Succeed())

		res, err := reconciler.Reconcile(ctx, reconcile.Request{NamespacedName: shootObjectKey})
		Expect(err).ToNot(HaveOccurred())
		Expect(res).To(Equal(ctrl.Result{}))

		Expect(fakeClient.Get(ctx, oidcObjectKey, oidc)).To(BeNotFoundError())
	})

	It("should do nothing because shoot is not trusted", func() {
		shoot.Annotations = map[string]string{
			"authentication.gardener.cloud/issuer": "managed",
		}
		Expect(fakeClient.Create(ctx, shoot)).To(Succeed())
		Expect(fakeClient.Get(ctx, shootObjectKey, shoot)).To(Succeed())

		res, err := reconciler.Reconcile(ctx, reconcile.Request{NamespacedName: shootObjectKey})
		Expect(err).ToNot(HaveOccurred())
		Expect(res).To(Equal(ctrl.Result{}))

		Expect(fakeClient.Get(ctx, oidcObjectKey, oidc)).To(BeNotFoundError())
	})

	It("should delete OIDC resource because shoot is not trusted", func() {
		shoot.Annotations = map[string]string{
			"authentication.gardener.cloud/issuer":  "managed",
			"authentication.gardener.cloud/trusted": "false",
		}
		Expect(fakeClient.Create(ctx, shoot)).To(Succeed())
		Expect(fakeClient.Get(ctx, shootObjectKey, shoot)).To(Succeed())
		// Create OIDC resource that should be deleted
		Expect(fakeClient.Create(ctx, oidc)).To(Succeed())
		Expect(fakeClient.Get(ctx, oidcObjectKey, oidc)).To(Succeed())

		res, err := reconciler.Reconcile(ctx, reconcile.Request{NamespacedName: shootObjectKey})
		Expect(err).ToNot(HaveOccurred())
		Expect(res).To(Equal(ctrl.Result{}))

		Expect(fakeClient.Get(ctx, oidcObjectKey, oidc)).To(BeNotFoundError())
	})

	It("should trigger deletion of OIDC resource because shoot is not trusted but OIDC is not found", func() {
		shoot.Annotations = map[string]string{
			"authentication.gardener.cloud/issuer":  "managed",
			"authentication.gardener.cloud/trusted": "false",
		}
		Expect(fakeClient.Create(ctx, shoot)).To(Succeed())
		Expect(fakeClient.Get(ctx, shootObjectKey, shoot)).To(Succeed())

		res, err := reconciler.Reconcile(ctx, reconcile.Request{NamespacedName: shootObjectKey})
		Expect(err).ToNot(HaveOccurred())
		Expect(res).To(Equal(ctrl.Result{}))

		Expect(fakeClient.Get(ctx, oidcObjectKey, oidc)).To(BeNotFoundError())
	})

	It("should delete OIDC resource because shoot is being deleted", func() {
		// Adding a finalizer to simulate that the shoot is being deleted and not yet fully deleted to trigger the shoot.DeletionTimestamp check
		shoot.Finalizers = []string{"some/finalizer"}
		Expect(fakeClient.Create(ctx, shoot)).To(Succeed())
		Expect(fakeClient.Get(ctx, shootObjectKey, shoot)).To(Succeed())
		// Create OIDC resource that should be deleted
		Expect(fakeClient.Create(ctx, oidc)).To(Succeed())
		Expect(fakeClient.Get(ctx, oidcObjectKey, oidc)).To(Succeed())

		res, err := reconciler.Reconcile(ctx, reconcile.Request{NamespacedName: shootObjectKey})
		Expect(err).ToNot(HaveOccurred())
		Expect(res).To(Equal(ctrl.Result{}))

		Expect(fakeClient.Delete(ctx, shoot)).To(Succeed())
		res, err = reconciler.Reconcile(ctx, reconcile.Request{NamespacedName: shootObjectKey})
		Expect(err).ToNot(HaveOccurred())
		Expect(res).To(Equal(ctrl.Result{}))

		Expect(fakeClient.Get(ctx, oidcObjectKey, oidc)).To(BeNotFoundError())
		Expect(fakeClient.Get(ctx, shootObjectKey, shoot)).To(Succeed())
	})

	It("should do nothing because shoot annotation is invalid", func() {
		shoot.Annotations = map[string]string{
			"authentication.gardener.cloud/issuer":  "managed",
			"authentication.gardener.cloud/trusted": "foo",
		}
		Expect(fakeClient.Create(ctx, shoot)).To(Succeed())
		Expect(fakeClient.Get(ctx, shootObjectKey, shoot)).To(Succeed())

		res, err := reconciler.Reconcile(ctx, reconcile.Request{NamespacedName: shootObjectKey})
		Expect(err).ToNot(HaveOccurred())
		Expect(res).To(Equal(ctrl.Result{}))

		Expect(fakeClient.Get(ctx, oidcObjectKey, oidc)).To(BeNotFoundError())
	})

	It("should trigger deletion of OIDC resource because shoot is missing", func() {
		res, err := reconciler.Reconcile(ctx, reconcile.Request{NamespacedName: types.NamespacedName{Name: shoot.Name, Namespace: shoot.Namespace}})
		Expect(err).ToNot(HaveOccurred())
		Expect(res).To(Equal(ctrl.Result{}))

		Expect(fakeClient.Get(ctx, oidcObjectKey, oidc)).To(BeNotFoundError())
	})

	It("should do nothing because shoot status.advertisedAddresses is empty", func() {
		shoot.Status.AdvertisedAddresses = nil
		Expect(fakeClient.Create(ctx, shoot)).To(Succeed())
		Expect(fakeClient.Get(ctx, shootObjectKey, shoot)).To(Succeed())

		res, err := reconciler.Reconcile(ctx, reconcile.Request{NamespacedName: shootObjectKey})
		Expect(err).ToNot(HaveOccurred())
		Expect(res).To(Equal(ctrl.Result{}))

		Expect(fakeClient.Get(ctx, oidcObjectKey, oidc)).To(BeNotFoundError())
	})

	It("should do nothing because shoot status.advertisedAddresses has no service account issuer", func() {
		shoot.Status.AdvertisedAddresses = []gardencorev1beta1.ShootAdvertisedAddress{
			{
				Name: "foo",
				URL:  "https://foo",
			},
		}
		Expect(fakeClient.Create(ctx, shoot)).To(Succeed())
		Expect(fakeClient.Get(ctx, shootObjectKey, shoot)).To(Succeed())

		res, err := reconciler.Reconcile(ctx, reconcile.Request{NamespacedName: shootObjectKey})
		Expect(err).ToNot(HaveOccurred())
		Expect(res).To(Equal(ctrl.Result{}))

		Expect(fakeClient.Get(ctx, oidcObjectKey, oidc)).To(BeNotFoundError())
	})
})
