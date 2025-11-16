package ingress

import (
	"github.com/fmendonca/openshift-mcp/internal/clients"
	"github.com/fmendonca/openshift-mcp/internal/server"
	"github.com/mark3labs/mcp-go/mcp"
)

func RegisterTools(srv *server.MCPServer, clients *clients.Clients) {
	srv.AddTool(&mcp.Tool{
		Name:        "list_ingresses",
		Description: "List all Ingresses in a namespace",
		InputSchema: mcp.ToolInputSchema{
			Type: "object",
			Properties: map[string]interface{}{
				"namespace": map[string]interface{}{
					"type":        "string",
					"description": "Namespace to list Ingresses from",
				},
			},
		},
	}, newListIngressesHandler(clients))

	srv.AddTool(&mcp.Tool{
		Name:        "get_ingress",
		Description: "Get detailed information about a specific Ingress",
		InputSchema: mcp.ToolInputSchema{
			Type: "object",
			Properties: map[string]interface{}{
				"name": map[string]interface{}{
					"type":        "string",
					"description": "Name of the Ingress",
				},
				"namespace": map[string]interface{}{
					"type":        "string",
					"description": "Namespace of the Ingress",
				},
			},
			Required: []string{"name", "namespace"},
		},
	}, newGetIngressHandler(clients))
}
