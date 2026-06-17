package tools_test

import (
	"testing"

	"github.com/cockroachdb/errors"

	"github.com/lexfrei/mcp-raker/internal/tools"
)

func TestPrintStart_RequiresFilename(t *testing.T) {
	t.Parallel()

	_, _, err := tools.NewPrintStartHandler(&mockAPI{})(t.Context(), nil, tools.PrintStartParams{})
	if !errors.Is(err, tools.ErrValidation) {
		t.Fatalf("err = %v, want ErrValidation", err)
	}
}

func TestPrintStart_SendsFilename(t *testing.T) {
	t.Parallel()

	mock := &mockAPI{result: okJSON}

	_, _, err := tools.NewPrintStartHandler(mock)(t.Context(), nil, tools.PrintStartParams{Filename: "benchy.gcode"})
	if err != nil {
		t.Fatalf("handler: %v", err)
	}

	assertCall(t, mock, methodPost, "/printer/print/start")

	if mock.lastQuery.Get("filename") != "benchy.gcode" {
		t.Errorf("filename = %q, want benchy.gcode", mock.lastQuery.Get("filename"))
	}
}

func TestPrintControls(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name    string
		handler func(*mockAPI) (any, error)
		path    string
	}{
		{
			name: "pause",
			path: "/printer/print/pause",
			handler: func(mock *mockAPI) (any, error) {
				_, out, err := tools.NewPrintPauseHandler(mock)(t.Context(), nil, tools.NoParams{})

				return out, err
			},
		},
		{
			name: "resume",
			path: "/printer/print/resume",
			handler: func(mock *mockAPI) (any, error) {
				_, out, err := tools.NewPrintResumeHandler(mock)(t.Context(), nil, tools.NoParams{})

				return out, err
			},
		},
		{
			name: "cancel",
			path: "/printer/print/cancel",
			handler: func(mock *mockAPI) (any, error) {
				_, out, err := tools.NewPrintCancelHandler(mock)(t.Context(), nil, tools.NoParams{})

				return out, err
			},
		},
	}

	for _, testCase := range tests {
		t.Run(testCase.name, func(t *testing.T) {
			t.Parallel()

			mock := &mockAPI{result: okJSON}

			_, err := testCase.handler(mock)
			if err != nil {
				t.Fatalf("handler: %v", err)
			}

			assertCall(t, mock, methodPost, testCase.path)
		})
	}
}
