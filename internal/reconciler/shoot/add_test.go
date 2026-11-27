// SPDX-FileCopyrightText: Contributors to the Gardener project
//
// SPDX-License-Identifier: Apache-2.0

package reconciler_test

import (
	"context"
	"time"

	gardencorev1beta1 "github.com/gardener/gardener/pkg/apis/core/v1beta1"
	v1beta1constants "github.com/gardener/gardener/pkg/apis/core/v1beta1/constants"
	authenticationv1alpha1 "github.com/gardener/oidc-webhook-authenticator/apis/authentication/v1alpha1"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/types"
	"sigs.k8s.io/controller-runtime/pkg/event"
	"sigs.k8s.io/controller-runtime/pkg/reconcile"

	shootcontroller "github.com/gardener/garden-shoot-trust-configurator/internal/reconciler/shoot"
	"github.com/gardener/garden-shoot-trust-configurator/pkg/apis/constants"
)

var _ = Describe("#ShootPredicate", func() {
	const (
		shootName      = "my-shoot"
		shootNamespace = "garden-abc"
	)

	var shoot *gardencorev1beta1.Shoot

	BeforeEach(func() {
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

	Describe("#IsRelevantShoot", func() {
		It("should return true for a shoot with the 'authentication.gardener.cloud/trusted' annotation set to 'true'", func() {
			Expect(shootcontroller.IsRelevantShoot(shoot)).To(BeTrue())
		})

		It("should return true for a shoot with the 'authentication.gardener.cloud/trusted' annotation set to 'True'", func() {
			shoot.Annotations["authentication.gardener.cloud/trusted"] = "True"
			Expect(shootcontroller.IsRelevantShoot(shoot)).To(BeTrue())
		})

		It("should return false for a shoot with the 'authentication.gardener.cloud/trusted' annotation set to 'false'", func() {
			shoot.Annotations["authentication.gardener.cloud/trusted"] = "false"
			Expect(shootcontroller.IsRelevantShoot(shoot)).To(BeFalse())
		})

		It("should return false for a shoot without the 'authentication.gardener.cloud/trusted' annotation", func() {
			shoot.Annotations = map[string]string{
				"authentication.gardener.cloud/issuer": "managed",
			}
			Expect(shootcontroller.IsRelevantShoot(shoot)).To(BeFalse())
		})

		It("should return false for a shoot which doesn't have managed issuer", func() {
			shoot.Annotations = map[string]string{}
			Expect(shootcontroller.IsRelevantShoot(shoot)).To(BeFalse())
		})

		It("should return false if object is not shoot", func() {
			nonShoot := &gardencorev1beta1.Seed{}
			Expect(shootcontroller.IsRelevantShoot(nonShoot)).To(BeFalse())
		})
	})

	Describe("#IsRelevantShootUpdate", func() {
		It("should return true if the shoot is updated to be trusted", func() {
			oldShoot := &gardencorev1beta1.Shoot{
				ObjectMeta: metav1.ObjectMeta{
					Name:      shootName,
					Namespace: shootNamespace,
				},
			}
			newShoot := shoot
			Expect(shootcontroller.IsRelevantShootUpdate(oldShoot, newShoot)).To(BeTrue())
		})

		It("should return true if the old shoot was trusted but the new one has the annotation removed", func() {
			oldShoot := shoot
			newShoot := shoot.DeepCopy()
			newShoot.Annotations = map[string]string{
				"authentication.gardener.cloud/issuer": "managed",
			}
			Expect(shootcontroller.IsRelevantShootUpdate(oldShoot, newShoot)).To(BeTrue())
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
			Expect(shootcontroller.IsRelevantShootUpdate(oldShoot, newShoot)).To(BeFalse())
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
			Expect(shootcontroller.IsRelevantShootUpdate(oldShoot, newShoot)).To(BeTrue())
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
			Expect(shootcontroller.IsRelevantShootUpdate(oldShoot, newShoot)).To(BeTrue())
		})

		It("should return false if the shoot's service-account-issuer has not been populated", func() {
			oldShoot := shoot.DeepCopy()
			oldShoot.Status.AdvertisedAddresses = []gardencorev1beta1.ShootAdvertisedAddress{}
			newShoot := shoot.DeepCopy()
			newShoot.Status.AdvertisedAddresses = []gardencorev1beta1.ShootAdvertisedAddress{}
			Expect(shootcontroller.IsRelevantShootUpdate(oldShoot, newShoot)).To(BeFalse())
		})

		It("should return true if a shoot is updated to be deleted", func() {
			oldShoot := shoot
			newShoot := shoot.DeepCopy()
			newShoot.Finalizers = []string{"some/finalizer"}
			newShoot.DeletionTimestamp = &metav1.Time{Time: time.Now()}
			Expect(shootcontroller.IsRelevantShootUpdate(oldShoot, newShoot)).To(BeTrue())
		})

		It("should return false if the shoot is updated but remains trusted", func() {
			oldShoot := shoot
			newShoot := shoot.DeepCopy()
			newShoot.Labels = map[string]string{"some": "label"}
			Expect(shootcontroller.IsRelevantShootUpdate(oldShoot, newShoot)).To(BeFalse())
		})

		It("should return false if new object is not shoot", func() {
			oldObj := &gardencorev1beta1.Shoot{}
			newObj := &gardencorev1beta1.Seed{}
			Expect(shootcontroller.IsRelevantShootUpdate(oldObj, newObj)).To(BeFalse())
		})

		It("should return false if both old and new objects are not shoots", func() {
			oldObj := &gardencorev1beta1.Seed{}
			newObj := &gardencorev1beta1.Seed{}
			Expect(shootcontroller.IsRelevantShootUpdate(oldObj, newObj)).To(BeFalse())
		})
	})
})

var _ = Describe("#MapOIDCToShoot", func() {
	It("should map a well-formed OIDC resource name to its owning shoot", func() {
		oidc := &authenticationv1alpha1.OpenIDConnect{
			ObjectMeta: metav1.ObjectMeta{
				Name: "garden-abc" + constants.Separator + "my-shoot" + constants.Separator + "some-uid",
			},
		}
		Expect(shootcontroller.MapOIDCToShoot(context.Background(), oidc)).To(ConsistOf(
			reconcile.Request{NamespacedName: types.NamespacedName{Namespace: "garden-abc", Name: "my-shoot"}},
		))
	})

	It("should return no requests for a name that does not follow the convention", func() {
		oidc := &authenticationv1alpha1.OpenIDConnect{
			ObjectMeta: metav1.ObjectMeta{Name: "garden-abc" + constants.Separator + "my-shoot"},
		}
		Expect(shootcontroller.MapOIDCToShoot(context.Background(), oidc)).To(BeEmpty())
	})
})

var _ = Describe("#ManagedByOIDCPredicate", func() {
	It("should admit an OIDC resource carrying the managed-by label", func() {
		oidc := &authenticationv1alpha1.OpenIDConnect{
			ObjectMeta: metav1.ObjectMeta{
				Labels: map[string]string{constants.LabelManagedByKey: constants.LabelManagedByValue},
			},
		}
		Expect(shootcontroller.ManagedByOIDCPredicate().Create(event.CreateEvent{Object: oidc})).To(BeTrue())
	})

	It("should reject an OIDC resource without the managed-by label", func() {
		oidc := &authenticationv1alpha1.OpenIDConnect{}
		Expect(shootcontroller.ManagedByOIDCPredicate().Create(event.CreateEvent{Object: oidc})).To(BeFalse())
	})

	It("should reject an OIDC resource with a different managed-by value", func() {
		oidc := &authenticationv1alpha1.OpenIDConnect{
			ObjectMeta: metav1.ObjectMeta{
				Labels: map[string]string{constants.LabelManagedByKey: "someone-else"},
			},
		}
		Expect(shootcontroller.ManagedByOIDCPredicate().Create(event.CreateEvent{Object: oidc})).To(BeFalse())
	})
})
