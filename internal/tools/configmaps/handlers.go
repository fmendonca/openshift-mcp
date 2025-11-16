package configmaps

import (
	"context"
	"fmt"

	"github.com/fmendonca/openshift-mcp/internal/clients"
	"github.com/fmendonca/openshift-mcp/internal/utils"
	"github.com/mark3labs/mcp-go/mcp"
	"github.com/mark3labs/mcp-go/server"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

func newListConfigMapsHandler(clients *clients.Clients) server.ToolHandlerFunc {
	return func(ctx context.Context, args map[string]interface{}) (*mcp.ToolResponse, error) {
		namespace := utils.GetStringArg(args, "namespace", "")

		configMaps, err := clients.Kubernetes.CoreV1().ConfigMaps(namespace).List(ctx, metav1.ListOptions{})
		if err != nil {
			return mcp.NewToolResponseError(fmt.Sprintf("Failed to list ConfigMaps: %v", err)), nil
		}

		return mcp.NewToolResponseText(formatConfigMapsList(configMaps)), nil
	}
}

func newGetConfigMapHandler(clients *clients.Clients) server.ToolHandlerFunc {
	return func(ctx context.Context, args map[string]interface{}) (*mcp.ToolResponse, error) {
		name, ok := args["name"].(string)
		if !ok {
			return mcp.NewToolResponseError("name is required"), nil
		}

		namespace, ok := args["namespace"].(string)
		if !ok {
			return mcp.NewToolResponseError("namespace is required"), nil
		}

		configMap, err := clients.Kubernetes.CoreV1().ConfigMaps(namespace).Get(ctx, name, metav1.GetOptions{})
		if err != nil {
			return mcp.NewToolResponseError(fmt.Sprintf("Failed to get ConfigMap: %v", err)), nil
		}

		return mcp.NewToolResponseText(formatConfigMapDetails(configMap)), nil
	}
}
