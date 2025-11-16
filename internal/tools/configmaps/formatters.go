package configmaps

import (
	"fmt"
	"strings"

	corev1 "k8s.io/api/core/v1"
)

func formatConfigMapsList(configMaps *corev1.ConfigMapList) string {
	var sb strings.Builder

	sb.WriteString(fmt.Sprintf("Total ConfigMaps: %d\n\n", len(configMaps.Items)))

	for _, cm := range configMaps.Items {
		sb.WriteString(fmt.Sprintf("Name: %s\n", cm.Name))
		sb.WriteString(fmt.Sprintf("Namespace: %s\n", cm.Namespace))
		sb.WriteString(fmt.Sprintf("Data Keys: %d\n", len(cm.Data)))
		sb.WriteString("\n---\n\n")
	}

	return sb.String()
}

func formatConfigMapDetails(cm *corev1.ConfigMap) string {
	var sb strings.Builder

	sb.WriteString(fmt.Sprintf("ConfigMap: %s\n", cm.Name))
	sb.WriteString(fmt.Sprintf("Namespace: %s\n", cm.Namespace))

	if len(cm.Data) > 0 {
		sb.WriteString("\nData:\n")
		for k, v := range cm.Data {
			sb.WriteString(fmt.Sprintf("  %s:\n", k))
			sb.WriteString(fmt.Sprintf("    %s\n", v))
		}
	}

	if len(cm.BinaryData) > 0 {
		sb.WriteString("\nBinary Data Keys:\n")
		for k := range cm.BinaryData {
			sb.WriteString(fmt.Sprintf("  - %s\n", k))
		}
	}

	return sb.String()
}
