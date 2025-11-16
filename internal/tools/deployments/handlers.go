package deployments

import (
	"context"
	"fmt"
	"time"

	"github.com/fmendonca/openshift-mcp/internal/clients"
	"github.com/fmendonca/openshift-mcp/internal/utils"
	"github.com/mark3labs/mcp-go/mcp"
	"github.com/mark3labs/mcp-go/server"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

func newListDeploymentsHandler(clients *clients.Clients) server.ToolHandlerFunc {
	return func(ctx context.Context, args map[string]interface{}) (*mcp.ToolResponse, error) {
		namespace := utils.GetStringArg(args, "namespace", "")

		deployments, err := clients.Kubernetes.AppsV1().Deployments(namespace).List(ctx, metav1.ListOptions{})
		if err != nil {
			return mcp.NewToolResponseError(fmt.Sprintf("Failed to list deployments: %v", err)), nil
		}

		return mcp.NewToolResponseText(formatDeploymentsList(deployments)), nil
	}
}

func newGetDeploymentHandler(clients *clients.Clients) server.ToolHandlerFunc {
	return func(ctx context.Context, args map[string]interface{}) (*mcp.ToolResponse, error) {
		name, ok := args["name"].(string)
		if !ok {
			return mcp.NewToolResponseError("name is required"), nil
		}

		namespace, ok := args["namespace"].(string)
		if !ok {
			return mcp.NewToolResponseError("namespace is required"), nil
		}

		deployment, err := clients.Kubernetes.AppsV1().Deployments(namespace).Get(ctx, name, metav1.GetOptions{})
		if err != nil {
			return mcp.NewToolResponseError(fmt.Sprintf("Failed to get deployment: %v", err)), nil
		}

		return mcp.NewToolResponseText(formatDeploymentDetails(deployment)), nil
	}
}

func newScaleDeploymentHandler(clients *clients.Clients) server.ToolHandlerFunc {
	return func(ctx context.Context, args map[string]interface{}) (*mcp.ToolResponse, error) {
		name, ok := args["name"].(string)
		if !ok {
			return mcp.NewToolResponseError("name is required"), nil
		}

		namespace, ok := args["namespace"].(string)
		if !ok {
			return mcp.NewToolResponseError("namespace is required"), nil
		}

		replicas := int32(utils.GetIntArg(args, "replicas", 1))

		deployment, err := clients.Kubernetes.AppsV1().Deployments(namespace).Get(ctx, name, metav1.GetOptions{})
		if err != nil {
			return mcp.NewToolResponseError(fmt.Sprintf("Failed to get deployment: %v", err)), nil
		}

		deployment.Spec.Replicas = &replicas

		_, err = clients.Kubernetes.AppsV1().Deployments(namespace).Update(ctx, deployment, metav1.UpdateOptions{})
		if err != nil {
			return mcp.NewToolResponseError(fmt.Sprintf("Failed to scale deployment: %v", err)), nil
		}

		return mcp.NewToolResponseText(fmt.Sprintf("Deployment %s scaled to %d replicas", name, replicas)), nil
	}
}

func newRestartDeploymentHandler(clients *clients.Clients) server.ToolHandlerFunc {
	return func(ctx context.Context, args map[string]interface{}) (*mcp.ToolResponse, error) {
		name, ok := args["name"].(string)
		if !ok {
			return mcp.NewToolResponseError("name is required"), nil
		}

		namespace, ok := args["namespace"].(string)
		if !ok {
			return mcp.NewToolResponseError("namespace is required"), nil
		}

		deployment, err := clients.Kubernetes.AppsV1().Deployments(namespace).Get(ctx, name, metav1.GetOptions{})
		if err != nil {
			return mcp.NewToolResponseError(fmt.Sprintf("Failed to get deployment: %v", err)), nil
		}

		if deployment.Spec.Template.ObjectMeta.Annotations == nil {
			deployment.Spec.Template.ObjectMeta.Annotations = make(map[string]string)
		}
		deployment.Spec.Template.ObjectMeta.Annotations["kubectl.kubernetes.io/restartedAt"] = time.Now().Format(time.RFC3339)

		_, err = clients.Kubernetes.AppsV1().Deployments(namespace).Update(ctx, deployment, metav1.UpdateOptions{})
		if err != nil {
			return mcp.NewToolResponseError(fmt.Sprintf("Failed to restart deployment: %v", err)), nil
		}

		return mcp.NewToolResponseText(fmt.Sprintf("Deployment %s restarted successfully", name)), nil
	}
}
