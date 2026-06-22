package tools_test

import (
	"encoding/json"
	"testing"

	"github.com/cockroachdb/errors"

	"github.com/lexfrei/mcp-raker/internal/tools"
)

// historyJobMetadata digs out the first job's metadata map from a history_list
// handler result, failing the test if the shape is unexpected.
func historyJobMetadata(t *testing.T, out map[string]any) map[string]any {
	t.Helper()

	jobs, ok := out["jobs"].([]any)
	if !ok || len(jobs) == 0 {
		t.Fatalf("jobs = %v, want at least one job", out["jobs"])
	}

	job, ok := jobs[0].(map[string]any)
	if !ok {
		t.Fatalf("job = %v, want an object", jobs[0])
	}

	metadata, ok := job["metadata"].(map[string]any)
	if !ok {
		t.Fatalf("metadata = %v, want an object", job["metadata"])
	}

	return metadata
}

// historyFixtureWithThumbnails is a one-job history result whose metadata carries
// thumbnails plus an ordinary field, so a test can tell stripping from wiping.
var historyFixtureWithThumbnails = json.RawMessage(
	`{"count":1,"jobs":[{"job_id":"1","metadata":{"slicer":"OrcaSlicer","thumbnails":[{"width":300}]}}]}`)

func TestHistoryList_DropsThumbnailsByDefault(t *testing.T) {
	t.Parallel()

	mock := &mockAPI{result: historyFixtureWithThumbnails}

	_, out, err := tools.NewHistoryListHandler(mock)(t.Context(), nil, tools.HistoryListParams{})
	if err != nil {
		t.Fatalf("handler: %v", err)
	}

	metadata := historyJobMetadata(t, out)

	if _, has := metadata["thumbnails"]; has {
		t.Error("thumbnails should be dropped by default")
	}

	if _, has := metadata["slicer"]; !has {
		t.Error("non-thumbnail metadata should be preserved")
	}
}

func TestHistoryList_KeepsThumbnailsWhenRequested(t *testing.T) {
	t.Parallel()

	mock := &mockAPI{result: historyFixtureWithThumbnails}

	params := tools.HistoryListParams{IncludeThumbnails: true}

	_, out, err := tools.NewHistoryListHandler(mock)(t.Context(), nil, params)
	if err != nil {
		t.Fatalf("handler: %v", err)
	}

	if _, has := historyJobMetadata(t, out)["thumbnails"]; !has {
		t.Error("thumbnails should be kept when include_thumbnails is set")
	}
}

func TestHistoryList_Filters(t *testing.T) {
	t.Parallel()

	mock := &mockAPI{result: okJSON}
	params := tools.HistoryListParams{Limit: 5, Order: "desc"}

	_, _, err := tools.NewHistoryListHandler(mock)(t.Context(), nil, params)
	if err != nil {
		t.Fatalf("handler: %v", err)
	}

	assertCall(t, mock, methodGet, "/server/history/list")

	if mock.lastQuery.Get("limit") != "5" || mock.lastQuery.Get("order") != "desc" {
		t.Errorf("query = %v, want limit=5 order=desc", mock.lastQuery)
	}
}

func TestHistoryTotals(t *testing.T) {
	t.Parallel()

	mock := &mockAPI{result: okJSON}

	_, _, err := tools.NewHistoryTotalsHandler(mock)(t.Context(), nil, tools.NoParams{})
	if err != nil {
		t.Fatalf("handler: %v", err)
	}

	assertCall(t, mock, methodGet, "/server/history/totals")
}

func TestHistoryJob_RequiresUID(t *testing.T) {
	t.Parallel()

	_, _, err := tools.NewHistoryJobHandler(&mockAPI{})(t.Context(), nil, tools.HistoryJobParams{})
	if !errors.Is(err, tools.ErrValidation) {
		t.Fatalf("err = %v, want ErrValidation", err)
	}
}

func TestHistoryResetTotals(t *testing.T) {
	t.Parallel()

	mock := &mockAPI{result: okJSON}

	_, _, err := tools.NewHistoryResetTotalsHandler(mock)(t.Context(), nil, tools.NoParams{})
	if err != nil {
		t.Fatalf("handler: %v", err)
	}

	assertCall(t, mock, methodPost, "/server/history/reset_totals")
}

func TestHistoryDeleteJob_RequiresUIDorAll(t *testing.T) {
	t.Parallel()

	_, _, err := tools.NewHistoryDeleteJobHandler(&mockAPI{})(t.Context(), nil, tools.HistoryDeleteJobParams{})
	if !errors.Is(err, tools.ErrValidation) {
		t.Fatalf("err = %v, want ErrValidation", err)
	}
}

func TestHistoryDeleteJob_RejectsBoth(t *testing.T) {
	t.Parallel()

	params := tools.HistoryDeleteJobParams{UID: "5", All: true}

	_, _, err := tools.NewHistoryDeleteJobHandler(&mockAPI{})(t.Context(), nil, params)
	if !errors.Is(err, tools.ErrValidation) {
		t.Fatalf("err = %v, want ErrValidation when uid and all are both set", err)
	}
}

func TestHistoryDeleteJob_All(t *testing.T) {
	t.Parallel()

	mock := &mockAPI{result: okJSON}

	_, _, err := tools.NewHistoryDeleteJobHandler(mock)(t.Context(), nil, tools.HistoryDeleteJobParams{All: true})
	if err != nil {
		t.Fatalf("handler: %v", err)
	}

	assertCall(t, mock, methodDelete, "/server/history/job")

	if mock.lastQuery.Get("all") != queryTrue {
		t.Errorf("all = %q, want true", mock.lastQuery.Get("all"))
	}
}
