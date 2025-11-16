package deployments

import (
	"fmt"
	"strings"

	appsv1 "k8s.io/api/apps/v1"
)

func formatDeploymentsList(deployments *appsv1.DeploymentList) string {
	var sb strings.Builder

	sb.WriteString(fmt.Sprintf("Total Deployments: %d\n\n", len(deployments.Items)))

	for _, deploy := range deployments.Items {
		sb.WriteString(fmt.Sprintf("Name: %s\n", deploy.Name))
		sb.WriteString(fmt.Sprintf("Namespace: %s\n", deploy.Namespace))
		sb.WriteString(fmt.Sprintf("Replicas: %d/%d\n", deploy.Status.ReadyReplicas, *deploy.Spec.Replicas))
		sb.WriteString(fmt.Sprintf("Available: %d\n", deploy.Status.AvailableReplicas))
		sb.WriteString(fmt.Sprintf("Updated: %d\n", deploy.Status.UpdatedReplicas))
		sb.WriteString("\n---\n\n")
	}

	return sb.String()
}

func formatDeploymentDetails(deploy *appsv1.Deployment) string {
	var sb strings.Builder

	sb.WriteString(fmt.Sprintf("Deployment: %s\n", deploy.Name))
	sb.WriteString(fmt.Sprintf("Namespace: %s\n", deploy.Namespace))
	sb.WriteString(fmt.Sprintf("Replicas: %d/%d\n", deploy.Status.ReadyReplicas, *deploy.Spec.Replicas))
	sb.WriteString(fmt.Sprintf("Strategy: %s\n", deploy.Spec.Strategy.Type))

	if len(deploy.Labels) > 0 {
		sb.WriteString("\nLabels:\n")
		for k, v := range deploy.Labels {
			sb.WriteString(fmt.Sprintf("  %s: %s\n", k, v))
		}
	}

	if len(deploy.Spec.Selector.MatchLabels) > 0 {
		sb.WriteString("\nSelector:\n")
		for k, v := range deploy.Spec.Selector.MatchLabels {
			sb.WriteString(fmt.Sprintf("  %s: %s\n", k, v))
		}
	}

	sb.WriteString("\nContainers:\n")
	for _, container := range deploy.Spec.Template.Spec.Containers {
		sb.WriteString(fmt.Sprintf("  Name: %s\n", container.Name))
		sb.WriteString(fmt.Sprintf("  Image: %s\n", container.Image))
	}

	if len(deploy.Status.Conditions) > 0 {
		sb.WriteString("\nConditions:\n")
		for _, cond := range deploy.Status.Conditions {
			sb.WriteString(fmt.Sprintf("  %s: %s (Reason: %s)\n",
				cond.Type, cond.Status, cond.Reason))
		}
	}

	return sb.String()
}
