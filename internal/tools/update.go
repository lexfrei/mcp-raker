package tools

import (
	"context"
	"net/url"

	"github.com/modelcontextprotocol/go-sdk/mcp"

	"github.com/lexfrei/mcp-raker/internal/moonraker"
)

// UpdateNameParams names an optional update target (a configured client,
// "moonraker", "klipper", "system", etc.).
type UpdateNameParams struct {
	Name string `json:"name" jsonschema:"Name of the update target; omit to apply to every configured target"`
}

// UpdateStatusTool returns the definition for moonraker_update_status.
func UpdateStatusTool() *mcp.Tool {
	return &mcp.Tool{
		Name:        "moonraker_update_status",
		Description: "Report installed versions and available updates for Moonraker, Klipper, and clients (GET /machine/update/status).",
		Annotations: readOnly("Update Status"),
	}
}

// NewUpdateStatusHandler creates the handler for moonraker_update_status.
func NewUpdateStatusHandler(api moonraker.API) mcp.ToolHandlerFor[NoParams, map[string]any] {
	return func(ctx context.Context, _ *mcp.CallToolRequest, _ NoParams) (*mcp.CallToolResult, map[string]any, error) {
		out, err := decodeResult(api.Get(ctx, "/machine/update/status", nil))

		return nil, out, err
	}
}

// UpdateRefreshTool returns the definition for moonraker_update_refresh (admin).
func UpdateRefreshTool() *mcp.Tool {
	return &mcp.Tool{
		Name:        "moonraker_update_refresh",
		Description: "Refresh the cached update status for a target (POST /machine/update/refresh).",
		Annotations: writeDestructive("Refresh Update Status"),
	}
}

// NewUpdateRefreshHandler creates the handler for moonraker_update_refresh.
func NewUpdateRefreshHandler(api moonraker.API) mcp.ToolHandlerFor[UpdateNameParams, map[string]any] {
	return func(ctx context.Context, _ *mcp.CallToolRequest, params UpdateNameParams) (*mcp.CallToolResult, map[string]any, error) {
		query := url.Values{}
		if params.Name != "" {
			query.Set(paramName, params.Name)
		}

		out, err := decodeResult(api.Post(ctx, "/machine/update/refresh", query, nil))

		return nil, out, err
	}
}

// UpdateUpgradeTool returns the definition for moonraker_update_upgrade (admin).
func UpdateUpgradeTool() *mcp.Tool {
	return &mcp.Tool{
		Name:        "moonraker_update_upgrade",
		Description: "Install available updates for a target, or all targets when omitted (POST /machine/update/upgrade).",
		Annotations: writeDestructive("Upgrade Software"),
	}
}

// NewUpdateUpgradeHandler creates the handler for moonraker_update_upgrade.
func NewUpdateUpgradeHandler(api moonraker.API) mcp.ToolHandlerFor[UpdateNameParams, map[string]any] {
	return func(ctx context.Context, _ *mcp.CallToolRequest, params UpdateNameParams) (*mcp.CallToolResult, map[string]any, error) {
		query := url.Values{}
		if params.Name != "" {
			query.Set(paramName, params.Name)
		}

		out, err := decodeResult(api.Post(ctx, "/machine/update/upgrade", query, nil))

		return nil, out, err
	}
}

// UpdateRecoverParams defines the parameters for moonraker_update_recover.
type UpdateRecoverParams struct {
	Name string `json:"name" jsonschema:"Name of the git-backed update target to repair"`
	Hard bool   `json:"hard" jsonschema:"When true, perform a hard recovery by re-cloning the repository"`
}

// UpdateRecoverTool returns the definition for moonraker_update_recover (admin).
func UpdateRecoverTool() *mcp.Tool {
	return &mcp.Tool{
		Name:        "moonraker_update_recover",
		Description: "Repair a corrupt git-backed update target (POST /machine/update/recover).",
		Annotations: writeDestructive("Recover Repo"),
	}
}

// NewUpdateRecoverHandler creates the handler for moonraker_update_recover.
func NewUpdateRecoverHandler(api moonraker.API) mcp.ToolHandlerFor[UpdateRecoverParams, map[string]any] {
	return func(ctx context.Context, _ *mcp.CallToolRequest, params UpdateRecoverParams) (*mcp.CallToolResult, map[string]any, error) {
		valErr := requireString(paramName, params.Name)
		if valErr != nil {
			return nil, map[string]any{}, valErr
		}

		query := url.Values{paramName: {params.Name}}
		if params.Hard {
			query.Set("hard", "true")
		}

		out, err := decodeResult(api.Post(ctx, "/machine/update/recover", query, nil))

		return nil, out, err
	}
}

// UpdateRollbackParams defines the parameters for moonraker_update_rollback.
// Unlike refresh and upgrade, rollback targets a single named repository, so
// the name is required.
type UpdateRollbackParams struct {
	Name string `json:"name" jsonschema:"Name of the update target to roll back (required)"`
}

// UpdateRollbackTool returns the definition for moonraker_update_rollback (admin).
func UpdateRollbackTool() *mcp.Tool {
	return &mcp.Tool{
		Name:        "moonraker_update_rollback",
		Description: "Roll a single target back to its previous version (POST /machine/update/rollback).",
		Annotations: writeDestructive("Roll Back Software"),
	}
}

// NewUpdateRollbackHandler creates the handler for moonraker_update_rollback.
func NewUpdateRollbackHandler(api moonraker.API) mcp.ToolHandlerFor[UpdateRollbackParams, map[string]any] {
	return func(ctx context.Context, _ *mcp.CallToolRequest, params UpdateRollbackParams) (*mcp.CallToolResult, map[string]any, error) {
		valErr := requireString(paramName, params.Name)
		if valErr != nil {
			return nil, map[string]any{}, valErr
		}

		out, err := decodeResult(api.Post(ctx, "/machine/update/rollback", url.Values{paramName: {params.Name}}, nil))

		return nil, out, err
	}
}
