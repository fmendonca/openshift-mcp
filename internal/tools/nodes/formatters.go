package nodes

import (
	"fmt"
	"strings"

	corev1 "k8s.io/api/core/v1"
)

func formatNodesList(nodes *corev1.NodeList) string {
	var sb strings.Builder

	sb.WriteString(fmt.Sprintf("Total Nodes: %d\n\n", len(nodes.Items)))

	for _, node := range nodes.Items {
		sb.WriteString(fmt.Sprintf("Name: %s\n", node.Name))

		for _, addr := range node.Status.Addresses {
			if addr.Type == corev1.NodeInternalIP {
				sb.WriteString(fmt.Sprintf("Internal IP: %s\n", addr.Address))
			}
		}

		for _, cond := range node.Status.Conditions {
			if cond.Type == corev1.NodeReady {
				sb.WriteString(fmt.Sprintf("Ready: %s\n", cond.Status))
			}
		}

		sb.WriteString(fmt.Sprintf("Kubelet Version: %s\n", node.Status.NodeInfo.KubeletVersion))
		sb.WriteString("\n---\n\n")
	}

	return sb.String()
}

func formatNodeDetails(node *corev1.Node) string {
	var sb strings.Builder

	sb.WriteString(fmt.Sprintf("Node: %s\n", node.Name))

	sb.WriteString("\nAddresses:\n")
	for _, addr := range node.Status.Addresses {
		sb.WriteString(fmt.Sprintf("  %s: %s\n", addr.Type, addr.Address))
	}

	sb.WriteString("\nSystem Info:\n")
	sb.WriteString(fmt.Sprintf("  OS: %s\n", node.Status.NodeInfo.OSImage))
	sb.WriteString(fmt.Sprintf("  Kernel: %s\n", node.Status.NodeInfo.KernelVersion))
	sb.WriteString(fmt.Sprintf("  Container Runtime: %s\n", node.Status.NodeInfo.ContainerRuntimeVersion))
	sb.WriteString(fmt.Sprintf("  Kubelet: %s\n", node.Status.NodeInfo.KubeletVersion))
	sb.WriteString(fmt.Sprintf("  Kube-Proxy: %s\n", node.Status.NodeInfo.KubeProxyVersion))

	sb.WriteString("\nCapacity:\n")
	sb.WriteString(fmt.Sprintf("  CPU: %s\n", node.Status.Capacity.Cpu()))
	sb.WriteString(fmt.Sprintf("  Memory: %s\n", node.Status.Capacity.Memory()))
	sb.WriteString(fmt.Sprintf("  Pods: %s\n", node.Status.Capacity.Pods()))

	sb.WriteString("\nAllocatable:\n")
	sb.WriteString(fmt.Sprintf("  CPU: %s\n", node.Status.Allocatable.Cpu()))
	sb.WriteString(fmt.Sprintf("  Memory: %s\n", node.Status.Allocatable.Memory()))
	sb.WriteString(fmt.Sprintf("  Pods: %s\n", node.Status.Allocatable.Pods()))

	sb.WriteString("\nConditions:\n")
	for _, cond := range node.Status.Conditions {
		sb.WriteString(fmt.Sprintf("  %s: %s (Reason: %s)\n", cond.Type, cond.Status, cond.Reason))
	}

	if len(node.Labels) > 0 {
		sb.WriteString("\nLabels:\n")
		for k, v := range node.Labels {
			sb.WriteString(fmt.Sprintf("  %s: %s\n", k, v))
		}
	}

	return sb.String()
}
