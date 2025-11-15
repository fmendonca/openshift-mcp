package mcpserver

import (
	"github.com/mark3labs/mcp-go/mcp"
	"github.com/mark3labs/mcp-go/server"
	"k8s.io/client-go/dynamic"
	"k8s.io/client-go/kubernetes"
)

type ServerContext struct {
	KubeClient *kubernetes.Clientset
	DynClient  dynamic.Interface
}

func NewServer(name, version string, kube *kubernetes.Clientset, dyn dynamic.Interface) *server.MCPServer {
	s := server.NewMCPServer(
		name,
		version,
		server.WithToolCapabilities(false),
	)

	ctx := &ServerContext{
		KubeClient: kube,
		DynClient:  dyn,
	}

	RegisterCoreTools(s, ctx)
	RegisterOpenShiftTools(s, ctx)
	RegisterKubeVirtTools(s, ctx)

	// opcional: registrar prompts, resources, etc.[web:35][web:74]

	return s
}

// Helper para criar um resultado de texto
func TextResult(text string) *mcp.CallToolResult {
	return mcp.NewToolResultText(text)
}
