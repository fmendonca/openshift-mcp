package secrets

import (
	"github.com/fmendonca/openshift-mcp/internal/clients"
	"github.com/fmendonca/openshift-mcp/internal/server"
	"github.com/mark3labs/mcp-go/mcp"
)

func RegisterTools(srv *server.MCPServer, clients *clients.Clients) {
	srv.AddTool(&mcp.Tool{
		Name:        "list_secrets",
		Description: "List all Secrets in a namespace",
		InputSchema: mcp.ToolInputSchema{
			Type: "object",
			Properties: map[string]interface{}{
				"namespace": map[string]interface{}{
					"type":        "string",
					"description": "Namespace to list Secrets from",
				},
			},
		},
	}, newListSecretsHandler(clients))

	srv.AddTool(&mcp.Tool{
		Name:        "get_secret",
		Description: "Get information about a specific Secret (data will be masked)",
		InputSchema: mcp.ToolInputSchema{
			Type: "object",
			Properties: map[string]interface{}{
				"name": map[string]interface{}{
					"type":        "string",
					"description": "Name of the Secret",
				},
				"namespace": map[string]interface{}{
					"type":        "string",
					"description": "Namespace of the Secret",
				},
			},
			Required: []string{"name", "namespace"},
		},
	}, newGetSecretHandler(clients))
}
