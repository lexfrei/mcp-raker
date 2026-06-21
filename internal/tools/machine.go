package tools

import (
	"context"
	"net/url"

	"github.com/modelcontextprotocol/go-sdk/mcp"

	"github.com/lexfrei/mcp-raker/internal/moonraker"
)

// SystemInfoTool returns the definition for moonraker_system_info.
func SystemInfoTool() *mcp.Tool {
	return &mcp.Tool{
		Name:        "moonraker_system_info",
		Description: "Get host system information: CPU, memory, distribution, and network (GET /machine/system_info).",
		Annotations: readOnly("System Info"),
	}
}

// NewSystemInfoHandler creates the handler for moonraker_system_info.
func NewSystemInfoHandler(api moonraker.API) mcp.ToolHandlerFor[NoParams, map[string]any] {
	return func(ctx context.Context, _ *mcp.CallToolRequest, _ NoParams) (*mcp.CallToolResult, map[string]any, error) {
		out, err := decodeResult(api.Get(ctx, "/machine/system_info", nil))

		return nil, out, err
	}
}

// ProcStatsTool returns the definition for moonraker_proc_stats.
func ProcStatsTool() *mcp.Tool {
	return &mcp.Tool{
		Name:        "moonraker_proc_stats",
		Description: "Get CPU and network usage, throttle state, and uptime (GET /machine/proc_stats).",
		Annotations: readOnly("Process Stats"),
	}
}

// NewProcStatsHandler creates the handler for moonraker_proc_stats.
func NewProcStatsHandler(api moonraker.API) mcp.ToolHandlerFor[NoParams, map[string]any] {
	return func(ctx context.Context, _ *mcp.CallToolRequest, _ NoParams) (*mcp.CallToolResult, map[string]any, error) {
		out, err := decodeResult(api.Get(ctx, "/machine/proc_stats", nil))

		return nil, out, err
	}
}

// SudoInfoTool returns the definition for moonraker_sudo_info.
func SudoInfoTool() *mcp.Tool {
	return &mcp.Tool{
		Name:        "moonraker_sudo_info",
		Description: "Report whether Moonraker has sudo access and any pending sudo requests (GET /machine/sudo/info).",
		Annotations: readOnly("Sudo Info"),
	}
}

// NewSudoInfoHandler creates the handler for moonraker_sudo_info.
func NewSudoInfoHandler(api moonraker.API) mcp.ToolHandlerFor[NoParams, map[string]any] {
	return func(ctx context.Context, _ *mcp.CallToolRequest, _ NoParams) (*mcp.CallToolResult, map[string]any, error) {
		out, err := decodeResult(api.Get(ctx, "/machine/sudo/info", nil))

		return nil, out, err
	}
}

// PeripheralsUSBTool returns the definition for moonraker_peripherals_usb.
func PeripheralsUSBTool() *mcp.Tool {
	return &mcp.Tool{
		Name:        "moonraker_peripherals_usb",
		Description: "List connected USB devices (GET /machine/peripherals/usb).",
		Annotations: readOnly("USB Peripherals"),
	}
}

// NewPeripheralsUSBHandler creates the handler for moonraker_peripherals_usb.
func NewPeripheralsUSBHandler(api moonraker.API) mcp.ToolHandlerFor[NoParams, map[string]any] {
	return func(ctx context.Context, _ *mcp.CallToolRequest, _ NoParams) (*mcp.CallToolResult, map[string]any, error) {
		out, err := decodeResult(api.Get(ctx, "/machine/peripherals/usb", nil))

		return nil, out, err
	}
}

// PeripheralsSerialTool returns the definition for moonraker_peripherals_serial.
func PeripheralsSerialTool() *mcp.Tool {
	return &mcp.Tool{
		Name:        "moonraker_peripherals_serial",
		Description: "List available serial devices (GET /machine/peripherals/serial).",
		Annotations: readOnly("Serial Peripherals"),
	}
}

// NewPeripheralsSerialHandler creates the handler for moonraker_peripherals_serial.
func NewPeripheralsSerialHandler(api moonraker.API) mcp.ToolHandlerFor[NoParams, map[string]any] {
	return func(ctx context.Context, _ *mcp.CallToolRequest, _ NoParams) (*mcp.CallToolResult, map[string]any, error) {
		out, err := decodeResult(api.Get(ctx, "/machine/peripherals/serial", nil))

		return nil, out, err
	}
}

// PeripheralsVideoTool returns the definition for moonraker_peripherals_video.
func PeripheralsVideoTool() *mcp.Tool {
	return &mcp.Tool{
		Name:        "moonraker_peripherals_video",
		Description: "List available video capture devices (GET /machine/peripherals/video).",
		Annotations: readOnly("Video Peripherals"),
	}
}

// NewPeripheralsVideoHandler creates the handler for moonraker_peripherals_video.
func NewPeripheralsVideoHandler(api moonraker.API) mcp.ToolHandlerFor[NoParams, map[string]any] {
	return func(ctx context.Context, _ *mcp.CallToolRequest, _ NoParams) (*mcp.CallToolResult, map[string]any, error) {
		out, err := decodeResult(api.Get(ctx, "/machine/peripherals/video", nil))

		return nil, out, err
	}
}

// PeripheralsCanbusParams defines the parameters for moonraker_peripherals_canbus.
type PeripheralsCanbusParams struct {
	Interface string `json:"interface" jsonschema:"CAN interface to scan, e.g. 'can0'; omit to use the default"`
}

// PeripheralsCanbusTool returns the definition for moonraker_peripherals_canbus.
func PeripheralsCanbusTool() *mcp.Tool {
	return &mcp.Tool{
		Name:        "moonraker_peripherals_canbus",
		Description: "List CAN bus devices on an interface (GET /machine/peripherals/canbus).",
		Annotations: readOnly("CAN Bus Peripherals"),
	}
}

// NewPeripheralsCanbusHandler creates the handler for moonraker_peripherals_canbus.
func NewPeripheralsCanbusHandler(api moonraker.API) mcp.ToolHandlerFor[PeripheralsCanbusParams, map[string]any] {
	return func(
		ctx context.Context,
		_ *mcp.CallToolRequest,
		params PeripheralsCanbusParams,
	) (*mcp.CallToolResult, map[string]any, error) {
		query := url.Values{}
		if params.Interface != "" {
			query.Set("interface", params.Interface)
		}

		out, err := decodeResult(api.Get(ctx, "/machine/peripherals/canbus", query))

		return nil, out, err
	}
}

// MachineShutdownTool returns the definition for moonraker_machine_shutdown (admin).
func MachineShutdownTool() *mcp.Tool {
	return &mcp.Tool{
		Name:        "moonraker_machine_shutdown",
		Description: "Shut down the host operating system (POST /machine/shutdown).",
		Annotations: writeDestructive("Shut Down Host"),
	}
}

// NewMachineShutdownHandler creates the handler for moonraker_machine_shutdown.
func NewMachineShutdownHandler(api moonraker.API) mcp.ToolHandlerFor[NoParams, map[string]any] {
	return func(ctx context.Context, _ *mcp.CallToolRequest, _ NoParams) (*mcp.CallToolResult, map[string]any, error) {
		out, err := decodeResult(api.Post(ctx, "/machine/shutdown", nil, nil))

		return nil, out, err
	}
}

// MachineRebootTool returns the definition for moonraker_machine_reboot (admin).
func MachineRebootTool() *mcp.Tool {
	return &mcp.Tool{
		Name:        "moonraker_machine_reboot",
		Description: "Reboot the host operating system (POST /machine/reboot).",
		Annotations: writeDestructive("Reboot Host"),
	}
}

// NewMachineRebootHandler creates the handler for moonraker_machine_reboot.
func NewMachineRebootHandler(api moonraker.API) mcp.ToolHandlerFor[NoParams, map[string]any] {
	return func(ctx context.Context, _ *mcp.CallToolRequest, _ NoParams) (*mcp.CallToolResult, map[string]any, error) {
		out, err := decodeResult(api.Post(ctx, "/machine/reboot", nil, nil))

		return nil, out, err
	}
}

// ServiceParams names a systemd service for the service control tools.
type ServiceParams struct {
	Service string `json:"service" jsonschema:"Name of the systemd service, e.g. 'klipper' or 'moonraker'"`
}

// ServiceStartTool returns the definition for moonraker_service_start (admin).
func ServiceStartTool() *mcp.Tool {
	return &mcp.Tool{
		Name:        "moonraker_service_start",
		Description: "Start a systemd service managed by Moonraker (POST /machine/services/start).",
		Annotations: writeDestructive("Start Service"),
	}
}

// NewServiceStartHandler creates the handler for moonraker_service_start.
func NewServiceStartHandler(api moonraker.API) mcp.ToolHandlerFor[ServiceParams, map[string]any] {
	return func(ctx context.Context, _ *mcp.CallToolRequest, params ServiceParams) (*mcp.CallToolResult, map[string]any, error) {
		valErr := requireString(paramService, params.Service)
		if valErr != nil {
			return nil, map[string]any{}, valErr
		}

		out, err := decodeResult(api.Post(ctx, "/machine/services/start", url.Values{paramService: {params.Service}}, nil))

		return nil, out, err
	}
}

// ServiceStopTool returns the definition for moonraker_service_stop (admin).
func ServiceStopTool() *mcp.Tool {
	return &mcp.Tool{
		Name:        "moonraker_service_stop",
		Description: "Stop a systemd service managed by Moonraker (POST /machine/services/stop).",
		Annotations: writeDestructive("Stop Service"),
	}
}

// NewServiceStopHandler creates the handler for moonraker_service_stop.
func NewServiceStopHandler(api moonraker.API) mcp.ToolHandlerFor[ServiceParams, map[string]any] {
	return func(ctx context.Context, _ *mcp.CallToolRequest, params ServiceParams) (*mcp.CallToolResult, map[string]any, error) {
		valErr := requireString(paramService, params.Service)
		if valErr != nil {
			return nil, map[string]any{}, valErr
		}

		out, err := decodeResult(api.Post(ctx, "/machine/services/stop", url.Values{paramService: {params.Service}}, nil))

		return nil, out, err
	}
}

// ServiceRestartTool returns the definition for moonraker_service_restart (admin).
func ServiceRestartTool() *mcp.Tool {
	return &mcp.Tool{
		Name:        "moonraker_service_restart",
		Description: "Restart a systemd service managed by Moonraker (POST /machine/services/restart).",
		Annotations: writeDestructive("Restart Service"),
	}
}

// NewServiceRestartHandler creates the handler for moonraker_service_restart.
func NewServiceRestartHandler(api moonraker.API) mcp.ToolHandlerFor[ServiceParams, map[string]any] {
	return func(ctx context.Context, _ *mcp.CallToolRequest, params ServiceParams) (*mcp.CallToolResult, map[string]any, error) {
		valErr := requireString(paramService, params.Service)
		if valErr != nil {
			return nil, map[string]any{}, valErr
		}

		out, err := decodeResult(api.Post(ctx, "/machine/services/restart", url.Values{paramService: {params.Service}}, nil))

		return nil, out, err
	}
}

// SudoPasswordParams defines the parameters for moonraker_sudo_password (admin).
type SudoPasswordParams struct {
	Password string `json:"password" jsonschema:"The sudo password to grant Moonraker elevated access"`
}

// SudoPasswordTool returns the definition for moonraker_sudo_password (admin).
func SudoPasswordTool() *mcp.Tool {
	return &mcp.Tool{
		Name:        "moonraker_sudo_password",
		Description: "Provide the host sudo password so Moonraker can run privileged actions (POST /machine/sudo/password).",
		Annotations: writeDestructive("Set Sudo Password"),
	}
}

// NewSudoPasswordHandler creates the handler for moonraker_sudo_password.
func NewSudoPasswordHandler(api moonraker.API) mcp.ToolHandlerFor[SudoPasswordParams, map[string]any] {
	return func(
		ctx context.Context,
		_ *mcp.CallToolRequest,
		params SudoPasswordParams,
	) (*mcp.CallToolResult, map[string]any, error) {
		valErr := requireString(paramPassword, params.Password)
		if valErr != nil {
			return nil, map[string]any{}, valErr
		}

		out, err := decodeResult(api.Post(ctx, "/machine/sudo/password", nil, map[string]any{paramPassword: params.Password}))

		return nil, out, err
	}
}
