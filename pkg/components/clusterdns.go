package components

// ClusterDNSPlacement describe the cluster DNS information, such as namespace and common label
type ClusterDNSPlacement struct {
	// Namespace is the Namespace name where cluster DNS endpoints reside.
	Namespace string
	// LabelSelectorKey is the cluster DNS pods label key, to be used by a label selector.
	LabelSelectorKey string
	// LabelSelectorValue is the cluster DNS pods label value, to be used by a label selector.
	LabelSelectorValue string
}
