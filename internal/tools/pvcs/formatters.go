package pvcs

import (
	"fmt"
	"strings"

	corev1 "k8s.io/api/core/v1"
)

func formatPVCsList(pvcs *corev1.PersistentVolumeClaimList) string {
	var sb strings.Builder

	sb.WriteString(fmt.Sprintf("Total PVCs: %d\n\n", len(pvcs.Items)))

	for _, pvc := range pvcs.Items {
		sb.WriteString(fmt.Sprintf("Name: %s\n", pvc.Name))
		sb.WriteString(fmt.Sprintf("Namespace: %s\n", pvc.Namespace))
		sb.WriteString(fmt.Sprintf("Status: %s\n", pvc.Status.Phase))
		if pvc.Spec.StorageClassName != nil {
			sb.WriteString(fmt.Sprintf("Storage Class: %s\n", *pvc.Spec.StorageClassName))
		}
		sb.WriteString(fmt.Sprintf("Capacity: %s\n", pvc.Status.Capacity.Storage()))
		sb.WriteString("\n---\n\n")
	}

	return sb.String()
}

func formatPVCDetails(pvc *corev1.PersistentVolumeClaim) string {
	var sb strings.Builder

	sb.WriteString(fmt.Sprintf("PVC: %s\n", pvc.Name))
	sb.WriteString(fmt.Sprintf("Namespace: %s\n", pvc.Namespace))
	sb.WriteString(fmt.Sprintf("Status: %s\n", pvc.Status.Phase))

	if pvc.Spec.StorageClassName != nil {
		sb.WriteString(fmt.Sprintf("Storage Class: %s\n", *pvc.Spec.StorageClassName))
	}

	sb.WriteString(fmt.Sprintf("Access Modes: %v\n", pvc.Spec.AccessModes))
	sb.WriteString(fmt.Sprintf("Requested: %s\n", pvc.Spec.Resources.Requests.Storage()))
	sb.WriteString(fmt.Sprintf("Capacity: %s\n", pvc.Status.Capacity.Storage()))

	if pvc.Spec.VolumeName != "" {
		sb.WriteString(fmt.Sprintf("Volume: %s\n", pvc.Spec.VolumeName))
	}

	return sb.String()
}
