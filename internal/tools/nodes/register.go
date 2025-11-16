package nodes

import (
	"github.com/mark3labs/mcp-go/mcp"
	"github.com/yourusername/openshift-k8s-mcp/internal/clients"
	"github.com/yourusername/openshift-k8s-mcp/internal/server"
)

func RegisterTools(srv *server.MCPServer, clients *clients.Clients) {
	srv.AddTool(&mcp.Tool{
		Name:        "list_nodes",
		Description: "List all nodes in the cluster",
		InputSchema: mcp.ToolInputSchema{
			Type:       "object",
			Properties: map[string]interface{}{},
		},
	}, newListNodesHandler(clients))

	srv.AddTool(&mcp.Tool{
		Name:        "get_node",
		Description: "Get detailed information about a specific node",
		InputSchema: mcp.ToolInputSchema{
			Type: "object",
			Properties: map[string]interface{}{
				"name": map[string]interface{}{
					"type":        "string",
					"description": "Name of the node",
				},
			},
			Required: []string{"name"},
		},
	}, newGetNodeHandler(clients))
}
