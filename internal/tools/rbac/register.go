package rbac

import (
	"github.com/mark3labs/mcp-go/mcp"
	"github.com/yourusername/openshift-k8s-mcp/internal/clients"
	"github.com/yourusername/openshift-k8s-mcp/internal/server"
)

func RegisterTools(srv *server.MCPServer, clients *clients.Clients) {
	srv.AddTool(&mcp.Tool{
		Name:        "list_roles",
		Description: "List all Roles in a namespace",
		InputSchema: mcp.ToolInputSchema{
			Type: "object",
			Properties: map[string]interface{}{
				"namespace": map[string]interface{}{
					"type":        "string",
					"description": "Namespace to list Roles from",
				},
			},
		},
	}, newListRolesHandler(clients))

	srv.AddTool(&mcp.Tool{
		Name:        "list_rolebindings",
		Description: "List all RoleBindings in a namespace",
		InputSchema: mcp.ToolInputSchema{
			Type: "object",
			Properties: map[string]interface{}{
				"namespace": map[string]interface{}{
					"type":        "string",
					"description": "Namespace to list RoleBindings from",
				},
			},
		},
	}, newListRoleBindingsHandler(clients))

	srv.AddTool(&mcp.Tool{
		Name:        "list_clusterroles",
		Description: "List all ClusterRoles",
		InputSchema: mcp.ToolInputSchema{
			Type:       "object",
			Properties: map[string]interface{}{},
		},
	}, newListClusterRolesHandler(clients))

	srv.AddTool(&mcp.Tool{
		Name:        "list_clusterrolebindings",
		Description: "List all ClusterRoleBindings",
		InputSchema: mcp.ToolInputSchema{
			Type:       "object",
			Properties: map[string]interface{}{},
		},
	}, newListClusterRoleBindingsHandler(clients))
}
