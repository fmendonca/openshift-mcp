package cluster

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/fmendonca/openshift-mcp/internal/clients"
	"github.com/mark3labs/mcp-go/mcp"
	"github.com/mark3labs/mcp-go/server"
)

// newClusterInfoHandler: informações gerais do cluster Kubernetes.
func newClusterInfoHandler(c *clients.Clients) server.ResourceHandlerFunc {
	return func(ctx context.Context, req mcp.ReadResourceRequest) (*mcp.ReadResourceResult, error) {
		sv, err := c.Kubernetes.Discovery().ServerVersion()
		if err != nil {
			return nil, fmt.Errorf("failed to get server version: %w", err)
		}

		groups, err := c.Kubernetes.Discovery().ServerGroups()
		if err != nil {
			return nil, fmt.Errorf("failed to get API groups: %w", err)
		}

		info := map[string]any{
			"gitVersion":    sv.GitVersion,
			"major":         sv.Major,
			"minor":         sv.Minor,
			"platform":      sv.Platform,
			"goVersion":     sv.GoVersion,
			"compiler":      sv.Compiler,
			"apiGroupCount": len(groups.Groups),
		}

		data, err := json.MarshalIndent(info, "", "  ")
		if err != nil {
			return nil, err
		}

		return &mcp.ReadResourceResult{
			Contents: []mcp.ResourceContent{
				mcp.TextResourceContent{
					URI:      req.Params.URI,
					MIMEType: "application/json",
					Text:     string(data),
				},
			},
		}, nil
	}
}

// newOpenShiftVersionHandler: por enquanto devolve um stub (sem usar openshift/client-go).
func newOpenShiftVersionHandler(c *clients.Clients) server.ResourceHandlerFunc {
	return func(ctx context.Context, req mcp.ReadResourceRequest) (*mcp.ReadResourceResult, error) {
		info := map[string]any{
			"openshift": false,
			"message":   "OpenShift client not wired; only Kubernetes info is available.",
		}

		data, err := json.MarshalIndent(info, "", "  ")
		if err != nil {
			return nil, err
		}

		return &mcp.ReadResourceResult{
			Contents: []mcp.ResourceContent{
				mcp.TextResourceContent{
					URI:      req.Params.URI,
					MIMEType: "application/json",
					Text:     string(data),
				},
			},
		}, nil
	}
}

// newAPIGroupsHandler: lista todos os API groups disponíveis no cluster.
func newAPIGroupsHandler(c *clients.Clients) server.ResourceHandlerFunc {
	return func(ctx context.Context, req mcp.ReadResourceRequest) (*mcp.ReadResourceResult, error) {
		groups, err := c.Kubernetes.Discovery().ServerGroups()
		if err != nil {
			return nil, fmt.Errorf("failed to get API groups: %w", err)
		}

		type groupInfo struct {
			Name      string   `json:"name"`
			Versions  []string `json:"versions"`
			Preferred string   `json:"preferredVersion,omitempty"`
		}

		out := make([]groupInfo, 0, len(groups.Groups))
		for _, g := range groups.Groups {
			gi := groupInfo{Name: g.Name}
			for _, v := range g.Versions {
				gi.Versions = append(gi.Versions, v.Version)
			}
			if g.PreferredVersion.Version != "" {
				gi.Preferred = g.PreferredVersion.Version
			}
			out = append(out, gi)
		}

		data, err := json.MarshalIndent(out, "", "  ")
		if err != nil {
			return nil, err
		}

		return &mcp.ReadResourceResult{
			Contents: []mcp.ResourceContent{
				mcp.TextResourceContent{
					URI:      req.Params.URI,
					MIMEType: "application/json",
					Text:     string(data),
				},
			},
		}, nil
	}
}
