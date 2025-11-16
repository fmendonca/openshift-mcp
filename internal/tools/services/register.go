package services
package services

import (
	"github.com/mark3labs/mcp-go/mcp"
	"github.com/fmendonca/openshift-mcp/internal/clients"
	"github.com/fmendonca/openshift-mcp/internal/server"
)

func RegisterTools(srv *server.MCPServer, clients *clients.Clients) {
	srv.AddTool(&mcp.Tool{
		Name:        "list_services",
		Description: "List all services in a namespace or across all namespaces",
		InputSchema: mcp.ToolInputSchema{
			Type: "object",
			Properties: map[string]interface{}{
				"namespace": map[string]interface{}{
					"type":        "string",
					"description": "Namespace to list services from",
				},
			},
		},
	}, newListServicesHandler(clients))

	srv.AddTool(&mcp.Tool{
		Name:        "get_service",
		Description: "Get detailed information about a specific service",
		InputSchema: mcp.ToolInputSchema{
			Type: "object",
			Properties: map[string]interface{}{
				"name": map[string]interface{}{
					"type":        "string",
					"description": "Name of the service",
				},
				"namespace": map[string]interface{}{
					"type":        "string",
					"description": "Namespace of the service",
				},
			},
			Required: []string{"name", "namespace"},
		},
	}, newGetServiceHandler(clients))
}
