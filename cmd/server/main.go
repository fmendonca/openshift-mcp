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

	// Tools (Kubernetes)
	"github.com/fmendonca/openshift-mcp/internal/tools/configmaps"
	"github.com/fmendonca/openshift-mcp/internal/tools/deployments"
	"github.com/fmendonca/openshift-mcp/internal/tools/imagestreams"
	"github.com/fmendonca/openshift-mcp/internal/tools/ingress"
	"github.com/fmendonca/openshift-mcp/internal/tools/nodes"
	"github.com/fmendonca/openshift-mcp/internal/tools/pods"
	"github.com/fmendonca/openshift-mcp/internal/tools/projects"
	"github.com/fmendonca/openshift-mcp/internal/tools/pvcs"
	"github.com/fmendonca/openshift-mcp/internal/tools/rbac"
	"github.com/fmendonca/openshift-mcp/internal/tools/routes"
	"github.com/fmendonca/openshift-mcp/internal/tools/secrets"
	"github.com/fmendonca/openshift-mcp/internal/tools/services"

	// Resources (MCP resources)
	clusterres "github.com/fmendonca/openshift-mcp/internal/resources/cluster"
	namespaceres "github.com/fmendonca/openshift-mcp/internal/resources/namespace"
)

func main() {
	// Context com cancelamento para desligar o servidor com SIGINT/SIGTERM
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	// Captura de sinal para shutdown gracioso
	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, os.Interrupt, syscall.SIGTERM)
	go func() {
		s := <-sigChan
		log.Printf("Received signal %s, shutting down gracefully...\n", s.String())
		cancel()
	}()

	// Inicializa os clients de Kubernetes/OpenShift
	k8sClients, err := clients.NewClients()
	if err != nil {
		log.Fatalf("Failed to initialize Kubernetes/OpenShift clients: %v", err)
	}

	// Cria o servidor MCP (STDIO)
	srv := mcpserver.NewServer("openshift-mcp", "1.0.0")

	// Registro de TOOLS (organizados por funcionalidade)

	// Pods: listar, descrever, logs, exec, delete
	pods.RegisterTools(srv, k8sClients)

	// Deployments: listar, descrever, scale, restart
	deployments.RegisterTools(srv, k8sClients)

	// Services: listar, descrever
	services.RegisterTools(srv, k8sClients)

	// Routes (OpenShift): listar, descrever
	routes.RegisterTools(srv, k8sClients)

	// ImageStreams (OpenShift): listar, descrever
	imagestreams.RegisterTools(srv, k8sClients)

	// Projects / Namespaces (OpenShift Projects API)
	projects.RegisterTools(srv, k8sClients)

	// Nodes: listar, descrever
	nodes.RegisterTools(srv, k8sClients)

	// ConfigMaps: listar, descrever
	configmaps.RegisterTools(srv, k8sClients)

	// Secrets: listar, descrever (dados mascarados)
	secrets.RegisterTools(srv, k8sClients)

	// PVCs: listar, descrever
	pvcs.RegisterTools(srv, k8sClients)

	// Ingress: listar, descrever
	ingress.RegisterTools(srv, k8sClients)

	// RBAC: roles, rolebindings, clusterroles, clusterrolebindings
	rbac.RegisterTools(srv, k8sClients)

	// Registro de RESOURCES (visões de leitura estruturadas)

	// Informações de cluster (versão K8s, versão OpenShift, API groups)
	clusterres.RegisterResources(srv, k8sClients)

	// Namespaces / Projects (lista e detalhe)
	namespaceres.RegisterResources(srv, k8sClients)

	// Inicializa o servidor MCP usando transporte STDIO
	log.Println("Starting OpenShift/Kubernetes MCP server over stdio...")
	if err := srv.Start(ctx); err != nil {
		// Erro ao servir (por exemplo, stdin/stdout fechados)
		log.Fatalf("MCP server error: %v", err)
	}

	log.Println("MCP server stopped.")
}
