// internal/resources/cluster/register.go
package cluster

import (
	"github.com/fmendonca/openshift-mcp/internal/clients"
	"github.com/mark3labs/mcp-go/mcp"
	"github.com/mark3labs/mcp-go/server"
)

func RegisterResources(srv *server.MCPServer, c *clients.Clients) {
	// cluster://info
	infoRes := mcp.NewResource(
		"cluster://info",
		"Cluster Information",
		mcp.WithResourceDescription("Kubernetes/OpenShift cluster version and basic info"),
		mcp.WithMIMEType("application/json"),
	)
	srv.AddResource(infoRes, newClusterInfoHandler(c))

	// cluster://openshift/version
	verRes := mcp.NewResource(
		"cluster://openshift/version",
		"OpenShift Version",
		mcp.WithResourceDescription("OpenShift cluster version and channel (if available)"),
		mcp.WithMIMEType("application/json"),
	)
	srv.AddResource(verRes, newOpenShiftVersionHandler(c))

	// cluster://apigroups
	groupsRes := mcp.NewResource(
		"cluster://apigroups",
		"API Groups",
		mcp.WithResourceDescription("List of available API groups in the cluster"),
		mcp.WithMIMEType("application/json"),
	)
	srv.AddResource(groupsRes, newAPIGroupsHandler(c))
}
