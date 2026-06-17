package tools

import (
	"context"
	"net/url"
	"strconv"

	"github.com/modelcontextprotocol/go-sdk/mcp"

	"github.com/lexfrei/mcp-raker/internal/moonraker"
)

// MCPVersion is the build information reported by moonraker_mcp_version.
type MCPVersion struct {
	Version   string `json:"version"`
	Revision  string `json:"revision"`
	GoVersion string `json:"go_version"`
}

// MCPVersionTool returns the definition for moonraker_mcp_version.
func MCPVersionTool() *mcp.Tool {
	return &mcp.Tool{
		Name:        "moonraker_mcp_version",
		Description: "Report the mcp-raker server build version (version, git revision, Go toolchain).",
		Annotations: readOnly("MCP Server Version"),
	}
}

// NewMCPVersionHandler creates the handler for moonraker_mcp_version. It does
// not call the Moonraker API; it returns the server's own build information.
func NewMCPVersionHandler(version, revision, goVersion string) mcp.ToolHandlerFor[NoParams, MCPVersion] {
	return func(
		_ context.Context,
		_ *mcp.CallToolRequest,
		_ NoParams,
	) (*mcp.CallToolResult, MCPVersion, error) {
		return nil, MCPVersion{Version: version, Revision: revision, GoVersion: goVersion}, nil
	}
}

// ServerInfoTool returns the definition for moonraker_server_info.
func ServerInfoTool() *mcp.Tool {
	return &mcp.Tool{
		Name:        "moonraker_server_info",
		Description: "Get Moonraker server state, version, and component status (GET /server/info).",
		Annotations: readOnly("Server Info"),
	}
}

// NewServerInfoHandler creates the handler for moonraker_server_info.
func NewServerInfoHandler(api moonraker.API) mcp.ToolHandlerFor[NoParams, moonraker.ServerInfo] {
	return func(
		ctx context.Context,
		_ *mcp.CallToolRequest,
		_ NoParams,
	) (*mcp.CallToolResult, moonraker.ServerInfo, error) {
		out, err := decodeTyped[moonraker.ServerInfo](api.Get(ctx, "/server/info", nil))

		return nil, out, err
	}
}

// ServerConfigTool returns the definition for moonraker_server_config.
func ServerConfigTool() *mcp.Tool {
	return &mcp.Tool{
		Name:        "moonraker_server_config",
		Description: "Get the parsed Moonraker server configuration (GET /server/config).",
		Annotations: readOnly("Server Config"),
	}
}

// NewServerConfigHandler creates the handler for moonraker_server_config.
func NewServerConfigHandler(api moonraker.API) mcp.ToolHandlerFor[NoParams, RawResult] {
	return func(
		ctx context.Context,
		_ *mcp.CallToolRequest,
		_ NoParams,
	) (*mcp.CallToolResult, RawResult, error) {
		out, err := decodeRaw(api.Get(ctx, "/server/config", nil))

		return nil, out, err
	}
}

// TemperatureStoreParams defines the parameters for moonraker_temperature_store.
type TemperatureStoreParams struct {
	IncludeMonitors bool `json:"include_monitors" jsonschema:"When true, include temperature monitors in addition to heaters and sensors"`
}

// TemperatureStoreTool returns the definition for moonraker_temperature_store.
func TemperatureStoreTool() *mcp.Tool {
	return &mcp.Tool{
		Name: "moonraker_temperature_store",
		Description: "Get the cached temperature history for heaters and sensors " +
			"(GET /server/temperature_store).",
		Annotations: readOnly("Temperature Store"),
	}
}

// NewTemperatureStoreHandler creates the handler for moonraker_temperature_store.
func NewTemperatureStoreHandler(api moonraker.API) mcp.ToolHandlerFor[TemperatureStoreParams, RawResult] {
	return func(
		ctx context.Context,
		_ *mcp.CallToolRequest,
		params TemperatureStoreParams,
	) (*mcp.CallToolResult, RawResult, error) {
		query := url.Values{}
		if params.IncludeMonitors {
			query.Set("include_monitors", "true")
		}

		out, err := decodeRaw(api.Get(ctx, "/server/temperature_store", query))

		return nil, out, err
	}
}

// GcodeStoreParams defines the parameters for moonraker_gcode_store.
type GcodeStoreParams struct {
	Count int `json:"count" jsonschema:"Number of recent G-code store entries to return; omit or 0 to use the server default"`
}

// GcodeStoreTool returns the definition for moonraker_gcode_store.
func GcodeStoreTool() *mcp.Tool {
	return &mcp.Tool{
		Name:        "moonraker_gcode_store",
		Description: "Get the cached G-code console output and command history (GET /server/gcode_store).",
		Annotations: readOnly("G-code Store"),
	}
}

// NewGcodeStoreHandler creates the handler for moonraker_gcode_store.
func NewGcodeStoreHandler(api moonraker.API) mcp.ToolHandlerFor[GcodeStoreParams, RawResult] {
	return func(
		ctx context.Context,
		_ *mcp.CallToolRequest,
		params GcodeStoreParams,
	) (*mcp.CallToolResult, RawResult, error) {
		query := url.Values{}
		if params.Count > 0 {
			query.Set("count", strconv.Itoa(params.Count))
		}

		out, err := decodeRaw(api.Get(ctx, "/server/gcode_store", query))

		return nil, out, err
	}
}

// LogsRolloverParams defines the parameters for moonraker_logs_rollover.
type LogsRolloverParams struct {
	Application string `json:"application" jsonschema:"Optional application log to roll over (e.g. 'moonraker' or 'klipper'); omit to roll over all logs"`
}

// LogsRolloverTool returns the definition for moonraker_logs_rollover (admin).
func LogsRolloverTool() *mcp.Tool {
	return &mcp.Tool{
		Name:        "moonraker_logs_rollover",
		Description: "Roll over Moonraker and Klipper log files (POST /server/logs/rollover).",
		Annotations: writeDestructive("Roll Over Logs"),
	}
}

// NewLogsRolloverHandler creates the handler for moonraker_logs_rollover.
func NewLogsRolloverHandler(api moonraker.API) mcp.ToolHandlerFor[LogsRolloverParams, RawResult] {
	return func(
		ctx context.Context,
		_ *mcp.CallToolRequest,
		params LogsRolloverParams,
	) (*mcp.CallToolResult, RawResult, error) {
		query := url.Values{}
		if params.Application != "" {
			query.Set("application", params.Application)
		}

		out, err := decodeRaw(api.Post(ctx, "/server/logs/rollover", query, nil))

		return nil, out, err
	}
}

// ServerRestartTool returns the definition for moonraker_server_restart (admin).
func ServerRestartTool() *mcp.Tool {
	return &mcp.Tool{
		Name:        "moonraker_server_restart",
		Description: "Restart the Moonraker server process (POST /server/restart).",
		Annotations: writeDestructive("Restart Moonraker"),
	}
}

// NewServerRestartHandler creates the handler for moonraker_server_restart.
func NewServerRestartHandler(api moonraker.API) mcp.ToolHandlerFor[NoParams, RawResult] {
	return func(
		ctx context.Context,
		_ *mcp.CallToolRequest,
		_ NoParams,
	) (*mcp.CallToolResult, RawResult, error) {
		out, err := decodeRaw(api.Post(ctx, "/server/restart", nil, nil))

		return nil, out, err
	}
}
