package main

import (
	"log"
	"os"

	"github.com/fmendonca/openshfit-mcp/internal/k8s"
	"github.com/fmendonca/openshfit-mcp/internal/mcpserver"
	"github.com/mark3labs/mcp-go/server"
)

func main() {
	cfg := k8s.Config{
		Kubeconfig: os.Getenv("KUBECONFIG"),
		InCluster:  os.Getenv("IN_CLUSTER") == "true",
	}

	restCfg, err := k8s.BuildRestConfig(cfg)
	if err != nil {
		log.Fatalf("erro ao criar rest.Config: %v", err)
	}

	kubeClient, dynClient, err := k8s.NewClients(restCfg)
	if err != nil {
		log.Fatalf("erro ao criar clients kubernetes/dynamic: %v", err)
	}

	s := mcpserver.NewServer("openshift-kubevirt-mcp-go", "0.1.0", kubeClient, dynClient)

	if err := server.ServeStdio(s); err != nil {
		log.Fatalf("erro ao iniciar MCP server: %v", err)
	}
}
