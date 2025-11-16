package nodes

import (
	"context"
	"fmt"

	"github.com/fmendonca/openshift-mcp/internal/clients"
	"github.com/mark3labs/mcp-go/mcp"
	"github.com/mark3labs/mcp-go/server"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

func newListNodesHandler(clients *clients.Clients) server.ToolHandlerFunc {
	return func(ctx context.Context, args map[string]interface{}) (*mcp.ToolResponse, error) {
		nodes, err := clients.Kubernetes.CoreV1().Nodes().List(ctx, metav1.ListOptions{})
		if err != nil {
			return mcp.NewToolResponseError(fmt.Sprintf("Failed to list nodes: %v", err)), nil
		}

		return mcp.NewToolResponseText(formatNodesList(nodes)), nil
	}
}

func newGetNodeHandler(clients *clients.Clients) server.ToolHandlerFunc {
	return func(ctx context.Context, args map[string]interface{}) (*mcp.ToolResponse, error) {
		name, ok := args["name"].(string)
		if !ok {
			return mcp.NewToolResponseError("name is required"), nil
		}

		node, err := clients.Kubernetes.CoreV1().Nodes().Get(ctx, name, metav1.GetOptions{})
		if err != nil {
			return mcp.NewToolResponseError(fmt.Sprintf("Failed to get node: %v", err)), nil
		}

		return mcp.NewToolResponseText(formatNodeDetails(node)), nil
	}
}
