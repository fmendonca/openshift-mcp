package server

import (
	"context"
	"log"

	"github.com/mark3labs/mcp-go/mcp"
	mcpsrv "github.com/mark3labs/mcp-go/server"
)

type MCPServer struct {
	server *mcpsrv.MCPServer
}

func NewServer(name, version string) *MCPServer {
	s := mcpsrv.NewMCPServer(
		name,
		version,
		mcpsrv.WithToolCapabilities(true),
		mcpsrv.WithResourceCapabilities(true, false),
		mcpsrv.WithRecovery(),
	)

	return &MCPServer{server: s}
}

func (s *MCPServer) AddTool(tool *mcp.Tool, handler mcpsrv.ToolHandlerFunc) {
	s.server.AddTool(*tool, handler)
}

func (s *MCPServer) AddResource(res *mcp.Resource, handler mcpsrv.ResourceHandlerFunc) {
	s.server.AddResource(*res, handler)
}

func (s *MCPServer) Start(ctx context.Context) error {
	log.Println("MCP Server is ready (stdio)")
	return mcpsrv.ServeStdio(s.server)
}

// Expondo o servidor interno do SDK para quem precisar
func (s *MCPServer) Inner() *mcpsrv.MCPServer {
	return s.server
}
