package ingress

import (
	"fmt"
	"strings"

	networkingv1 "k8s.io/api/networking/v1"
)

func formatIngressesList(ingresses *networkingv1.IngressList) string {
	var sb strings.Builder

	sb.WriteString(fmt.Sprintf("Total Ingresses: %d\n\n", len(ingresses.Items)))

	for _, ing := range ingresses.Items {
		sb.WriteString(fmt.Sprintf("Name: %s\n", ing.Name))
		sb.WriteString(fmt.Sprintf("Namespace: %s\n", ing.Namespace))
		if ing.Spec.IngressClassName != nil {
			sb.WriteString(fmt.Sprintf("Class: %s\n", *ing.Spec.IngressClassName))
		}
		sb.WriteString("\n---\n\n")
	}

	return sb.String()
}

func formatIngressDetails(ing *networkingv1.Ingress) string {
	var sb strings.Builder

	sb.WriteString(fmt.Sprintf("Ingress: %s\n", ing.Name))
	sb.WriteString(fmt.Sprintf("Namespace: %s\n", ing.Namespace))

	if ing.Spec.IngressClassName != nil {
		sb.WriteString(fmt.Sprintf("Class: %s\n", *ing.Spec.IngressClassName))
	}

	if len(ing.Spec.Rules) > 0 {
		sb.WriteString("\nRules:\n")
		for _, rule := range ing.Spec.Rules {
			sb.WriteString(fmt.Sprintf("  Host: %s\n", rule.Host))
			if rule.HTTP != nil {
				for _, path := range rule.HTTP.Paths {
					sb.WriteString(fmt.Sprintf("    Path: %s\n", path.Path))
					sb.WriteString(fmt.Sprintf("    Backend: %s:%d\n",
						path.Backend.Service.Name,
						path.Backend.Service.Port.Number))
				}
			}
		}
	}

	if len(ing.Spec.TLS) > 0 {
		sb.WriteString("\nTLS:\n")
		for _, tls := range ing.Spec.TLS {
			sb.WriteString(fmt.Sprintf("  Hosts: %v\n", tls.Hosts))
			sb.WriteString(fmt.Sprintf("  Secret: %s\n", tls.SecretName))
		}
	}

	return sb.String()
}
