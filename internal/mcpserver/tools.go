package mcpserver

import (
	"context"

	contextx "github.com/fmendonca/openshfit-mcp/internal/context"
	"github.com/fmendonca/openshfit-mcp/internal/k8s"
	"github.com/mark3labs/mcp-go/mcp"
	"github.com/mark3labs/mcp-go/server"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"
	"k8s.io/apimachinery/pkg/runtime/schema"
	"k8s.io/client-go/dynamic"
)

// helpers para ler args
func getStringArg(req mcp.CallToolRequest, key, def string) string {
	args, ok := req.Params.Arguments.(map[string]any)
	if !ok {
		return def
	}
	if v, ok := args[key]; ok {
		if s, ok2 := v.(string); ok2 {
			return s
		}
	}
	return def
}

func getIntArg(req mcp.CallToolRequest, key string, def int64) int64 {
	args, ok := req.Params.Arguments.(map[string]any)
	if !ok {
		return def
	}
	if v, ok := args[key]; ok {
		switch n := v.(type) {
		case float64:
			return int64(n)
		case int64:
			return n
		}
	}
	return def
}

// ---------- Registro das tools core ----------

func RegisterCoreTools(s *server.MCPServer, ctx *contextx.ServerContext) {
	registerApisListTool(s, ctx)
	registerResourcesListTool(s, ctx)
	registerResourcesGetTool(s, ctx)
	registerResourcesApplyTool(s, ctx)
	registerResourcesDeleteTool(s, ctx)
}

// ---------- apis_list ----------

func registerApisListTool(s *server.MCPServer, ctx *contextx.ServerContext) {
	tool := mcp.NewTool(
		"apis_list",
		mcp.WithDescription("Lista todos os API groups/versions/resources disponíveis (inclui OpenShift e KubeVirt)."),
	)

	s.AddTool(tool, func(c context.Context, req mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		out, err := k8s.ListServerGroupsAndResources(ctx.KubeClient.Discovery())
		if err != nil {
			return mcp.NewToolResultErrorFromErr("erro ao listar APIs", err), nil
		}
		return mcp.NewToolResultText(out), nil
	})
}

// ---------- resources_list ----------

type ResourcesListInput struct {
	Group         string `json:"group"`
	Version       string `json:"version"`
	Resource      string `json:"resource"`
	Namespace     string `json:"namespace,omitempty"`
	LabelSelector string `json:"labelSelector,omitempty"`
	FieldSelector string `json:"fieldSelector,omitempty"`
	Limit         int64  `json:"limit,omitempty"`
}

func registerResourcesListTool(s *server.MCPServer, ctx *contextx.ServerContext) {
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
		in := ResourcesListInput{
			Group:         getStringArg(req, "group", ""),
			Version:       getStringArg(req, "version", ""),
			Resource:      getStringArg(req, "resource", ""),
			Namespace:     getStringArg(req, "namespace", ""),
			LabelSelector: getStringArg(req, "labelSelector", ""),
			FieldSelector: getStringArg(req, "fieldSelector", ""),
			Limit:         getIntArg(req, "limit", 0),
		}
		if in.Version == "" || in.Resource == "" {
			return mcp.NewToolResultError("version e resource são obrigatórios"), nil
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

		return mcp.NewToolResultText(string(b)), nil
	})
}

// ---------- resources_get ----------

type ResourcesGetInput struct {
	Group     string `json:"group"`
	Version   string `json:"version"`
	Resource  string `json:"resource"`
	Namespace string `json:"namespace,omitempty"`
	Name      string `json:"name"`
}

func registerResourcesGetTool(s *server.MCPServer, ctx *contextx.ServerContext) {
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
		in := ResourcesGetInput{
			Group:     getStringArg(req, "group", ""),
			Version:   getStringArg(req, "version", ""),
			Resource:  getStringArg(req, "resource", ""),
			Namespace: getStringArg(req, "namespace", ""),
			Name:      getStringArg(req, "name", ""),
		}
		if in.Version == "" || in.Resource == "" || in.Name == "" {
			return mcp.NewToolResultError("version, resource e name são obrigatórios"), nil
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

		return mcp.NewToolResultText(string(b)), nil
	})
}

// ---------- resources_apply ----------

type ResourcesApplyInput struct {
	// Para simplicidade, vamos aceitar os campos básicos como argumentos
	APIVersion string                 `json:"apiVersion"`
	Kind       string                 `json:"kind"`
	Namespace  string                 `json:"namespace,omitempty"`
	Name       string                 `json:"name"`
	Resource   string                 `json:"resource"` // plural
	Object     map[string]interface{} `json:"object"`   // conteúdo completo
}

func registerResourcesApplyTool(s *server.MCPServer, ctx *contextx.ServerContext) {
	tool := mcp.NewTool(
		"resources_apply",
		mcp.WithDescription("Cria ou atualiza qualquer recurso Kubernetes/OpenShift/KubeVirt via objeto JSON genérico."),
		mcp.WithString("apiVersion", mcp.Required(), mcp.Description("apiVersion do recurso.")),
		mcp.WithString("kind", mcp.Required(), mcp.Description("Kind do recurso.")),
		mcp.WithString("namespace", mcp.Description("Namespace (opcional).")),
		mcp.WithString("name", mcp.Required(), mcp.Description("Nome do recurso.")),
		mcp.WithString("resource", mcp.Required(), mcp.Description("Nome plural (GVR) do recurso.")),
		// object virá como arguments["object"], tratado manualmente
	)

	s.AddTool(tool, func(c context.Context, req mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		args, ok := req.Params.Arguments.(map[string]any)
		if !ok {
			return mcp.NewToolResultError("arguments deve ser um objeto"), nil
		}

		in := ResourcesApplyInput{
			APIVersion: getStringArg(req, "apiVersion", ""),
			Kind:       getStringArg(req, "kind", ""),
			Namespace:  getStringArg(req, "namespace", ""),
			Name:       getStringArg(req, "name", ""),
			Resource:   getStringArg(req, "resource", ""),
		}
		objRaw, ok := args["object"]
		if !ok {
			return mcp.NewToolResultError("campo object é obrigatório"), nil
		}
		objMap, ok := objRaw.(map[string]any)
		if !ok {
			return mcp.NewToolResultError("object deve ser um JSON object"), nil
		}
		in.Object = objMap

		if in.APIVersion == "" || in.Kind == "" || in.Name == "" || in.Resource == "" {
			return mcp.NewToolResultError("apiVersion, kind, name e resource são obrigatórios"), nil
		}

		gv, err := schema.ParseGroupVersion(in.APIVersion)
		if err != nil {
			return mcp.NewToolResultErrorFromErr("apiVersion inválido", err), nil
		}

		gvr := gv.WithResource(in.Resource)

		u := &unstructured.Unstructured{Object: in.Object}
		u.SetAPIVersion(in.APIVersion)
		u.SetKind(in.Kind)
		u.SetName(in.Name)
		if in.Namespace != "" {
			u.SetNamespace(in.Namespace)
		}

		var ri dynamic.ResourceInterface
		if in.Namespace != "" {
			ri = ctx.DynClient.Resource(gvr).Namespace(in.Namespace)
		} else {
			ri = ctx.DynClient.Resource(gvr)
		}

		current, err := ri.Get(c, in.Name, metav1.GetOptions{})
		if err == nil {
			u.SetResourceVersion(current.GetResourceVersion())
			updated, err := ri.Update(c, u, metav1.UpdateOptions{})
			if err != nil {
				return mcp.NewToolResultErrorFromErr("erro ao atualizar recurso", err), nil
			}
			b, _ := updated.MarshalJSON()
			return mcp.NewToolResultText(string(b)), nil
		}

		created, err := ri.Create(c, u, metav1.CreateOptions{})
		if err != nil {
			return mcp.NewToolResultErrorFromErr("erro ao criar recurso", err), nil
		}
		b, _ := created.MarshalJSON()
		return mcp.NewToolResultText(string(b)), nil
	})
}

// ---------- resources_delete ----------

type ResourcesDeleteInput struct {
	Group     string `json:"group"`
	Version   string `json:"version"`
	Resource  string `json:"resource"`
	Namespace string `json:"namespace,omitempty"`
	Name      string `json:"name"`
}

func registerResourcesDeleteTool(s *server.MCPServer, ctx *contextx.ServerContext) {
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
		in := ResourcesDeleteInput{
			Group:     getStringArg(req, "group", ""),
			Version:   getStringArg(req, "version", ""),
			Resource:  getStringArg(req, "resource", ""),
			Namespace: getStringArg(req, "namespace", ""),
			Name:      getStringArg(req, "name", ""),
		}
		if in.Version == "" || in.Resource == "" || in.Name == "" {
			return mcp.NewToolResultError("version, resource e name são obrigatórios"), nil
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

		return mcp.NewToolResultText("deleted"), nil
	})
}
