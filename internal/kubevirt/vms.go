package kubevirt

import (
	"context"
	"encoding/json"
	"fmt"

	contextx "github.com/fmendonca/openshfit-mcp/internal/context"
	"github.com/mark3labs/mcp-go/mcp"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime/schema"
	"k8s.io/apimachinery/pkg/types"
	"k8s.io/client-go/dynamic"
)

// Registra tools específicas de KubeVirt
func RegisterKubeVirtTools(reg *mcp.ToolRegistry, ctx *contextx.ServerContext) {
	registerVMsListTool(reg, ctx)
	registerVMStartTool(reg, ctx)
	registerVMStopTool(reg, ctx)
	registerVMRestartTool(reg, ctx)
}

// ---------- List VMs ----------

type VMsListInput struct {
	Namespace string `json:"namespace,omitempty"`
}

func registerVMsListTool(reg *mcp.ToolRegistry, ctx *contextx.ServerContext) {
	reg.RegisterTool(&mcp.Tool{
		Name:        "kubevirt_vms_list",
		Description: "Lista VirtualMachines do KubeVirt (kubevirt.io/v1, resource 'virtualmachines').",
		InputSchema: &mcp.JSONSchema{
			Type: "object",
			Properties: map[string]*mcp.JSONSchema{
				"namespace": {Type: "string"},
			},
		},
	}, func(c context.Context, req *mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		var in VMsListInput
		if err := json.Unmarshal(req.Params, &in); err != nil {
			return mcp.NewErrorToolResult("falha ao decodificar argumentos", err), nil
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
			return mcp.NewErrorToolResult("erro ao listar VMs", err), nil
		}

		b, _ := vms.MarshalJSON()
		return mcp.NewToolResultText(string(b)), nil
	})
}

// ---------- VM start/stop/restart ----------

type VMActionInput struct {
	Namespace string `json:"namespace"`
	Name      string `json:"name"`
}

func vmGVR() schema.GroupVersionResource {
	return schema.GroupVersionResource{
		Group:    "kubevirt.io",
		Version:  "v1",
		Resource: "virtualmachines",
	}
}

func registerVMStartTool(reg *mcp.ToolRegistry, ctx *contextx.ServerContext) {
	reg.RegisterTool(&mcp.Tool{
		Name:        "kubevirt_vm_start",
		Description: "Inicia uma VirtualMachine do KubeVirt ajustando runStrategy para 'Always'.",
		InputSchema: &mcp.JSONSchema{
			Type: "object",
			Properties: map[string]*mcp.JSONSchema{
				"namespace": {Type: "string"},
				"name":      {Type: "string"},
			},
			Required: []string{"namespace", "name"},
		},
	}, func(c context.Context, req *mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		var in VMActionInput
		if err := json.Unmarshal(req.Params, &in); err != nil {
			return mcp.NewErrorToolResult("falha ao decodificar argumentos", err), nil
		}

		patch := []byte(`{"spec":{"runStrategy":"Always"}}`)

		ri := ctx.DynClient.Resource(vmGVR()).Namespace(in.Namespace)
		_, err := ri.Patch(c, in.Name, types.MergePatchType, patch, metav1.PatchOptions{})
		if err != nil {
			return mcp.NewErrorToolResult("erro ao iniciar VM", err), nil
		}

		return mcp.NewToolResultText(fmt.Sprintf("VM %s/%s iniciada (runStrategy=Always)", in.Namespace, in.Name)), nil
	})
}

func registerVMStopTool(reg *mcp.ToolRegistry, ctx *contextx.ServerContext) {
	reg.RegisterTool(&mcp.Tool{
		Name:        "kubevirt_vm_stop",
		Description: "Desliga uma VirtualMachine do KubeVirt ajustando runStrategy para 'Halted'.",
		InputSchema: &mcp.JSONSchema{
			Type: "object",
			Properties: map[string]*mcp.JSONSchema{
				"namespace": {Type: "string"},
				"name":      {Type: "string"},
			},
			Required: []string{"namespace", "name"},
		},
	}, func(c context.Context, req *mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		var in VMActionInput
		if err := json.Unmarshal(req.Params, &in); err != nil {
			return mcp.NewErrorToolResult("falha ao decodificar argumentos", err), nil
		}

		patch := []byte(`{"spec":{"runStrategy":"Halted"}}`)

		ri := ctx.DynClient.Resource(vmGVR()).Namespace(in.Namespace)
		_, err := ri.Patch(c, in.Name, types.MergePatchType, patch, metav1.PatchOptions{})
		if err != nil {
			return mcp.NewErrorToolResult("erro ao parar VM", err), nil
		}

		return mcp.NewToolResultText(fmt.Sprintf("VM %s/%s parada (runStrategy=Halted)", in.Namespace, in.Name)), nil
	})
}

func registerVMRestartTool(reg *mcp.ToolRegistry, ctx *contextx.ServerContext) {
	reg.RegisterTool(&mcp.Tool{
		Name:        "kubevirt_vm_restart",
		Description: "Reinicia uma VirtualMachine do KubeVirt adicionando um stateChangeRequest 'Restart'.",
		InputSchema: &mcp.JSONSchema{
			Type: "object",
			Properties: map[string]*mcp.JSONSchema{
				"namespace": {Type: "string"},
				"name":      {Type: "string"},
			},
			Required: []string{"namespace", "name"},
		},
	}, func(c context.Context, req *mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		var in VMActionInput
		if err := json.Unmarshal(req.Params, &in); err != nil {
			return mcp.NewErrorToolResult("falha ao decodificar argumentos", err), nil
		}

		patch := []byte(`{"spec":{"runStrategy":"Always","stateChangeRequests":[{"action":"Restart"}]}}`)

		ri := ctx.DynClient.Resource(vmGVR()).Namespace(in.Namespace)
		_, err := ri.Patch(c, in.Name, types.MergePatchType, patch, metav1.PatchOptions{})
		if err != nil {
			return mcp.NewErrorToolResult("erro ao reiniciar VM", err), nil
		}

		return mcp.NewToolResultText(fmt.Sprintf("VM %s/%s reiniciada (stateChangeRequests=Restart)", in.Namespace, in.Name)), nil
	})
}
