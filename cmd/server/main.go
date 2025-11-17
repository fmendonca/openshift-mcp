// cmd/server/main.go
package main

import (
	"context"
	"log"
	"os"
	"os/signal"
	"syscall"

	"github.com/fmendonca/openshift-mcp/internal/clients"
	mcpserver "github.com/fmendonca/openshift-mcp/internal/server"

	// Handlers unificados (todos os tools em um único arquivo)
	"github.com/fmendonca/openshift-mcp/internal/handlers"
	// Resources (cluster e namespaces)
)

func main() {
	// Contexto com cancel para shutdown gracioso
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	// Captura SIGINT/SIGTERM para encerrar o servidor de forma limpa
	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, os.Interrupt, syscall.SIGTERM)
	go func() {
		sig := <-sigChan
		log.Printf("Received signal %s, shutting down gracefully...\n", sig.String())
		cancel()
	}()

	// Inicializa clients Kubernetes (in-cluster ou kubeconfig local)
	k8sClients, err := clients.NewClients()
	if err != nil {
		log.Fatalf("Failed to initialize Kubernetes clients: %v", err)
	}

	// Cria MCP server com suporte a tools e resources
	srv := mcpserver.NewServer("openshift-mcp", "1.0.0")

	// Registra TODOS os tools em um único lugar (pods, services, etc.)
	handlers.RegisterAllTools(srv.Inner(), k8sClients)

	// Registra resources (cluster://..., namespaces://...)
	//clusterres.RegisterResources(srv, k8sClients)
	//namespaceres.RegisterResources(srv, k8sClients)

	log.Println("Starting OpenShift/Kubernetes MCP server over stdio...")

	// Inicia o servidor usando transporte stdio (Claude, VS Code, etc.)
	if err := srv.Start(ctx); err != nil {
		log.Fatalf("MCP server error: %v", err)
	}

	log.Println("MCP server stopped.")
}
