package tools

import (
	"context"
	"net/url"
	"strings"

	"github.com/modelcontextprotocol/go-sdk/mcp"

	"github.com/lexfrei/mcp-raker/internal/moonraker"
)

// PrinterInfoTool returns the definition for moonraker_printer_info.
func PrinterInfoTool() *mcp.Tool {
	return &mcp.Tool{
		Name:        "moonraker_printer_info",
		Description: "Get Klippy host state, hostname, and software version (GET /printer/info).",
		Annotations: readOnly("Printer Info"),
	}
}

// NewPrinterInfoHandler creates the handler for moonraker_printer_info.
func NewPrinterInfoHandler(api moonraker.API) mcp.ToolHandlerFor[NoParams, moonraker.PrinterInfo] {
	return func(
		ctx context.Context,
		_ *mcp.CallToolRequest,
		_ NoParams,
	) (*mcp.CallToolResult, moonraker.PrinterInfo, error) {
		out, err := decodeTyped[moonraker.PrinterInfo](api.Get(ctx, "/printer/info", nil))

		return nil, out, err
	}
}

// ObjectsListTool returns the definition for moonraker_objects_list.
func ObjectsListTool() *mcp.Tool {
	return &mcp.Tool{
		Name:        "moonraker_objects_list",
		Description: "List every printer object that can be queried (GET /printer/objects/list).",
		Annotations: readOnly("List Printer Objects"),
	}
}

// NewObjectsListHandler creates the handler for moonraker_objects_list.
func NewObjectsListHandler(api moonraker.API) mcp.ToolHandlerFor[NoParams, moonraker.ObjectsList] {
	return func(
		ctx context.Context,
		_ *mcp.CallToolRequest,
		_ NoParams,
	) (*mcp.CallToolResult, moonraker.ObjectsList, error) {
		out, err := decodeTyped[moonraker.ObjectsList](api.Get(ctx, "/printer/objects/list", nil))

		return nil, out, err
	}
}

// ObjectsQueryParams defines the parameters for moonraker_objects_query.
type ObjectsQueryParams struct {
	Objects map[string][]string `json:"objects" jsonschema:"Map of printer object name to the fields to return; use an empty list to return all fields of that object. Example: {\"extruder\": [\"temperature\", \"target\"], \"heater_bed\": [], \"print_stats\": []}"`
}

// ObjectsQueryTool returns the definition for moonraker_objects_query.
func ObjectsQueryTool() *mcp.Tool {
	return &mcp.Tool{
		Name: "moonraker_objects_query",
		Description: "Query the current values of printer objects such as extruder, heater_bed, toolhead, " +
			"print_stats, and virtual_sdcard (GET /printer/objects/query). Call moonraker_objects_list first " +
			"to discover available object names.",
		Annotations: readOnly("Query Printer Objects"),
	}
}

// NewObjectsQueryHandler creates the handler for moonraker_objects_query.
func NewObjectsQueryHandler(api moonraker.API) mcp.ToolHandlerFor[ObjectsQueryParams, moonraker.ObjectsQuery] {
	return func(
		ctx context.Context,
		_ *mcp.CallToolRequest,
		params ObjectsQueryParams,
	) (*mcp.CallToolResult, moonraker.ObjectsQuery, error) {
		valErr := requirePresent("objects", len(params.Objects))
		if valErr != nil {
			return nil, moonraker.ObjectsQuery{}, valErr
		}

		query := url.Values{}
		for name, fields := range params.Objects {
			query.Set(name, strings.Join(fields, ","))
		}

		out, err := decodeTyped[moonraker.ObjectsQuery](api.Get(ctx, "/printer/objects/query", query))

		return nil, out, err
	}
}

// QueryEndstopsTool returns the definition for moonraker_query_endstops.
func QueryEndstopsTool() *mcp.Tool {
	return &mcp.Tool{
		Name:        "moonraker_query_endstops",
		Description: "Report the triggered state of each endstop (GET /printer/query_endstops/status).",
		Annotations: readOnly("Query Endstops"),
	}
}

// NewQueryEndstopsHandler creates the handler for moonraker_query_endstops.
func NewQueryEndstopsHandler(api moonraker.API) mcp.ToolHandlerFor[NoParams, map[string]any] {
	return func(
		ctx context.Context,
		_ *mcp.CallToolRequest,
		_ NoParams,
	) (*mcp.CallToolResult, map[string]any, error) {
		out, err := decodeResult(api.Get(ctx, "/printer/query_endstops/status", nil))

		return nil, out, err
	}
}

// GcodeScriptParams defines the parameters for moonraker_gcode_script.
type GcodeScriptParams struct {
	Script string `json:"script" jsonschema:"G-code command(s) to run, e.g. 'G28' to home or 'M104 S200' to set the hotend temperature"`
}

// GcodeScriptTool returns the definition for moonraker_gcode_script.
func GcodeScriptTool() *mcp.Tool {
	return &mcp.Tool{
		Name: "moonraker_gcode_script",
		Description: "Run one or more raw G-code commands on the printer (POST /printer/gcode/script). " +
			"G-code can move axes, heat the hotend/bed, or disable steppers, so this is treated as destructive.",
		Annotations: writeDestructive("Run G-code"),
	}
}

// NewGcodeScriptHandler creates the handler for moonraker_gcode_script.
func NewGcodeScriptHandler(api moonraker.API) mcp.ToolHandlerFor[GcodeScriptParams, map[string]any] {
	return func(
		ctx context.Context,
		_ *mcp.CallToolRequest,
		params GcodeScriptParams,
	) (*mcp.CallToolResult, map[string]any, error) {
		valErr := requireString(paramScript, params.Script)
		if valErr != nil {
			return nil, map[string]any{}, valErr
		}

		query := url.Values{paramScript: {params.Script}}

		out, err := decodeResult(api.Post(ctx, "/printer/gcode/script", query, nil))

		return nil, out, err
	}
}

// GcodeHelpTool returns the definition for moonraker_gcode_help.
func GcodeHelpTool() *mcp.Tool {
	return &mcp.Tool{
		Name:        "moonraker_gcode_help",
		Description: "List the registered G-code commands and their descriptions (GET /printer/gcode/help).",
		Annotations: readOnly("G-code Help"),
	}
}

// NewGcodeHelpHandler creates the handler for moonraker_gcode_help.
func NewGcodeHelpHandler(api moonraker.API) mcp.ToolHandlerFor[NoParams, map[string]any] {
	return func(
		ctx context.Context,
		_ *mcp.CallToolRequest,
		_ NoParams,
	) (*mcp.CallToolResult, map[string]any, error) {
		out, err := decodeResult(api.Get(ctx, "/printer/gcode/help", nil))

		return nil, out, err
	}
}

// EmergencyStopTool returns the definition for moonraker_emergency_stop.
func EmergencyStopTool() *mcp.Tool {
	return &mcp.Tool{
		Name: "moonraker_emergency_stop",
		Description: "Immediately halt the printer and shut down Klipper (POST /printer/emergency_stop). " +
			"Use in an emergency; a firmware restart is required afterwards to resume.",
		Annotations: writeDestructive("Emergency Stop"),
	}
}

// NewEmergencyStopHandler creates the handler for moonraker_emergency_stop.
func NewEmergencyStopHandler(api moonraker.API) mcp.ToolHandlerFor[NoParams, map[string]any] {
	return func(
		ctx context.Context,
		_ *mcp.CallToolRequest,
		_ NoParams,
	) (*mcp.CallToolResult, map[string]any, error) {
		out, err := decodeResult(api.Post(ctx, "/printer/emergency_stop", nil, nil))

		return nil, out, err
	}
}

// PrinterRestartTool returns the definition for moonraker_printer_restart.
func PrinterRestartTool() *mcp.Tool {
	return &mcp.Tool{
		Name:        "moonraker_printer_restart",
		Description: "Restart the Klipper host software (POST /printer/restart).",
		Annotations: writeDestructive("Restart Klipper"),
	}
}

// NewPrinterRestartHandler creates the handler for moonraker_printer_restart.
func NewPrinterRestartHandler(api moonraker.API) mcp.ToolHandlerFor[NoParams, map[string]any] {
	return func(
		ctx context.Context,
		_ *mcp.CallToolRequest,
		_ NoParams,
	) (*mcp.CallToolResult, map[string]any, error) {
		out, err := decodeResult(api.Post(ctx, "/printer/restart", nil, nil))

		return nil, out, err
	}
}

// FirmwareRestartTool returns the definition for moonraker_firmware_restart.
func FirmwareRestartTool() *mcp.Tool {
	return &mcp.Tool{
		Name: "moonraker_firmware_restart",
		Description: "Restart the Klipper host and reset the micro-controllers (POST /printer/firmware_restart). " +
			"Required to recover after an emergency stop or shutdown.",
		Annotations: writeDestructive("Firmware Restart"),
	}
}

// NewFirmwareRestartHandler creates the handler for moonraker_firmware_restart.
func NewFirmwareRestartHandler(api moonraker.API) mcp.ToolHandlerFor[NoParams, map[string]any] {
	return func(
		ctx context.Context,
		_ *mcp.CallToolRequest,
		_ NoParams,
	) (*mcp.CallToolResult, map[string]any, error) {
		out, err := decodeResult(api.Post(ctx, "/printer/firmware_restart", nil, nil))

		return nil, out, err
	}
}
