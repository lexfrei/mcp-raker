package tools_test

import (
	"testing"

	"github.com/cockroachdb/errors"

	"github.com/lexfrei/mcp-raker/internal/tools"
)

func TestSpoolmanStatusAndGet(t *testing.T) {
	t.Parallel()

	status := &mockAPI{result: okJSON}

	_, _, statusErr := tools.NewSpoolmanStatusHandler(status)(t.Context(), nil, tools.NoParams{})
	if statusErr != nil {
		t.Fatalf("status: %v", statusErr)
	}

	assertCall(t, status, methodGet, "/server/spoolman/status")

	get := &mockAPI{result: okJSON}

	_, _, getErr := tools.NewSpoolmanGetSpoolHandler(get)(t.Context(), nil, tools.NoParams{})
	if getErr != nil {
		t.Fatalf("get: %v", getErr)
	}

	assertCall(t, get, methodGet, "/server/spoolman/spool_id")
}

func TestSpoolmanSetSpool_RequiresPositive(t *testing.T) {
	t.Parallel()

	_, _, err := tools.NewSpoolmanSetSpoolHandler(&mockAPI{})(t.Context(), nil, tools.SpoolmanSetSpoolParams{})
	if !errors.Is(err, tools.ErrValidation) {
		t.Fatalf("err = %v, want ErrValidation", err)
	}
}

func TestSpoolmanSetSpool(t *testing.T) {
	t.Parallel()

	mock := &mockAPI{result: okJSON}

	_, _, err := tools.NewSpoolmanSetSpoolHandler(mock)(t.Context(), nil, tools.SpoolmanSetSpoolParams{SpoolID: 7})
	if err != nil {
		t.Fatalf("handler: %v", err)
	}

	assertCall(t, mock, methodPost, "/server/spoolman/spool_id")
}

func TestSpoolmanProxy_RequiresMethodAndPath(t *testing.T) {
	t.Parallel()

	_, _, err := tools.NewSpoolmanProxyHandler(&mockAPI{})(t.Context(), nil, tools.SpoolmanProxyParams{Path: "/v1/spool"})
	if !errors.Is(err, tools.ErrValidation) {
		t.Fatalf("err = %v, want ErrValidation", err)
	}
}

func TestSpoolmanProxy_ForwardsMethod(t *testing.T) {
	t.Parallel()

	mock := &mockAPI{result: okJSON}
	params := tools.SpoolmanProxyParams{RequestMethod: "DELETE", Path: "/v1/spool/3"}

	_, _, err := tools.NewSpoolmanProxyHandler(mock)(t.Context(), nil, params)
	if err != nil {
		t.Fatalf("handler: %v", err)
	}

	assertCall(t, mock, methodPost, "/server/spoolman/proxy")

	body, ok := mock.lastBody.(map[string]any)
	if !ok {
		t.Fatalf("body type = %T, want map[string]any", mock.lastBody)
	}

	if body["request_method"] != "DELETE" {
		t.Errorf("request_method = %v, want DELETE forwarded verbatim", body["request_method"])
	}
}
