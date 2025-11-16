package projects

import (
	"fmt"
	"strings"

	projectv1 "github.com/openshift/api/project/v1"
)

func formatProjectsList(projects *projectv1.ProjectList) string {
	var sb strings.Builder

	sb.WriteString(fmt.Sprintf("Total Projects: %d\n\n", len(projects.Items)))

	for _, proj := range projects.Items {
		sb.WriteString(fmt.Sprintf("Name: %s\n", proj.Name))
		sb.WriteString(fmt.Sprintf("Display Name: %s\n", proj.Annotations["openshift.io/display-name"]))
		sb.WriteString(fmt.Sprintf("Status: %s\n", proj.Status.Phase))
		sb.WriteString("\n---\n\n")
	}

	return sb.String()
}

func formatProjectDetails(proj *projectv1.Project) string {
	var sb strings.Builder

	sb.WriteString(fmt.Sprintf("Project: %s\n", proj.Name))
	sb.WriteString(fmt.Sprintf("Display Name: %s\n", proj.Annotations["openshift.io/display-name"]))
	sb.WriteString(fmt.Sprintf("Description: %s\n", proj.Annotations["openshift.io/description"]))
	sb.WriteString(fmt.Sprintf("Status: %s\n", proj.Status.Phase))

	if len(proj.Annotations) > 0 {
		sb.WriteString("\nAnnotations:\n")
		for k, v := range proj.Annotations {
			sb.WriteString(fmt.Sprintf("  %s: %s\n", k, v))
		}
	}

	return sb.String()
}
