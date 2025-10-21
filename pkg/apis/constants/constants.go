package constants

const (
	// AnnotationTrustedShoot is the annotation that marks a Shoot to be trusted in the Garden cluster.
	AnnotationTrustedShoot = "authentication.gardener.cloud/trusted"
	// LabelManagedByKey is a constant for a key of a label on an OIDC resource describing who is managing it.
	LabelManagedByKey = "app.kubernetes.io/managed-by"
	// LabelManagedByValue is a constant for a value of a label on a OIDC describing the value 'garden-shoot-trust-configurator'.
	LabelManagedByValue = "garden-shoot-trust-configurator"
	// Separator is the separator used in the OIDC resource name to separate namespace, name and uid of the shoot.
	Separator = "--"
)
