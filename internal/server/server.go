package server

import (
	"context"
	"log"

	"github.com/mark3labs/mcp-go/mcp"
	"github.com/mark3labs/mcp-go/server"
)

type MCPServer struct {
	server *server.MCPServer
}

func NewServer(name, version string) *MCPServer {
	mcpServer := server.NewMCPServer(
		name,
		version,
		server.WithToolCapabilities(true),
		server.WithResourceCapabilities(true, false),
	)

	return &MCPServer{
		server: mcpServer,
	}
}

func (s *MCPServer) AddTool(tool *mcp.Tool, handler server.ToolHandlerFunc) {
	s.server.AddTool(tool, handler)
}

func (s *MCPServer) AddResource(resource *mcp.Resource, handler server.ResourceHandlerFunc) {
	s.server.AddResource(resource, handler)
}

func (s *MCPServer) Start(ctx context.Context) error {
	log.Println("MCP Server is ready")
	return s.server.Serve(ctx, server.StdioTransport())
}
