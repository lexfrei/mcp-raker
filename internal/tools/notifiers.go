package tools

import (
	"context"

	"github.com/modelcontextprotocol/go-sdk/mcp"

	"github.com/lexfrei/mcp-raker/internal/moonraker"
)

// NotifiersListTool returns the definition for moonraker_notifiers_list.
func NotifiersListTool() *mcp.Tool {
	return &mcp.Tool{
		Name:        "moonraker_notifiers_list",
		Description: "List configured Apprise notifiers (GET /server/notifiers/list).",
		Annotations: readOnly("List Notifiers"),
	}
}

// NewNotifiersListHandler creates the handler for moonraker_notifiers_list.
func NewNotifiersListHandler(api moonraker.API) mcp.ToolHandlerFor[NoParams, map[string]any] {
	return func(ctx context.Context, _ *mcp.CallToolRequest, _ NoParams) (*mcp.CallToolResult, map[string]any, error) {
		out, err := decodeResult(api.Get(ctx, "/server/notifiers/list", nil))

		return nil, out, err
	}
}
