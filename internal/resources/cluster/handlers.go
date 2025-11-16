package cluster

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/fmendonca/openshift-mcp/internal/clients"
	"github.com/mark3labs/mcp-go/mcp"
	"github.com/mark3labs/mcp-go/server"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

func newClusterInfoHandler(c *clients.Clients) server.ResourceHandlerFunc {
	return func(ctx context.Context, uri string) (*mcp.ResourceResponse, error) {
		sv, err := c.Kubernetes.Discovery().ServerVersion()
		if err != nil {
			return nil, fmt.Errorf("failed to get server version: %w", err)
		}

		groups, err := c.Kubernetes.Discovery().ServerGroups()
		if err != nil {
			return nil, fmt.Errorf("failed to get server groups: %w", err)
		}

		info := map[string]any{
			"gitVersion":    sv.GitVersion,
			"major":         sv.Major,
			"minor":         sv.Minor,
			"platform":      sv.Platform,
			"buildDate":     sv.BuildDate,
			"compiler":      sv.Compiler,
			"goVersion":     sv.GoVersion,
			"apiGroupCount": len(groups.Groups),
		}

		data, err := json.MarshalIndent(info, "", "  ")
		if err != nil {
			return nil, err
		}

		return mcp.NewResourceResponseBytes("application/json", data), nil
	}
}

func newOpenShiftVersionHandler(c *clients.Clients) server.ResourceHandlerFunc {
	return func(ctx context.Context, uri string) (*mcp.ResourceResponse, error) {
		// Para clusters Kubernetes puros, apenas retorna vazio.
		if c.Config == nil {
			empty := map[string]any{
				"openshift": false,
				"message":   "OpenShift config client not available",
			}
			data, _ := json.MarshalIndent(empty, "", "  ")
			return mcp.NewResourceResponseBytes("application/json", data), nil
		}

		cv, err := c.Config.ConfigV1().ClusterVersions().Get(ctx, "version", metav1.GetOptions{})
		if err != nil {
			info := map[string]any{
				"openshift": true,
				"error":     err.Error(),
			}
			data, _ := json.MarshalIndent(info, "", "  ")
			return mcp.NewResourceResponseBytes("application/json", data), nil
		}

		channel := cv.Spec.Channel
		desired := cv.Status.Desired.Version

		hist := []map[string]any{}
		for _, h := range cv.Status.History {
			hist = append(hist, map[string]any{
				"version":  h.Version,
				"state":    string(h.State),
				"started":  h.StartedTime,
				"verified": h.Verified,
			})
		}

		info := map[string]any{
			"openshift":      true,
			"desiredVersion": desired,
			"channel":        channel,
			"history":        hist,
		}

		data, err := json.MarshalIndent(info, "", "  ")
		if err != nil {
			return nil, err
		}

		return mcp.NewResourceResponseBytes("application/json", data), nil
	}
}

func newAPIGroupsHandler(c *clients.Clients) server.ResourceHandlerFunc {
	return func(ctx context.Context, uri string) (*mcp.ResourceResponse, error) {
		groups, err := c.Kubernetes.Discovery().ServerGroups()
		if err != nil {
			return nil, fmt.Errorf("failed to get API groups: %w", err)
		}

		type groupInfo struct {
			Name      string   `json:"name"`
			Versions  []string `json:"versions"`
			Preferred string   `json:"preferredVersion,omitempty"`
		}

		list := []groupInfo{}
		for _, g := range groups.Groups {
			gi := groupInfo{Name: g.Name}
			for _, v := range g.Versions {
				gi.Versions = append(gi.Versions, v.Version)
			}
			if g.PreferredVersion.Version != "" {
				gi.Preferred = g.PreferredVersion.Version
			}
			list = append(list, gi)
		}

		data, err := json.MarshalIndent(list, "", "  ")
		if err != nil {
			return nil, err
		}

		return mcp.NewResourceResponseBytes("application/json", data), nil
	}
}
