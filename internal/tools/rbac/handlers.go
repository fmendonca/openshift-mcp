package rbac

import (
	"context"
	"fmt"

	"github.com/mark3labs/mcp-go/mcp"
	"github.com/mark3labs/mcp-go/server"
	"github.com/yourusername/openshift-k8s-mcp/internal/clients"
	"github.com/yourusername/openshift-k8s-mcp/internal/utils"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

func newListRolesHandler(clients *clients.Clients) server.ToolHandlerFunc {
	return func(ctx context.Context, args map[string]interface{}) (*mcp.ToolResponse, error) {
		namespace := utils.GetStringArg(args, "namespace", "")

		roles, err := clients.Kubernetes.RbacV1().Roles(namespace).List(ctx, metav1.ListOptions{})
		if err != nil {
			return mcp.NewToolResponseError(fmt.Sprintf("Failed to list Roles: %v", err)), nil
		}

		return mcp.NewToolResponseText(formatRolesList(roles)), nil
	}
}

func newListRoleBindingsHandler(clients *clients.Clients) server.ToolHandlerFunc {
	return func(ctx context.Context, args map[string]interface{}) (*mcp.ToolResponse, error) {
		namespace := utils.GetStringArg(args, "namespace", "")

		roleBindings, err := clients.Kubernetes.RbacV1().RoleBindings(namespace).List(ctx, metav1.ListOptions{})
		if err != nil {
			return mcp.NewToolResponseError(fmt.Sprintf("Failed to list RoleBindings: %v", err)), nil
		}

		return mcp.NewToolResponseText(formatRoleBindingsList(roleBindings)), nil
	}
}

func newListClusterRolesHandler(clients *clients.Clients) server.ToolHandlerFunc {
	return func(ctx context.Context, args map[string]interface{}) (*mcp.ToolResponse, error) {
		clusterRoles, err := clients.Kubernetes.RbacV1().ClusterRoles().List(ctx, metav1.ListOptions{})
		if err != nil {
			return mcp.NewToolResponseError(fmt.Sprintf("Failed to list ClusterRoles: %v", err)), nil
		}

		return mcp.NewToolResponseText(formatClusterRolesList(clusterRoles)), nil
	}
}

func newListClusterRoleBindingsHandler(clients *clients.Clients) server.ToolHandlerFunc {
	return func(ctx context.Context, args map[string]interface{}) (*mcp.ToolResponse, error) {
		clusterRoleBindings, err := clients.Kubernetes.RbacV1().ClusterRoleBindings().List(ctx, metav1.ListOptions{})
		if err != nil {
			return mcp.NewToolResponseError(fmt.Sprintf("Failed to list ClusterRoleBindings: %v", err)), nil
		}

		return mcp.NewToolResponseText(formatClusterRoleBindingsList(clusterRoleBindings)), nil
	}
}
