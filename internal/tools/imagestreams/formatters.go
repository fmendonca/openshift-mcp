package imagestreams

import (
	"fmt"
	"strings"

	imagev1 "github.com/openshift/api/image/v1"
)

func formatImageStreamsList(imageStreams *imagev1.ImageStreamList) string {
	var sb strings.Builder

	sb.WriteString(fmt.Sprintf("Total ImageStreams: %d\n\n", len(imageStreams.Items)))

	for _, is := range imageStreams.Items {
		sb.WriteString(fmt.Sprintf("Name: %s\n", is.Name))
		sb.WriteString(fmt.Sprintf("Namespace: %s\n", is.Namespace))
		sb.WriteString(fmt.Sprintf("Docker Repository: %s\n", is.Status.DockerImageRepository))
		sb.WriteString(fmt.Sprintf("Tags: %d\n", len(is.Status.Tags)))
		sb.WriteString("\n---\n\n")
	}

	return sb.String()
}

func formatImageStreamDetails(is *imagev1.ImageStream) string {
	var sb strings.Builder

	sb.WriteString(fmt.Sprintf("ImageStream: %s\n", is.Name))
	sb.WriteString(fmt.Sprintf("Namespace: %s\n", is.Namespace))
	sb.WriteString(fmt.Sprintf("Docker Repository: %s\n", is.Status.DockerImageRepository))
	sb.WriteString(fmt.Sprintf("Public Docker Repository: %s\n", is.Status.PublicDockerImageRepository))

	if len(is.Spec.Tags) > 0 {
		sb.WriteString("\nSpec Tags:\n")
		for _, tag := range is.Spec.Tags {
			sb.WriteString(fmt.Sprintf("  - %s\n", tag.Name))
			if tag.From != nil {
				sb.WriteString(fmt.Sprintf("    From: %s/%s\n", tag.From.Kind, tag.From.Name))
			}
		}
	}

	if len(is.Status.Tags) > 0 {
		sb.WriteString("\nStatus Tags:\n")
		for _, tag := range is.Status.Tags {
			sb.WriteString(fmt.Sprintf("  - %s (%d items)\n", tag.Tag, len(tag.Items)))
			if len(tag.Items) > 0 {
				sb.WriteString(fmt.Sprintf("    Latest: %s\n", tag.Items[0].DockerImageReference))
			}
		}
	}

	return sb.String()
}
