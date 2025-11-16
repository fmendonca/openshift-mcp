package main

import (
	"context"
	"log"
	"os"
	"os/signal"
	"syscall"

	"github.com/yourusername/openshift-k8s-mcp/internal/clients"
	"github.com/yourusername/openshift-k8s-mcp/internal/server"
	"github.com/yourusername/openshift-k8s-mcp/internal/tools/configmaps"
	"github.com/yourusername/openshift-k8s-mcp/internal/tools/deployments"
	"github.com/yourusername/openshift-k8s-mcp/internal/tools/imagestreams"
	"github.com/yourusername/openshift-k8s-mcp/internal/tools/ingress"
	"github.com/yourusername/openshift-k8s-mcp/internal/tools/nodes"
	"github.com/yourusername/openshift-k8s-mcp/internal/tools/pods"
	"github.com/yourusername/openshift-k8s-mcp/internal/tools/projects"
	"github.com/yourusername/openshift-k8s-mcp/internal/tools/pvcs"
	"github.com/yourusername/openshift-k8s-mcp/internal/tools/rbac"
	"github.com/yourusername/openshift-k8s-mcp/internal/tools/routes"
	"github.com/yourusername/openshift-k8s-mcp/internal/tools/secrets"
	"github.com/yourusername/openshift-k8s-mcp/internal/tools/services"
)

func main() {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	// Handle graceful shutdown
	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, os.Interrupt, syscall.SIGTERM)
	go func() {
		<-sigChan
		log.Println("Shutting down gracefully...")
		cancel()
	}()

	// Initialize Kubernetes/OpenShift clients
	k8sClients, err := clients.NewClients()
	if err != nil {
		log.Fatalf("Failed to initialize clients: %v", err)
	}

	// Create MCP server
	mcpServer := server.NewServer("openshift-k8s-mcp", "1.0.0")

	// Register all tools
	pods.RegisterTools(mcpServer, k8sClients)
	deployments.RegisterTools(mcpServer, k8sClients)
	services.RegisterTools(mcpServer, k8sClients)
	routes.RegisterTools(mcpServer, k8sClients)
	imagestreams.RegisterTools(mcpServer, k8sClients)
	projects.RegisterTools(mcpServer, k8sClients)
	nodes.RegisterTools(mcpServer, k8sClients)
	configmaps.RegisterTools(mcpServer, k8sClients)
	secrets.RegisterTools(mcpServer, k8sClients)
	pvcs.RegisterTools(mcpServer, k8sClients)
	ingress.RegisterTools(mcpServer, k8sClients)
	rbac.RegisterTools(mcpServer, k8sClients)

	// Start server
	log.Println("Starting MCP server...")
	if err := mcpServer.Start(ctx); err != nil {
		log.Fatalf("Server error: %v", err)
	}
}
