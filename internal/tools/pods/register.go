package pods

import (
	"github.com/fmendonca/openshift-mcp/internal/clients"
	"github.com/mark3labs/mcp-go/mcp"
	"github.com/mark3labs/mcp-go/server"
)

func RegisterTools(srv *server.MCPServer, c *clients.Clients) {
	list := mcp.NewTool(
		"list_pods",
		"List all pods in a namespace or across all namespaces",
		mcp.WithInputSchema(map[string]any{
			"type": "object",
			"properties": map[string]any{
				"namespace": map[string]any{
					"type":        "string",
					"description": "Namespace to list pods from (empty for all namespaces)",
				},
				"labelSelector": map[string]any{
					"type":        "string",
					"description": "Label selector to filter pods",
				},
			},
		}),
	)
	srv.AddTool(list, newListPodsHandler(c))

	get := mcp.NewTool(
		"get_pod",
		"Get detailed information about a specific pod",
		mcp.WithInputSchema(map[string]any{
			"type": "object",
			"properties": map[string]any{
				"name": map[string]any{
					"type":        "string",
					"description": "Name of the pod",
				},
				"namespace": map[string]any{
					"type":        "string",
					"description": "Namespace of the pod",
				},
			},
			"required": []string{"name", "namespace"},
		}),
	)
	srv.AddTool(get, newGetPodHandler(c))

	logs := mcp.NewTool(
		"get_pod_logs",
		"Get logs from a pod container",
		mcp.WithInputSchema(map[string]any{
			"type": "object",
			"properties": map[string]any{
				"name": map[string]any{
					"type":        "string",
					"description": "Name of the pod",
				},
				"namespace": map[string]any{
					"type":        "string",
					"description": "Namespace of the pod",
				},
				"container": map[string]any{
					"type":        "string",
					"description": "Container name (optional, first container if omitted)",
				},
				"tailLines": map[string]any{
					"type":        "integer",
					"description": "Number of log lines from the end",
					"default":     100,
				},
				"previous": map[string]any{
					"type":        "boolean",
					"description": "Get logs from previous container instance",
					"default":     false,
				},
			},
			"required": []string{"name", "namespace"},
		}),
	)
	srv.AddTool(logs, newGetPodLogsHandler(c))

	del := mcp.NewTool(
		"delete_pod",
		"Delete a specific pod",
		mcp.WithInputSchema(map[string]any{
			"type": "object",
			"properties": map[string]any{
				"name": map[string]any{
					"type":        "string",
					"description": "Name of the pod",
				},
				"namespace": map[string]any{
					"type":        "string",
					"description": "Namespace of the pod",
				},
			},
			"required": []string{"name", "namespace"},
		}),
	)
	srv.AddTool(del, newDeletePodHandler(c))

	execTool := mcp.NewTool(
		"exec_pod",
		"Execute a command in a pod container",
		mcp.WithInputSchema(map[string]any{
			"type": "object",
			"properties": map[string]any{
				"name": map[string]any{
					"type":        "string",
					"description": "Name of the pod",
				},
				"namespace": map[string]any{
					"type":        "string",
					"description": "Namespace of the pod",
				},
				"container": map[string]any{
					"type":        "string",
					"description": "Container name (optional)",
				},
				"command": map[string]any{
					"type":        "array",
					"description": "Command to execute (array of strings)",
					"items": map[string]any{
						"type": "string",
					},
				},
			},
			"required": []string{"name", "namespace", "command"},
		}),
	)
	srv.AddTool(execTool, newExecPodHandler(c))
}
