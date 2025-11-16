package namespace
package namespace

import (
    "github.com/mark3labs/mcp-go/mcp"
    mcpserver "github.com/fmendonca/openshift-mcp/internal/server"
    "github.com/fmendonca/openshift-mcp/internal/clients"
)

func RegisterResources(srv *mcpserver.MCPServer, c *clients.Clients) {
    // Lista de namespaces / projetos
    srv.AddResource(&mcp.Resource{
        URI:         "namespaces://all",
        Name:        "Namespaces List",
        Description: "List all namespaces (and OpenShift projects if available)",
        MimeType:    "application/json",
    }, newNamespacesListHandler(c))

    // Detalhes de um namespace específico (via query na URI)
    // Ex: namespaces://detail?name=my-namespace
    srv.AddResource(&mcp.Resource{
        URI:         "namespaces://detail",
        Name:        "Namespace Detail",
        Description: "Detailed information for a single namespace",
        MimeType:    "application/json",
    }, newNamespaceDetailHandler(c))
}
