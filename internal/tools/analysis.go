package tools

import (
	"context"

	"github.com/modelcontextprotocol/go-sdk/mcp"

	"github.com/lexfrei/mcp-raker/internal/moonraker"
)

// AnalysisStatusTool returns the definition for moonraker_analysis_status.
func AnalysisStatusTool() *mcp.Tool {
	return &mcp.Tool{
		Name:        "moonraker_analysis_status",
		Description: "Report whether the Klipper estimator is available and configured (GET /server/analysis/status).",
		Annotations: readOnly("Analysis Status"),
	}
}

// NewAnalysisStatusHandler creates the handler for moonraker_analysis_status.
func NewAnalysisStatusHandler(api moonraker.API) mcp.ToolHandlerFor[NoParams, map[string]any] {
	return func(ctx context.Context, _ *mcp.CallToolRequest, _ NoParams) (*mcp.CallToolResult, map[string]any, error) {
		out, err := decodeResult(api.Get(ctx, "/server/analysis/status", nil))

		return nil, out, err
	}
}

// AnalysisEstimateParams defines the parameters for moonraker_analysis_estimate.
type AnalysisEstimateParams struct {
	Filename        string `json:"filename"                   jsonschema:"Gcode file to estimate, relative to the gcodes root"`
	EstimatorConfig string `json:"estimator_config,omitempty" jsonschema:"Optional estimator configuration name to use"`
	UpdateMetadata  bool   `json:"update_metadata,omitempty"  jsonschema:"When true, write the estimate back into the file's metadata"`
}

// AnalysisEstimateTool returns the definition for moonraker_analysis_estimate.
func AnalysisEstimateTool() *mcp.Tool {
	return &mcp.Tool{
		Name:        "moonraker_analysis_estimate",
		Description: "Estimate the print duration of a gcode file (POST /server/analysis/estimate).",
		Annotations: write("Estimate Print"),
	}
}

// NewAnalysisEstimateHandler creates the handler for moonraker_analysis_estimate.
func NewAnalysisEstimateHandler(api moonraker.API) mcp.ToolHandlerFor[AnalysisEstimateParams, map[string]any] {
	return func(ctx context.Context, _ *mcp.CallToolRequest, params AnalysisEstimateParams) (*mcp.CallToolResult, map[string]any, error) {
		valErr := requireString(paramFilename, params.Filename)
		if valErr != nil {
			return nil, map[string]any{}, valErr
		}

		body := map[string]any{paramFilename: params.Filename, "update_metadata": params.UpdateMetadata}
		if params.EstimatorConfig != "" {
			body["estimator_config"] = params.EstimatorConfig
		}

		out, err := decodeResult(api.Post(ctx, "/server/analysis/estimate", nil, body))

		return nil, out, err
	}
}

// AnalysisProcessTool returns the definition for moonraker_analysis_process.
func AnalysisProcessTool() *mcp.Tool {
	return &mcp.Tool{
		Name:        "moonraker_analysis_process",
		Description: "Post-process a gcode file, injecting estimated time data (POST /server/analysis/process).",
		Annotations: write("Post-Process G-code"),
	}
}

// NewAnalysisProcessHandler creates the handler for moonraker_analysis_process.
func NewAnalysisProcessHandler(api moonraker.API) mcp.ToolHandlerFor[FilenameParams, map[string]any] {
	return func(ctx context.Context, _ *mcp.CallToolRequest, params FilenameParams) (*mcp.CallToolResult, map[string]any, error) {
		valErr := requireString(paramFilename, params.Filename)
		if valErr != nil {
			return nil, map[string]any{}, valErr
		}

		out, err := decodeResult(api.Post(ctx, "/server/analysis/process", nil, map[string]any{paramFilename: params.Filename}))

		return nil, out, err
	}
}

// AnalysisDumpConfigParams defines the parameters for moonraker_analysis_dump_config.
type AnalysisDumpConfigParams struct {
	DestConfig string `json:"dest_config,omitempty" jsonschema:"Optional destination path to write the estimator configuration to"`
}

// AnalysisDumpConfigTool returns the definition for moonraker_analysis_dump_config.
func AnalysisDumpConfigTool() *mcp.Tool {
	return &mcp.Tool{
		Name:        "moonraker_analysis_dump_config",
		Description: "Export the active Klipper estimator configuration (POST /server/analysis/dump_config).",
		Annotations: write("Dump Estimator Config"),
	}
}

// NewAnalysisDumpConfigHandler creates the handler for moonraker_analysis_dump_config.
func NewAnalysisDumpConfigHandler(api moonraker.API) mcp.ToolHandlerFor[AnalysisDumpConfigParams, map[string]any] {
	return func(ctx context.Context, _ *mcp.CallToolRequest, params AnalysisDumpConfigParams) (*mcp.CallToolResult, map[string]any, error) {
		var body map[string]any
		if params.DestConfig != "" {
			body = map[string]any{"dest_config": params.DestConfig}
		}

		out, err := decodeResult(api.Post(ctx, "/server/analysis/dump_config", nil, body))

		return nil, out, err
	}
}
