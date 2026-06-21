package tools

import (
	"context"
	"net/url"

	"github.com/modelcontextprotocol/go-sdk/mcp"

	"github.com/lexfrei/mcp-raker/internal/moonraker"
)

// DBListTool returns the definition for moonraker_db_list.
func DBListTool() *mcp.Tool {
	return &mcp.Tool{
		Name:        "moonraker_db_list",
		Description: "List the namespaces stored in Moonraker's database (GET /server/database/list).",
		Annotations: readOnly("List DB Namespaces"),
	}
}

// NewDBListHandler creates the handler for moonraker_db_list.
func NewDBListHandler(api moonraker.API) mcp.ToolHandlerFor[NoParams, map[string]any] {
	return func(ctx context.Context, _ *mcp.CallToolRequest, _ NoParams) (*mcp.CallToolResult, map[string]any, error) {
		out, err := decodeResult(api.Get(ctx, "/server/database/list", nil))

		return nil, out, err
	}
}

// DBGetItemParams defines the parameters for moonraker_db_get_item.
type DBGetItemParams struct {
	Namespace string `json:"namespace"     jsonschema:"Database namespace to read from"`
	Key       string `json:"key,omitempty" jsonschema:"Dotted key within the namespace; omit to return the whole namespace"`
}

// DBGetItemTool returns the definition for moonraker_db_get_item.
func DBGetItemTool() *mcp.Tool {
	return &mcp.Tool{
		Name:        "moonraker_db_get_item",
		Description: "Read a value from Moonraker's key-value database (GET /server/database/item).",
		Annotations: readOnly("Get DB Item"),
	}
}

// NewDBGetItemHandler creates the handler for moonraker_db_get_item.
func NewDBGetItemHandler(api moonraker.API) mcp.ToolHandlerFor[DBGetItemParams, map[string]any] {
	return func(ctx context.Context, _ *mcp.CallToolRequest, params DBGetItemParams) (*mcp.CallToolResult, map[string]any, error) {
		valErr := requireString(paramNamespace, params.Namespace)
		if valErr != nil {
			return nil, map[string]any{}, valErr
		}

		query := url.Values{paramNamespace: {params.Namespace}}
		if params.Key != "" {
			query.Set(paramKey, params.Key)
		}

		out, err := decodeResult(api.Get(ctx, "/server/database/item", query))

		return nil, out, err
	}
}

// DBPostItemParams defines the parameters for moonraker_db_post_item.
type DBPostItemParams struct {
	Namespace string `json:"namespace"       jsonschema:"Database namespace to write to"`
	Key       string `json:"key"             jsonschema:"Dotted key within the namespace"`
	Value     any    `json:"value,omitempty" jsonschema:"Value to store; may be any JSON type"`
}

// DBPostItemTool returns the definition for moonraker_db_post_item.
func DBPostItemTool() *mcp.Tool {
	return &mcp.Tool{
		Name:        "moonraker_db_post_item",
		Description: "Store a value in Moonraker's key-value database (POST /server/database/item).",
		Annotations: write("Set DB Item"),
	}
}

// NewDBPostItemHandler creates the handler for moonraker_db_post_item.
func NewDBPostItemHandler(api moonraker.API) mcp.ToolHandlerFor[DBPostItemParams, map[string]any] {
	return func(ctx context.Context, _ *mcp.CallToolRequest, params DBPostItemParams) (*mcp.CallToolResult, map[string]any, error) {
		valErr := requireString(paramNamespace, params.Namespace)
		if valErr != nil {
			return nil, map[string]any{}, valErr
		}

		keyErr := requireString(paramKey, params.Key)
		if keyErr != nil {
			return nil, map[string]any{}, keyErr
		}

		body := map[string]any{paramNamespace: params.Namespace, paramKey: params.Key, paramValue: params.Value}

		out, err := decodeResult(api.Post(ctx, "/server/database/item", nil, body))

		return nil, out, err
	}
}

// DBDeleteItemParams defines the parameters for moonraker_db_delete_item.
type DBDeleteItemParams struct {
	Namespace string `json:"namespace" jsonschema:"Database namespace to delete from"`
	Key       string `json:"key"       jsonschema:"Dotted key within the namespace to delete"`
}

// DBDeleteItemTool returns the definition for moonraker_db_delete_item.
func DBDeleteItemTool() *mcp.Tool {
	return &mcp.Tool{
		Name:        "moonraker_db_delete_item",
		Description: "Delete a value from Moonraker's key-value database (DELETE /server/database/item).",
		Annotations: writeDestructive("Delete DB Item"),
	}
}

// NewDBDeleteItemHandler creates the handler for moonraker_db_delete_item.
func NewDBDeleteItemHandler(api moonraker.API) mcp.ToolHandlerFor[DBDeleteItemParams, map[string]any] {
	return func(ctx context.Context, _ *mcp.CallToolRequest, params DBDeleteItemParams) (*mcp.CallToolResult, map[string]any, error) {
		valErr := requireString(paramNamespace, params.Namespace)
		if valErr != nil {
			return nil, map[string]any{}, valErr
		}

		keyErr := requireString(paramKey, params.Key)
		if keyErr != nil {
			return nil, map[string]any{}, keyErr
		}

		query := url.Values{paramNamespace: {params.Namespace}, paramKey: {params.Key}}

		out, err := decodeResult(api.Delete(ctx, "/server/database/item", query))

		return nil, out, err
	}
}

// DBBackupParams defines the parameters for moonraker_db_backup.
type DBBackupParams struct {
	Filename string `json:"filename,omitempty" jsonschema:"Optional backup filename; omit to use the server default"`
}

// DBBackupTool returns the definition for moonraker_db_backup.
func DBBackupTool() *mcp.Tool {
	return &mcp.Tool{
		Name:        "moonraker_db_backup",
		Description: "Create a backup of Moonraker's database (POST /server/database/backup).",
		Annotations: write("Backup Database"),
	}
}

// NewDBBackupHandler creates the handler for moonraker_db_backup.
func NewDBBackupHandler(api moonraker.API) mcp.ToolHandlerFor[DBBackupParams, map[string]any] {
	return func(ctx context.Context, _ *mcp.CallToolRequest, params DBBackupParams) (*mcp.CallToolResult, map[string]any, error) {
		query := url.Values{}
		if params.Filename != "" {
			query.Set(paramFilename, params.Filename)
		}

		out, err := decodeResult(api.Post(ctx, "/server/database/backup", query, nil))

		return nil, out, err
	}
}

// DBCompactTool returns the definition for moonraker_db_compact.
func DBCompactTool() *mcp.Tool {
	return &mcp.Tool{
		Name:        "moonraker_db_compact",
		Description: "Defragment and compact Moonraker's database (POST /server/database/compact).",
		Annotations: write("Compact Database"),
	}
}

// NewDBCompactHandler creates the handler for moonraker_db_compact.
func NewDBCompactHandler(api moonraker.API) mcp.ToolHandlerFor[NoParams, map[string]any] {
	return func(ctx context.Context, _ *mcp.CallToolRequest, _ NoParams) (*mcp.CallToolResult, map[string]any, error) {
		out, err := decodeResult(api.Post(ctx, "/server/database/compact", nil, nil))

		return nil, out, err
	}
}
