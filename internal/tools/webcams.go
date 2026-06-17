package tools

import (
	"context"
	"maps"
	"net/url"

	"github.com/modelcontextprotocol/go-sdk/mcp"

	"github.com/lexfrei/mcp-raker/internal/moonraker"
)

// WebcamsListTool returns the definition for moonraker_webcams_list.
func WebcamsListTool() *mcp.Tool {
	return &mcp.Tool{
		Name:        "moonraker_webcams_list",
		Description: "List configured webcams (GET /server/webcams/list).",
		Annotations: readOnly("List Webcams"),
	}
}

// NewWebcamsListHandler creates the handler for moonraker_webcams_list.
func NewWebcamsListHandler(api moonraker.API) mcp.ToolHandlerFor[NoParams, RawResult] {
	return func(ctx context.Context, _ *mcp.CallToolRequest, _ NoParams) (*mcp.CallToolResult, RawResult, error) {
		out, err := decodeRaw(api.Get(ctx, "/server/webcams/list", nil))

		return nil, out, err
	}
}

// WebcamNameParams names a webcam configuration.
type WebcamNameParams struct {
	Name string `json:"name" jsonschema:"Name of the webcam configuration"`
}

// WebcamsGetTool returns the definition for moonraker_webcams_get.
func WebcamsGetTool() *mcp.Tool {
	return &mcp.Tool{
		Name:        "moonraker_webcams_get",
		Description: "Get a single webcam configuration (GET /server/webcams/item).",
		Annotations: readOnly("Get Webcam"),
	}
}

// NewWebcamsGetHandler creates the handler for moonraker_webcams_get.
func NewWebcamsGetHandler(api moonraker.API) mcp.ToolHandlerFor[WebcamNameParams, RawResult] {
	return func(ctx context.Context, _ *mcp.CallToolRequest, params WebcamNameParams) (*mcp.CallToolResult, RawResult, error) {
		valErr := requireString(paramName, params.Name)
		if valErr != nil {
			return nil, RawResult{}, valErr
		}

		out, err := decodeRaw(api.Get(ctx, "/server/webcams/item", url.Values{paramName: {params.Name}}))

		return nil, out, err
	}
}

// WebcamsAddParams defines the parameters for moonraker_webcams_add.
type WebcamsAddParams struct {
	Name     string         `json:"name"     jsonschema:"Name of the webcam to create or update"`
	Settings map[string]any `json:"settings" jsonschema:"Webcam fields such as stream_url, snapshot_url, rotation, flip_horizontal, flip_vertical, target_fps"`
}

// WebcamsAddTool returns the definition for moonraker_webcams_add.
func WebcamsAddTool() *mcp.Tool {
	return &mcp.Tool{
		Name:        "moonraker_webcams_add",
		Description: "Create or update a webcam configuration (POST /server/webcams/item).",
		Annotations: write("Add/Update Webcam"),
	}
}

// NewWebcamsAddHandler creates the handler for moonraker_webcams_add.
func NewWebcamsAddHandler(api moonraker.API) mcp.ToolHandlerFor[WebcamsAddParams, RawResult] {
	return func(ctx context.Context, _ *mcp.CallToolRequest, params WebcamsAddParams) (*mcp.CallToolResult, RawResult, error) {
		valErr := requireString(paramName, params.Name)
		if valErr != nil {
			return nil, RawResult{}, valErr
		}

		body := make(map[string]any, len(params.Settings)+1)
		maps.Copy(body, params.Settings)
		body[paramName] = params.Name

		out, err := decodeRaw(api.Post(ctx, "/server/webcams/item", nil, body))

		return nil, out, err
	}
}

// WebcamsDeleteTool returns the definition for moonraker_webcams_delete.
func WebcamsDeleteTool() *mcp.Tool {
	return &mcp.Tool{
		Name:        "moonraker_webcams_delete",
		Description: "Delete a webcam configuration (DELETE /server/webcams/item).",
		Annotations: writeDestructive("Delete Webcam"),
	}
}

// NewWebcamsDeleteHandler creates the handler for moonraker_webcams_delete.
func NewWebcamsDeleteHandler(api moonraker.API) mcp.ToolHandlerFor[WebcamNameParams, RawResult] {
	return func(ctx context.Context, _ *mcp.CallToolRequest, params WebcamNameParams) (*mcp.CallToolResult, RawResult, error) {
		valErr := requireString(paramName, params.Name)
		if valErr != nil {
			return nil, RawResult{}, valErr
		}

		out, err := decodeRaw(api.Delete(ctx, "/server/webcams/item", url.Values{paramName: {params.Name}}))

		return nil, out, err
	}
}

// WebcamsTestTool returns the definition for moonraker_webcams_test.
func WebcamsTestTool() *mcp.Tool {
	return &mcp.Tool{
		Name:        "moonraker_webcams_test",
		Description: "Test whether a webcam's URLs are reachable (POST /server/webcams/test).",
		Annotations: write("Test Webcam"),
	}
}

// NewWebcamsTestHandler creates the handler for moonraker_webcams_test.
func NewWebcamsTestHandler(api moonraker.API) mcp.ToolHandlerFor[WebcamNameParams, RawResult] {
	return func(ctx context.Context, _ *mcp.CallToolRequest, params WebcamNameParams) (*mcp.CallToolResult, RawResult, error) {
		valErr := requireString(paramName, params.Name)
		if valErr != nil {
			return nil, RawResult{}, valErr
		}

		out, err := decodeRaw(api.Post(ctx, "/server/webcams/test", url.Values{paramName: {params.Name}}, nil))

		return nil, out, err
	}
}
