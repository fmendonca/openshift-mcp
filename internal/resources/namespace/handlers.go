package namespace

import (
	"context"
	"encoding/json"
	"fmt"
	"net/url"
	"strings"

	"github.com/fmendonca/openshift-mcp/internal/clients"
	"github.com/mark3labs/mcp-go/mcp"
	"github.com/mark3labs/mcp-go/server"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

func newNamespacesListHandler(c *clients.Clients) server.ResourceHandlerFunc {
	return func(ctx context.Context, req mcp.ReadResourceRequest) (*mcp.ReadResourceResult, error) {
		nsList, err := c.Kubernetes.CoreV1().Namespaces().List(ctx, metav1.ListOptions{})
		if err != nil {
			return nil, fmt.Errorf("failed to list namespaces: %w", err)
		}

		type nsInfo struct {
			Name   string            `json:"name"`
			Phase  string            `json:"phase"`
			Labels map[string]string `json:"labels,omitempty"`
		}

		out := make([]nsInfo, 0, len(nsList.Items))
		for _, ns := range nsList.Items {
			out = append(out, nsInfo{
				Name:   ns.Name,
				Phase:  string(ns.Status.Phase),
				Labels: ns.Labels,
			})
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

func parseQuery(raw string) url.Values {
	parts := strings.SplitN(raw, "?", 2)
	if len(parts) != 2 {
		return url.Values{}
	}
	v, _ := url.ParseQuery(parts[1])
	return v
}

func newNamespaceDetailHandler(c *clients.Clients) server.ResourceHandlerFunc {
	return func(ctx context.Context, req mcp.ReadResourceRequest) (*mcp.ReadResourceResult, error) {
		q := parseQuery(req.Params.URI)
		name := q.Get("name")
		if name == "" {
			body := map[string]any{
				"error":  "missing namespace name",
				"usage":  "namespaces://detail?name=<namespace>",
				"status": "bad_request",
			}
			data, _ := json.MarshalIndent(body, "", "  ")
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

		ns, err := c.Kubernetes.CoreV1().Namespaces().Get(ctx, name, metav1.GetOptions{})
		if err != nil {
			body := map[string]any{
				"error":  err.Error(),
				"status": "not_found",
				"name":   name,
			}
			data, _ := json.MarshalIndent(body, "", "  ")
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

		info := map[string]any{
			"name":        ns.Name,
			"phase":       ns.Status.Phase,
			"labels":      ns.Labels,
			"annotations": ns.Annotations,
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
