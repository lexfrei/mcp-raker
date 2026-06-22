package tools

import (
	"context"

	"github.com/modelcontextprotocol/go-sdk/mcp"

	"github.com/lexfrei/mcp-raker/internal/moonraker"
)

// ExtensionsListTool returns the definition for moonraker_extensions_list.
func ExtensionsListTool() *mcp.Tool {
	return &mcp.Tool{
		Name:        "moonraker_extensions_list",
		Description: "List connected agents and extensions (GET /server/extensions/list).",
		Annotations: readOnly("List Extensions"),
	}
}

// NewExtensionsListHandler creates the handler for moonraker_extensions_list.
func NewExtensionsListHandler(api moonraker.API) mcp.ToolHandlerFor[NoParams, map[string]any] {
	return func(ctx context.Context, _ *mcp.CallToolRequest, _ NoParams) (*mcp.CallToolResult, map[string]any, error) {
		out, err := decodeResult(api.Get(ctx, "/server/extensions/list", nil))

		return nil, out, err
	}
}

// ExtensionsRequestParams defines the parameters for moonraker_extensions_request.
type ExtensionsRequestParams struct {
	Agent     string `json:"agent"               jsonschema:"Name of the registered agent or extension to call"`
	Method    string `json:"method"              jsonschema:"Method the agent exposes"`
	Arguments any    `json:"arguments,omitempty" jsonschema:"Optional arguments to pass to the agent method"`
}

// ExtensionsRequestTool returns the definition for moonraker_extensions_request.
func ExtensionsRequestTool() *mcp.Tool {
	return &mcp.Tool{
		Name:        "moonraker_extensions_request",
		Description: "Invoke a method on a connected agent or extension (POST /server/extensions/request).",
		Annotations: write("Call Extension"),
	}
}

// NewExtensionsRequestHandler creates the handler for moonraker_extensions_request.
// The agent's response shape is arbitrary, so it is returned verbatim under a
// "response" key — the wrapper keeps the structured content an object (required
// by MCP) for scalar or array payloads.
func NewExtensionsRequestHandler(api moonraker.API) mcp.ToolHandlerFor[ExtensionsRequestParams, map[string]any] {
	return func(ctx context.Context, _ *mcp.CallToolRequest, params ExtensionsRequestParams) (*mcp.CallToolResult, map[string]any, error) {
		agentErr := requireString("agent", params.Agent)
		if agentErr != nil {
			return nil, nil, agentErr
		}

		methodErr := requireString("method", params.Method)
		if methodErr != nil {
			return nil, nil, methodErr
		}

		body := map[string]any{"agent": params.Agent, "method": params.Method}
		if params.Arguments != nil {
			body["arguments"] = params.Arguments
		}

		value, err := decodePassthrough(api.Post(ctx, "/server/extensions/request", nil, body))
		if err != nil {
			return nil, nil, err
		}

		return nil, map[string]any{"response": value}, nil
	}
}
