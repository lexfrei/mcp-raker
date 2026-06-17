package tools_test

import (
	"testing"

	"github.com/cockroachdb/errors"

	"github.com/lexfrei/mcp-raker/internal/tools"
)

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
