package routes

import (
	"fmt"
	"strings"

	routev1 "github.com/openshift/api/route/v1"
)

func formatRoutesList(routes *routev1.RouteList) string {
	var sb strings.Builder

	sb.WriteString(fmt.Sprintf("Total Routes: %d\n\n", len(routes.Items)))

	for _, route := range routes.Items {
		sb.WriteString(fmt.Sprintf("Name: %s\n", route.Name))
		sb.WriteString(fmt.Sprintf("Namespace: %s\n", route.Namespace))
		sb.WriteString(fmt.Sprintf("Host: %s\n", route.Spec.Host))
		sb.WriteString(fmt.Sprintf("Path: %s\n", route.Spec.Path))
		sb.WriteString(fmt.Sprintf("Service: %s\n", route.Spec.To.Name))
		if route.Spec.TLS != nil {
			sb.WriteString(fmt.Sprintf("TLS: %s\n", route.Spec.TLS.Termination))
		}
		sb.WriteString("\n---\n\n")
	}

	return sb.String()
}

func formatRouteDetails(route *routev1.Route) string {
	var sb strings.Builder

	sb.WriteString(fmt.Sprintf("Route: %s\n", route.Name))
	sb.WriteString(fmt.Sprintf("Namespace: %s\n", route.Namespace))
	sb.WriteString(fmt.Sprintf("Host: %s\n", route.Spec.Host))
	sb.WriteString(fmt.Sprintf("Path: %s\n", route.Spec.Path))

	sb.WriteString("\nTarget:\n")
	sb.WriteString(fmt.Sprintf("  Kind: %s\n", route.Spec.To.Kind))
	sb.WriteString(fmt.Sprintf("  Name: %s\n", route.Spec.To.Name))
	sb.WriteString(fmt.Sprintf("  Weight: %d\n", *route.Spec.To.Weight))

	if route.Spec.Port != nil {
		sb.WriteString(fmt.Sprintf("\nPort: %s\n", route.Spec.Port.TargetPort.String()))
	}

	if route.Spec.TLS != nil {
		sb.WriteString("\nTLS:\n")
		sb.WriteString(fmt.Sprintf("  Termination: %s\n", route.Spec.TLS.Termination))
		sb.WriteString(fmt.Sprintf("  Insecure Edge Termination Policy: %s\n", route.Spec.TLS.InsecureEdgeTerminationPolicy))
	}

	if len(route.Status.Ingress) > 0 {
		sb.WriteString("\nIngress Status:\n")
		for _, ing := range route.Status.Ingress {
			sb.WriteString(fmt.Sprintf("  Host: %s\n", ing.Host))
			sb.WriteString(fmt.Sprintf("  Router Name: %s\n", ing.RouterName))
		}
	}

	return sb.String()
}
