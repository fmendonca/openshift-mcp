package contextx

import (
	"github.com/mark3labs/mcp-go/server"
	"k8s.io/client-go/dynamic"
	"k8s.io/client-go/kubernetes"
)

// ServerContext agrega os clientes usados por todas as tools.
type ServerContext struct {
	KubeClient *kubernetes.Clientset
	DynClient  dynamic.Interface
}

// Registrar é a assinatura genérica para funções que registram tools.
type Registrar func(s *server.MCPServer, ctx *ServerContext)
