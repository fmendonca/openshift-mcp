package cluster

import (
	"github.com/fmendonca/openshift-mcp/internal/clients"
	mcpserver "github.com/fmendonca/openshift-mcp/internal/server"
	"github.com/mark3labs/mcp-go/mcp"
)

func RegisterResources(srv *mcpserver.MCPServer, c *clients.Clients) {
	// Informações gerais do cluster (versão, API groups, etc.)
	srv.AddResource(&mcp.Resource{
		URI:         "cluster://info",
		Name:        "Cluster Information",
		Description: "Kubernetes/OpenShift cluster version and basic info",
		MimeType:    "application/json",
	}, newClusterInfoHandler(c))

	// Versão do OpenShift (se disponível)
	srv.AddResource(&mcp.Resource{
		URI:         "cluster://openshift/version",
		Name:        "OpenShift Version",
		Description: "OpenShift cluster version and channel (if available)",
		MimeType:    "application/json",
	}, newOpenShiftVersionHandler(c))

	// API groups disponíveis
	srv.AddResource(&mcp.Resource{
		URI:         "cluster://apigroups",
		Name:        "API Groups",
		Description: "List of available API groups in the cluster",
		MimeType:    "application/json",
	}, newAPIGroupsHandler(c))
}
