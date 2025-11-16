package pods

import (
	"github.com/mark3labs/mcp-go/mcp"
	"github.com/yourusername/openshift-k8s-mcp/internal/clients"
	"github.com/yourusername/openshift-k8s-mcp/internal/server"
)

func RegisterTools(srv *server.MCPServer, clients *clients.Clients) {
	srv.AddTool(&mcp.Tool{
		Name:        "list_pods",
		Description: "List all pods in a namespace or across all namespaces",
		InputSchema: mcp.ToolInputSchema{
			Type: "object",
			Properties: map[string]interface{}{
				"namespace": map[string]interface{}{
					"type":        "string",
					"description": "Namespace to list pods from (empty for all namespaces)",
				},
				"labelSelector": map[string]interface{}{
					"type":        "string",
					"description": "Label selector to filter pods (e.g., 'app=nginx')",
				},
			},
		},
	}, newListPodsHandler(clients))

	srv.AddTool(&mcp.Tool{
		Name:        "get_pod",
		Description: "Get detailed information about a specific pod",
		InputSchema: mcp.ToolInputSchema{
			Type: "object",
			Properties: map[string]interface{}{
				"name": map[string]interface{}{
					"type":        "string",
					"description": "Name of the pod",
				},
				"namespace": map[string]interface{}{
					"type":        "string",
					"description": "Namespace of the pod",
				},
			},
			Required: []string{"name", "namespace"},
		},
	}, newGetPodHandler(clients))

	srv.AddTool(&mcp.Tool{
		Name:        "get_pod_logs",
		Description: "Get logs from a pod container",
		InputSchema: mcp.ToolInputSchema{
			Type: "object",
			Properties: map[string]interface{}{
				"name": map[string]interface{}{
					"type":        "string",
					"description": "Name of the pod",
				},
				"namespace": map[string]interface{}{
					"type":        "string",
					"description": "Namespace of the pod",
				},
				"container": map[string]interface{}{
					"type":        "string",
					"description": "Container name (optional, first container if not specified)",
				},
				"tailLines": map[string]interface{}{
					"type":        "integer",
					"description": "Number of lines to retrieve from the end",
					"default":     100,
				},
				"previous": map[string]interface{}{
					"type":        "boolean",
					"description": "Get logs from previous container instance",
					"default":     false,
				},
			},
			Required: []string{"name", "namespace"},
		},
	}, newGetPodLogsHandler(clients))

	srv.AddTool(&mcp.Tool{
		Name:        "delete_pod",
		Description: "Delete a specific pod",
		InputSchema: mcp.ToolInputSchema{
			Type: "object",
			Properties: map[string]interface{}{
				"name": map[string]interface{}{
					"type":        "string",
					"description": "Name of the pod to delete",
				},
				"namespace": map[string]interface{}{
					"type":        "string",
					"description": "Namespace of the pod",
				},
			},
			Required: []string{"name", "namespace"},
		},
	}, newDeletePodHandler(clients))

	srv.AddTool(&mcp.Tool{
		Name:        "exec_pod",
		Description: "Execute a command in a pod container",
		InputSchema: mcp.ToolInputSchema{
			Type: "object",
			Properties: map[string]interface{}{
				"name": map[string]interface{}{
					"type":        "string",
					"description": "Name of the pod",
				},
				"namespace": map[string]interface{}{
					"type":        "string",
					"description": "Namespace of the pod",
				},
				"container": map[string]interface{}{
					"type":        "string",
					"description": "Container name (optional)",
				},
				"command": map[string]interface{}{
					"type":        "array",
					"description": "Command to execute as array of strings",
					"items": map[string]interface{}{
						"type": "string",
					},
				},
			},
			Required: []string{"name", "namespace", "command"},
		},
	}, newExecPodHandler(clients))
}
