package openshift

import (
	"context"
	"encoding/json"
	"fmt"

	contextx "github.com/fmendonca/openshfit-mcp/internal/context"
	"github.com/mark3labs/mcp-go/mcp"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"
	"k8s.io/apimachinery/pkg/runtime/schema"
	"k8s.io/apimachinery/pkg/types"
	"k8s.io/client-go/dynamic"
)

// Registra todas as tools específicas de OpenShift
func RegisterOpenShiftTools(reg *mcp.ToolRegistry, ctx *contextx.ServerContext) {
	registerRoutesListTool(reg, ctx)
	registerBuildConfigsListTool(reg, ctx)
	registerImageStreamsListTool(reg, ctx)
	registerProjectsListTool(reg, ctx)
	registerDeploymentConfigsListTool(reg, ctx)

	// ações
	registerBuildConfigStartBuildTool(reg, ctx)
	registerDeploymentConfigRolloutTool(reg, ctx)
	registerImageStreamPromoteTagTool(reg, ctx)
}

// ---------- Routes ----------

type RoutesListInput struct {
	Namespace string `json:"namespace,omitempty"`
}

func registerRoutesListTool(reg *mcp.ToolRegistry, ctx *contextx.ServerContext) {
	reg.RegisterTool(&mcp.Tool{
		Name:        "routes_list",
		Description: "Lista Routes do OpenShift (route.openshift.io/v1, resource 'routes').",
		InputSchema: &mcp.JSONSchema{
			Type: "object",
			Properties: map[string]*mcp.JSONSchema{
				"namespace": {Type: "string"},
			},
		},
	}, func(c context.Context, req *mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		var in RoutesListInput
		if err := json.Unmarshal(req.Params, &in); err != nil {
			return mcp.NewErrorToolResult("falha ao decodificar argumentos", err), nil
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
			return mcp.NewErrorToolResult("erro ao listar routes", err), nil
		}

		b, _ := routes.MarshalJSON()
		return mcp.NewToolResultText(string(b)), nil
	})
}

// ---------- BuildConfig ----------

type BuildConfigsListInput struct {
	Namespace string `json:"namespace,omitempty"`
}

func registerBuildConfigsListTool(reg *mcp.ToolRegistry, ctx *contextx.ServerContext) {
	reg.RegisterTool(&mcp.Tool{
		Name:        "buildconfigs_list",
		Description: "Lista BuildConfigs (build.openshift.io/v1, resource 'buildconfigs').",
		InputSchema: &mcp.JSONSchema{
			Type: "object",
			Properties: map[string]*mcp.JSONSchema{
				"namespace": {Type: "string"},
			},
		},
	}, func(c context.Context, req *mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		var in BuildConfigsListInput
		if err := json.Unmarshal(req.Params, &in); err != nil {
			return mcp.NewErrorToolResult("falha ao decodificar argumentos", err), nil
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
			return mcp.NewErrorToolResult("erro ao listar BuildConfigs", err), nil
		}

		b, _ := bcs.MarshalJSON()
		return mcp.NewToolResultText(string(b)), nil
	})
}

// ---------- ImageStream ----------

type ImageStreamsListInput struct {
	Namespace string `json:"namespace,omitempty"`
}

func registerImageStreamsListTool(reg *mcp.ToolRegistry, ctx *contextx.ServerContext) {
	reg.RegisterTool(&mcp.Tool{
		Name:        "imagestreams_list",
		Description: "Lista ImageStreams (image.openshift.io/v1, resource 'imagestreams').",
		InputSchema: &mcp.JSONSchema{
			Type: "object",
			Properties: map[string]*mcp.JSONSchema{
				"namespace": {Type: "string"},
			},
		},
	}, func(c context.Context, req *mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		var in ImageStreamsListInput
		if err := json.Unmarshal(req.Params, &in); err != nil {
			return mcp.NewErrorToolResult("falha ao decodificar argumentos", err), nil
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
			return mcp.NewErrorToolResult("erro ao listar ImageStreams", err), nil
		}

		b, _ := iss.MarshalJSON()
		return mcp.NewToolResultText(string(b)), nil
	})
}

// ---------- Project ----------

type ProjectsListInput struct{}

func registerProjectsListTool(reg *mcp.ToolRegistry, ctx *contextx.ServerContext) {
	reg.RegisterTool(&mcp.Tool{
		Name:        "projects_list",
		Description: "Lista Projects (project.openshift.io/v1, resource 'projects').",
	}, func(c context.Context, req *mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		gvr := schema.GroupVersionResource{
			Group:    "project.openshift.io",
			Version:  "v1",
			Resource: "projects",
		}

		ri := ctx.DynClient.Resource(gvr)

		projects, err := ri.List(c, metav1.ListOptions{})
		if err != nil {
			return mcp.NewErrorToolResult("erro ao listar Projects", err), nil
		}

		b, _ := projects.MarshalJSON()
		return mcp.NewToolResultText(string(b)), nil
	})
}

// ---------- DeploymentConfig ----------

type DeploymentConfigsListInput struct {
	Namespace string `json:"namespace,omitempty"`
}

func registerDeploymentConfigsListTool(reg *mcp.ToolRegistry, ctx *contextx.ServerContext) {
	reg.RegisterTool(&mcp.Tool{
		Name:        "deploymentconfigs_list",
		Description: "Lista DeploymentConfigs (apps.openshift.io/v1, resource 'deploymentconfigs').",
		InputSchema: &mcp.JSONSchema{
			Type: "object",
			Properties: map[string]*mcp.JSONSchema{
				"namespace": {Type: "string"},
			},
		},
	}, func(c context.Context, req *mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		var in DeploymentConfigsListInput
		if err := json.Unmarshal(req.Params, &in); err != nil {
			return mcp.NewErrorToolResult("falha ao decodificar argumentos", err), nil
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
			return mcp.NewErrorToolResult("erro ao listar DeploymentConfigs", err), nil
		}

		b, _ := dcs.MarshalJSON()
		return mcp.NewToolResultText(string(b)), nil
	})
}

// ---------- BuildConfig actions ----------

type BuildConfigStartBuildInput struct {
	Namespace string `json:"namespace"`
	Name      string `json:"name"`
}

func registerBuildConfigStartBuildTool(reg *mcp.ToolRegistry, ctx *contextx.ServerContext) {
	reg.RegisterTool(&mcp.Tool{
		Name:        "buildconfig_start_build",
		Description: "Dispara um build a partir de um BuildConfig (equivalente a 'oc start-build').",
		InputSchema: &mcp.JSONSchema{
			Type: "object",
			Properties: map[string]*mcp.JSONSchema{
				"namespace": {Type: "string"},
				"name":      {Type: "string"},
			},
			Required: []string{"namespace", "name"},
		},
	}, func(c context.Context, req *mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		var in BuildConfigStartBuildInput
		if err := json.Unmarshal(req.Params, &in); err != nil {
			return mcp.NewErrorToolResult("falha ao decodificar argumentos", err), nil
		}

		bcGVR := schema.GroupVersionResource{
			Group:    "build.openshift.io",
			Version:  "v1",
			Resource: "buildconfigs",
		}

		bc, err := ctx.DynClient.Resource(bcGVR).Namespace(in.Namespace).Get(c, in.Name, metav1.GetOptions{})
		if err != nil {
			return mcp.NewErrorToolResult("erro ao obter BuildConfig", err), nil
		}

		spec, ok := bc.Object["spec"].(map[string]interface{})
		if !ok {
			return mcp.NewErrorToolResult("BuildConfig sem spec válido", nil), nil
		}

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
			return mcp.NewErrorToolResult("erro ao criar Build", err), nil
		}

		b, _ := created.MarshalJSON()
		return mcp.NewToolResultText(string(b)), nil
	})
}

// ---------- DeploymentConfig actions ----------

type DeploymentConfigRolloutInput struct {
	Namespace string `json:"namespace"`
	Name      string `json:"name"`
}

func registerDeploymentConfigRolloutTool(reg *mcp.ToolRegistry, ctx *contextx.ServerContext) {
	reg.RegisterTool(&mcp.Tool{
		Name:        "deploymentconfig_rollout_latest",
		Description: "Dispara um rollout manual da DeploymentConfig (equivalente a 'oc rollout latest').",
		InputSchema: &mcp.JSONSchema{
			Type: "object",
			Properties: map[string]*mcp.JSONSchema{
				"namespace": {Type: "string"},
				"name":      {Type: "string"},
			},
			Required: []string{"namespace", "name"},
		},
	}, func(c context.Context, req *mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		var in DeploymentConfigRolloutInput
		if err := json.Unmarshal(req.Params, &in); err != nil {
			return mcp.NewErrorToolResult("falha ao decodificar argumentos", err), nil
		}

		gvr := schema.GroupVersionResource{
			Group:    "apps.openshift.io",
			Version:  "v1",
			Resource: "deploymentconfigs",
		}

		ri := ctx.DynClient.Resource(gvr).Namespace(in.Namespace)

		dc, err := ri.Get(c, in.Name, metav1.GetOptions{})
		if err != nil {
			return mcp.NewErrorToolResult("erro ao obter DeploymentConfig", err), nil
		}

		latest, found, _ := unstructured.NestedInt64(dc.Object, "status", "latestVersion")
		if !found {
			latest = 0
		}
		newLatest := latest + 1

		patch := []byte(fmt.Sprintf(`{"status":{"latestVersion":%d}}`, newLatest))

		updated, err := ri.Patch(c, in.Name, types.MergePatchType, patch, metav1.PatchOptions{})
		if err != nil {
			return mcp.NewErrorToolResult("erro ao acionar rollout", err), nil
		}

		b, _ := updated.MarshalJSON()
		return mcp.NewToolResultText(string(b)), nil
	})
}

// ---------- ImageStream actions ----------

type ImageStreamPromoteTagInput struct {
	Namespace    string `json:"namespace"`
	ImageStream  string `json:"imageStream"`
	SourceTag    string `json:"sourceTag"`
	TargetTag    string `json:"targetTag"`
	TargetIsCopy bool   `json:"targetIsCopy"`
}

func registerImageStreamPromoteTagTool(reg *mcp.ToolRegistry, ctx *contextx.ServerContext) {
	reg.RegisterTool(&mcp.Tool{
		Name:        "imagestream_promote_tag",
		Description: "Promove uma tag de ImageStream (ex: 'app:dev' -> 'app:prod') usando ImageStreamTag.",
		InputSchema: &mcp.JSONSchema{
			Type: "object",
			Properties: map[string]*mcp.JSONSchema{
				"namespace":    {Type: "string"},
				"imageStream":  {Type: "string"},
				"sourceTag":    {Type: "string"},
				"targetTag":    {Type: "string"},
				"targetIsCopy": {Type: "boolean"},
			},
			Required: []string{"namespace", "imageStream", "sourceTag", "targetTag"},
		},
	}, func(c context.Context, req *mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		var in ImageStreamPromoteTagInput
		if err := json.Unmarshal(req.Params, &in); err != nil {
			return mcp.NewErrorToolResult("falha ao decodificar argumentos", err), nil
		}

		istGVR := schema.GroupVersionResource{
			Group:    "image.openshift.io",
			Version:  "v1",
			Resource: "imagestreamtags",
		}

		srcName := fmt.Sprintf("%s:%s", in.ImageStream, in.SourceTag)
		dstName := fmt.Sprintf("%s:%s", in.ImageStream, in.TargetTag)

		ri := ctx.DynClient.Resource(istGVR).Namespace(in.Namespace)

		src, err := ri.Get(c, srcName, metav1.GetOptions{})
		if err != nil {
			return mcp.NewErrorToolResult("erro ao obter ImageStreamTag de origem", err), nil
		}

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
				return mcp.NewErrorToolResult("não foi possível determinar dockerImageReference de origem", nil), nil
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

		dst, err := ri.Get(c, dstName, metav1.GetOptions{})
		if err == nil {
			obj["metadata"].(map[string]interface{})["resourceVersion"] = dst.GetResourceVersion()
			updated, err := ri.Update(c, &unstructured.Unstructured{Object: obj}, metav1.UpdateOptions{})
			if err != nil {
				return mcp.NewErrorToolResult("erro ao atualizar ImageStreamTag de destino", err), nil
			}
			b, _ := updated.MarshalJSON()
			return mcp.NewToolResultText(string(b)), nil
		}

		created, err := ri.Create(c, &unstructured.Unstructured{Object: obj}, metav1.CreateOptions{})
		if err != nil {
			return mcp.NewErrorToolResult("erro ao criar ImageStreamTag de destino", err), nil
		}

		b, _ := created.MarshalJSON()
		return mcp.NewToolResultText(string(b)), nil
	})
}
