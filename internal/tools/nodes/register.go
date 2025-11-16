package nodes

import (
	"github.com/fmendonca/openshift-mcp/internal/clients"
	"github.com/fmendonca/openshift-mcp/internal/server"
	"github.com/mark3labs/mcp-go/mcp"
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
