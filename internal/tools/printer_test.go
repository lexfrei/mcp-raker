package tools_test

import (
	"encoding/json"
	"testing"

	"github.com/cockroachdb/errors"

	"github.com/lexfrei/mcp-raker/internal/tools"
)

func TestPrinterInfo(t *testing.T) {
	t.Parallel()

	mock := &mockAPI{result: json.RawMessage(`{"state":"ready","hostname":"voron"}`)}
	handler := tools.NewPrinterInfoHandler(mock)

	_, out, err := handler(t.Context(), nil, tools.NoParams{})
	if err != nil {
		t.Fatalf("handler: %v", err)
	}

	if out.State != "ready" || out.Hostname != "voron" {
		t.Errorf("out = %+v, want state=ready hostname=voron", out)
	}

	if mock.lastMethod != methodGet || mock.lastPath != "/printer/info" {
		t.Errorf("call = %s %s, want GET /printer/info", mock.lastMethod, mock.lastPath)
	}
}

func TestObjectsQuery_RequiresObjects(t *testing.T) {
	t.Parallel()

	handler := tools.NewObjectsQueryHandler(&mockAPI{})

	_, _, err := handler(t.Context(), nil, tools.ObjectsQueryParams{})
	if !errors.Is(err, tools.ErrValidation) {
		t.Fatalf("err = %v, want ErrValidation", err)
	}
}

func TestObjectsQuery_BuildsQuery(t *testing.T) {
	t.Parallel()

	mock := &mockAPI{result: json.RawMessage(`{"eventtime":1.0,"status":{"extruder":{"temperature":24.5}}}`)}
	handler := tools.NewObjectsQueryHandler(mock)

	params := tools.ObjectsQueryParams{Objects: map[string][]string{
		"extruder":   {"temperature", "target"},
		"heater_bed": {},
	}}

	_, out, err := handler(t.Context(), nil, params)
	if err != nil {
		t.Fatalf("handler: %v", err)
	}

	if mock.lastPath != "/printer/objects/query" {
		t.Errorf("path = %q, want /printer/objects/query", mock.lastPath)
	}

	if got := mock.lastQuery.Get("extruder"); got != "temperature,target" {
		t.Errorf("extruder query = %q, want temperature,target", got)
	}

	if _, ok := mock.lastQuery["heater_bed"]; !ok {
		t.Error("heater_bed not present in query")
	}

	if out.Status["extruder"] == nil {
		t.Error("status.extruder missing from result")
	}
}

func TestGcodeScript_RequiresScript(t *testing.T) {
	t.Parallel()

	handler := tools.NewGcodeScriptHandler(&mockAPI{})

	_, _, err := handler(t.Context(), nil, tools.GcodeScriptParams{})
	if !errors.Is(err, tools.ErrValidation) {
		t.Fatalf("err = %v, want ErrValidation", err)
	}
}

func TestGcodeScript_SendsScriptQuery(t *testing.T) {
	t.Parallel()

	mock := &mockAPI{result: json.RawMessage(`"ok"`)}
	handler := tools.NewGcodeScriptHandler(mock)

	_, _, err := handler(t.Context(), nil, tools.GcodeScriptParams{Script: gcodeG28})
	if err != nil {
		t.Fatalf("handler: %v", err)
	}

	if mock.lastMethod != methodPost || mock.lastPath != "/printer/gcode/script" {
		t.Errorf("call = %s %s, want POST /printer/gcode/script", mock.lastMethod, mock.lastPath)
	}

	if got := mock.lastQuery.Get("script"); got != gcodeG28 {
		t.Errorf("script query = %q, want %s", got, gcodeG28)
	}
}

func TestEmergencyStop(t *testing.T) {
	t.Parallel()

	mock := &mockAPI{result: json.RawMessage(`"ok"`)}
	handler := tools.NewEmergencyStopHandler(mock)

	_, _, err := handler(t.Context(), nil, tools.NoParams{})
	if err != nil {
		t.Fatalf("handler: %v", err)
	}

	if mock.lastMethod != methodPost || mock.lastPath != "/printer/emergency_stop" {
		t.Errorf("call = %s %s, want POST /printer/emergency_stop", mock.lastMethod, mock.lastPath)
	}
}

func TestMoonrakerError_Wrapped(t *testing.T) {
	t.Parallel()

	mock := &mockAPI{err: errStub}
	handler := tools.NewPrinterInfoHandler(mock)

	_, _, err := handler(t.Context(), nil, tools.NoParams{})
	if !errors.Is(err, tools.ErrMoonraker) {
		t.Fatalf("err = %v, want ErrMoonraker", err)
	}
}
