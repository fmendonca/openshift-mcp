// internal/resources/namespace/register.go
package namespace

import (
	"github.com/fmendonca/openshift-mcp/internal/clients"
	"github.com/mark3labs/mcp-go/mcp"
	"github.com/mark3labs/mcp-go/server"
)

func RegisterResources(srv *server.MCPServer, c *clients.Clients) {
	// namespaces://all
	nsList := mcp.NewResource(
		"namespaces://all",
		"Namespaces List",
		mcp.WithResourceDescription("List all namespaces in the cluster"),
		mcp.WithMIMEType("application/json"),
	)
	srv.AddResource(nsList, newNamespacesListHandler(c))

	// namespaces://detail?name=<ns>
	nsDetail := mcp.NewResource(
		"namespaces://detail",
		"Namespace Detail",
		mcp.WithResourceDescription("Detailed information for a single namespace (name query param)"),
		mcp.WithMIMEType("application/json"),
	)
	srv.AddResource(nsDetail, newNamespaceDetailHandler(c))
}
