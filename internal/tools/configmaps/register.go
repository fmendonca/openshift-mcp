package configmaps

import (
	"github.com/mark3labs/mcp-go/mcp"
	"github.com/yourusername/openshift-k8s-mcp/internal/clients"
	"github.com/yourusername/openshift-k8s-mcp/internal/server"
)

func RegisterTools(srv *server.MCPServer, clients *clients.Clients) {
	srv.AddTool(&mcp.Tool{
		Name:        "list_configmaps",
		Description: "List all ConfigMaps in a namespace",
		InputSchema: mcp.ToolInputSchema{
			Type: "object",
			Properties: map[string]interface{}{
				"namespace": map[string]interface{}{
					"type":        "string",
					"description": "Namespace to list ConfigMaps from",
				},
			},
		},
	}, newListConfigMapsHandler(clients))

	srv.AddTool(&mcp.Tool{
		Name:        "get_configmap",
		Description: "Get detailed information about a specific ConfigMap",
		InputSchema: mcp.ToolInputSchema{
			Type: "object",
			Properties: map[string]interface{}{
				"name": map[string]interface{}{
					"type":        "string",
					"description": "Name of the ConfigMap",
				},
				"namespace": map[string]interface{}{
					"type":        "string",
					"description": "Namespace of the ConfigMap",
				},
			},
			Required: []string{"name", "namespace"},
		},
	}, newGetConfigMapHandler(clients))
}
