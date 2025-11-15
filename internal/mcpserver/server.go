package mcpserver

import (
	contextx "github.com/fmendonca/openshfit-mcp/internal/context"
	"github.com/fmendonca/openshfit-mcp/internal/kubevirt"
	"github.com/fmendonca/openshfit-mcp/internal/openshift"
	"github.com/mark3labs/mcp-go/mcp"
	"github.com/mark3labs/mcp-go/server"
	"k8s.io/client-go/dynamic"
	"k8s.io/client-go/kubernetes"
)

func NewServer(name, version string, kube *kubernetes.Clientset, dyn dynamic.Interface) *server.MCPServer {
	ctx := &contextx.ServerContext{
		KubeClient: kube,
		DynClient:  dyn,
	}

	s := server.NewMCPServer(
		name,
		version,
		server.WithTools(func(reg *mcp.ToolRegistry) {
			// tools genéricas
			RegisterCoreTools(reg, ctx)
			// OpenShift
			openshift.RegisterOpenShiftTools(reg, ctx)
			// KubeVirt
			kubevirt.RegisterKubeVirtTools(reg, ctx)
		}),
	)

	return s
}
