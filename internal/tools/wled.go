package tools

import (
	"context"
	"net/url"
	"strconv"

	"github.com/modelcontextprotocol/go-sdk/mcp"

	"github.com/lexfrei/mcp-raker/internal/moonraker"
)

// StripParams names a WLED strip.
type StripParams struct {
	Strip string `json:"strip" jsonschema:"Name of the WLED strip as configured in Moonraker"`
}

// WLEDStripsTool returns the definition for moonraker_wled_strips.
func WLEDStripsTool() *mcp.Tool {
	return &mcp.Tool{
		Name:        "moonraker_wled_strips",
		Description: "List configured WLED strips (GET /machine/wled/strips).",
		Annotations: readOnly("List WLED Strips"),
	}
}

// NewWLEDStripsHandler creates the handler for moonraker_wled_strips.
func NewWLEDStripsHandler(api moonraker.API) mcp.ToolHandlerFor[NoParams, RawResult] {
	return func(ctx context.Context, _ *mcp.CallToolRequest, _ NoParams) (*mcp.CallToolResult, RawResult, error) {
		out, err := decodeRaw(api.Get(ctx, "/machine/wled/strips", nil))

		return nil, out, err
	}
}

// WLEDStatusTool returns the definition for moonraker_wled_status.
func WLEDStatusTool() *mcp.Tool {
	return &mcp.Tool{
		Name:        "moonraker_wled_status",
		Description: "Get the current state of a WLED strip (GET /machine/wled/status).",
		Annotations: readOnly("WLED Status"),
	}
}

// NewWLEDStatusHandler creates the handler for moonraker_wled_status.
func NewWLEDStatusHandler(api moonraker.API) mcp.ToolHandlerFor[StripParams, RawResult] {
	return func(ctx context.Context, _ *mcp.CallToolRequest, params StripParams) (*mcp.CallToolResult, RawResult, error) {
		valErr := requireString(paramStrip, params.Strip)
		if valErr != nil {
			return nil, RawResult{}, valErr
		}

		out, err := decodeRaw(api.Get(ctx, "/machine/wled/status", url.Values{paramStrip: {params.Strip}}))

		return nil, out, err
	}
}

// WLEDOnTool returns the definition for moonraker_wled_on.
func WLEDOnTool() *mcp.Tool {
	return &mcp.Tool{
		Name:        "moonraker_wled_on",
		Description: "Turn a WLED strip on, restoring its configured preset (POST /machine/wled/on).",
		Annotations: write("WLED On"),
	}
}

// NewWLEDOnHandler creates the handler for moonraker_wled_on.
func NewWLEDOnHandler(api moonraker.API) mcp.ToolHandlerFor[StripParams, RawResult] {
	return func(ctx context.Context, _ *mcp.CallToolRequest, params StripParams) (*mcp.CallToolResult, RawResult, error) {
		valErr := requireString(paramStrip, params.Strip)
		if valErr != nil {
			return nil, RawResult{}, valErr
		}

		out, err := decodeRaw(api.Post(ctx, "/machine/wled/on", url.Values{paramStrip: {params.Strip}}, nil))

		return nil, out, err
	}
}

// WLEDOffTool returns the definition for moonraker_wled_off.
func WLEDOffTool() *mcp.Tool {
	return &mcp.Tool{
		Name:        "moonraker_wled_off",
		Description: "Turn a WLED strip off (POST /machine/wled/off).",
		Annotations: write("WLED Off"),
	}
}

// NewWLEDOffHandler creates the handler for moonraker_wled_off.
func NewWLEDOffHandler(api moonraker.API) mcp.ToolHandlerFor[StripParams, RawResult] {
	return func(ctx context.Context, _ *mcp.CallToolRequest, params StripParams) (*mcp.CallToolResult, RawResult, error) {
		valErr := requireString(paramStrip, params.Strip)
		if valErr != nil {
			return nil, RawResult{}, valErr
		}

		out, err := decodeRaw(api.Post(ctx, "/machine/wled/off", url.Values{paramStrip: {params.Strip}}, nil))

		return nil, out, err
	}
}

// WLEDToggleTool returns the definition for moonraker_wled_toggle.
func WLEDToggleTool() *mcp.Tool {
	return &mcp.Tool{
		Name:        "moonraker_wled_toggle",
		Description: "Toggle a WLED strip on or off (POST /machine/wled/toggle).",
		Annotations: write("WLED Toggle"),
	}
}

// NewWLEDToggleHandler creates the handler for moonraker_wled_toggle.
func NewWLEDToggleHandler(api moonraker.API) mcp.ToolHandlerFor[StripParams, RawResult] {
	return func(ctx context.Context, _ *mcp.CallToolRequest, params StripParams) (*mcp.CallToolResult, RawResult, error) {
		valErr := requireString(paramStrip, params.Strip)
		if valErr != nil {
			return nil, RawResult{}, valErr
		}

		out, err := decodeRaw(api.Post(ctx, "/machine/wled/toggle", url.Values{paramStrip: {params.Strip}}, nil))

		return nil, out, err
	}
}

// WLEDSetParams defines the parameters for moonraker_wled_set. Preset,
// intensity, and speed are pointers because 0 is a valid value in their ranges
// (preset index 0, intensity/speed 0-255), so it must be distinguishable from
// "unset". Brightness is documented 1-255, so its zero value safely means unset.
type WLEDSetParams struct {
	Strip      string `json:"strip"      jsonschema:"Name of the WLED strip"`
	Preset     *int   `json:"preset"     jsonschema:"Preset index to apply; omit to leave unchanged"`
	Brightness int    `json:"brightness" jsonschema:"Brightness 1-255; omit or 0 to leave unchanged"`
	Intensity  *int   `json:"intensity"  jsonschema:"Effect intensity 0-255; omit to leave unchanged"`
	Speed      *int   `json:"speed"      jsonschema:"Effect speed 0-255; omit to leave unchanged"`
}

// WLEDSetTool returns the definition for moonraker_wled_set.
func WLEDSetTool() *mcp.Tool {
	return &mcp.Tool{
		Name:        "moonraker_wled_set",
		Description: "Set a WLED strip's preset, brightness, intensity, or speed (POST /machine/wled/strip).",
		Annotations: write("Set WLED State"),
	}
}

// WLED brightness, intensity, and speed bounds, validated client-side so an
// out-of-range value fails fast with a clear message.
const (
	wledValueMin  = 0
	wledValueMax  = 255
	brightnessMin = 1
)

// validateWLEDRanges checks the optional numeric fields against their documented
// ranges. Brightness uses 0 to mean "unset"; the pointer fields use nil.
func validateWLEDRanges(params WLEDSetParams) error {
	if params.Brightness != 0 {
		brightErr := requireRange("brightness", params.Brightness, brightnessMin, wledValueMax)
		if brightErr != nil {
			return brightErr
		}
	}

	if params.Intensity != nil {
		intensityErr := requireRange("intensity", *params.Intensity, wledValueMin, wledValueMax)
		if intensityErr != nil {
			return intensityErr
		}
	}

	if params.Speed != nil {
		speedErr := requireRange("speed", *params.Speed, wledValueMin, wledValueMax)
		if speedErr != nil {
			return speedErr
		}
	}

	return nil
}

// NewWLEDSetHandler creates the handler for moonraker_wled_set.
func NewWLEDSetHandler(api moonraker.API) mcp.ToolHandlerFor[WLEDSetParams, RawResult] {
	return func(ctx context.Context, _ *mcp.CallToolRequest, params WLEDSetParams) (*mcp.CallToolResult, RawResult, error) {
		valErr := requireString(paramStrip, params.Strip)
		if valErr != nil {
			return nil, RawResult{}, valErr
		}

		rangeErr := validateWLEDRanges(params)
		if rangeErr != nil {
			return nil, RawResult{}, rangeErr
		}

		query := url.Values{paramStrip: {params.Strip}}
		setIntPtr(query, "preset", params.Preset)
		setPositive(query, "brightness", params.Brightness)
		setIntPtr(query, "intensity", params.Intensity)
		setIntPtr(query, "speed", params.Speed)

		out, err := decodeRaw(api.Post(ctx, "/machine/wled/strip", query, nil))

		return nil, out, err
	}
}

// setPositive adds key=value to query when value is positive.
func setPositive(query url.Values, key string, value int) {
	if value > 0 {
		query.Set(key, strconv.Itoa(value))
	}
}

// setIntPtr adds key=value to query when value is non-nil, so an explicit 0 is
// sent rather than treated as unset.
func setIntPtr(query url.Values, key string, value *int) {
	if value != nil {
		query.Set(key, strconv.Itoa(*value))
	}
}
