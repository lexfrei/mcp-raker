package tools

import (
	"context"

	"github.com/modelcontextprotocol/go-sdk/mcp"

	"github.com/lexfrei/mcp-raker/internal/moonraker"
)

// SpoolmanStatusTool returns the definition for moonraker_spoolman_status.
func SpoolmanStatusTool() *mcp.Tool {
	return &mcp.Tool{
		Name:        "moonraker_spoolman_status",
		Description: "Report the Spoolman connection state and any pending usage reports (GET /server/spoolman/status).",
		Annotations: readOnly("Spoolman Status"),
	}
}

// NewSpoolmanStatusHandler creates the handler for moonraker_spoolman_status.
func NewSpoolmanStatusHandler(api moonraker.API) mcp.ToolHandlerFor[NoParams, map[string]any] {
	return func(ctx context.Context, _ *mcp.CallToolRequest, _ NoParams) (*mcp.CallToolResult, map[string]any, error) {
		out, err := decodeResult(api.Get(ctx, "/server/spoolman/status", nil))

		return nil, out, err
	}
}

// SpoolmanGetSpoolTool returns the definition for moonraker_spoolman_get_spool.
func SpoolmanGetSpoolTool() *mcp.Tool {
	return &mcp.Tool{
		Name:        "moonraker_spoolman_get_spool",
		Description: "Get the currently active Spoolman spool id (GET /server/spoolman/spool_id).",
		Annotations: readOnly("Get Active Spool"),
	}
}

// NewSpoolmanGetSpoolHandler creates the handler for moonraker_spoolman_get_spool.
func NewSpoolmanGetSpoolHandler(api moonraker.API) mcp.ToolHandlerFor[NoParams, map[string]any] {
	return func(ctx context.Context, _ *mcp.CallToolRequest, _ NoParams) (*mcp.CallToolResult, map[string]any, error) {
		out, err := decodeResult(api.Get(ctx, "/server/spoolman/spool_id", nil))

		return nil, out, err
	}
}

// SpoolmanSetSpoolParams defines the parameters for moonraker_spoolman_set_spool.
type SpoolmanSetSpoolParams struct {
	SpoolID int `json:"spool_id" jsonschema:"Spoolman spool id to mark as active"`
}

// SpoolmanSetSpoolTool returns the definition for moonraker_spoolman_set_spool.
func SpoolmanSetSpoolTool() *mcp.Tool {
	return &mcp.Tool{
		Name:        "moonraker_spoolman_set_spool",
		Description: "Set the active Spoolman spool (POST /server/spoolman/spool_id).",
		Annotations: write("Set Active Spool"),
	}
}

// NewSpoolmanSetSpoolHandler creates the handler for moonraker_spoolman_set_spool.
func NewSpoolmanSetSpoolHandler(api moonraker.API) mcp.ToolHandlerFor[SpoolmanSetSpoolParams, map[string]any] {
	return func(ctx context.Context, _ *mcp.CallToolRequest, params SpoolmanSetSpoolParams) (*mcp.CallToolResult, map[string]any, error) {
		valErr := requirePositive("spool_id", params.SpoolID)
		if valErr != nil {
			return nil, map[string]any{}, valErr
		}

		out, err := decodeResult(api.Post(ctx, "/server/spoolman/spool_id", nil, map[string]any{"spool_id": params.SpoolID}))

		return nil, out, err
	}
}

// SpoolmanProxyParams defines the parameters for moonraker_spoolman_proxy.
type SpoolmanProxyParams struct {
	RequestMethod string         `json:"request_method" jsonschema:"HTTP method to forward to the Spoolman server, e.g. 'GET' or 'POST'"`
	Path          string         `json:"path"           jsonschema:"Spoolman API path to call, e.g. '/v1/spool'"`
	Query         map[string]any `json:"query"          jsonschema:"Optional query parameters to forward"`
	Body          any            `json:"body"           jsonschema:"Optional request body to forward"`
}

// SpoolmanProxyTool returns the definition for moonraker_spoolman_proxy.
func SpoolmanProxyTool() *mcp.Tool {
	return &mcp.Tool{
		Name: "moonraker_spoolman_proxy",
		Description: "Forward an arbitrary request (any HTTP method, including DELETE) to the configured " +
			"Spoolman server (POST /server/spoolman/proxy).",
		Annotations: writeDestructive("Spoolman Proxy"),
	}
}

// NewSpoolmanProxyHandler creates the handler for moonraker_spoolman_proxy.
// The proxied response shape varies by request, so it passes through unchanged
// with no output schema.
func NewSpoolmanProxyHandler(api moonraker.API) mcp.ToolHandlerFor[SpoolmanProxyParams, any] {
	return func(ctx context.Context, _ *mcp.CallToolRequest, params SpoolmanProxyParams) (*mcp.CallToolResult, any, error) {
		methodErr := requireString("request_method", params.RequestMethod)
		if methodErr != nil {
			return nil, nil, methodErr
		}

		pathErr := requireString(paramPath, params.Path)
		if pathErr != nil {
			return nil, nil, pathErr
		}

		body := map[string]any{
			"request_method":  params.RequestMethod,
			"path":            params.Path,
			"use_v2_response": true,
		}
		if params.Query != nil {
			body["query"] = params.Query
		}

		if params.Body != nil {
			body["body"] = params.Body
		}

		out, err := decodePassthrough(api.Post(ctx, "/server/spoolman/proxy", nil, body))

		return nil, out, err
	}
}
