package components

import (
	corev1 "k8s.io/api/core/v1"
	netv1 "k8s.io/api/networking/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/util/intstr"
	"k8s.io/utils/ptr"
)

func GetNetworkPolices(
	targetNamespace string,
	cnaoLabelSelector map[string]string,
	clusterDNSPlacement ClusterDNSPlacement,
) []netv1.NetworkPolicy {
	return []netv1.NetworkPolicy{
		{
			TypeMeta: metav1.TypeMeta{
				APIVersion: "networking.k8s.io/v1",
				Kind:       "NetworkPolicy",
			},
			ObjectMeta: metav1.ObjectMeta{
				Name:      "allow-egress-to-dns",
				Namespace: targetNamespace,
			},
			Spec: netv1.NetworkPolicySpec{
				PodSelector: metav1.LabelSelector{
					MatchLabels: cnaoLabelSelector,
				},
				PolicyTypes: []netv1.PolicyType{
					netv1.PolicyTypeEgress,
				},
				Egress: []netv1.NetworkPolicyEgressRule{
					{
						To: []netv1.NetworkPolicyPeer{
							{
								NamespaceSelector: &metav1.LabelSelector{
									MatchLabels: map[string]string{
										corev1.LabelMetadataName: clusterDNSPlacement.Namespace,
									},
								},
								PodSelector: &metav1.LabelSelector{
									MatchLabels: map[string]string{
										clusterDNSPlacement.LabelSelectorKey: clusterDNSPlacement.LabelSelectorValue,
									},
								},
							},
						},
						Ports: []netv1.NetworkPolicyPort{
							{
								Protocol: ptr.To(corev1.ProtocolTCP),
								Port:     ptr.To(intstr.FromString("dns-tcp")),
							},
							{
								Protocol: ptr.To(corev1.ProtocolUDP),
								Port:     ptr.To(intstr.FromString("dns")),
							},
						},
					},
				},
			},
		},
		{
			TypeMeta: metav1.TypeMeta{
				APIVersion: "networking.k8s.io/v1",
				Kind:       "NetworkPolicy",
			},
			ObjectMeta: metav1.ObjectMeta{
				Name:      "allow-egress-to-kube-apiserver",
				Namespace: targetNamespace,
			},
			Spec: netv1.NetworkPolicySpec{
				PodSelector: metav1.LabelSelector{
					MatchLabels: cnaoLabelSelector,
				},
				PolicyTypes: []netv1.PolicyType{
					netv1.PolicyTypeEgress,
				},
				Egress: []netv1.NetworkPolicyEgressRule{
					{
						Ports: []netv1.NetworkPolicyPort{
							{
								Protocol: ptr.To(corev1.ProtocolTCP),
								Port:     ptr.To(intstr.FromInt32(6443)),
							},
						},
					},
				},
			},
		},
		{
			TypeMeta: metav1.TypeMeta{
				APIVersion: "networking.k8s.io/v1",
				Kind:       "NetworkPolicy",
			},
			ObjectMeta: metav1.ObjectMeta{
				Name:      "allow-ingress-to-metrics-endpoint",
				Namespace: targetNamespace,
			},
			Spec: netv1.NetworkPolicySpec{
				PodSelector: metav1.LabelSelector{
					MatchLabels: cnaoLabelSelector,
				},
				PolicyTypes: []netv1.PolicyType{
					netv1.PolicyTypeIngress,
				},
				Ingress: []netv1.NetworkPolicyIngressRule{
					{
						Ports: []netv1.NetworkPolicyPort{
							{
								Protocol: ptr.To(corev1.ProtocolTCP),
								Port:     ptr.To(intstr.FromString("metrics")),
							},
						},
					},
				},
			},
		},
	}
}
