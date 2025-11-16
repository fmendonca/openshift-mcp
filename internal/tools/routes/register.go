package routes

import (
	"github.com/fmendonca/openshift-mcp/internal/clients"
	"github.com/fmendonca/openshift-mcp/internal/server"
	"github.com/mark3labs/mcp-go/mcp"
)

func RegisterTools(srv *server.MCPServer, clients *clients.Clients) {
	srv.AddTool(&mcp.Tool{
		Name:        "list_routes",
		Description: "List all OpenShift routes in a namespace",
		InputSchema: mcp.ToolInputSchema{
			Type: "object",
			Properties: map[string]interface{}{
				"namespace": map[string]interface{}{
					"type":        "string",
					"description": "Namespace to list routes from",
				},
			},
		},
	}, newListRoutesHandler(clients))

	srv.AddTool(&mcp.Tool{
		Name:        "get_route",
		Description: "Get detailed information about a specific OpenShift route",
		InputSchema: mcp.ToolInputSchema{
			Type: "object",
			Properties: map[string]interface{}{
				"name": map[string]interface{}{
					"type":        "string",
					"description": "Name of the route",
				},
				"namespace": map[string]interface{}{
					"type":        "string",
					"description": "Namespace of the route",
				},
			},
			Required: []string{"name", "namespace"},
		},
	}, newGetRouteHandler(clients))
}
