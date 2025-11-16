package projects

import (
	"context"
	"fmt"

	"github.com/fmendonca/openshift-mcp/internal/clients"
	"github.com/mark3labs/mcp-go/mcp"
	"github.com/mark3labs/mcp-go/server"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

func newListProjectsHandler(clients *clients.Clients) server.ToolHandlerFunc {
	return func(ctx context.Context, args map[string]interface{}) (*mcp.ToolResponse, error) {
		projects, err := clients.Project.ProjectV1().Projects().List(ctx, metav1.ListOptions{})
		if err != nil {
			return mcp.NewToolResponseError(fmt.Sprintf("Failed to list projects: %v", err)), nil
		}

		return mcp.NewToolResponseText(formatProjectsList(projects)), nil
	}
}

func newGetProjectHandler(clients *clients.Clients) server.ToolHandlerFunc {
	return func(ctx context.Context, args map[string]interface{}) (*mcp.ToolResponse, error) {
		name, ok := args["name"].(string)
		if !ok {
			return mcp.NewToolResponseError("name is required"), nil
		}

		project, err := clients.Project.ProjectV1().Projects().Get(ctx, name, metav1.GetOptions{})
		if err != nil {
			return mcp.NewToolResponseError(fmt.Sprintf("Failed to get project: %v", err)), nil
		}

		return mcp.NewToolResponseText(formatProjectDetails(project)), nil
	}
}
