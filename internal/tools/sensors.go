package tools

import (
	"context"
	"net/url"

	"github.com/cockroachdb/errors"
	"github.com/modelcontextprotocol/go-sdk/mcp"

	"github.com/lexfrei/mcp-raker/internal/moonraker"
)

// SensorsListParams defines the parameters for moonraker_sensors_list.
type SensorsListParams struct {
	Extended bool `json:"extended,omitempty" jsonschema:"When true, include each sensor's configuration and parameters"`
}

// SensorsListTool returns the definition for moonraker_sensors_list.
func SensorsListTool() *mcp.Tool {
	return &mcp.Tool{
		Name:        "moonraker_sensors_list",
		Description: "List configured sensors and their latest values (GET /machine/sensors/list).",
		Annotations: readOnly("List Sensors"),
	}
}

// NewSensorsListHandler creates the handler for moonraker_sensors_list.
func NewSensorsListHandler(api moonraker.API) mcp.ToolHandlerFor[SensorsListParams, map[string]any] {
	return func(ctx context.Context, _ *mcp.CallToolRequest, params SensorsListParams) (*mcp.CallToolResult, map[string]any, error) {
		query := url.Values{}
		if params.Extended {
			query.Set("extended", "true")
		}

		raw, reqErr := api.Get(ctx, "/machine/sensors/list", query)

		// Moonraker replies 404 when no [sensor] section is configured. Treat that
		// as "no sensors" rather than an error, so the common case reads cleanly.
		if errors.Is(reqErr, moonraker.ErrNotFound) {
			return nil, map[string]any{"sensors": map[string]any{}}, nil
		}

		out, err := decodeResult(raw, reqErr)

		return nil, out, err
	}
}

// SensorParams names a single sensor.
type SensorParams struct {
	Sensor string `json:"sensor" jsonschema:"Name of the sensor as configured in Moonraker"`
}

// SensorsInfoTool returns the definition for moonraker_sensors_info.
func SensorsInfoTool() *mcp.Tool {
	return &mcp.Tool{
		Name:        "moonraker_sensors_info",
		Description: "Get a single sensor's configuration and latest values (GET /machine/sensors/info).",
		Annotations: readOnly("Sensor Info"),
	}
}

// NewSensorsInfoHandler creates the handler for moonraker_sensors_info.
func NewSensorsInfoHandler(api moonraker.API) mcp.ToolHandlerFor[SensorParams, map[string]any] {
	return func(ctx context.Context, _ *mcp.CallToolRequest, params SensorParams) (*mcp.CallToolResult, map[string]any, error) {
		valErr := requireString(paramSensor, params.Sensor)
		if valErr != nil {
			return nil, map[string]any{}, valErr
		}

		out, err := decodeResult(api.Get(ctx, "/machine/sensors/info", url.Values{paramSensor: {params.Sensor}}))

		return nil, out, err
	}
}

// SensorsMeasurementsParams defines the parameters for moonraker_sensors_measurements.
type SensorsMeasurementsParams struct {
	Sensor string `json:"sensor,omitempty" jsonschema:"Name of a single sensor; omit to return measurements for all sensors"`
}

// SensorsMeasurementsTool returns the definition for moonraker_sensors_measurements.
func SensorsMeasurementsTool() *mcp.Tool {
	return &mcp.Tool{
		Name:        "moonraker_sensors_measurements",
		Description: "Get stored measurement history for one or all sensors (GET /machine/sensors/measurements).",
		Annotations: readOnly("Sensor Measurements"),
	}
}

// NewSensorsMeasurementsHandler creates the handler for moonraker_sensors_measurements.
func NewSensorsMeasurementsHandler(api moonraker.API) mcp.ToolHandlerFor[SensorsMeasurementsParams, map[string]any] {
	return func(ctx context.Context, _ *mcp.CallToolRequest, params SensorsMeasurementsParams) (*mcp.CallToolResult, map[string]any, error) {
		query := url.Values{}
		if params.Sensor != "" {
			query.Set(paramSensor, params.Sensor)
		}

		raw, reqErr := api.Get(ctx, "/machine/sensors/measurements", query)

		// Like sensors_list, the endpoint 404s when no [sensor] section exists.
		// Measurements are keyed by sensor name, so "none configured" is {}.
		if errors.Is(reqErr, moonraker.ErrNotFound) {
			return nil, map[string]any{}, nil
		}

		out, err := decodeResult(raw, reqErr)

		return nil, out, err
	}
}
