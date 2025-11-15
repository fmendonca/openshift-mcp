package kubevirt

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/mark3labs/mcp-go/mcp"
	"github.com/mark3labs/mcp-go/server"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime/schema"
	"k8s.io/apimachinery/pkg/types"
	"k8s.io/client-go/dynamic"
)

func RegisterKubeVirtTools(s *server.MCPServer, ctx *mcpserver.ServerContext) {
	registerVMsListTool(s, ctx)
	registerVMStartTool(s, ctx)
	registerVMStopTool(s, ctx)
	registerVMRestartTool(s, ctx)
}

// ---------- List VMs ----------

type VMsListInput struct {
	Namespace string `json:"namespace,omitempty"`
}

func registerVMsListTool(s *server.MCPServer, ctx *mcpserver.ServerContext) {
	tool := mcp.NewTool(
		"kubevirt_vms_list",
		mcp.WithDescription("Lista VirtualMachines do KubeVirt (kubevirt.io/v1, resource 'virtualmachines')."),
		mcp.WithString("namespace", mcp.Description("Namespace (opcional).")),
	)

	s.AddTool(tool, func(c context.Context, req mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		var in VMsListInput
		if err := json.Unmarshal(req.Arguments, &in); err != nil {
			return mcp.NewToolResultError("falha ao decodificar argumentos"), nil
		}

		gvr := schema.GroupVersionResource{
			Group:    "kubevirt.io",
			Version:  "v1",
			Resource: "virtualmachines",
		}

		var ri dynamic.ResourceInterface
		if in.Namespace != "" {
			ri = ctx.DynClient.Resource(gvr).Namespace(in.Namespace)
		} else {
			ri = ctx.DynClient.Resource(gvr)
		}

		vms, err := ri.List(c, metav1.ListOptions{})
		if err != nil {
			return mcp.NewToolResultErrorFromErr("erro ao listar VMs", err), nil
		}

		b, _ := vms.MarshalJSON()
		return mcpserver.TextResult(string(b)), nil
	})
}

// ---------- VM start/stop/restart via runStrategy ----------

type VMActionInput struct {
	Namespace string `json:"namespace"`
	Name      string `json:"name"`
}

func vmGVR() schema.GroupVersionResource {
	return schema.GroupVersionResource{
		Group:    "kubevirt.io",
		Version:  "v1",
		Resource: "virtualmachines",
	} // VirtualMachine CRD[web:45][web:51]
}

// Start: define runStrategy = "Always" (ou running=true, dependendo da sua política)
func registerVMStartTool(s *server.MCPServer, ctx *mcpserver.ServerContext) {
	tool := mcp.NewTool(
		"kubevirt_vm_start",
		mcp.WithDescription("Inicia uma VirtualMachine do KubeVirt ajustando runStrategy para 'Always'."),
		mcp.WithString("namespace", mcp.Required(), mcp.Description("Namespace da VM.")),
		mcp.WithString("name", mcp.Required(), mcp.Description("Nome da VirtualMachine.")),
	)

	s.AddTool(tool, func(c context.Context, req mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		var in VMActionInput
		if err := json.Unmarshal(req.Arguments, &in); err != nil {
			return mcp.NewToolResultError("falha ao decodificar argumentos"), nil
		}

		patch := []byte(`{"spec":{"runStrategy":"Always"}}`) // liga a VM[web:95][web:99]

		ri := ctx.DynClient.Resource(vmGVR()).Namespace(in.Namespace)
		_, err := ri.Patch(c, in.Name, types.MergePatchType, patch, metav1.PatchOptions{})
		if err != nil {
			return mcp.NewToolResultErrorFromErr("erro ao iniciar VM", err), nil
		}

		return mcpserver.TextResult(fmt.Sprintf("VM %s/%s iniciada (runStrategy=Always)", in.Namespace, in.Name)), nil
	})
}

// Stop: define runStrategy = "Halted"
func registerVMStopTool(s *server.MCPServer, ctx *mcpserver.ServerContext) {
	tool := mcp.NewTool(
		"kubevirt_vm_stop",
		mcp.WithDescription("Desliga uma VirtualMachine do KubeVirt ajustando runStrategy para 'Halted'."),
		mcp.WithString("namespace", mcp.Required(), mcp.Description("Namespace da VM.")),
		mcp.WithString("name", mcp.Required(), mcp.Description("Nome da VirtualMachine.")),
	)

	s.AddTool(tool, func(c context.Context, req mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		var in VMActionInput
		if err := json.Unmarshal(req.Arguments, &in); err != nil {
			return mcp.NewToolResultError("falha ao decodificar argumentos"), nil
		}

		patch := []byte(`{"spec":{"runStrategy":"Halted"}}`) // desliga a VM[web:95][web:103]

		ri := ctx.DynClient.Resource(vmGVR()).Namespace(in.Namespace)
		_, err := ri.Patch(c, in.Name, types.MergePatchType, patch, metav1.PatchOptions{})
		if err != nil {
			return mcp.NewToolResultErrorFromErr("erro ao parar VM", err), nil
		}

		return mcpserver.TextResult(fmt.Sprintf("VM %s/%s parada (runStrategy=Halted)", in.Namespace, in.Name)), nil
	})
}

// Restart: adiciona um stateChangeRequest restart, conforme comportamento do virtctl/manual[web:99][web:103]
func registerVMRestartTool(s *server.MCPServer, ctx *mcpserver.ServerContext) {
	tool := mcp.NewTool(
		"kubevirt_vm_restart",
		mcp.WithDescription("Reinicia uma VirtualMachine do KubeVirt adicionando um stateChangeRequest 'Restart'."),
		mcp.WithString("namespace", mcp.Required(), mcp.Description("Namespace da VM.")),
		mcp.WithString("name", mcp.Required(), mcp.Description("Nome da VirtualMachine.")),
	)

	s.AddTool(tool, func(c context.Context, req mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		var in VMActionInput
		if err := json.Unmarshal(req.Arguments, &in); err != nil {
			return mcp.NewToolResultError("falha ao decodificar argumentos"), nil
		}

		// Padrão de restart: adicionar StateChangeRequest "Restart"
		patch := []byte(`{"spec":{"runStrategy":"Always","stateChangeRequests":[{"action":"Restart"}]}}`) // ilustra comportamento restart[web:99]

		ri := ctx.DynClient.Resource(vmGVR()).Namespace(in.Namespace)
		_, err := ri.Patch(c, in.Name, types.MergePatchType, patch, metav1.PatchOptions{})
		if err != nil {
			return mcp.NewToolResultErrorFromErr("erro ao reiniciar VM", err), nil
		}

		return mcpserver.TextResult(fmt.Sprintf("VM %s/%s reiniciada (stateChangeRequests=Restart)", in.Namespace, in.Name)), nil
	})
}
