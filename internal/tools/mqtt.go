package tools

import (
	"context"

	"github.com/modelcontextprotocol/go-sdk/mcp"

	"github.com/lexfrei/mcp-raker/internal/moonraker"
)

// mqttTopic is the body key naming the MQTT topic.
const mqttTopic = "topic"

// MQTT quality-of-service bounds, validated client-side.
const (
	qosMin = 0
	qosMax = 2
)

// MQTTPublishParams defines the parameters for moonraker_mqtt_publish.
type MQTTPublishParams struct {
	Topic   string `json:"topic"            jsonschema:"MQTT topic to publish to"`
	Payload any    `json:"payload"          jsonschema:"Payload to publish; may be any JSON type, including null or \"\" to clear a retained message"`
	QOS     int    `json:"qos,omitempty"    jsonschema:"MQTT quality-of-service level 0-2"`
	Retain  bool   `json:"retain,omitempty" jsonschema:"When true, the broker retains the message"`
}

// MQTTPublishTool returns the definition for moonraker_mqtt_publish.
func MQTTPublishTool() *mcp.Tool {
	return &mcp.Tool{
		Name:        "moonraker_mqtt_publish",
		Description: "Publish a payload to an MQTT topic (POST /server/mqtt/publish).",
		Annotations: write("MQTT Publish"),
	}
}

// NewMQTTPublishHandler creates the handler for moonraker_mqtt_publish.
func NewMQTTPublishHandler(api moonraker.API) mcp.ToolHandlerFor[MQTTPublishParams, map[string]any] {
	return func(ctx context.Context, _ *mcp.CallToolRequest, params MQTTPublishParams) (*mcp.CallToolResult, map[string]any, error) {
		valErr := requireString(mqttTopic, params.Topic)
		if valErr != nil {
			return nil, map[string]any{}, valErr
		}

		qosErr := requireRange("qos", params.QOS, qosMin, qosMax)
		if qosErr != nil {
			return nil, map[string]any{}, qosErr
		}

		body := map[string]any{
			mqttTopic: params.Topic,
			"payload": params.Payload,
			"qos":     params.QOS,
			"retain":  params.Retain,
		}

		out, err := decodeResult(api.Post(ctx, "/server/mqtt/publish", nil, body))

		return nil, out, err
	}
}

// MQTTSubscribeParams defines the parameters for moonraker_mqtt_subscribe.
type MQTTSubscribeParams struct {
	Topic   string  `json:"topic"             jsonschema:"MQTT topic to read a single value from"`
	QOS     int     `json:"qos,omitempty"     jsonschema:"MQTT quality-of-service level 0-2"`
	Timeout float64 `json:"timeout,omitempty" jsonschema:"Seconds to wait for a message before giving up"`
}

// MQTTSubscribeTool returns the definition for moonraker_mqtt_subscribe.
func MQTTSubscribeTool() *mcp.Tool {
	return &mcp.Tool{
		Name:        "moonraker_mqtt_subscribe",
		Description: "Read the next value published to an MQTT topic (POST /server/mqtt/subscribe).",
		Annotations: readOnly("MQTT Subscribe"),
	}
}

// NewMQTTSubscribeHandler creates the handler for moonraker_mqtt_subscribe.
// An MQTT payload can be any JSON value (a scalar, array, or object), so it is
// returned verbatim under a "payload" key. The wrapper keeps the structured
// content an object (required by MCP) while a scalar payload is not collapsed to
// an acknowledgement.
func NewMQTTSubscribeHandler(api moonraker.API) mcp.ToolHandlerFor[MQTTSubscribeParams, map[string]any] {
	return func(ctx context.Context, _ *mcp.CallToolRequest, params MQTTSubscribeParams) (*mcp.CallToolResult, map[string]any, error) {
		valErr := requireString(mqttTopic, params.Topic)
		if valErr != nil {
			return nil, nil, valErr
		}

		qosErr := requireRange("qos", params.QOS, qosMin, qosMax)
		if qosErr != nil {
			return nil, nil, qosErr
		}

		body := map[string]any{mqttTopic: params.Topic, "qos": params.QOS, "timeout": params.Timeout}

		value, err := decodePassthrough(api.Post(ctx, "/server/mqtt/subscribe", nil, body))
		if err != nil {
			return nil, nil, err
		}

		return nil, map[string]any{"payload": value}, nil
	}
}
