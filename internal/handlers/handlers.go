// internal/handlers/handlers.go
package handlers

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"

	"github.com/fmendonca/openshift-mcp/internal/clients"
	"github.com/fmendonca/openshift-mcp/internal/utils"
	"github.com/mark3labs/mcp-go/mcp"
	mcpsrv "github.com/mark3labs/mcp-go/server"

	corev1 "k8s.io/api/core/v1"
	storagev1 "k8s.io/api/storage/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"
	"k8s.io/apimachinery/pkg/runtime/schema"
	"k8s.io/apimachinery/pkg/types"

	"k8s.io/client-go/kubernetes/scheme"
	"k8s.io/client-go/tools/remotecommand"
)

///////////////////////////////////////////////////////////////////////////////
// REGISTRO GERAL
///////////////////////////////////////////////////////////////////////////////

// RegisterAllTools registra todos os tools do MCP em um único lugar.
func RegisterAllTools(srv *mcpsrv.MCPServer, c *clients.Clients) {
	registerPodTools(srv, c)
	registerServiceTools(srv, c)
	registerClusterTools(srv, c)
	registerKubeVirtTools(srv, c)
}

///////////////////////////////////////////////////////////////////////////////
// KUBEVIRT VIRTUAL MACHINES via dynamic client (sem kubevirt.io/client-go)
///////////////////////////////////////////////////////////////////////////////

var vmGVR = schema.GroupVersionResource{
	Group:    "kubevirt.io",
	Version:  "v1",
	Resource: "virtualmachines",
}

func registerKubeVirtTools(srv *mcpsrv.MCPServer, c *clients.Clients) {
	listVMTool := mcp.NewTool(
		"list_virtualmachines",
		mcp.WithDescription("List KubeVirt VirtualMachines. Args: namespace (string, optional)."),
	)
	srv.AddTool(listVMTool, listVirtualMachinesHandler(c))

	startVMTool := mcp.NewTool(
		"start_virtualmachine",
		mcp.WithDescription("Start a VirtualMachine by setting spec.runStrategy=Always. Args: name (string), namespace (string)."),
	)
	srv.AddTool(startVMTool, startVirtualMachineHandler(c))

	stopVMTool := mcp.NewTool(
		"stop_virtualmachine",
		mcp.WithDescription("Stop a VirtualMachine by setting spec.runStrategy=Halted. Args: name (string), namespace (string)."),
	)
	srv.AddTool(stopVMTool, stopVirtualMachineHandler(c))

	restartVMTool := mcp.NewTool(
		"restart_virtualmachine",
		mcp.WithDescription("Restart a VirtualMachine by toggling spec.runStrategy. Args: name (string), namespace (string)."),
	)
	srv.AddTool(restartVMTool, restartVirtualMachineHandler(c))

	editVMResTool := mcp.NewTool(
		"edit_virtualmachine_resources",
		mcp.WithDescription("Edit CPU and memory resources of a VirtualMachine. Args: name, namespace, cpu (string, optional), memory (string, optional)."),
	)
	srv.AddTool(editVMResTool, editVirtualMachineResourcesHandler(c))
}

func listVirtualMachinesHandler(c *clients.Clients) mcpsrv.ToolHandlerFunc {
	return func(ctx context.Context, req mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		raw := req.Params.Arguments
		args, ok := raw.(map[string]any)
		if !ok {
			return mcp.NewToolResultError("invalid arguments payload (expected object)"), nil
		}

		ns := utils.GetStringArg(args, "namespace", "")

		var (
			list *unstructured.UnstructuredList
			err  error
		)

		if ns == "" {
			// todas as namespaces
			list, err = c.Dynamic.Resource(vmGVR).Namespace(metav1.NamespaceAll).List(ctx, metav1.ListOptions{})
		} else {
			list, err = c.Dynamic.Resource(vmGVR).Namespace(ns).List(ctx, metav1.ListOptions{})
		}

		if err != nil {
			return mcp.NewToolResultError(fmt.Sprintf("Failed to list VirtualMachines: %v", err)), nil
		}

		var buf bytes.Buffer
		fmt.Fprintf(&buf, "Total VirtualMachines: %d\n\n", len(list.Items))
		for _, item := range list.Items {
			name := item.GetName()
			namespace := item.GetNamespace()
			spec, _ := item.Object["spec"].(map[string]any)
			runStrategy := ""
			if spec != nil {
				if rs, ok := spec["runStrategy"].(string); ok {
					runStrategy = rs
				}
			}
			fmt.Fprintf(&buf, "Name: %s\nNamespace: %s\nRunStrategy: %s\n\n---\n\n",
				name, namespace, runStrategy)
		}

		return mcp.NewToolResultText(buf.String()), nil
	}
}

func startVirtualMachineHandler(c *clients.Clients) mcpsrv.ToolHandlerFunc {
	return func(ctx context.Context, req mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		raw := req.Params.Arguments
		args, ok := raw.(map[string]any)
		if !ok {
			return mcp.NewToolResultError("invalid arguments payload (expected object)"), nil
		}

		name, _ := args["name"].(string)
		ns, _ := args["namespace"].(string)
		if name == "" || ns == "" {
			return mcp.NewToolResultError("name and namespace are required"), nil
		}

		patch := map[string]any{
			"spec": map[string]any{
				"runStrategy": "Always",
			},
		}
		data, _ := json.Marshal(patch)

		_, err := c.Dynamic.Resource(vmGVR).Namespace(ns).Patch(
			ctx,
			name,
			types.MergePatchType,
			data,
			metav1.PatchOptions{},
		)
		if err != nil {
			return mcp.NewToolResultError(fmt.Sprintf("Failed to start VirtualMachine: %v", err)), nil
		}

		return mcp.NewToolResultText(fmt.Sprintf("VirtualMachine %s/%s started (runStrategy=Always)", ns, name)), nil
	}
}

func stopVirtualMachineHandler(c *clients.Clients) mcpsrv.ToolHandlerFunc {
	return func(ctx context.Context, req mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		raw := req.Params.Arguments
		args, ok := raw.(map[string]any)
		if !ok {
			return mcp.NewToolResultError("invalid arguments payload (expected object)"), nil
		}

		name, _ := args["name"].(string)
		ns, _ := args["namespace"].(string)
		if name == "" || ns == "" {
			return mcp.NewToolResultError("name and namespace are required"), nil
		}

		patch := map[string]any{
			"spec": map[string]any{
				"runStrategy": "Halted",
			},
		}
		data, _ := json.Marshal(patch)

		_, err := c.Dynamic.Resource(vmGVR).Namespace(ns).Patch(
			ctx,
			name,
			types.MergePatchType,
			data,
			metav1.PatchOptions{},
		)
		if err != nil {
			return mcp.NewToolResultError(fmt.Sprintf("Failed to stop VirtualMachine: %v", err)), nil
		}

		return mcp.NewToolResultText(fmt.Sprintf("VirtualMachine %s/%s stopped (runStrategy=Halted)", ns, name)), nil
	}
}

func restartVirtualMachineHandler(c *clients.Clients) mcpsrv.ToolHandlerFunc {
	return func(ctx context.Context, req mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		raw := req.Params.Arguments
		args, ok := raw.(map[string]any)
		if !ok {
			return mcp.NewToolResultError("invalid arguments payload (expected object)"), nil
		}

		name, _ := args["name"].(string)
		ns, _ := args["namespace"].(string)
		if name == "" || ns == "" {
			return mcp.NewToolResultError("name and namespace are required"), nil
		}

		patchHalted := map[string]any{
			"spec": map[string]any{
				"runStrategy": "Halted",
			},
		}
		dataHalted, _ := json.Marshal(patchHalted)
		_, err := c.Dynamic.Resource(vmGVR).Namespace(ns).Patch(
			ctx,
			name,
			types.MergePatchType,
			dataHalted,
			metav1.PatchOptions{},
		)
		if err != nil {
			return mcp.NewToolResultError(fmt.Sprintf("Failed to set runStrategy=Halted: %v", err)), nil
		}

		patchRun := map[string]any{
			"spec": map[string]any{
				"runStrategy": "Always",
			},
		}
		dataRun, _ := json.Marshal(patchRun)
		_, err = c.Dynamic.Resource(vmGVR).Namespace(ns).Patch(
			ctx,
			name,
			types.MergePatchType,
			dataRun,
			metav1.PatchOptions{},
		)
		if err != nil {
			return mcp.NewToolResultError(fmt.Sprintf("Failed to set runStrategy=Always: %v", err)), nil
		}

		return mcp.NewToolResultText(fmt.Sprintf("VirtualMachine %s/%s restarted via runStrategy toggle", ns, name)), nil
	}
}

func editVirtualMachineResourcesHandler(c *clients.Clients) mcpsrv.ToolHandlerFunc {
	return func(ctx context.Context, req mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		raw := req.Params.Arguments
		args, ok := raw.(map[string]any)
		if !ok {
			return mcp.NewToolResultError("invalid arguments payload (expected object)"), nil
		}

		name, _ := args["name"].(string)
		ns, _ := args["namespace"].(string)
		if name == "" || ns == "" {
			return mcp.NewToolResultError("name and namespace are required"), nil
		}

		cpuStr := utils.GetStringArg(args, "cpu", "")
		memStr := utils.GetStringArg(args, "memory", "")

		vm, err := c.Dynamic.Resource(vmGVR).Namespace(ns).Get(ctx, name, metav1.GetOptions{})
		if err != nil {
			return mcp.NewToolResultError(fmt.Sprintf("Failed to get VirtualMachine: %v", err)), nil
		}

		spec, ok := vm.Object["spec"].(map[string]any)
		if !ok {
			spec = map[string]any{}
			vm.Object["spec"] = spec
		}
		tpl, ok := spec["template"].(map[string]any)
		if !ok {
			tpl = map[string]any{}
			spec["template"] = tpl
		}
		tplSpec, ok := tpl["spec"].(map[string]any)
		if !ok {
			tplSpec = map[string]any{}
			tpl["spec"] = tplSpec
		}
		domain, ok := tplSpec["domain"].(map[string]any)
		if !ok {
			domain = map[string]any{}
			tplSpec["domain"] = domain
		}
		resources, ok := domain["resources"].(map[string]any)
		if !ok {
			resources = map[string]any{}
			domain["resources"] = resources
		}
		requests, ok := resources["requests"].(map[string]any)
		if !ok {
			requests = map[string]any{}
			resources["requests"] = requests
		}

		if cpuStr != "" {
			requests["cpu"] = cpuStr
		}
		if memStr != "" {
			requests["memory"] = memStr
		}

		_, err = c.Dynamic.Resource(vmGVR).Namespace(ns).Update(ctx, vm, metav1.UpdateOptions{})
		if err != nil {
			return mcp.NewToolResultError(fmt.Sprintf("Failed to update VirtualMachine resources: %v", err)), nil
		}

		return mcp.NewToolResultText(
			fmt.Sprintf("VirtualMachine %s/%s resources updated (cpu=%s, memory=%s)", ns, name, cpuStr, memStr),
		), nil
	}
}

///////////////////////////////////////////////////////////////////////////
// CLUSTER INVENTORY (namespaces, storage, ingress, RBAC)
///////////////////////////////////////////////////////////////////////////

func registerClusterTools(srv *mcpsrv.MCPServer, c *clients.Clients) {
	nsTool := mcp.NewTool(
		"list_namespaces",
		mcp.WithDescription("List all namespaces (projects) in the cluster."),
	)
	srv.AddTool(nsTool, listNamespacesHandler(c))

	scTool := mcp.NewTool(
		"list_storageclasses",
		mcp.WithDescription("List all StorageClasses in the cluster."),
	)
	srv.AddTool(scTool, listStorageClassesHandler(c))

	ingTool := mcp.NewTool(
		"list_ingresses",
		mcp.WithDescription("List ingresses. Args: namespace (string, optional)."),
	)
	srv.AddTool(ingTool, listIngressesHandler(c))

	rolesTool := mcp.NewTool(
		"list_rbac_roles",
		mcp.WithDescription("List Roles and RoleBindings in a namespace. Args: namespace (string, required)."),
	)
	srv.AddTool(rolesTool, listRBACRolesHandler(c))

	crTool := mcp.NewTool(
		"list_cluster_roles",
		mcp.WithDescription("List ClusterRoles and ClusterRoleBindings."),
	)
	srv.AddTool(crTool, listClusterRolesHandler(c))
}

func listNamespacesHandler(c *clients.Clients) mcpsrv.ToolHandlerFunc {
	return func(ctx context.Context, req mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		nsList, err := c.Kubernetes.CoreV1().Namespaces().List(ctx, metav1.ListOptions{})
		if err != nil {
			return mcp.NewToolResultError(fmt.Sprintf("Failed to list namespaces: %v", err)), nil
		}

		var buf bytes.Buffer
		fmt.Fprintf(&buf, "Total namespaces: %d\n\n", len(nsList.Items))
		for _, ns := range nsList.Items {
			fmt.Fprintf(&buf, "Name: %s\nStatus: %s\n\n", ns.Name, ns.Status.Phase)
		}

		return mcp.NewToolResultText(buf.String()), nil
	}
}

func listStorageClassesHandler(c *clients.Clients) mcpsrv.ToolHandlerFunc {
	return func(ctx context.Context, req mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		scList, err := c.Kubernetes.StorageV1().StorageClasses().List(ctx, metav1.ListOptions{})
		if err != nil {
			return mcp.NewToolResultError(fmt.Sprintf("Failed to list StorageClasses: %v", err)), nil
		}

		var buf bytes.Buffer
		fmt.Fprintf(&buf, "Total StorageClasses: %d\n\n", len(scList.Items))
		for _, sc := range scList.Items {
			allow := "false"
			if sc.AllowVolumeExpansion != nil && *sc.AllowVolumeExpansion {
				allow = "true"
			}
			fmt.Fprintf(&buf,
				"Name: %s\nProvisioner: %s\nAllowVolumeExpansion: %s\nDefault: %t\n\n---\n\n",
				sc.Name,
				sc.Provisioner,
				allow,
				isDefaultStorageClass(&sc),
			)
		}
		return mcp.NewToolResultText(buf.String()), nil
	}
}

func isDefaultStorageClass(sc *storagev1.StorageClass) bool {
	for k, v := range sc.Annotations {
		if k == "storageclass.kubernetes.io/is-default-class" && v == "true" {
			return true
		}
	}
	return false
}

func listIngressesHandler(c *clients.Clients) mcpsrv.ToolHandlerFunc {
	return func(ctx context.Context, req mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		raw := req.Params.Arguments
		args, ok := raw.(map[string]any)
		if !ok {
			return mcp.NewToolResultError("invalid arguments payload (expected object)"), nil
		}

		ns := utils.GetStringArg(args, "namespace", "")

		ingList, err := c.Kubernetes.NetworkingV1().Ingresses(ns).List(ctx, metav1.ListOptions{})
		if err != nil {
			return mcp.NewToolResultError(fmt.Sprintf("Failed to list ingresses: %v", err)), nil
		}

		var buf bytes.Buffer
		fmt.Fprintf(&buf, "Total Ingresses: %d\n\n", len(ingList.Items))
		for _, ing := range ingList.Items {
			fmt.Fprintf(&buf, "Name: %s\nNamespace: %s\nClass: %s\n",
				ing.Name, ing.Namespace, ptrToString(ing.Spec.IngressClassName))
			for _, rule := range ing.Spec.Rules {
				fmt.Fprintf(&buf, "  Host: %s\n", rule.Host)
			}
			fmt.Fprintln(&buf, "\n---\n")
		}

		return mcp.NewToolResultText(buf.String()), nil
	}
}

func ptrToString(p *string) string {
	if p == nil {
		return ""
	}
	return *p
}

func listRBACRolesHandler(c *clients.Clients) mcpsrv.ToolHandlerFunc {
	return func(ctx context.Context, req mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		raw := req.Params.Arguments
		args, ok := raw.(map[string]any)
		if !ok {
			return mcp.NewToolResultError("invalid arguments payload (expected object)"), nil
		}

		ns := utils.GetStringArg(args, "namespace", "")
		if ns == "" {
			return mcp.NewToolResultError("namespace is required"), nil
		}

		roles, err := c.Kubernetes.RbacV1().Roles(ns).List(ctx, metav1.ListOptions{})
		if err != nil {
			return mcp.NewToolResultError(fmt.Sprintf("Failed to list Roles: %v", err)), nil
		}

		rbs, err := c.Kubernetes.RbacV1().RoleBindings(ns).List(ctx, metav1.ListOptions{})
		if err != nil {
			return mcp.NewToolResultError(fmt.Sprintf("Failed to list RoleBindings: %v", err)), nil
		}

		var buf bytes.Buffer
		fmt.Fprintf(&buf, "Namespace: %s\n\nRoles: %d\n\n", ns, len(roles.Items))
		for _, r := range roles.Items {
			fmt.Fprintf(&buf, "- Role: %s\n", r.Name)
		}

		fmt.Fprintf(&buf, "\nRoleBindings: %d\n\n", len(rbs.Items))
		for _, rb := range rbs.Items {
			fmt.Fprintf(&buf, "- RoleBinding: %s (role: %s)\n", rb.Name, rb.RoleRef.Name)
			for _, subj := range rb.Subjects {
				fmt.Fprintf(&buf, "    Subject: %s %s (%s)\n", subj.Kind, subj.Name, subj.Namespace)
			}
		}

		return mcp.NewToolResultText(buf.String()), nil
	}
}

func listClusterRolesHandler(c *clients.Clients) mcpsrv.ToolHandlerFunc {
	return func(ctx context.Context, req mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		crs, err := c.Kubernetes.RbacV1().ClusterRoles().List(ctx, metav1.ListOptions{})
		if err != nil {
			return mcp.NewToolResultError(fmt.Sprintf("Failed to list ClusterRoles: %v", err)), nil
		}

		crbs, err := c.Kubernetes.RbacV1().ClusterRoleBindings().List(ctx, metav1.ListOptions{})
		if err != nil {
			return mcp.NewToolResultError(fmt.Sprintf("Failed to list ClusterRoleBindings: %v", err)), nil
		}

		var buf bytes.Buffer
		fmt.Fprintf(&buf, "ClusterRoles: %d\n\n", len(crs.Items))
		for _, r := range crs.Items {
			fmt.Fprintf(&buf, "- ClusterRole: %s\n", r.Name)
		}

		fmt.Fprintf(&buf, "\nClusterRoleBindings: %d\n\n", len(crbs.Items))
		for _, rb := range crbs.Items {
			fmt.Fprintf(&buf, "- ClusterRoleBinding: %s (role: %s)\n", rb.Name, rb.RoleRef.Name)
			for _, subj := range rb.Subjects {
				fmt.Fprintf(&buf, "    Subject: %s %s (%s)\n", subj.Kind, subj.Name, subj.Namespace)
			}
		}

		return mcp.NewToolResultText(buf.String()), nil
	}
}

///////////////////////////////////////////////////////////////////////////////
// PODS
///////////////////////////////////////////////////////////////////////////////

func registerPodTools(srv *mcpsrv.MCPServer, c *clients.Clients) {
	listPodsTool := mcp.NewTool(
		"list_pods",
		mcp.WithDescription("List all pods. Args: namespace (string, optional), labelSelector (string, optional)."),
	)
	srv.AddTool(listPodsTool, listPodsHandler(c))

	getPodTool := mcp.NewTool(
		"get_pod",
		mcp.WithDescription("Get pod details. Args: name (string), namespace (string)."),
	)
	srv.AddTool(getPodTool, getPodHandler(c))

	logsTool := mcp.NewTool(
		"get_pod_logs",
		mcp.WithDescription("Get pod logs. Args: name (string), namespace (string), container (string, optional), tailLines (int, optional), previous (bool, optional)."),
	)
	srv.AddTool(logsTool, getPodLogsHandler(c))

	deleteTool := mcp.NewTool(
		"delete_pod",
		mcp.WithDescription("Delete a pod. Args: name (string), namespace (string)."),
	)
	srv.AddTool(deleteTool, deletePodHandler(c))

	execTool := mcp.NewTool(
		"exec_pod",
		mcp.WithDescription("Exec command in pod container. Args: name (string), namespace (string), container (string, optional), command ([]string)."),
	)
	srv.AddTool(execTool, execPodHandler(c))
}

func listPodsHandler(c *clients.Clients) mcpsrv.ToolHandlerFunc {
	return func(ctx context.Context, req mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		raw := req.Params.Arguments
		args, ok := raw.(map[string]any)
		if !ok {
			return mcp.NewToolResultError("invalid arguments payload (expected object)"), nil
		}

		namespace := utils.GetStringArg(args, "namespace", "")
		labelSelector := utils.GetStringArg(args, "labelSelector", "")

		pods, err := c.Kubernetes.CoreV1().Pods(namespace).List(ctx, metav1.ListOptions{
			LabelSelector: labelSelector,
		})
		if err != nil {
			return mcp.NewToolResultError(fmt.Sprintf("Failed to list pods: %v", err)), nil
		}

		return mcp.NewToolResultText(formatPodsList(pods)), nil
	}
}

func getPodHandler(c *clients.Clients) mcpsrv.ToolHandlerFunc {
	return func(ctx context.Context, req mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		raw := req.Params.Arguments
		args, ok := raw.(map[string]any)
		if !ok {
			return mcp.NewToolResultError("invalid arguments payload (expected object)"), nil
		}

		name, _ := args["name"].(string)
		if name == "" {
			return mcp.NewToolResultError("name is required"), nil
		}
		ns, _ := args["namespace"].(string)
		if ns == "" {
			return mcp.NewToolResultError("namespace is required"), nil
		}

		pod, err := c.Kubernetes.CoreV1().Pods(ns).Get(ctx, name, metav1.GetOptions{})
		if err != nil {
			return mcp.NewToolResultError(fmt.Sprintf("Failed to get pod: %v", err)), nil
		}

		return mcp.NewToolResultText(formatPodDetails(pod)), nil
	}
}

func getPodLogsHandler(c *clients.Clients) mcpsrv.ToolHandlerFunc {
	return func(ctx context.Context, req mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		raw := req.Params.Arguments
		args, ok := raw.(map[string]any)
		if !ok {
			return mcp.NewToolResultError("invalid arguments payload (expected object)"), nil
		}

		name, _ := args["name"].(string)
		if name == "" {
			return mcp.NewToolResultError("name is required"), nil
		}
		ns, _ := args["namespace"].(string)
		if ns == "" {
			return mcp.NewToolResultError("namespace is required"), nil
		}

		container := utils.GetStringArg(args, "container", "")
		tailLines := int64(utils.GetIntArg(args, "tailLines", 100))
		previous := utils.GetBoolArg(args, "previous", false)

		opts := &corev1.PodLogOptions{
			Container: container,
			TailLines: &tailLines,
			Previous:  previous,
		}

		stream, err := c.Kubernetes.CoreV1().Pods(ns).GetLogs(name, opts).Stream(ctx)
		if err != nil {
			return mcp.NewToolResultError(fmt.Sprintf("Failed to get pod logs: %v", err)), nil
		}
		defer stream.Close()

		buf := new(bytes.Buffer)
		if _, err := io.Copy(buf, stream); err != nil {
			return mcp.NewToolResultError(fmt.Sprintf("Failed to read logs: %v", err)), nil
		}

		return mcp.NewToolResultText(buf.String()), nil
	}
}

func deletePodHandler(c *clients.Clients) mcpsrv.ToolHandlerFunc {
	return func(ctx context.Context, req mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		raw := req.Params.Arguments
		args, ok := raw.(map[string]any)
		if !ok {
			return mcp.NewToolResultError("invalid arguments payload (expected object)"), nil
		}

		name, _ := args["name"].(string)
		if name == "" {
			return mcp.NewToolResultError("name is required"), nil
		}
		ns, _ := args["namespace"].(string)
		if ns == "" {
			return mcp.NewToolResultError("namespace is required"), nil
		}

		if err := c.Kubernetes.CoreV1().Pods(ns).Delete(ctx, name, metav1.DeleteOptions{}); err != nil {
			return mcp.NewToolResultError(fmt.Sprintf("Failed to delete pod: %v", err)), nil
		}

		return mcp.NewToolResultText(fmt.Sprintf("Pod %s in namespace %s deleted successfully", name, ns)), nil
	}
}

func execPodHandler(c *clients.Clients) mcpsrv.ToolHandlerFunc {
	return func(ctx context.Context, req mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		raw := req.Params.Arguments
		args, ok := raw.(map[string]any)
		if !ok {
			return mcp.NewToolResultError("invalid arguments payload (expected object)"), nil
		}

		name, _ := args["name"].(string)
		if name == "" {
			return mcp.NewToolResultError("name is required"), nil
		}
		ns, _ := args["namespace"].(string)
		if ns == "" {
			return mcp.NewToolResultError("namespace is required"), nil
		}

		container := utils.GetStringArg(args, "container", "")
		cmdRaw, ok := args["command"]
		if !ok {
			return mcp.NewToolResultError("command is required"), nil
		}
		command := utils.InterfaceSliceToStringSlice(cmdRaw)

		reqExec := c.Kubernetes.CoreV1().RESTClient().Post().
			Resource("pods").
			Name(name).
			Namespace(ns).
			SubResource("exec")

		execOpts := &corev1.PodExecOptions{
			Container: container,
			Command:   command,
			Stdout:    true,
			Stderr:    true,
		}

		reqExec.VersionedParams(execOpts, scheme.ParameterCodec)

		executor, err := remotecommand.NewSPDYExecutor(c.RestConfig, "POST", reqExec.URL())
		if err != nil {
			return mcp.NewToolResultError(fmt.Sprintf("Failed to create executor: %v", err)), nil
		}

		var stdout, stderr bytes.Buffer
		err = executor.StreamWithContext(ctx, remotecommand.StreamOptions{
			Stdout: &stdout,
			Stderr: &stderr,
		})
		if err != nil {
			return mcp.NewToolResultError(fmt.Sprintf("Exec failed: %v\nStderr: %s", err, stderr.String())), nil
		}

		result := fmt.Sprintf("Stdout:\n%s\n\nStderr:\n%s", stdout.String(), stderr.String())
		return mcp.NewToolResultText(result), nil
	}
}

func formatPodsList(pods *corev1.PodList) string {
	var buf bytes.Buffer
	fmt.Fprintf(&buf, "Total Pods: %d\n\n", len(pods.Items))
	for _, p := range pods.Items {
		fmt.Fprintf(&buf, "Name: %s\nNamespace: %s\nStatus: %s\nNode: %s\nIP: %s\n\n---\n\n",
			p.Name, p.Namespace, p.Status.Phase, p.Spec.NodeName, p.Status.PodIP)
	}
	return buf.String()
}

func formatPodDetails(pod *corev1.Pod) string {
	var buf bytes.Buffer
	fmt.Fprintf(&buf, "Pod: %s\nNamespace: %s\nStatus: %s\nNode: %s\nPod IP: %s\nHost IP: %s\n",
		pod.Name, pod.Namespace, pod.Status.Phase, pod.Spec.NodeName, pod.Status.PodIP, pod.Status.HostIP)
	return buf.String()
}

///////////////////////////////////////////////////////////////////////////////
// SERVICES
///////////////////////////////////////////////////////////////////////////////

func registerServiceTools(srv *mcpsrv.MCPServer, c *clients.Clients) {
	listSvc := mcp.NewTool(
		"list_services",
		mcp.WithDescription("List all services. Args: namespace (string, optional)."),
	)
	srv.AddTool(listSvc, listServicesHandler(c))

	getSvc := mcp.NewTool(
		"get_service",
		mcp.WithDescription("Get service details. Args: name (string), namespace (string)."),
	)
	srv.AddTool(getSvc, getServiceHandler(c))
}

func listServicesHandler(c *clients.Clients) mcpsrv.ToolHandlerFunc {
	return func(ctx context.Context, req mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		raw := req.Params.Arguments
		args, ok := raw.(map[string]any)
		if !ok {
			return mcp.NewToolResultError("invalid arguments payload (expected object)"), nil
		}

		ns := utils.GetStringArg(args, "namespace", "")

		svcs, err := c.Kubernetes.CoreV1().Services(ns).List(ctx, metav1.ListOptions{})
		if err != nil {
			return mcp.NewToolResultError(fmt.Sprintf("Failed to list services: %v", err)), nil
		}

		var buf bytes.Buffer
		fmt.Fprintf(&buf, "Total Services: %d\n\n", len(svcs.Items))
		for _, s := range svcs.Items {
			fmt.Fprintf(&buf, "Name: %s\nNamespace: %s\nType: %s\nClusterIP: %s\n\n---\n\n",
				s.Name, s.Namespace, s.Spec.Type, s.Spec.ClusterIP)
		}

		return mcp.NewToolResultText(buf.String()), nil
	}
}

func getServiceHandler(c *clients.Clients) mcpsrv.ToolHandlerFunc {
	return func(ctx context.Context, req mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		raw := req.Params.Arguments
		args, ok := raw.(map[string]any)
		if !ok {
			return mcp.NewToolResultError("invalid arguments payload (expected object)"), nil
		}

		name, _ := args["name"].(string)
		if name == "" {
			return mcp.NewToolResultError("name is required"), nil
		}
		ns, _ := args["namespace"].(string)
		if ns == "" {
			return mcp.NewToolResultError("namespace is required"), nil
		}

		svc, err := c.Kubernetes.CoreV1().Services(ns).Get(ctx, name, metav1.GetOptions{})
		if err != nil {
			return mcp.NewToolResultError(fmt.Sprintf("Failed to get service: %v", err)), nil
		}

		var buf bytes.Buffer
		fmt.Fprintf(&buf, "Service: %s\nNamespace: %s\nType: %s\nClusterIP: %s\n",
			svc.Name, svc.Namespace, svc.Spec.Type, svc.Spec.ClusterIP)

		return mcp.NewToolResultText(buf.String()), nil
	}
}
