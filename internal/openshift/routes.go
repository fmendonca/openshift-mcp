package openshift

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/mark3labs/mcp-go/mcp"
	"github.com/mark3labs/mcp-go/server"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"
	"k8s.io/apimachinery/pkg/runtime/schema"
	"k8s.io/apimachinery/pkg/types"
	"k8s.io/client-go/dynamic"
)

// Registra todas as tools específicas de OpenShift
func RegisterOpenShiftTools(s *server.MCPServer, ctx *mcpserver.ServerContext) {
	registerRoutesListTool(s, ctx)
	registerBuildConfigsListTool(s, ctx)
	registerImageStreamsListTool(s, ctx)
	registerProjectsListTool(s, ctx)
	registerDeploymentConfigsListTool(s, ctx)

	// ações
	registerBuildConfigStartBuildTool(s, ctx)
	registerDeploymentConfigRolloutTool(s, ctx)
	registerImageStreamPromoteTagTool(s, ctx)
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
		}

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
		}

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
		}

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
		}

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

		// BuildConfig GVR
		bcGVR := schema.GroupVersionResource{
			Group:    "build.openshift.io",
			Version:  "v1",
			Resource: "buildconfigs",
		}

		// 1) pega o BuildConfig
		bc, err := ctx.DynClient.Resource(bcGVR).Namespace(in.Namespace).Get(c, in.Name, metav1.GetOptions{})
		if err != nil {
			return mcp.NewToolResultErrorFromErr("erro ao obter BuildConfig", err), nil
		}

		spec, ok := bc.Object["spec"].(map[string]interface{})
		if !ok {
			return mcp.NewToolResultError("BuildConfig sem spec válido"), nil
		}

		// 2) monta um Build baseado no BuildConfig (estilo oc start-build)
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
				"source":   spec["source"],
				"strategy": spec["strategy"],
				"output":   spec["output"],
			},
		}
		// serviceAccount é opcional; só coloca se existir
		if sa, exists := spec["serviceAccount"]; exists {
			build["spec"].(map[string]interface{})["serviceAccount"] = sa
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

// ---------- DeploymentConfig actions ----------

type DeploymentConfigRolloutInput struct {
	Namespace string `json:"namespace"`
	Name      string `json:"name"`
}

func registerDeploymentConfigRolloutTool(s *server.MCPServer, ctx *mcpserver.ServerContext) {
	tool := mcp.NewTool(
		"deploymentconfig_rollout_latest",
		mcp.WithDescription("Dispara um rollout manual da DeploymentConfig (equivalente a 'oc rollout latest')."),
		mcp.WithString("namespace", mcp.Required(), mcp.Description("Namespace da DeploymentConfig.")),
		mcp.WithString("name", mcp.Required(), mcp.Description("Nome da DeploymentConfig.")),
	)

	s.AddTool(tool, func(c context.Context, req mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		var in DeploymentConfigRolloutInput
		if err := json.Unmarshal(req.Arguments, &in); err != nil {
			return mcp.NewToolResultError("falha ao decodificar argumentos"), nil
		}

		gvr := schema.GroupVersionResource{
			Group:    "apps.openshift.io",
			Version:  "v1",
			Resource: "deploymentconfigs",
		}

		ri := ctx.DynClient.Resource(gvr).Namespace(in.Namespace)

		dc, err := ri.Get(c, in.Name, metav1.GetOptions{})
		if err != nil {
			return mcp.NewToolResultErrorFromErr("erro ao obter DeploymentConfig", err), nil
		}

		latest, found, _ := unstructured.NestedInt64(dc.Object, "status", "latestVersion")
		if !found {
			latest = 0
		}
		newLatest := latest + 1

		patch := []byte(fmt.Sprintf(`{"status":{"latestVersion":%d}}`, newLatest))

		updated, err := ri.Patch(c, in.Name, types.MergePatchType, patch, metav1.PatchOptions{})
		if err != nil {
			return mcp.NewToolResultErrorFromErr("erro ao acionar rollout", err), nil
		}

		b, _ := updated.MarshalJSON()
		return mcpserver.TextResult(string(b)), nil
	})
}

// ---------- ImageStream actions ----------

type ImageStreamPromoteTagInput struct {
	Namespace    string `json:"namespace"`
	ImageStream  string `json:"imageStream"`
	SourceTag    string `json:"sourceTag"`    // ex: "dev"
	TargetTag    string `json:"targetTag"`    // ex: "prod"
	TargetIsCopy bool   `json:"targetIsCopy"` // se true, copia a imagem atual; se false, referencia a mesma fonte
}

func registerImageStreamPromoteTagTool(s *server.MCPServer, ctx *mcpserver.ServerContext) {
	tool := mcp.NewTool(
		"imagestream_promote_tag",
		mcp.WithDescription("Promove uma tag de ImageStream (ex: 'app:dev' -> 'app:prod') usando ImageStreamTag."),
		mcp.WithString("namespace", mcp.Required(), mcp.Description("Namespace do ImageStream.")),
		mcp.WithString("imageStream", mcp.Required(), mcp.Description("Nome do ImageStream.")),
		mcp.WithString("sourceTag", mcp.Required(), mcp.Description("Tag de origem (ex: dev).")),
		mcp.WithString("targetTag", mcp.Required(), mcp.Description("Tag de destino (ex: prod).")),
		mcp.WithBool("targetIsCopy", mcp.Description("Se true, copia a imagem; se false, referencia a mesma origem.")),
	)

	s.AddTool(tool, func(c context.Context, req mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		var in ImageStreamPromoteTagInput
		if err := json.Unmarshal(req.Arguments, &in); err != nil {
			return mcp.NewToolResultError("falha ao decodificar argumentos"), nil
		}

		istGVR := schema.GroupVersionResource{
			Group:    "image.openshift.io",
			Version:  "v1",
			Resource: "imagestreamtags",
		}

		srcName := fmt.Sprintf("%s:%s", in.ImageStream, in.SourceTag)
		dstName := fmt.Sprintf("%s:%s", in.ImageStream, in.TargetTag)

		ri := ctx.DynClient.Resource(istGVR).Namespace(in.Namespace)

		// 1) lê o ImageStreamTag de origem
		src, err := ri.Get(c, srcName, metav1.GetOptions{})
		if err != nil {
			return mcp.NewToolResultErrorFromErr("erro ao obter ImageStreamTag de origem", err), nil
		}

		// 2) prepara o objeto de destino
		obj := map[string]interface{}{
			"apiVersion": "image.openshift.io/v1",
			"kind":       "ImageStreamTag",
			"metadata": map[string]interface{}{
				"name":      dstName,
				"namespace": in.Namespace,
			},
		}

		if in.TargetIsCopy {
			image, _, _ := unstructured.NestedString(src.Object, "image", "dockerImageReference")
			if image == "" {
				return mcp.NewToolResultError("não foi possível determinar dockerImageReference de origem"), nil
			}

			obj["tag"] = map[string]interface{}{
				"name": in.TargetTag,
				"from": map[string]interface{}{
					"kind": "DockerImage",
					"name": image,
				},
			}
		} else {
			obj["tag"] = map[string]interface{}{
				"name": in.TargetTag,
				"from": map[string]interface{}{
					"kind": "ImageStreamTag",
					"name": srcName,
				},
			}
		}

		// 3) cria ou substitui o ImageStreamTag de destino
		dst, err := ri.Get(c, dstName, metav1.GetOptions{})
		if err == nil {
			obj["metadata"].(map[string]interface{})["resourceVersion"] = dst.GetResourceVersion()
			updated, err := ri.Update(c, &unstructured.Unstructured{Object: obj}, metav1.UpdateOptions{})
			if err != nil {
				return mcp.NewToolResultErrorFromErr("erro ao atualizar ImageStreamTag de destino", err), nil
			}
			b, _ := updated.MarshalJSON()
			return mcpserver.TextResult(string(b)), nil
		}

		created, err := ri.Create(c, &unstructured.Unstructured{Object: obj}, metav1.CreateOptions{})
		if err != nil {
			return mcp.NewToolResultErrorFromErr("erro ao criar ImageStreamTag de destino", err), nil
		}

		b, _ := created.MarshalJSON()
		return mcpserver.TextResult(string(b)), nil
	})
}
