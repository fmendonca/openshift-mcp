package services

import (
	"context"
	"fmt"

	"github.com/mark3labs/mcp-go/mcp"
	"github.com/mark3labs/mcp-go/server"
	"github.com/yourusername/openshift-k8s-mcp/internal/clients"
	"github.com/yourusername/openshift-k8s-mcp/internal/utils"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

func newListServicesHandler(clients *clients.Clients) server.ToolHandlerFunc {
	return func(ctx context.Context, args map[string]interface{}) (*mcp.ToolResponse, error) {
		namespace := utils.GetStringArg(args, "namespace", "")

		services, err := clients.Kubernetes.CoreV1().Services(namespace).List(ctx, metav1.ListOptions{})
		if err != nil {
			return mcp.NewToolResponseError(fmt.Sprintf("Failed to list services: %v", err)), nil
		}

		return mcp.NewToolResponseText(formatServicesList(services)), nil
	}
}

func newGetServiceHandler(clients *clients.Clients) server.ToolHandlerFunc {
	return func(ctx context.Context, args map[string]interface{}) (*mcp.ToolResponse, error) {
		name, ok := args["name"].(string)
		if !ok {
			return mcp.NewToolResponseError("name is required"), nil
		}

		namespace, ok := args["namespace"].(string)
		if !ok {
			return mcp.NewToolResponseError("namespace is required"), nil
		}

		service, err := clients.Kubernetes.CoreV1().Services(namespace).Get(ctx, name, metav1.GetOptions{})
		if err != nil {
			return mcp.NewToolResponseError(fmt.Sprintf("Failed to get service: %v", err)), nil
		}

		return mcp.NewToolResponseText(formatServiceDetails(service)), nil
	}
}
