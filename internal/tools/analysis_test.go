package tools_test

import (
	"testing"

	"github.com/cockroachdb/errors"

	"github.com/lexfrei/mcp-raker/internal/tools"
)

func TestAnalysisStatus(t *testing.T) {
	t.Parallel()

	mock := &mockAPI{result: okJSON}

	_, _, err := tools.NewAnalysisStatusHandler(mock)(t.Context(), nil, tools.NoParams{})
	if err != nil {
		t.Fatalf("handler: %v", err)
	}

	assertCall(t, mock, methodGet, "/server/analysis/status")
}

func TestAnalysisEstimate_RequiresFilename(t *testing.T) {
	t.Parallel()

	_, _, err := tools.NewAnalysisEstimateHandler(&mockAPI{})(t.Context(), nil, tools.AnalysisEstimateParams{})
	if !errors.Is(err, tools.ErrValidation) {
		t.Fatalf("err = %v, want ErrValidation", err)
	}
}

func TestAnalysisEstimate(t *testing.T) {
	t.Parallel()

	mock := &mockAPI{result: okJSON}

	_, _, err := tools.NewAnalysisEstimateHandler(mock)(t.Context(), nil, tools.AnalysisEstimateParams{Filename: testGcodeFile})
	if err != nil {
		t.Fatalf("handler: %v", err)
	}

	assertCall(t, mock, methodPost, "/server/analysis/estimate")
}

func TestAnalysisProcess_RequiresFilename(t *testing.T) {
	t.Parallel()

	_, _, err := tools.NewAnalysisProcessHandler(&mockAPI{})(t.Context(), nil, tools.FilenameParams{})
	if !errors.Is(err, tools.ErrValidation) {
		t.Fatalf("err = %v, want ErrValidation", err)
	}
}

func TestAnalysisDumpConfig(t *testing.T) {
	t.Parallel()

	mock := &mockAPI{result: okJSON}

	_, _, err := tools.NewAnalysisDumpConfigHandler(mock)(t.Context(), nil, tools.AnalysisDumpConfigParams{})
	if err != nil {
		t.Fatalf("handler: %v", err)
	}

	assertCall(t, mock, methodPost, "/server/analysis/dump_config")
}
