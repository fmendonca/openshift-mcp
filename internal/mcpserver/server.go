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
	s := server.NewMCPServer(
		name,
		version,
		server.WithToolCapabilities(false),
	)

	ctx := &contextx.ServerContext{
		KubeClient: kube,
		DynClient:  dyn,
	}

	// tools genéricas
	RegisterCoreTools(s, ctx)
	// toolsets específicos
	openshift.RegisterOpenShiftTools(s, ctx)
	kubevirt.RegisterKubeVirtTools(s, ctx)

	return s
}

func TextResult(text string) *mcp.CallToolResult {
	return mcp.NewToolResultText(text)
}
