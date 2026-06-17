package tools

import (
	"context"
	"net/url"

	"github.com/modelcontextprotocol/go-sdk/mcp"

	"github.com/lexfrei/mcp-raker/internal/moonraker"
)

// PowerDeviceParams names a power device managed by Moonraker.
type PowerDeviceParams struct {
	Device string `json:"device" jsonschema:"Name of the power device as configured in Moonraker"`
}

// PowerDevicesTool returns the definition for moonraker_power_devices.
func PowerDevicesTool() *mcp.Tool {
	return &mcp.Tool{
		Name:        "moonraker_power_devices",
		Description: "List the power devices Moonraker can control (GET /machine/device_power/devices).",
		Annotations: readOnly("List Power Devices"),
	}
}

// NewPowerDevicesHandler creates the handler for moonraker_power_devices.
func NewPowerDevicesHandler(api moonraker.API) mcp.ToolHandlerFor[NoParams, RawResult] {
	return func(ctx context.Context, _ *mcp.CallToolRequest, _ NoParams) (*mcp.CallToolResult, RawResult, error) {
		out, err := decodeRaw(api.Get(ctx, "/machine/device_power/devices", nil))

		return nil, out, err
	}
}

// PowerStatusTool returns the definition for moonraker_power_status.
func PowerStatusTool() *mcp.Tool {
	return &mcp.Tool{
		Name:        "moonraker_power_status",
		Description: "Get the on/off state of a power device (GET /machine/device_power/device).",
		Annotations: readOnly("Power Device Status"),
	}
}

// NewPowerStatusHandler creates the handler for moonraker_power_status.
func NewPowerStatusHandler(api moonraker.API) mcp.ToolHandlerFor[PowerDeviceParams, RawResult] {
	return func(ctx context.Context, _ *mcp.CallToolRequest, params PowerDeviceParams) (*mcp.CallToolResult, RawResult, error) {
		valErr := requireString(paramDevice, params.Device)
		if valErr != nil {
			return nil, RawResult{}, valErr
		}

		out, err := decodeRaw(api.Get(ctx, "/machine/device_power/device", url.Values{paramDevice: {params.Device}}))

		return nil, out, err
	}
}

// devicePath is the single-device endpoint that takes a device + action.
const devicePath = "/machine/device_power/device"

// powerActionBody builds the JSON body for the single-device action endpoint.
func powerActionBody(device, action string) map[string]any {
	return map[string]any{paramDevice: device, paramAction: action}
}

// PowerOnTool returns the definition for moonraker_power_on.
func PowerOnTool() *mcp.Tool {
	return &mcp.Tool{
		Name:        "moonraker_power_on",
		Description: "Switch a power device on (POST /machine/device_power/device, action=on).",
		Annotations: write("Power On Device"),
	}
}

// NewPowerOnHandler creates the handler for moonraker_power_on.
func NewPowerOnHandler(api moonraker.API) mcp.ToolHandlerFor[PowerDeviceParams, RawResult] {
	return func(ctx context.Context, _ *mcp.CallToolRequest, params PowerDeviceParams) (*mcp.CallToolResult, RawResult, error) {
		valErr := requireString(paramDevice, params.Device)
		if valErr != nil {
			return nil, RawResult{}, valErr
		}

		out, err := decodeRaw(api.Post(ctx, devicePath, nil, powerActionBody(params.Device, "on")))

		return nil, out, err
	}
}

// PowerOffTool returns the definition for moonraker_power_off.
func PowerOffTool() *mcp.Tool {
	return &mcp.Tool{
		Name: "moonraker_power_off",
		Description: "Switch a power device off (POST /machine/device_power/device, action=off). " +
			"May interrupt an active print.",
		Annotations: writeDestructive("Power Off Device"),
	}
}

// NewPowerOffHandler creates the handler for moonraker_power_off.
func NewPowerOffHandler(api moonraker.API) mcp.ToolHandlerFor[PowerDeviceParams, RawResult] {
	return func(ctx context.Context, _ *mcp.CallToolRequest, params PowerDeviceParams) (*mcp.CallToolResult, RawResult, error) {
		valErr := requireString(paramDevice, params.Device)
		if valErr != nil {
			return nil, RawResult{}, valErr
		}

		out, err := decodeRaw(api.Post(ctx, devicePath, nil, powerActionBody(params.Device, "off")))

		return nil, out, err
	}
}

// PowerToggleTool returns the definition for moonraker_power_toggle.
func PowerToggleTool() *mcp.Tool {
	return &mcp.Tool{
		Name:        "moonraker_power_toggle",
		Description: "Toggle a power device's state (POST /machine/device_power/device, action=toggle).",
		Annotations: writeDestructive("Toggle Power Device"),
	}
}

// NewPowerToggleHandler creates the handler for moonraker_power_toggle.
func NewPowerToggleHandler(api moonraker.API) mcp.ToolHandlerFor[PowerDeviceParams, RawResult] {
	return func(ctx context.Context, _ *mcp.CallToolRequest, params PowerDeviceParams) (*mcp.CallToolResult, RawResult, error) {
		valErr := requireString(paramDevice, params.Device)
		if valErr != nil {
			return nil, RawResult{}, valErr
		}

		out, err := decodeRaw(api.Post(ctx, devicePath, nil, powerActionBody(params.Device, "toggle")))

		return nil, out, err
	}
}
