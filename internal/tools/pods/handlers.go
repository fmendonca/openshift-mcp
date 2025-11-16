package pods

import (
	"bytes"
	"context"
	"fmt"
	"io"

	"github.com/fmendonca/openshift-mcp/internal/clients"
	"github.com/fmendonca/openshift-mcp/internal/utils"
	"github.com/mark3labs/mcp-go/mcp"
	"github.com/mark3labs/mcp-go/server"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes/scheme"
	"k8s.io/client-go/tools/remotecommand"
)

func newListPodsHandler(clients *clients.Clients) server.ToolHandlerFunc {
	return func(ctx context.Context, args map[string]interface{}) (*mcp.ToolResponse, error) {
		namespace := utils.GetStringArg(args, "namespace", "")
		labelSelector := utils.GetStringArg(args, "labelSelector", "")

		pods, err := clients.Kubernetes.CoreV1().Pods(namespace).List(ctx, metav1.ListOptions{
			LabelSelector: labelSelector,
		})
		if err != nil {
			return mcp.NewToolResponseError(fmt.Sprintf("Failed to list pods: %v", err)), nil
		}

		return mcp.NewToolResponseText(formatPodsList(pods)), nil
	}
}

func newGetPodHandler(clients *clients.Clients) server.ToolHandlerFunc {
	return func(ctx context.Context, args map[string]interface{}) (*mcp.ToolResponse, error) {
		name, ok := args["name"].(string)
		if !ok {
			return mcp.NewToolResponseError("name is required"), nil
		}

		namespace, ok := args["namespace"].(string)
		if !ok {
			return mcp.NewToolResponseError("namespace is required"), nil
		}

		pod, err := clients.Kubernetes.CoreV1().Pods(namespace).Get(ctx, name, metav1.GetOptions{})
		if err != nil {
			return mcp.NewToolResponseError(fmt.Sprintf("Failed to get pod: %v", err)), nil
		}

		return mcp.NewToolResponseText(formatPodDetails(pod)), nil
	}
}

func newGetPodLogsHandler(clients *clients.Clients) server.ToolHandlerFunc {
	return func(ctx context.Context, args map[string]interface{}) (*mcp.ToolResponse, error) {
		name, ok := args["name"].(string)
		if !ok {
			return mcp.NewToolResponseError("name is required"), nil
		}

		namespace, ok := args["namespace"].(string)
		if !ok {
			return mcp.NewToolResponseError("namespace is required"), nil
		}

		container := utils.GetStringArg(args, "container", "")
		tailLines := int64(utils.GetIntArg(args, "tailLines", 100))
		previous := utils.GetBoolArg(args, "previous", false)

		logOptions := &corev1.PodLogOptions{
			Container: container,
			TailLines: &tailLines,
			Previous:  previous,
		}

		req := clients.Kubernetes.CoreV1().Pods(namespace).GetLogs(name, logOptions)
		podLogs, err := req.Stream(ctx)
		if err != nil {
			return mcp.NewToolResponseError(fmt.Sprintf("Failed to get pod logs: %v", err)), nil
		}
		defer podLogs.Close()

		buf := new(bytes.Buffer)
		_, err = io.Copy(buf, podLogs)
		if err != nil {
			return mcp.NewToolResponseError(fmt.Sprintf("Failed to read logs: %v", err)), nil
		}

		return mcp.NewToolResponseText(buf.String()), nil
	}
}

func newDeletePodHandler(clients *clients.Clients) server.ToolHandlerFunc {
	return func(ctx context.Context, args map[string]interface{}) (*mcp.ToolResponse, error) {
		name, ok := args["name"].(string)
		if !ok {
			return mcp.NewToolResponseError("name is required"), nil
		}

		namespace, ok := args["namespace"].(string)
		if !ok {
			return mcp.NewToolResponseError("namespace is required"), nil
		}

		err := clients.Kubernetes.CoreV1().Pods(namespace).Delete(ctx, name, metav1.DeleteOptions{})
		if err != nil {
			return mcp.NewToolResponseError(fmt.Sprintf("Failed to delete pod: %v", err)), nil
		}

		return mcp.NewToolResponseText(fmt.Sprintf("Pod %s in namespace %s deleted successfully", name, namespace)), nil
	}
}

func newExecPodHandler(clients *clients.Clients) server.ToolHandlerFunc {
	return func(ctx context.Context, args map[string]interface{}) (*mcp.ToolResponse, error) {
		name, ok := args["name"].(string)
		if !ok {
			return mcp.NewToolResponseError("name is required"), nil
		}

		namespace, ok := args["namespace"].(string)
		if !ok {
			return mcp.NewToolResponseError("namespace is required"), nil
		}

		container := utils.GetStringArg(args, "container", "")

		commandInterface, ok := args["command"]
		if !ok {
			return mcp.NewToolResponseError("command is required"), nil
		}

		command := utils.InterfaceSliceToStringSlice(commandInterface)

		req := clients.Kubernetes.CoreV1().RESTClient().Post().
			Resource("pods").
			Name(name).
			Namespace(namespace).
			SubResource("exec")

		execOptions := &corev1.PodExecOptions{
			Container: container,
			Command:   command,
			Stdout:    true,
			Stderr:    true,
		}

		req.VersionedParams(execOptions, scheme.ParameterCodec)

		exec, err := remotecommand.NewSPDYExecutor(clients.RestConfig, "POST", req.URL())
		if err != nil {
			return mcp.NewToolResponseError(fmt.Sprintf("Failed to create executor: %v", err)), nil
		}

		var stdout, stderr bytes.Buffer
		err = exec.StreamWithContext(ctx, remotecommand.StreamOptions{
			Stdout: &stdout,
			Stderr: &stderr,
		})

		if err != nil {
			return mcp.NewToolResponseError(fmt.Sprintf("Exec failed: %v\nStderr: %s", err, stderr.String())), nil
		}

		result := fmt.Sprintf("Stdout:\n%s\n\nStderr:\n%s", stdout.String(), stderr.String())
		return mcp.NewToolResponseText(result), nil
	}
}
