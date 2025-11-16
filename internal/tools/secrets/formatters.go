package secrets

import (
	"fmt"
	"strings"

	corev1 "k8s.io/api/core/v1"
)

func formatSecretsList(secrets *corev1.SecretList) string {
	var sb strings.Builder

	sb.WriteString(fmt.Sprintf("Total Secrets: %d\n\n", len(secrets.Items)))

	for _, secret := range secrets.Items {
		sb.WriteString(fmt.Sprintf("Name: %s\n", secret.Name))
		sb.WriteString(fmt.Sprintf("Namespace: %s\n", secret.Namespace))
		sb.WriteString(fmt.Sprintf("Type: %s\n", secret.Type))
		sb.WriteString(fmt.Sprintf("Data Keys: %d\n", len(secret.Data)))
		sb.WriteString("\n---\n\n")
	}

	return sb.String()
}

func formatSecretDetails(secret *corev1.Secret) string {
	var sb strings.Builder

	sb.WriteString(fmt.Sprintf("Secret: %s\n", secret.Name))
	sb.WriteString(fmt.Sprintf("Namespace: %s\n", secret.Namespace))
	sb.WriteString(fmt.Sprintf("Type: %s\n", secret.Type))

	if len(secret.Data) > 0 {
		sb.WriteString("\nData Keys (values masked for security):\n")
		for k := range secret.Data {
			sb.WriteString(fmt.Sprintf("  - %s\n", k))
		}
	}

	return sb.String()
}
