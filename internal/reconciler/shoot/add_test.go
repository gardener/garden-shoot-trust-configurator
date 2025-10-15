// SPDX-FileCopyrightText: SAP SE or an SAP affiliate company and Gardener contributors
//
// SPDX-License-Identifier: Apache-2.0

package reconciler_test

import (
	"time"

	gardencorev1beta1 "github.com/gardener/gardener/pkg/apis/core/v1beta1"
	v1beta1constants "github.com/gardener/gardener/pkg/apis/core/v1beta1/constants"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"

	shootcontroller "github.com/gardener/garden-shoot-trust-configurator/internal/reconciler/shoot"
)

var _ = Describe("ShootPredicate", func() {
	const (
		shootName      = "my-shoot"
		shootNamespace = "garden-abc"
	)

	var (
		reconciler *shootcontroller.Reconciler
		shoot      *gardencorev1beta1.Shoot
	)

	BeforeEach(func() {
		reconciler = &shootcontroller.Reconciler{}
		shoot = &gardencorev1beta1.Shoot{
			ObjectMeta: metav1.ObjectMeta{
				Name:      shootName,
				Namespace: shootNamespace,
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

		It("should return true if the old shoot was trusted but the new one has the annotation removed", func() {
			oldShoot := shoot
			newShoot := shoot.DeepCopy()
			newShoot.Annotations = map[string]string{
				"authentication.gardener.cloud/issuer": "managed",
			}
			Expect(reconciler.IsRelevantShootUpdate(oldShoot, newShoot)).To(BeTrue())
		})

		It("should return false if neither the old nor the new shoot have the expected annotation", func() {
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

		It("should return true if the shoot's service-account-issuer has been added", func() {
			oldShoot := &gardencorev1beta1.Shoot{
				ObjectMeta: metav1.ObjectMeta{
					Name:      shootName,
					Namespace: shootNamespace,
					Annotations: map[string]string{
						"authentication.gardener.cloud/issuer":  "managed",
						"authentication.gardener.cloud/trusted": "true",
					},
				},
				Status: gardencorev1beta1.ShootStatus{
					AdvertisedAddresses: []gardencorev1beta1.ShootAdvertisedAddress{},
				},
			}
			newShoot := shoot
			Expect(reconciler.IsRelevantShootUpdate(oldShoot, newShoot)).To(BeTrue())
		})

		It("should return true if the shoot's service-account-issuer has been changed", func() {
			oldShoot := shoot
			newShoot := shoot.DeepCopy()
			newShoot.Status.AdvertisedAddresses = []gardencorev1beta1.ShootAdvertisedAddress{
				{
					Name: v1beta1constants.AdvertisedAddressServiceAccountIssuer,
					URL:  "https://shoot/new-issuer",
				},
			}
			Expect(reconciler.IsRelevantShootUpdate(oldShoot, newShoot)).To(BeTrue())
		})

		It("should return false if the shoot's service-account-issuer has not been populated", func() {
			oldShoot := shoot.DeepCopy()
			oldShoot.Status.AdvertisedAddresses = []gardencorev1beta1.ShootAdvertisedAddress{}
			newShoot := shoot.DeepCopy()
			newShoot.Status.AdvertisedAddresses = []gardencorev1beta1.ShootAdvertisedAddress{}
			Expect(reconciler.IsRelevantShootUpdate(oldShoot, newShoot)).To(BeFalse())
		})

		It("should return true if a shoot is updated to be deleted", func() {
			oldShoot := shoot
			newShoot := shoot.DeepCopy()
			newShoot.Finalizers = []string{"some/finalizer"}
			newShoot.DeletionTimestamp = &metav1.Time{Time: time.Now()}
			Expect(reconciler.IsRelevantShootUpdate(oldShoot, newShoot)).To(BeTrue())
		})

		It("should return false if the shoot is updated but remains trusted", func() {
			oldShoot := shoot
			newShoot := shoot.DeepCopy()
			newShoot.Labels = map[string]string{"some": "label"}
			Expect(reconciler.IsRelevantShootUpdate(oldShoot, newShoot)).To(BeFalse())
		})

		It("should return false if new object is not shoot", func() {
			oldObj := &gardencorev1beta1.Shoot{}
			newObj := &gardencorev1beta1.Seed{}
			Expect(reconciler.IsRelevantShootUpdate(oldObj, newObj)).To(BeFalse())
		})

		It("should return false if both old and new objects are not shoots", func() {
			oldObj := &gardencorev1beta1.Seed{}
			newObj := &gardencorev1beta1.Seed{}
			Expect(reconciler.IsRelevantShootUpdate(oldObj, newObj)).To(BeFalse())
		})
	})
})
