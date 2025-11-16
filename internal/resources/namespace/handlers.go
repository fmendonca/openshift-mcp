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
	return func(ctx context.Context, uri string) (*mcp.ResourceResponse, error) {
		nsList, err := c.Kubernetes.CoreV1().Namespaces().List(ctx, metav1.ListOptions{})
		if err != nil {
			return nil, fmt.Errorf("failed to list namespaces: %w", err)
		}

		type nsInfo struct {
			Name   string            `json:"name"`
			Phase  string            `json:"phase"`
			Labels map[string]string `json:"labels,omitempty"`
		}

		items := []nsInfo{}
		for _, ns := range nsList.Items {
			items = append(items, nsInfo{
				Name:   ns.Name,
				Phase:  string(ns.Status.Phase),
				Labels: ns.Labels,
			})
		}

		data, err := json.MarshalIndent(items, "", "  ")
		if err != nil {
			return nil, err
		}

		return mcp.NewResourceResponseBytes("application/json", data), nil
	}
}

// parseQuery extrai querystring de uma URI resource MCP simples
func parseQuery(raw string) url.Values {
	// exemplo uri: namespaces://detail?name=my-namespace
	parts := strings.SplitN(raw, "?", 2)
	if len(parts) != 2 {
		return url.Values{}
	}
	v, _ := url.ParseQuery(parts[1])
	return v
}

func newNamespaceDetailHandler(c *clients.Clients) server.ResourceHandlerFunc {
	return func(ctx context.Context, uri string) (*mcp.ResourceResponse, error) {
		q := parseQuery(uri)
		name := q.Get("name")
		if name == "" {
			// sem nome, retorna erro "lógico" via conteúdo
			errBody := map[string]any{
				"error":  "missing namespace name",
				"usage":  "namespaces://detail?name=<namespace>",
				"status": "bad_request",
			}
			data, _ := json.MarshalIndent(errBody, "", "  ")
			return mcp.NewResourceResponseBytes("application/json", data), nil
		}

		ns, err := c.Kubernetes.CoreV1().Namespaces().Get(ctx, name, metav1.GetOptions{})
		if err != nil {
			errBody := map[string]any{
				"error":  err.Error(),
				"status": "not_found",
				"name":   name,
			}
			data, _ := json.MarshalIndent(errBody, "", "  ")
			return mcp.NewResourceResponseBytes("application/json", data), nil
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

		return mcp.NewResourceResponseBytes("application/json", data), nil
	}
}
