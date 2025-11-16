package imagestreams

import (
	"github.com/mark3labs/mcp-go/mcp"
	"github.com/yourusername/openshift-k8s-mcp/internal/clients"
	"github.com/yourusername/openshift-k8s-mcp/internal/server"
)

func RegisterTools(srv *server.MCPServer, clients *clients.Clients) {
	srv.AddTool(&mcp.Tool{
		Name:        "list_imagestreams",
		Description: "List all ImageStreams in a namespace",
		InputSchema: mcp.ToolInputSchema{
			Type: "object",
			Properties: map[string]interface{}{
				"namespace": map[string]interface{}{
					"type":        "string",
					"description": "Namespace to list ImageStreams from",
				},
			},
		},
	}, newListImageStreamsHandler(clients))

	srv.AddTool(&mcp.Tool{
		Name:        "get_imagestream",
		Description: "Get detailed information about a specific ImageStream",
		InputSchema: mcp.ToolInputSchema{
			Type: "object",
			Properties: map[string]interface{}{
				"name": map[string]interface{}{
					"type":        "string",
					"description": "Name of the ImageStream",
				},
				"namespace": map[string]interface{}{
					"type":        "string",
					"description": "Namespace of the ImageStream",
				},
			},
			Required: []string{"name", "namespace"},
		},
	}, newGetImageStreamHandler(clients))
}
