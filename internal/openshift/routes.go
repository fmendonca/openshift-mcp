package openshift

import (
	"context"
	"encoding/json"

	"github.com/fmendonca/openshfit-mcp/internal/mcpserver"
	"github.com/mark3labs/mcp-go/mcp"
	"github.com/mark3labs/mcp-go/server"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime/schema"
	"k8s.io/client-go/dynamic"
)

func RegisterOpenShiftTools(s *server.MCPServer, ctx *mcpserver.ServerContext) {
	registerRoutesListTool(s, ctx)
	registerBuildConfigsListTool(s, ctx)
	registerImageStreamsListTool(s, ctx)
	registerProjectsListTool(s, ctx)
	registerDeploymentConfigsListTool(s, ctx)
}

// ---------- Routes ----------

type RoutesListInput struct {
	Namespace string `json:"namespace,omitempty"`
}

func registerRoutesListTool(s *server.MCPServer, ctx *mcpserver.ServerContext) {
	tool := mcp.NewTool(
		"routes_list",
		mcp.WithDescription("Lista Routes do OpenShift (route.openshift.io/v1, resource 'routes')."),
		mcp.WithString("namespace", mcp.Description("Namespace (opcional).")),
	)

	s.AddTool(tool, func(c context.Context, req mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		var in RoutesListInput
		if err := json.Unmarshal(req.Arguments, &in); err != nil {
			return mcp.NewToolResultError("falha ao decodificar argumentos"), nil
		}

		gvr := schema.GroupVersionResource{
			Group:    "route.openshift.io",
			Version:  "v1",
			Resource: "routes",
		}

		var ri dynamic.ResourceInterface
		if in.Namespace != "" {
			ri = ctx.DynClient.Resource(gvr).Namespace(in.Namespace)
		} else {
			ri = ctx.DynClient.Resource(gvr)
		}

		routes, err := ri.List(c, metav1.ListOptions{})
		if err != nil {
			return mcp.NewToolResultErrorFromErr("erro ao listar routes", err), nil
		}

		b, _ := routes.MarshalJSON()
		return mcpserver.TextResult(string(b)), nil
	})
}

// ---------- BuildConfig ----------

type BuildConfigsListInput struct {
	Namespace string `json:"namespace,omitempty"`
}

func registerBuildConfigsListTool(s *server.MCPServer, ctx *mcpserver.ServerContext) {
	tool := mcp.NewTool(
		"buildconfigs_list",
		mcp.WithDescription("Lista BuildConfigs (build.openshift.io/v1, resource 'buildconfigs')."),
		mcp.WithString("namespace", mcp.Description("Namespace (opcional).")),
	)

	s.AddTool(tool, func(c context.Context, req mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		var in BuildConfigsListInput
		if err := json.Unmarshal(req.Arguments, &in); err != nil {
			return mcp.NewToolResultError("falha ao decodificar argumentos"), nil
		}

		gvr := schema.GroupVersionResource{
			Group:    "build.openshift.io",
			Version:  "v1",
			Resource: "buildconfigs",
		} // BuildConfig API[web:88][web:91]

		var ri dynamic.ResourceInterface
		if in.Namespace != "" {
			ri = ctx.DynClient.Resource(gvr).Namespace(in.Namespace)
		} else {
			ri = ctx.DynClient.Resource(gvr)
		}

		bcs, err := ri.List(c, metav1.ListOptions{})
		if err != nil {
			return mcp.NewToolResultErrorFromErr("erro ao listar BuildConfigs", err), nil
		}

		b, _ := bcs.MarshalJSON()
		return mcpserver.TextResult(string(b)), nil
	})
}

// ---------- ImageStream ----------

type ImageStreamsListInput struct {
	Namespace string `json:"namespace,omitempty"`
}

func registerImageStreamsListTool(s *server.MCPServer, ctx *mcpserver.ServerContext) {
	tool := mcp.NewTool(
		"imagestreams_list",
		mcp.WithDescription("Lista ImageStreams (image.openshift.io/v1, resource 'imagestreams')."),
		mcp.WithString("namespace", mcp.Description("Namespace (opcional).")),
	)

	s.AddTool(tool, func(c context.Context, req mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		var in ImageStreamsListInput
		if err := json.Unmarshal(req.Arguments, &in); err != nil {
			return mcp.NewToolResultError("falha ao decodificar argumentos"), nil
		}

		gvr := schema.GroupVersionResource{
			Group:    "image.openshift.io",
			Version:  "v1",
			Resource: "imagestreams",
		} // ImageStream API[web:94][web:98]

		var ri dynamic.ResourceInterface
		if in.Namespace != "" {
			ri = ctx.DynClient.Resource(gvr).Namespace(in.Namespace)
		} else {
			ri = ctx.DynClient.Resource(gvr)
		}

		iss, err := ri.List(c, metav1.ListOptions{})
		if err != nil {
			return mcp.NewToolResultErrorFromErr("erro ao listar ImageStreams", err), nil
		}

		b, _ := iss.MarshalJSON()
		return mcpserver.TextResult(string(b)), nil
	})
}

// ---------- Project ----------

type ProjectsListInput struct{}

func registerProjectsListTool(s *server.MCPServer, ctx *mcpserver.ServerContext) {
	tool := mcp.NewTool(
		"projects_list",
		mcp.WithDescription("Lista Projects (project.openshift.io/v1, resource 'projects')."),
	)

	s.AddTool(tool, func(c context.Context, req mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		gvr := schema.GroupVersionResource{
			Group:    "project.openshift.io",
			Version:  "v1",
			Resource: "projects",
		} // Project API[web:56]

		ri := ctx.DynClient.Resource(gvr)

		projects, err := ri.List(c, metav1.ListOptions{})
		if err != nil {
			return mcp.NewToolResultErrorFromErr("erro ao listar Projects", err), nil
		}

		b, _ := projects.MarshalJSON()
		return mcpserver.TextResult(string(b)), nil
	})
}

// ---------- DeploymentConfig ----------

type DeploymentConfigsListInput struct {
	Namespace string `json:"namespace,omitempty"`
}

func registerDeploymentConfigsListTool(s *server.MCPServer, ctx *mcpserver.ServerContext) {
	tool := mcp.NewTool(
		"deploymentconfigs_list",
		mcp.WithDescription("Lista DeploymentConfigs (apps.openshift.io/v1, resource 'deploymentconfigs')."),
		mcp.WithString("namespace", mcp.Description("Namespace (opcional).")),
	)

	s.AddTool(tool, func(c context.Context, req mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		var in DeploymentConfigsListInput
		if err := json.Unmarshal(req.Arguments, &in); err != nil {
			return mcp.NewToolResultError("falha ao decodificar argumentos"), nil
		}

		gvr := schema.GroupVersionResource{
			Group:    "apps.openshift.io",
			Version:  "v1",
			Resource: "deploymentconfigs",
		} // DeploymentConfig API[web:93][web:97]

		var ri dynamic.ResourceInterface
		if in.Namespace != "" {
			ri = ctx.DynClient.Resource(gvr).Namespace(in.Namespace)
		} else {
			ri = ctx.DynClient.Resource(gvr)
		}

		dcs, err := ri.List(c, metav1.ListOptions{})
		if err != nil {
			return mcp.NewToolResultErrorFromErr("erro ao listar DeploymentConfigs", err), nil
		}

		b, _ := dcs.MarshalJSON()
		return mcpserver.TextResult(string(b)), nil
	})
}

// ---------- BuildConfig actions ----------

type BuildConfigStartBuildInput struct {
	Namespace string `json:"namespace"`
	Name      string `json:"name"` // nome do BuildConfig
}

func registerBuildConfigStartBuildTool(s *server.MCPServer, ctx *mcpserver.ServerContext) {
	tool := mcp.NewTool(
		"buildconfig_start_build",
		mcp.WithDescription("Dispara um build a partir de um BuildConfig (equivalente a 'oc start-build')."),
		mcp.WithString("namespace", mcp.Required(), mcp.Description("Namespace do BuildConfig.")),
		mcp.WithString("name", mcp.Required(), mcp.Description("Nome do BuildConfig.")),
	)

	s.AddTool(tool, func(c context.Context, req mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		var in BuildConfigStartBuildInput
		if err := json.Unmarshal(req.Arguments, &in); err != nil {
			return mcp.NewToolResultError("falha ao decodificar argumentos"), nil
		}

		// BuildRequest CR de build.openshift.io/v1[web:106][web:115]
		gvr := schema.GroupVersionResource{
			Group:    "build.openshift.io",
			Version:  "v1",
			Resource: "buildconfigs",
		}

		// Subresource "instantiate" do BuildConfig não é diretamente exposto pelo client dinâmico,
		// mas podemos criar um BuildRequest no recurso "builds" com referência ao BuildConfig.[web:107]

		buildReq := map[string]interface{}{
			"apiVersion": "build.openshift.io/v1",
			"kind":       "BuildRequest",
			"metadata": map[string]interface{}{
				"name":      in.Name,
				"namespace": in.Namespace,
			},
			"kindRef": map[string]interface{}{
				"name": in.Name,
			},
		}

		// Muitos clusters aceitam POST em:
		// /apis/build.openshift.io/v1/namespaces/{ns}/buildconfigs/{name}/instantiate[web:107][web:112]
		// Como o dynamic não tem API direta para subresource, usamos o recurso 'builds' com spec.from, se necessário.
		// Para simplificar, vamos usar create em 'builds' com spec.source e spec.strategy herdados do BuildConfig.

		// 1) pega o BuildConfig
		bc, err := ctx.DynClient.Resource(gvr).Namespace(in.Namespace).Get(c, in.Name, metav1.GetOptions{})
		if err != nil {
			return mcp.NewToolResultErrorFromErr("erro ao obter BuildConfig", err), nil
		}

		// 2) monta um Build baseado no BuildConfig (estilo oc start-build)[web:115]
		build := map[string]interface{}{
			"apiVersion": "build.openshift.io/v1",
			"kind":       "Build",
			"metadata": map[string]interface{}{
				"generateName": in.Name + "-",
				"namespace":    in.Namespace,
				"labels": map[string]interface{}{
					"buildconfig": in.Name,
				},
			},
			"spec": map[string]interface{}{
				"serviceAccount": bc.Object["spec"].(map[string]interface{})["serviceAccount"],
				"source":         bc.Object["spec"].(map[string]interface{})["source"],
				"strategy":       bc.Object["spec"].(map[string]interface{})["strategy"],
				"output":         bc.Object["spec"].(map[string]interface{})["output"],
			},
		}

		buildsGVR := schema.GroupVersionResource{
			Group:    "build.openshift.io",
			Version:  "v1",
			Resource: "builds",
		}

		created, err := ctx.DynClient.Resource(buildsGVR).Namespace(in.Namespace).Create(
			c,
			&unstructured.Unstructured{Object: build},
			metav1.CreateOptions{},
		)
		if err != nil {
			return mcp.NewToolResultErrorFromErr("erro ao criar Build", err), nil
		}

		b, _ := created.MarshalJSON()
		return mcpserver.TextResult(string(b)), nil
	})
}
