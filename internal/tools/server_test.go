package tools_test

import (
	"encoding/json"
	"testing"

	"github.com/lexfrei/mcp-raker/internal/tools"
)

func TestMCPVersion(t *testing.T) {
	t.Parallel()

	handler := tools.NewMCPVersionHandler("1.2.3", "abc123", "go1.26.4")

	_, out, err := handler(t.Context(), nil, tools.NoParams{})
	if err != nil {
		t.Fatalf("handler: %v", err)
	}

	if out.Version != "1.2.3" || out.Revision != "abc123" || out.GoVersion != "go1.26.4" {
		t.Errorf("out = %+v, want 1.2.3/abc123/go1.26.4", out)
	}
}

func TestServerInfo(t *testing.T) {
	t.Parallel()

	mock := &mockAPI{result: json.RawMessage(`{"klippy_state":"ready","moonraker_version":"v0.9.0"}`)}
	handler := tools.NewServerInfoHandler(mock)

	_, out, err := handler(t.Context(), nil, tools.NoParams{})
	if err != nil {
		t.Fatalf("handler: %v", err)
	}

	if out.KlippyState != "ready" {
		t.Errorf("KlippyState = %q, want ready", out.KlippyState)
	}

	assertCall(t, mock, methodGet, "/server/info")
}

func TestServerConfig(t *testing.T) {
	t.Parallel()

	mock := &mockAPI{result: okJSON}

	_, _, err := tools.NewServerConfigHandler(mock)(t.Context(), nil, tools.NoParams{})
	if err != nil {
		t.Fatalf("handler: %v", err)
	}

	assertCall(t, mock, methodGet, "/server/config")
}

func TestTemperatureStore_IncludeMonitors(t *testing.T) {
	t.Parallel()

	mock := &mockAPI{result: okJSON}
	params := tools.TemperatureStoreParams{IncludeMonitors: true}

	_, _, err := tools.NewTemperatureStoreHandler(mock)(t.Context(), nil, params)
	if err != nil {
		t.Fatalf("handler: %v", err)
	}

	assertCall(t, mock, methodGet, "/server/temperature_store")

	if mock.lastQuery.Get("include_monitors") != queryTrue {
		t.Errorf("include_monitors = %q, want true", mock.lastQuery.Get("include_monitors"))
	}
}

func TestGcodeStore_Count(t *testing.T) {
	t.Parallel()

	mock := &mockAPI{result: okJSON}

	_, _, err := tools.NewGcodeStoreHandler(mock)(t.Context(), nil, tools.GcodeStoreParams{Count: 25})
	if err != nil {
		t.Fatalf("handler: %v", err)
	}

	assertCall(t, mock, methodGet, "/server/gcode_store")

	if mock.lastQuery.Get("count") != "25" {
		t.Errorf("count = %q, want 25", mock.lastQuery.Get("count"))
	}
}

func TestLogsRollover(t *testing.T) {
	t.Parallel()

	mock := &mockAPI{result: okJSON}

	_, _, err := tools.NewLogsRolloverHandler(mock)(t.Context(), nil, tools.LogsRolloverParams{Application: svcKlipper})
	if err != nil {
		t.Fatalf("handler: %v", err)
	}

	assertCall(t, mock, methodPost, "/server/logs/rollover")

	if mock.lastQuery.Get("application") != svcKlipper {
		t.Errorf("application = %q, want %s", mock.lastQuery.Get("application"), svcKlipper)
	}
}

func TestServerRestart(t *testing.T) {
	t.Parallel()

	mock := &mockAPI{result: okJSON}

	_, _, err := tools.NewServerRestartHandler(mock)(t.Context(), nil, tools.NoParams{})
	if err != nil {
		t.Fatalf("handler: %v", err)
	}

	assertCall(t, mock, methodPost, "/server/restart")
}
