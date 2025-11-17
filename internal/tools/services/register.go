// internal/tools/services/register.go
package services

import (
	"github.com/fmendonca/openshift-mcp/internal/clients"
	"github.com/mark3labs/mcp-go/mcp"
	"github.com/mark3labs/mcp-go/server"
)

func RegisterTools(srv *server.MCPServer, c *clients.Clients) {
	// list_services
	listTool := mcp.NewTool(
		"list_services",
		"List all services in a namespace or across all namespaces",
		mcp.WithInputSchema(map[string]any{
			"type": "object",
			"properties": map[string]any{
				"namespace": map[string]any{
					"type":        "string",
					"description": "Namespace to list services from (empty for all namespaces)",
				},
			},
		}),
	)
	srv.AddTool(listTool, newListServicesHandler(c))

	// get_service
	getTool := mcp.NewTool(
		"get_service",
		"Get detailed information about a specific Service",
		mcp.WithInputSchema(map[string]any{
			"type": "object",
			"properties": map[string]any{
				"name": map[string]any{
					"type":        "string",
					"description": "Name of the Service",
				},
				"namespace": map[string]any{
					"type":        "string",
					"description": "Namespace of the Service",
				},
			},
			"required": []string{"name", "namespace"},
		}),
	)
	srv.AddTool(getTool, newGetServiceHandler(c))
}
