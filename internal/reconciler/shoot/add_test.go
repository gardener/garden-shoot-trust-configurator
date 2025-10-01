// SPDX-FileCopyrightText: SAP SE or an SAP affiliate company and Gardener contributors
//
// SPDX-License-Identifier: Apache-2.0

package reconciler_test

import (
	gardencorev1beta1 "github.com/gardener/gardener/pkg/apis/core/v1beta1"
	"github.com/gardener/gardener/pkg/client/kubernetes"
	authenticationv1alpha1 "github.com/gardener/oidc-webhook-authenticator/apis/authentication/v1alpha1"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"

	shootreconciler "github.com/gardener/garden-shoot-trust-configurator/internal/reconciler/shoot"
)

var _ = Describe("ShootPredicate", func() {
	const (
		shootName      = "my-shoot"
		shootNamespace = "garden-abc"
	)

	var (
		reconciler *shootreconciler.Reconciler

		shoot *gardencorev1beta1.Shoot
	)

	BeforeEach(func() {
		scheme := runtime.NewScheme()
		Expect(kubernetes.AddGardenSchemeToScheme(scheme)).To(Succeed())
		Expect(authenticationv1alpha1.AddToScheme(scheme)).To(Succeed())

		reconciler = &shootreconciler.Reconciler{}

		shoot = &gardencorev1beta1.Shoot{
			ObjectMeta: metav1.ObjectMeta{
				Name:      shootName,
				Namespace: shootNamespace,
				Annotations: map[string]string{
					"authentication.gardener.cloud/issuer":  "managed",
					"authentication.gardener.cloud/trusted": "true",
				},
			},
		}
	})

	Describe("IsRelevantShoot", func() {
		It("should return true for a shoot with the 'authentication.gardener.cloud/trusted' annotation set to 'true'", func() {
			Expect(reconciler.IsRelevantShoot(shoot)).To(BeTrue())
		})

		It("should return true for a shoot with the 'authentication.gardener.cloud/trusted' annotation set to 'True'", func() {
			shoot.Annotations["authentication.gardener.cloud/trusted"] = "True"
			Expect(reconciler.IsRelevantShoot(shoot)).To(BeTrue())
		})

		It("should return false for a shoot with the 'authentication.gardener.cloud/trusted' annotation set to 'false'", func() {
			shoot.Annotations["authentication.gardener.cloud/trusted"] = "false"
			Expect(reconciler.IsRelevantShoot(shoot)).To(BeFalse())
		})

		It("should return false for a shoot without the 'authentication.gardener.cloud/trusted' annotation", func() {
			shoot.Annotations = map[string]string{
				"authentication.gardener.cloud/issuer": "managed",
			}
			Expect(reconciler.IsRelevantShoot(shoot)).To(BeFalse())
		})

		It("should return false for a shoot which doesn't have managed issuer", func() {
			shoot.Annotations = map[string]string{}
			Expect(reconciler.IsRelevantShoot(shoot)).To(BeFalse())
		})

		It("should return false if object is not shoot", func() {
			nonShoot := &gardencorev1beta1.Seed{}
			Expect(reconciler.IsRelevantShoot(nonShoot)).To(BeFalse())
		})
	})

	Describe("IsRelevantShootUpdate", func() {
		It("should return true if the shoot is updated to be trusted", func() {
			oldShoot := &gardencorev1beta1.Shoot{
				ObjectMeta: metav1.ObjectMeta{
					Name:      shootName,
					Namespace: shootNamespace,
				},
			}
			newShoot := shoot
			Expect(reconciler.IsRelevantShootUpdate(oldShoot, newShoot)).To(BeTrue())
		})

		It("should return true if the old object was trusted but the new one has the annotation removed", func() {
			oldShoot := shoot
			newShoot := shoot.DeepCopy()
			newShoot.Annotations = map[string]string{
				"authentication.gardener.cloud/issuer": "managed",
			}
			Expect(reconciler.IsRelevantShootUpdate(oldShoot, newShoot)).To(BeTrue())
		})

		It("should return false if neither the old nor the new object have the expected annotation", func() {
			oldShoot := &gardencorev1beta1.Shoot{
				ObjectMeta: metav1.ObjectMeta{
					Name:      shootName,
					Namespace: shootNamespace,
				},
			}
			newShoot := &gardencorev1beta1.Shoot{
				ObjectMeta: metav1.ObjectMeta{
					Name:      shootName,
					Namespace: shootNamespace,
				},
			}
			Expect(reconciler.IsRelevantShootUpdate(oldShoot, newShoot)).To(BeFalse())
		})

		It("should return false if both old and new objects are not shoots", func() {
			oldObj := &gardencorev1beta1.Seed{}
			newObj := &gardencorev1beta1.Seed{}
			Expect(reconciler.IsRelevantShootUpdate(oldObj, newObj)).To(BeFalse())
		})
	})
})
