package services

import (
	"fmt"
	"strings"

	corev1 "k8s.io/api/core/v1"
)

func formatServicesList(services *corev1.ServiceList) string {
	var sb strings.Builder

	sb.WriteString(fmt.Sprintf("Total Services: %d\n\n", len(services.Items)))

	for _, svc := range services.Items {
		sb.WriteString(fmt.Sprintf("Name: %s\n", svc.Name))
		sb.WriteString(fmt.Sprintf("Namespace: %s\n", svc.Namespace))
		sb.WriteString(fmt.Sprintf("Type: %s\n", svc.Spec.Type))
		sb.WriteString(fmt.Sprintf("Cluster IP: %s\n", svc.Spec.ClusterIP))
		if len(svc.Spec.Ports) > 0 {
			sb.WriteString("Ports:\n")
			for _, port := range svc.Spec.Ports {
				sb.WriteString(fmt.Sprintf("  - %s:%d/%s\n", port.Name, port.Port, port.Protocol))
			}
		}
		sb.WriteString("\n---\n\n")
	}

	return sb.String()
}

func formatServiceDetails(svc *corev1.Service) string {
	var sb strings.Builder

	sb.WriteString(fmt.Sprintf("Service: %s\n", svc.Name))
	sb.WriteString(fmt.Sprintf("Namespace: %s\n", svc.Namespace))
	sb.WriteString(fmt.Sprintf("Type: %s\n", svc.Spec.Type))
	sb.WriteString(fmt.Sprintf("Cluster IP: %s\n", svc.Spec.ClusterIP))

	if len(svc.Spec.ExternalIPs) > 0 {
		sb.WriteString(fmt.Sprintf("External IPs: %v\n", svc.Spec.ExternalIPs))
	}

	if len(svc.Labels) > 0 {
		sb.WriteString("\nLabels:\n")
		for k, v := range svc.Labels {
			sb.WriteString(fmt.Sprintf("  %s: %s\n", k, v))
		}
	}

	if len(svc.Spec.Selector) > 0 {
		sb.WriteString("\nSelector:\n")
		for k, v := range svc.Spec.Selector {
			sb.WriteString(fmt.Sprintf("  %s: %s\n", k, v))
		}
	}

	if len(svc.Spec.Ports) > 0 {
		sb.WriteString("\nPorts:\n")
		for _, port := range svc.Spec.Ports {
			sb.WriteString(fmt.Sprintf("  Name: %s\n", port.Name))
			sb.WriteString(fmt.Sprintf("  Protocol: %s\n", port.Protocol))
			sb.WriteString(fmt.Sprintf("  Port: %d\n", port.Port))
			sb.WriteString(fmt.Sprintf("  TargetPort: %v\n", port.TargetPort))
			if port.NodePort != 0 {
				sb.WriteString(fmt.Sprintf("  NodePort: %d\n", port.NodePort))
			}
			sb.WriteString("\n")
		}
	}

	return sb.String()
}
