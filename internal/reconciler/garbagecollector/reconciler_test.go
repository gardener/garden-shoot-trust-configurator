// SPDX-FileCopyrightText: SAP SE or an SAP affiliate company and Gardener contributors
//
// SPDX-License-Identifier: Apache-2.0

package garbagecollector_test

import (
	"context"
	"time"

	"github.com/gardener/gardener/pkg/client/kubernetes"
	authenticationv1alpha1 "github.com/gardener/oidc-webhook-authenticator/apis/authentication/v1alpha1"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/client/fake"
	logf "sigs.k8s.io/controller-runtime/pkg/log"
	logzap "sigs.k8s.io/controller-runtime/pkg/log/zap"
	"sigs.k8s.io/controller-runtime/pkg/reconcile"

	garbagecollectorcontroller "github.com/gardener/garden-shoot-trust-configurator/internal/reconciler/garbagecollector"
	configv1alpha1 "github.com/gardener/garden-shoot-trust-configurator/pkg/apis/config/v1alpha1"
)

var _ = Describe("Collector", func() {
	var (
		ctx = logf.IntoContext(context.Background(), logzap.New(logzap.WriteTo(GinkgoWriter)))

		gc         *garbagecollectorcontroller.Reconciler
		fakeClient client.Client
	)

	BeforeEach(func() {
		scheme := runtime.NewScheme()
		Expect(kubernetes.AddGardenSchemeToScheme(scheme)).To(Succeed())
		Expect(authenticationv1alpha1.AddToScheme(scheme)).To(Succeed())

		fakeClient = fake.NewClientBuilder().WithScheme(scheme).Build()
		gc = &garbagecollectorcontroller.Reconciler{
			Client: fakeClient,
			Config: configv1alpha1.GarbageCollectorControllerConfig{
				SyncPeriod:            &metav1.Duration{Duration: time.Hour},
				MinimumObjectLifetime: &metav1.Duration{Duration: time.Minute},
			},
		}
	})

	// TODO(theoddora): add tests for actual garbage collection logic
	It("should do nothing", func() {
		res, err := gc.Reconcile(ctx, reconcile.Request{})
		Expect(err).ToNot(HaveOccurred())
		Expect(res).To(Equal(ctrl.Result{RequeueAfter: gc.Config.SyncPeriod.Duration, Requeue: true}))
	})
})
