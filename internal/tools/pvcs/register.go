package pvcs

import (
	"github.com/mark3labs/mcp-go/mcp"
	"github.com/yourusername/openshift-k8s-mcp/internal/clients"
	"github.com/yourusername/openshift-k8s-mcp/internal/server"
)

func RegisterTools(srv *server.MCPServer, clients *clients.Clients) {
	srv.AddTool(&mcp.Tool{
		Name:        "list_pvcs",
		Description: "List all PersistentVolumeClaims in a namespace",
		InputSchema: mcp.ToolInputSchema{
			Type: "object",
			Properties: map[string]interface{}{
				"namespace": map[string]interface{}{
					"type":        "string",
					"description": "Namespace to list PVCs from",
				},
			},
		},
	}, newListPVCsHandler(clients))

	srv.AddTool(&mcp.Tool{
		Name:        "get_pvc",
		Description: "Get detailed information about a specific PersistentVolumeClaim",
		InputSchema: mcp.ToolInputSchema{
			Type: "object",
			Properties: map[string]interface{}{
				"name": map[string]interface{}{
					"type":        "string",
					"description": "Name of the PVC",
				},
				"namespace": map[string]interface{}{
					"type":        "string",
					"description": "Namespace of the PVC",
				},
			},
			Required: []string{"name", "namespace"},
		},
	}, newGetPVCHandler(clients))
}
