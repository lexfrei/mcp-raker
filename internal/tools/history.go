package tools

import (
	"context"
	"net/url"
	"strconv"

	"github.com/modelcontextprotocol/go-sdk/mcp"

	"github.com/lexfrei/mcp-raker/internal/moonraker"
)

// HistoryListParams defines the parameters for moonraker_history_list.
type HistoryListParams struct {
	Limit  int     `json:"limit"  jsonschema:"Maximum number of jobs to return"`
	Start  int     `json:"start"  jsonschema:"Number of jobs to skip from the start"`
	Before float64 `json:"before" jsonschema:"Only include jobs that ended before this Unix timestamp"`
	Since  float64 `json:"since"  jsonschema:"Only include jobs that started after this Unix timestamp"`
	Order  string  `json:"order"  jsonschema:"Sort order: 'asc' or 'desc'"`
}

// HistoryListTool returns the definition for moonraker_history_list.
func HistoryListTool() *mcp.Tool {
	return &mcp.Tool{
		Name:        "moonraker_history_list",
		Description: "List recorded print jobs with optional filtering and paging (GET /server/history/list).",
		Annotations: readOnly("List Print History"),
	}
}

// NewHistoryListHandler creates the handler for moonraker_history_list.
func NewHistoryListHandler(api moonraker.API) mcp.ToolHandlerFor[HistoryListParams, RawResult] {
	return func(ctx context.Context, _ *mcp.CallToolRequest, params HistoryListParams) (*mcp.CallToolResult, RawResult, error) {
		query := url.Values{}
		if params.Limit > 0 {
			query.Set("limit", strconv.Itoa(params.Limit))
		}

		if params.Start > 0 {
			query.Set("start", strconv.Itoa(params.Start))
		}

		if params.Before > 0 {
			query.Set("before", strconv.FormatFloat(params.Before, 'f', -1, 64))
		}

		if params.Since > 0 {
			query.Set("since", strconv.FormatFloat(params.Since, 'f', -1, 64))
		}

		if params.Order != "" {
			query.Set("order", params.Order)
		}

		out, err := decodeRaw(api.Get(ctx, "/server/history/list", query))

		return nil, out, err
	}
}

// HistoryTotalsTool returns the definition for moonraker_history_totals.
func HistoryTotalsTool() *mcp.Tool {
	return &mcp.Tool{
		Name:        "moonraker_history_totals",
		Description: "Get aggregated print statistics: total jobs, time, and filament used (GET /server/history/totals).",
		Annotations: readOnly("History Totals"),
	}
}

// NewHistoryTotalsHandler creates the handler for moonraker_history_totals.
func NewHistoryTotalsHandler(api moonraker.API) mcp.ToolHandlerFor[NoParams, RawResult] {
	return func(ctx context.Context, _ *mcp.CallToolRequest, _ NoParams) (*mcp.CallToolResult, RawResult, error) {
		out, err := decodeRaw(api.Get(ctx, "/server/history/totals", nil))

		return nil, out, err
	}
}

// HistoryJobParams names a single recorded job by its unique id.
type HistoryJobParams struct {
	UID string `json:"uid" jsonschema:"Unique identifier of the history job"`
}

// HistoryJobTool returns the definition for moonraker_history_job.
func HistoryJobTool() *mcp.Tool {
	return &mcp.Tool{
		Name:        "moonraker_history_job",
		Description: "Get a single recorded print job by its unique id (GET /server/history/job).",
		Annotations: readOnly("History Job"),
	}
}

// NewHistoryJobHandler creates the handler for moonraker_history_job.
func NewHistoryJobHandler(api moonraker.API) mcp.ToolHandlerFor[HistoryJobParams, RawResult] {
	return func(ctx context.Context, _ *mcp.CallToolRequest, params HistoryJobParams) (*mcp.CallToolResult, RawResult, error) {
		valErr := requireString(paramUID, params.UID)
		if valErr != nil {
			return nil, RawResult{}, valErr
		}

		out, err := decodeRaw(api.Get(ctx, "/server/history/job", url.Values{paramUID: {params.UID}}))

		return nil, out, err
	}
}

// HistoryResetTotalsTool returns the definition for moonraker_history_reset_totals.
func HistoryResetTotalsTool() *mcp.Tool {
	return &mcp.Tool{
		Name:        "moonraker_history_reset_totals",
		Description: "Reset the aggregated print statistics to zero (POST /server/history/reset_totals).",
		Annotations: writeDestructive("Reset History Totals"),
	}
}

// NewHistoryResetTotalsHandler creates the handler for moonraker_history_reset_totals.
func NewHistoryResetTotalsHandler(api moonraker.API) mcp.ToolHandlerFor[NoParams, RawResult] {
	return func(ctx context.Context, _ *mcp.CallToolRequest, _ NoParams) (*mcp.CallToolResult, RawResult, error) {
		out, err := decodeRaw(api.Post(ctx, "/server/history/reset_totals", nil, nil))

		return nil, out, err
	}
}

// HistoryDeleteJobParams defines the parameters for moonraker_history_delete_job.
type HistoryDeleteJobParams struct {
	UID string `json:"uid" jsonschema:"Unique id of the job to delete"`
	All bool   `json:"all" jsonschema:"When true, delete every recorded job instead of a single uid"`
}

// HistoryDeleteJobTool returns the definition for moonraker_history_delete_job.
func HistoryDeleteJobTool() *mcp.Tool {
	return &mcp.Tool{
		Name:        "moonraker_history_delete_job",
		Description: "Delete one recorded job by uid, or all jobs (DELETE /server/history/job).",
		Annotations: writeDestructive("Delete History Job"),
	}
}

// NewHistoryDeleteJobHandler creates the handler for moonraker_history_delete_job.
func NewHistoryDeleteJobHandler(api moonraker.API) mcp.ToolHandlerFor[HistoryDeleteJobParams, RawResult] {
	return func(ctx context.Context, _ *mcp.CallToolRequest, params HistoryDeleteJobParams) (*mcp.CallToolResult, RawResult, error) {
		if params.UID != "" && params.All {
			return nil, RawResult{}, mutuallyExclusive(paramUID, "all")
		}

		if params.UID == "" && !params.All {
			return nil, RawResult{}, requireString(paramUID, params.UID)
		}

		query := url.Values{}
		if params.All {
			query.Set("all", "true")
		} else {
			query.Set(paramUID, params.UID)
		}

		out, err := decodeRaw(api.Delete(ctx, "/server/history/job", query))

		return nil, out, err
	}
}
