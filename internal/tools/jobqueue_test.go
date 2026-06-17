package tools_test

import (
	"testing"

	"github.com/cockroachdb/errors"

	"github.com/lexfrei/mcp-raker/internal/tools"
)

func TestJobQueueStatus(t *testing.T) {
	t.Parallel()

	mock := &mockAPI{result: okJSON}

	_, _, err := tools.NewJobQueueStatusHandler(mock)(t.Context(), nil, tools.NoParams{})
	if err != nil {
		t.Fatalf("handler: %v", err)
	}

	assertCall(t, mock, methodGet, "/server/job_queue/status")
}

func TestJobQueueEnqueue(t *testing.T) {
	t.Parallel()

	_, _, missing := tools.NewJobQueueEnqueueHandler(&mockAPI{})(t.Context(), nil, tools.JobQueueEnqueueParams{})
	if !errors.Is(missing, tools.ErrValidation) {
		t.Errorf("missing filenames err = %v, want ErrValidation", missing)
	}

	mock := &mockAPI{result: okJSON}
	params := tools.JobQueueEnqueueParams{Filenames: []string{testGcodeFile}}

	_, _, err := tools.NewJobQueueEnqueueHandler(mock)(t.Context(), nil, params)
	if err != nil {
		t.Fatalf("handler: %v", err)
	}

	assertCall(t, mock, methodPost, "/server/job_queue/job")
}

func TestJobQueueRemove_RequiresSelection(t *testing.T) {
	t.Parallel()

	_, _, err := tools.NewJobQueueRemoveHandler(&mockAPI{})(t.Context(), nil, tools.JobQueueRemoveParams{})
	if !errors.Is(err, tools.ErrValidation) {
		t.Fatalf("err = %v, want ErrValidation", err)
	}
}

func TestJobQueueRemove_RejectsBoth(t *testing.T) {
	t.Parallel()

	params := tools.JobQueueRemoveParams{JobIDs: []string{"id1"}, All: true}

	_, _, err := tools.NewJobQueueRemoveHandler(&mockAPI{})(t.Context(), nil, params)
	if !errors.Is(err, tools.ErrValidation) {
		t.Fatalf("err = %v, want ErrValidation when job_ids and all are both set", err)
	}
}

func TestJobQueueRemove_ByIDs(t *testing.T) {
	t.Parallel()

	mock := &mockAPI{result: okJSON}
	params := tools.JobQueueRemoveParams{JobIDs: []string{"id1", "id2"}}

	_, _, err := tools.NewJobQueueRemoveHandler(mock)(t.Context(), nil, params)
	if err != nil {
		t.Fatalf("handler: %v", err)
	}

	assertCall(t, mock, methodDelete, "/server/job_queue/job")

	if len(mock.lastQuery["job_ids"]) != 2 {
		t.Errorf("job_ids = %v, want 2 entries", mock.lastQuery["job_ids"])
	}
}

func TestJobQueuePauseStart(t *testing.T) {
	t.Parallel()

	pause := &mockAPI{result: okJSON}

	_, _, pauseErr := tools.NewJobQueuePauseHandler(pause)(t.Context(), nil, tools.NoParams{})
	if pauseErr != nil {
		t.Fatalf("pause: %v", pauseErr)
	}

	assertCall(t, pause, methodPost, "/server/job_queue/pause")

	start := &mockAPI{result: okJSON}

	_, _, startErr := tools.NewJobQueueStartHandler(start)(t.Context(), nil, tools.NoParams{})
	if startErr != nil {
		t.Fatalf("start: %v", startErr)
	}

	assertCall(t, start, methodPost, "/server/job_queue/start")
}

func TestJobQueueJump_RequiresJobID(t *testing.T) {
	t.Parallel()

	_, _, err := tools.NewJobQueueJumpHandler(&mockAPI{})(t.Context(), nil, tools.JobQueueJumpParams{})
	if !errors.Is(err, tools.ErrValidation) {
		t.Fatalf("err = %v, want ErrValidation", err)
	}
}
