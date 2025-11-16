package mcpserver

import (
	contextx "github.com/fmendonca/openshfit-mcp/internal/context"
	"github.com/fmendonca/openshfit-mcp/internal/kubevirt"
	"github.com/fmendonca/openshfit-mcp/internal/openshift"
	"github.com/mark3labs/mcp-go/server"
	"k8s.io/client-go/dynamic"
	"k8s.io/client-go/kubernetes"
)

// NewServer cria o MCP server e registra todas as tools.
func NewServer(name, version string, kube *kubernetes.Clientset, dyn dynamic.Interface) *server.MCPServer {
	s := server.NewMCPServer(
		name,
		version,
	)

	ctx := &contextx.ServerContext{
		KubeClient: kube,
		DynClient:  dyn,
	}

	RegisterCoreTools(s, ctx)
	openshift.RegisterOpenShiftTools(s, ctx)
	kubevirt.RegisterKubeVirtTools(s, ctx)

	return s
}
