package pods

import (
	"fmt"
	"strings"

	corev1 "k8s.io/api/core/v1"
)

func formatPodsList(pods *corev1.PodList) string {
	var sb strings.Builder

	sb.WriteString(fmt.Sprintf("Total Pods: %d\n\n", len(pods.Items)))

	for _, pod := range pods.Items {
		sb.WriteString(fmt.Sprintf("Name: %s\n", pod.Name))
		sb.WriteString(fmt.Sprintf("Namespace: %s\n", pod.Namespace))
		sb.WriteString(fmt.Sprintf("Status: %s\n", pod.Status.Phase))
		sb.WriteString(fmt.Sprintf("Node: %s\n", pod.Spec.NodeName))
		sb.WriteString(fmt.Sprintf("IP: %s\n", pod.Status.PodIP))

		if len(pod.Status.ContainerStatuses) > 0 {
			sb.WriteString("Containers:\n")
			for _, cs := range pod.Status.ContainerStatuses {
				sb.WriteString(fmt.Sprintf("  - %s (Ready: %v, Restarts: %d)\n",
					cs.Name, cs.Ready, cs.RestartCount))
			}
		}

		sb.WriteString("\n---\n\n")
	}

	return sb.String()
}

func formatPodDetails(pod *corev1.Pod) string {
	var sb strings.Builder

	sb.WriteString(fmt.Sprintf("Pod: %s\n", pod.Name))
	sb.WriteString(fmt.Sprintf("Namespace: %s\n", pod.Namespace))
	sb.WriteString(fmt.Sprintf("Status: %s\n", pod.Status.Phase))
	sb.WriteString(fmt.Sprintf("Node: %s\n", pod.Spec.NodeName))
	sb.WriteString(fmt.Sprintf("Pod IP: %s\n", pod.Status.PodIP))
	sb.WriteString(fmt.Sprintf("Host IP: %s\n", pod.Status.HostIP))
	sb.WriteString(fmt.Sprintf("QoS Class: %s\n", pod.Status.QOSClass))

	if len(pod.Labels) > 0 {
		sb.WriteString("\nLabels:\n")
		for k, v := range pod.Labels {
			sb.WriteString(fmt.Sprintf("  %s: %s\n", k, v))
		}
	}

	if len(pod.Annotations) > 0 {
		sb.WriteString("\nAnnotations:\n")
		for k, v := range pod.Annotations {
			sb.WriteString(fmt.Sprintf("  %s: %s\n", k, v))
		}
	}

	sb.WriteString("\nContainers:\n")
	for _, container := range pod.Spec.Containers {
		sb.WriteString(fmt.Sprintf("  Name: %s\n", container.Name))
		sb.WriteString(fmt.Sprintf("  Image: %s\n", container.Image))
		if len(container.Ports) > 0 {
			sb.WriteString(fmt.Sprintf("  Ports: %v\n", container.Ports))
		}
	}

	if len(pod.Status.ContainerStatuses) > 0 {
		sb.WriteString("\nContainer Statuses:\n")
		for _, cs := range pod.Status.ContainerStatuses {
			sb.WriteString(fmt.Sprintf("  %s:\n", cs.Name))
			sb.WriteString(fmt.Sprintf("    Ready: %v\n", cs.Ready))
			sb.WriteString(fmt.Sprintf("    Restart Count: %d\n", cs.RestartCount))
			sb.WriteString(fmt.Sprintf("    Image: %s\n", cs.Image))
		}
	}

	if len(pod.Status.Conditions) > 0 {
		sb.WriteString("\nConditions:\n")
		for _, cond := range pod.Status.Conditions {
			sb.WriteString(fmt.Sprintf("  %s: %s (Reason: %s)\n",
				cond.Type, cond.Status, cond.Reason))
		}
	}

	return sb.String()
}
