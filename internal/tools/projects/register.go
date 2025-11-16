package projects

import (
	"github.com/fmendonca/openshift-mcp/internal/clients"
	"github.com/fmendonca/openshift-mcp/internal/server"
	"github.com/mark3labs/mcp-go/mcp"
)

func RegisterTools(srv *server.MCPServer, clients *clients.Clients) {
	srv.AddTool(&mcp.Tool{
		Name:        "list_projects",
		Description: "List all OpenShift projects (namespaces)",
		InputSchema: mcp.ToolInputSchema{
			Type:       "object",
			Properties: map[string]interface{}{},
		},
	}, newListProjectsHandler(clients))

	srv.AddTool(&mcp.Tool{
		Name:        "get_project",
		Description: "Get detailed information about a specific OpenShift project",
		InputSchema: mcp.ToolInputSchema{
			Type: "object",
			Properties: map[string]interface{}{
				"name": map[string]interface{}{
					"type":        "string",
					"description": "Name of the project",
				},
			},
			Required: []string{"name"},
		},
	}, newGetProjectHandler(clients))
}
