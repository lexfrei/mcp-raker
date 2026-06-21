package tools

import (
	"context"
	"net/url"

	"github.com/modelcontextprotocol/go-sdk/mcp"

	"github.com/lexfrei/mcp-raker/internal/moonraker"
)

// jobIDsKey is the query/body key for job-queue identifiers.
const jobIDsKey = "job_ids"

// JobQueueStatusTool returns the definition for moonraker_jobqueue_status.
func JobQueueStatusTool() *mcp.Tool {
	return &mcp.Tool{
		Name:        "moonraker_jobqueue_status",
		Description: "Get the job queue state and its pending jobs (GET /server/job_queue/status).",
		Annotations: readOnly("Job Queue Status"),
	}
}

// NewJobQueueStatusHandler creates the handler for moonraker_jobqueue_status.
func NewJobQueueStatusHandler(api moonraker.API) mcp.ToolHandlerFor[NoParams, map[string]any] {
	return func(ctx context.Context, _ *mcp.CallToolRequest, _ NoParams) (*mcp.CallToolResult, map[string]any, error) {
		out, err := decodeResult(api.Get(ctx, "/server/job_queue/status", nil))

		return nil, out, err
	}
}

// JobQueueEnqueueParams defines the parameters for moonraker_jobqueue_enqueue.
type JobQueueEnqueueParams struct {
	Filenames []string `json:"filenames"       jsonschema:"Gcode filenames (relative to the gcodes root) to append to the queue"`
	Reset     bool     `json:"reset,omitempty" jsonschema:"When true, clear the existing queue before adding these jobs"`
}

// JobQueueEnqueueTool returns the definition for moonraker_jobqueue_enqueue.
func JobQueueEnqueueTool() *mcp.Tool {
	return &mcp.Tool{
		Name:        "moonraker_jobqueue_enqueue",
		Description: "Add one or more gcode files to the job queue (POST /server/job_queue/job).",
		Annotations: write("Enqueue Jobs"),
	}
}

// NewJobQueueEnqueueHandler creates the handler for moonraker_jobqueue_enqueue.
func NewJobQueueEnqueueHandler(api moonraker.API) mcp.ToolHandlerFor[JobQueueEnqueueParams, map[string]any] {
	return func(ctx context.Context, _ *mcp.CallToolRequest, params JobQueueEnqueueParams) (*mcp.CallToolResult, map[string]any, error) {
		valErr := requirePresent("filenames", len(params.Filenames))
		if valErr != nil {
			return nil, map[string]any{}, valErr
		}

		body := map[string]any{"filenames": params.Filenames, "reset": params.Reset}

		out, err := decodeResult(api.Post(ctx, "/server/job_queue/job", nil, body))

		return nil, out, err
	}
}

// JobQueueRemoveParams defines the parameters for moonraker_jobqueue_remove.
type JobQueueRemoveParams struct {
	JobIDs []string `json:"job_ids,omitempty" jsonschema:"Queue job identifiers to remove"`
	All    bool     `json:"all,omitempty"     jsonschema:"When true, clear the entire queue instead of specific ids"`
}

// JobQueueRemoveTool returns the definition for moonraker_jobqueue_remove.
func JobQueueRemoveTool() *mcp.Tool {
	return &mcp.Tool{
		Name:        "moonraker_jobqueue_remove",
		Description: "Remove jobs from the queue by id, or clear it entirely (DELETE /server/job_queue/job).",
		Annotations: writeDestructive("Remove Queued Jobs"),
	}
}

// NewJobQueueRemoveHandler creates the handler for moonraker_jobqueue_remove.
func NewJobQueueRemoveHandler(api moonraker.API) mcp.ToolHandlerFor[JobQueueRemoveParams, map[string]any] {
	return func(ctx context.Context, _ *mcp.CallToolRequest, params JobQueueRemoveParams) (*mcp.CallToolResult, map[string]any, error) {
		if len(params.JobIDs) > 0 && params.All {
			return nil, map[string]any{}, mutuallyExclusive(jobIDsKey, "all")
		}

		if len(params.JobIDs) == 0 && !params.All {
			return nil, map[string]any{}, requirePresent(jobIDsKey, len(params.JobIDs))
		}

		query := url.Values{}
		if params.All {
			query.Set("all", "true")
		} else {
			for _, id := range params.JobIDs {
				query.Add(jobIDsKey, id)
			}
		}

		out, err := decodeResult(api.Delete(ctx, "/server/job_queue/job", query))

		return nil, out, err
	}
}

// JobQueuePauseTool returns the definition for moonraker_jobqueue_pause.
func JobQueuePauseTool() *mcp.Tool {
	return &mcp.Tool{
		Name:        "moonraker_jobqueue_pause",
		Description: "Pause processing of the job queue (POST /server/job_queue/pause).",
		Annotations: write("Pause Job Queue"),
	}
}

// NewJobQueuePauseHandler creates the handler for moonraker_jobqueue_pause.
func NewJobQueuePauseHandler(api moonraker.API) mcp.ToolHandlerFor[NoParams, map[string]any] {
	return func(ctx context.Context, _ *mcp.CallToolRequest, _ NoParams) (*mcp.CallToolResult, map[string]any, error) {
		out, err := decodeResult(api.Post(ctx, "/server/job_queue/pause", nil, nil))

		return nil, out, err
	}
}

// JobQueueStartTool returns the definition for moonraker_jobqueue_start.
func JobQueueStartTool() *mcp.Tool {
	return &mcp.Tool{
		Name:        "moonraker_jobqueue_start",
		Description: "Start or resume processing of the job queue (POST /server/job_queue/start).",
		Annotations: write("Start Job Queue"),
	}
}

// NewJobQueueStartHandler creates the handler for moonraker_jobqueue_start.
func NewJobQueueStartHandler(api moonraker.API) mcp.ToolHandlerFor[NoParams, map[string]any] {
	return func(ctx context.Context, _ *mcp.CallToolRequest, _ NoParams) (*mcp.CallToolResult, map[string]any, error) {
		out, err := decodeResult(api.Post(ctx, "/server/job_queue/start", nil, nil))

		return nil, out, err
	}
}

// JobQueueJumpParams names a job to move to the front of the queue.
type JobQueueJumpParams struct {
	JobID string `json:"job_id" jsonschema:"Identifier of the queued job to move to the front"`
}

// JobQueueJumpTool returns the definition for moonraker_jobqueue_jump.
func JobQueueJumpTool() *mcp.Tool {
	return &mcp.Tool{
		Name:        "moonraker_jobqueue_jump",
		Description: "Move a queued job to the front of the queue (POST /server/job_queue/jump).",
		Annotations: write("Jump Queue Job"),
	}
}

// NewJobQueueJumpHandler creates the handler for moonraker_jobqueue_jump.
func NewJobQueueJumpHandler(api moonraker.API) mcp.ToolHandlerFor[JobQueueJumpParams, map[string]any] {
	return func(ctx context.Context, _ *mcp.CallToolRequest, params JobQueueJumpParams) (*mcp.CallToolResult, map[string]any, error) {
		valErr := requireString(paramJobID, params.JobID)
		if valErr != nil {
			return nil, map[string]any{}, valErr
		}

		out, err := decodeResult(api.Post(ctx, "/server/job_queue/jump", url.Values{paramJobID: {params.JobID}}, nil))

		return nil, out, err
	}
}
