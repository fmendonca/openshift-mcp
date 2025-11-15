package mcpserver

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/fmendonca/openshfit-mcp/internal/k8s"
	"github.com/mark3labs/mcp-go/mcp"
	"github.com/mark3labs/mcp-go/server"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"
	"k8s.io/apimachinery/pkg/runtime/schema"
	"k8s.io/client-go/dynamic"
)

func RegisterCoreTools(s *server.MCPServer, ctx *ServerContext) {
	registerApisListTool(s, ctx)
	registerResourcesListTool(s, ctx)
	registerResourcesGetTool(s, ctx)
	registerResourcesApplyTool(s, ctx)
	registerResourcesDeleteTool(s, ctx)
}

func registerApisListTool(s *server.MCPServer, ctx *ServerContext) {
	tool := mcp.NewTool(
		"apis_list",
		mcp.WithDescription("Lista todos os API groups/versions/resources disponíveis (inclui OpenShift e KubeVirt)."),
	)

	s.AddTool(tool, func(c context.Context, req mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		out, err := k8s.ListServerGroupsAndResources(ctx.KubeClient.Discovery())
		if err != nil {
			return mcp.NewToolResultErrorFromErr("erro ao listar APIs", err), nil
		}
		return TextResult(out), nil
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

func registerResourcesListTool(s *server.MCPServer, ctx *ServerContext) {
	tool := mcp.NewTool(
		"resources_list",
		mcp.WithDescription("Lista qualquer recurso Kubernetes/OpenShift/KubeVirt por GVR (group/version/resource)."),
		mcp.WithString("group", mcp.Description("API group, vazio para core.")),
		mcp.WithString("version", mcp.Required(), mcp.Description("API version, ex: v1.")),
		mcp.WithString("resource", mcp.Required(), mcp.Description("Nome plural do recurso, ex: pods, routes.")),
		mcp.WithString("namespace", mcp.Description("Namespace (opcional).")),
		mcp.WithString("labelSelector", mcp.Description("Label selector (opcional).")),
		mcp.WithString("fieldSelector", mcp.Description("Field selector (opcional).")),
		mcp.WithNumber("limit", mcp.Description("Limite de itens (opcional).")),
	)

	s.AddTool(tool, func(c context.Context, req mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		var in ResourcesListInput
		if err := json.Unmarshal(req.Arguments, &in); err != nil {
			return mcp.NewToolResultError("falha ao decodificar argumentos"), nil
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
			return mcp.NewToolResultErrorFromErr("erro ao listar recursos", err), nil
		}

		b, err := list.MarshalJSON()
		if err != nil {
			return mcp.NewToolResultErrorFromErr("erro ao serializar lista", err), nil
		}

		return TextResult(string(b)), nil
	})
}

type ResourcesGetInput struct {
	Group     string `json:"group"`
	Version   string `json:"version"`
	Resource  string `json:"resource"`
	Namespace string `json:"namespace,omitempty"`
	Name      string `json:"name"`
}

func registerResourcesGetTool(s *server.MCPServer, ctx *ServerContext) {
	tool := mcp.NewTool(
		"resources_get",
		mcp.WithDescription("Obtém um recurso específico por Group/Version/Resource, namespace e nome."),
		mcp.WithString("group", mcp.Description("API group, vazio para core.")),
		mcp.WithString("version", mcp.Required(), mcp.Description("API version.")),
		mcp.WithString("resource", mcp.Required(), mcp.Description("Nome plural do recurso.")),
		mcp.WithString("namespace", mcp.Description("Namespace (opcional).")),
		mcp.WithString("name", mcp.Required(), mcp.Description("Nome do recurso.")),
	)

	s.AddTool(tool, func(c context.Context, req mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		var in ResourcesGetInput
		if err := json.Unmarshal(req.Arguments, &in); err != nil {
			return mcp.NewToolResultError("falha ao decodificar argumentos"), nil
		}

		gvr := schema.GroupVersionResource{
			Group:    in.Group,
			Version:  in.Version,
			Resource: in.Resource,
		}

		var ri dynamic.ResourceInterface
		if in.Namespace != "" {
			ri = ctx.DynClient.Resource(gvr).Namespace(in.Namespace)
		} else {
			ri = ctx.DynClient.Resource(gvr)
		}

		obj, err := ri.Get(c, in.Name, metav1.GetOptions{})
		if err != nil {
			return mcp.NewToolResultErrorFromErr("erro ao obter recurso", err), nil
		}

		b, err := obj.MarshalJSON()
		if err != nil {
			return mcp.NewToolResultErrorFromErr("erro ao serializar recurso", err), nil
		}

		return TextResult(string(b)), nil
	})
}

type ResourcesApplyInput struct {
	Object map[string]interface{} `json:"object"`
}

func registerResourcesApplyTool(s *server.MCPServer, ctx *ServerContext) {
	tool := mcp.NewTool(
		"resources_apply",
		mcp.WithDescription("Cria ou atualiza qualquer recurso Kubernetes/OpenShift/KubeVirt via objeto JSON genérico."),
		mcp.WithObject("object", mcp.Required(), mcp.Description("Objeto completo com apiVersion, kind, metadata..., mais campo __resource com o plural.")),
	)

	s.AddTool(tool, func(c context.Context, req mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		var in ResourcesApplyInput
		if err := json.Unmarshal(req.Arguments, &in); err != nil {
			return mcp.NewToolResultError("falha ao decodificar argumentos"), nil
		}

		u := &unstructured.Unstructured{Object: in.Object}

		apiVersion := u.GetAPIVersion()
		kind := u.GetKind()
		name := u.GetName()
		ns := u.GetNamespace()

		if apiVersion == "" || kind == "" || name == "" {
			return mcp.NewToolResultError("object deve ter apiVersion, kind e metadata.name"), nil
		}

		gv, err := schema.ParseGroupVersion(apiVersion)
		if err != nil {
			return mcp.NewToolResultErrorFromErr("apiVersion inválido", err), nil
		}

		plural, ok := in.Object["__resource"].(string)
		if !ok || plural == "" {
			return mcp.NewToolResultError("campo __resource (plural) é obrigatório"), nil
		}

		gvr := gv.WithResource(plural)

		var ri dynamic.ResourceInterface
		if ns != "" {
			ri = ctx.DynClient.Resource(gvr).Namespace(ns)
		} else {
			ri = ctx.DynClient.Resource(gvr)
		}

		// Tenta get; se existir, faz update; se não, create
		current, err := ri.Get(c, name, metav1.GetOptions{})
		if err == nil {
			u.SetResourceVersion(current.GetResourceVersion())
			updated, err := ri.Update(c, u, metav1.UpdateOptions{})
			if err != nil {
				return mcp.NewToolResultErrorFromErr("erro ao atualizar recurso", err), nil
			}
			b, _ := updated.MarshalJSON()
			return TextResult(string(b)), nil
		}

		created, err := ri.Create(c, u, metav1.CreateOptions{})
		if err != nil {
			return mcp.NewToolResultErrorFromErr("erro ao criar recurso", err), nil
		}
		b, _ := created.MarshalJSON()
		return TextResult(string(b)), nil
	})
}

type ResourcesDeleteInput struct {
	Group     string `json:"group"`
	Version   string `json:"version"`
	Resource  string `json:"resource"`
	Namespace string `json:"namespace,omitempty"`
	Name      string `json:"name"`
}

func registerResourcesDeleteTool(s *server.MCPServer, ctx *ServerContext) {
	tool := mcp.NewTool(
		"resources_delete",
		mcp.WithDescription("Deleta qualquer recurso Kubernetes/OpenShift/KubeVirt via GVR + nome."),
		mcp.WithString("group", mcp.Description("API group, vazio para core.")),
		mcp.WithString("version", mcp.Required(), mcp.Description("API version.")),
		mcp.WithString("resource", mcp.Required(), mcp.Description("Nome plural do recurso.")),
		mcp.WithString("namespace", mcp.Description("Namespace (opcional).")),
		mcp.WithString("name", mcp.Required(), mcp.Description("Nome do recurso.")),
	)

	s.AddTool(tool, func(c context.Context, req mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		var in ResourcesDeleteInput
		if err := json.Unmarshal(req.Arguments, &in); err != nil {
			return mcp.NewToolResultError("falha ao decodificar argumentos"), nil
		}

		gvr := schema.GroupVersionResource{
			Group:    in.Group,
			Version:  in.Version,
			Resource: in.Resource,
		}

		var ri dynamic.ResourceInterface
		if in.Namespace != "" {
			ri = ctx.DynClient.Resource(gvr).Namespace(in.Namespace)
		} else {
			ri = ctx.DynClient.Resource(gvr)
		}

		if err := ri.Delete(c, in.Name, metav1.DeleteOptions{}); err != nil {
			return mcp.NewToolResultErrorFromErr("erro ao deletar recurso", err), nil
		}

		return TextResult(fmt.Sprintf("deleted %s/%s", in.Resource, in.Name)), nil
	})
}
