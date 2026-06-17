package tools

import (
	"context"
	"net/url"

	"github.com/modelcontextprotocol/go-sdk/mcp"

	"github.com/lexfrei/mcp-raker/internal/moonraker"
)

// PrintStartParams defines the parameters for moonraker_print_start.
type PrintStartParams struct {
	Filename string `json:"filename" jsonschema:"Path of the G-code file to print, relative to the gcodes root (e.g. 'benchy.gcode' or 'subdir/part.gcode')"`
}

// PrintStartTool returns the definition for moonraker_print_start.
func PrintStartTool() *mcp.Tool {
	return &mcp.Tool{
		Name:        "moonraker_print_start",
		Description: "Start printing a G-code file from the gcodes root (POST /printer/print/start).",
		Annotations: writeDestructive("Start Print"),
	}
}

// NewPrintStartHandler creates the handler for moonraker_print_start.
func NewPrintStartHandler(api moonraker.API) mcp.ToolHandlerFor[PrintStartParams, RawResult] {
	return func(
		ctx context.Context,
		_ *mcp.CallToolRequest,
		params PrintStartParams,
	) (*mcp.CallToolResult, RawResult, error) {
		valErr := requireString(paramFilename, params.Filename)
		if valErr != nil {
			return nil, RawResult{}, valErr
		}

		query := url.Values{paramFilename: {params.Filename}}

		out, err := decodeRaw(api.Post(ctx, "/printer/print/start", query, nil))

		return nil, out, err
	}
}

// PrintPauseTool returns the definition for moonraker_print_pause.
func PrintPauseTool() *mcp.Tool {
	return &mcp.Tool{
		Name:        "moonraker_print_pause",
		Description: "Pause the current print (POST /printer/print/pause).",
		Annotations: write("Pause Print"),
	}
}

// NewPrintPauseHandler creates the handler for moonraker_print_pause.
func NewPrintPauseHandler(api moonraker.API) mcp.ToolHandlerFor[NoParams, RawResult] {
	return func(
		ctx context.Context,
		_ *mcp.CallToolRequest,
		_ NoParams,
	) (*mcp.CallToolResult, RawResult, error) {
		out, err := decodeRaw(api.Post(ctx, "/printer/print/pause", nil, nil))

		return nil, out, err
	}
}

// PrintResumeTool returns the definition for moonraker_print_resume.
func PrintResumeTool() *mcp.Tool {
	return &mcp.Tool{
		Name:        "moonraker_print_resume",
		Description: "Resume a paused print (POST /printer/print/resume).",
		Annotations: write("Resume Print"),
	}
}

// NewPrintResumeHandler creates the handler for moonraker_print_resume.
func NewPrintResumeHandler(api moonraker.API) mcp.ToolHandlerFor[NoParams, RawResult] {
	return func(
		ctx context.Context,
		_ *mcp.CallToolRequest,
		_ NoParams,
	) (*mcp.CallToolResult, RawResult, error) {
		out, err := decodeRaw(api.Post(ctx, "/printer/print/resume", nil, nil))

		return nil, out, err
	}
}

// PrintCancelTool returns the definition for moonraker_print_cancel.
func PrintCancelTool() *mcp.Tool {
	return &mcp.Tool{
		Name:        "moonraker_print_cancel",
		Description: "Cancel the current print (POST /printer/print/cancel).",
		Annotations: writeDestructive("Cancel Print"),
	}
}

// NewPrintCancelHandler creates the handler for moonraker_print_cancel.
func NewPrintCancelHandler(api moonraker.API) mcp.ToolHandlerFor[NoParams, RawResult] {
	return func(
		ctx context.Context,
		_ *mcp.CallToolRequest,
		_ NoParams,
	) (*mcp.CallToolResult, RawResult, error) {
		out, err := decodeRaw(api.Post(ctx, "/printer/print/cancel", nil, nil))

		return nil, out, err
	}
}
