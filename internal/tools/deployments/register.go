package deployments

import (
	"github.com/fmendonca/openshift-mcp/internal/clients"
	"github.com/fmendonca/openshift-mcp/internal/server"
	"github.com/mark3labs/mcp-go/mcp"
)

func RegisterTools(srv *server.MCPServer, clients *clients.Clients) {
	srv.AddTool(&mcp.Tool{
		Name:        "list_deployments",
		Description: "List all deployments in a namespace or across all namespaces",
		InputSchema: mcp.ToolInputSchema{
			Type: "object",
			Properties: map[string]interface{}{
				"namespace": map[string]interface{}{
					"type":        "string",
					"description": "Namespace to list deployments from",
				},
			},
		},
	}, newListDeploymentsHandler(clients))

	srv.AddTool(&mcp.Tool{
		Name:        "get_deployment",
		Description: "Get detailed information about a specific deployment",
		InputSchema: mcp.ToolInputSchema{
			Type: "object",
			Properties: map[string]interface{}{
				"name": map[string]interface{}{
					"type":        "string",
					"description": "Name of the deployment",
				},
				"namespace": map[string]interface{}{
					"type":        "string",
					"description": "Namespace of the deployment",
				},
			},
			Required: []string{"name", "namespace"},
		},
	}, newGetDeploymentHandler(clients))

	srv.AddTool(&mcp.Tool{
		Name:        "scale_deployment",
		Description: "Scale a deployment to specified number of replicas",
		InputSchema: mcp.ToolInputSchema{
			Type: "object",
			Properties: map[string]interface{}{
				"name": map[string]interface{}{
					"type":        "string",
					"description": "Name of the deployment",
				},
				"namespace": map[string]interface{}{
					"type":        "string",
					"description": "Namespace of the deployment",
				},
				"replicas": map[string]interface{}{
					"type":        "integer",
					"description": "Target number of replicas",
					"minimum":     0,
				},
			},
			Required: []string{"name", "namespace", "replicas"},
		},
	}, newScaleDeploymentHandler(clients))

	srv.AddTool(&mcp.Tool{
		Name:        "restart_deployment",
		Description: "Restart a deployment by updating its annotation",
		InputSchema: mcp.ToolInputSchema{
			Type: "object",
			Properties: map[string]interface{}{
				"name": map[string]interface{}{
					"type":        "string",
					"description": "Name of the deployment",
				},
				"namespace": map[string]interface{}{
					"type":        "string",
					"description": "Namespace of the deployment",
				},
			},
			Required: []string{"name", "namespace"},
		},
	}, newRestartDeploymentHandler(clients))
}
