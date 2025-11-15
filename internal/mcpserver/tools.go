package mcpserver

import (
	"context"
	"encoding/json"

	contextx "github.com/fmendonca/openshfit-mcp/internal/context"
	"github.com/fmendonca/openshfit-mcp/internal/k8s"
	"github.com/mark3labs/mcp-go/mcp"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime/schema"
	"k8s.io/client-go/dynamic"
)

func RegisterCoreTools(reg *mcp.ToolRegistry, ctx *contextx.ServerContext) {
	registerApisListTool(reg, ctx)
	registerResourcesListTool(reg, ctx)
	registerResourcesGetTool(reg, ctx)
	registerResourcesApplyTool(reg, ctx)
	registerResourcesDeleteTool(reg, ctx)
}

func registerApisListTool(reg *mcp.ToolRegistry, ctx *contextx.ServerContext) {
	reg.RegisterTool(&mcp.Tool{
		Name:        "apis_list",
		Description: "Lista todos os API groups/versions/resources disponíveis (inclui OpenShift e KubeVirt).",
	}, func(c context.Context, req *mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		out, err := k8s.ListServerGroupsAndResources(ctx.KubeClient.Discovery())
		if err != nil {
			return mcp.NewErrorToolResult("erro ao listar APIs", err), nil
		}
		return mcp.NewToolResultText(out), nil
	})
}

type ResourcesListInput struct {
	Group         string `json:"group"`
	Version       string `json:"version"`
	Resource      string `json:"resource"`
	Namespace     string `json:"namespace,omitempty"`
	LabelSelector string `json:"labelSelector,omitempty"`
	FieldSelector string `json:"fieldSelector,omitempty"`
	Limit         int64  `json:"limit,omitempty"`
}

func registerResourcesListTool(reg *mcp.ToolRegistry, ctx *contextx.ServerContext) {
	reg.RegisterTool(&mcp.Tool{
		Name:        "resources_list",
		Description: "Lista qualquer recurso Kubernetes/OpenShift/KubeVirt por GVR (group/version/resource).",
		InputSchema: &mcp.JSONSchema{
			Type: "object",
			Properties: map[string]*mcp.JSONSchema{
				"group":         {Type: "string"},
				"version":       {Type: "string"},
				"resource":      {Type: "string"},
				"namespace":     {Type: "string"},
				"labelSelector": {Type: "string"},
				"fieldSelector": {Type: "string"},
				"limit":         {Type: "number"},
			},
			Required: []string{"version", "resource"},
		},
	}, func(c context.Context, req *mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		var in ResourcesListInput
		if err := json.Unmarshal(req.Params, &in); err != nil {
			return mcp.NewErrorToolResult("falha ao decodificar argumentos", err), nil
		}

		gvr := schema.GroupVersionResource{
			Group:    in.Group,
			Version:  in.Version,
			Resource: in.Resource,
		}

		opts := metav1.ListOptions{
			LabelSelector: in.LabelSelector,
			FieldSelector: in.FieldSelector,
		}
		if in.Limit > 0 {
			opts.Limit = in.Limit
		}

		var ri dynamic.ResourceInterface
		if in.Namespace != "" {
			ri = ctx.DynClient.Resource(gvr).Namespace(in.Namespace)
		} else {
			ri = ctx.DynClient.Resource(gvr)
		}

		list, err := ri.List(c, opts)
		if err != nil {
			return mcp.NewErrorToolResult("erro ao listar recursos", err), nil
		}

		b, err := list.MarshalJSON()
		if err != nil {
			return mcp.NewErrorToolResult("erro ao serializar lista", err), nil
		}

		return mcp.NewToolResultText(string(b)), nil
	})
}
