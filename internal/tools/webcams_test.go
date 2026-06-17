package tools_test

import (
	"testing"

	"github.com/cockroachdb/errors"

	"github.com/lexfrei/mcp-raker/internal/tools"
)

const testCam = "cam1"

func TestWebcamsList(t *testing.T) {
	t.Parallel()

	mock := &mockAPI{result: okJSON}

	_, _, err := tools.NewWebcamsListHandler(mock)(t.Context(), nil, tools.NoParams{})
	if err != nil {
		t.Fatalf("handler: %v", err)
	}

	assertCall(t, mock, methodGet, "/server/webcams/list")
}

func TestWebcamsGet_RequiresName(t *testing.T) {
	t.Parallel()

	_, _, err := tools.NewWebcamsGetHandler(&mockAPI{})(t.Context(), nil, tools.WebcamNameParams{})
	if !errors.Is(err, tools.ErrValidation) {
		t.Fatalf("err = %v, want ErrValidation", err)
	}
}

func TestWebcamsAdd_MergesSettings(t *testing.T) {
	t.Parallel()

	mock := &mockAPI{result: okJSON}
	params := tools.WebcamsAddParams{
		Name:     testCam,
		Settings: map[string]any{"stream_url": "/webcam/stream"},
	}

	_, _, err := tools.NewWebcamsAddHandler(mock)(t.Context(), nil, params)
	if err != nil {
		t.Fatalf("handler: %v", err)
	}

	assertCall(t, mock, methodPost, "/server/webcams/item")

	body, ok := mock.lastBody.(map[string]any)
	if !ok {
		t.Fatalf("body type = %T, want map[string]any", mock.lastBody)
	}

	if body["name"] != testCam || body["stream_url"] != "/webcam/stream" {
		t.Errorf("body = %v, want name and stream_url merged", body)
	}
}

func TestWebcamsDeleteAndTest(t *testing.T) {
	t.Parallel()

	del := &mockAPI{result: okJSON}

	_, _, delErr := tools.NewWebcamsDeleteHandler(del)(t.Context(), nil, tools.WebcamNameParams{Name: testCam})
	if delErr != nil {
		t.Fatalf("delete: %v", delErr)
	}

	assertCall(t, del, methodDelete, "/server/webcams/item")

	test := &mockAPI{result: okJSON}

	_, _, testErr := tools.NewWebcamsTestHandler(test)(t.Context(), nil, tools.WebcamNameParams{Name: testCam})
	if testErr != nil {
		t.Fatalf("test: %v", testErr)
	}

	assertCall(t, test, methodPost, "/server/webcams/test")
}
