package tools_test

import (
	"encoding/json"
	"testing"

	"github.com/lexfrei/mcp-raker/internal/tools"
)

// TestContract_DataReadsPreserveTopLevel feeds every object-returning data-read
// tool a representative object fixture and asserts the payload's fields survive
// at the top level with no "result" envelope. Unlike the per-area call tests
// (which use a scalar "ok" fixture and so would pass even if a handler silently
// dropped a real payload), this pins the actual output contract across the whole
// data-read surface: re-typing any handler back to a collapsing path fails here.
func TestContract_DataReadsPreserveTopLevel(t *testing.T) {
	t.Parallel()

	const probe = "probe"

	fixture := json.RawMessage(`{"probe":42}`)

	tests := []struct {
		name string
		call func(*mockAPI) (map[string]any, error)
	}{
		{"system_info", func(m *mockAPI) (map[string]any, error) {
			_, o, e := tools.NewSystemInfoHandler(m)(t.Context(), nil, tools.NoParams{})

			return o, e
		}},
		{"proc_stats", func(m *mockAPI) (map[string]any, error) {
			_, o, e := tools.NewProcStatsHandler(m)(t.Context(), nil, tools.ProcStatsParams{})

			return o, e
		}},
		{"sudo_info", func(m *mockAPI) (map[string]any, error) {
			_, o, e := tools.NewSudoInfoHandler(m)(t.Context(), nil, tools.NoParams{})

			return o, e
		}},
		{"peripherals_usb", func(m *mockAPI) (map[string]any, error) {
			_, o, e := tools.NewPeripheralsUSBHandler(m)(t.Context(), nil, tools.NoParams{})

			return o, e
		}},
		{"peripherals_canbus", func(m *mockAPI) (map[string]any, error) {
			_, o, e := tools.NewPeripheralsCanbusHandler(m)(t.Context(), nil, tools.PeripheralsCanbusParams{})

			return o, e
		}},
		{"update_status", func(m *mockAPI) (map[string]any, error) {
			_, o, e := tools.NewUpdateStatusHandler(m)(t.Context(), nil, tools.NoParams{})

			return o, e
		}},
		{"power_devices", func(m *mockAPI) (map[string]any, error) {
			_, o, e := tools.NewPowerDevicesHandler(m)(t.Context(), nil, tools.NoParams{})

			return o, e
		}},
		{"server_config", func(m *mockAPI) (map[string]any, error) {
			_, o, e := tools.NewServerConfigHandler(m)(t.Context(), nil, tools.NoParams{})

			return o, e
		}},
		{"temperature_store", func(m *mockAPI) (map[string]any, error) {
			_, o, e := tools.NewTemperatureStoreHandler(m)(t.Context(), nil, tools.TemperatureStoreParams{})

			return o, e
		}},
		{"gcode_store", func(m *mockAPI) (map[string]any, error) {
			_, o, e := tools.NewGcodeStoreHandler(m)(t.Context(), nil, tools.GcodeStoreParams{})

			return o, e
		}},
		{"query_endstops", func(m *mockAPI) (map[string]any, error) {
			_, o, e := tools.NewQueryEndstopsHandler(m)(t.Context(), nil, tools.NoParams{})

			return o, e
		}},
		{"gcode_help", func(m *mockAPI) (map[string]any, error) {
			_, o, e := tools.NewGcodeHelpHandler(m)(t.Context(), nil, tools.NoParams{})

			return o, e
		}},
		{"history_totals", func(m *mockAPI) (map[string]any, error) {
			_, o, e := tools.NewHistoryTotalsHandler(m)(t.Context(), nil, tools.NoParams{})

			return o, e
		}},
		{"history_list", func(m *mockAPI) (map[string]any, error) {
			_, o, e := tools.NewHistoryListHandler(m)(t.Context(), nil, tools.HistoryListParams{})

			return o, e
		}},
		{"history_job", func(m *mockAPI) (map[string]any, error) {
			_, o, e := tools.NewHistoryJobHandler(m)(t.Context(), nil, tools.HistoryJobParams{UID: "1"})

			return o, e
		}},
		{"jobqueue_status", func(m *mockAPI) (map[string]any, error) {
			_, o, e := tools.NewJobQueueStatusHandler(m)(t.Context(), nil, tools.NoParams{})

			return o, e
		}},
		{"db_list", func(m *mockAPI) (map[string]any, error) {
			_, o, e := tools.NewDBListHandler(m)(t.Context(), nil, tools.NoParams{})

			return o, e
		}},
		{"db_get_item", func(m *mockAPI) (map[string]any, error) {
			_, o, e := tools.NewDBGetItemHandler(m)(t.Context(), nil, tools.DBGetItemParams{Namespace: "ns"})

			return o, e
		}},
		{"access_user_info", func(m *mockAPI) (map[string]any, error) {
			_, o, e := tools.NewAccessUserInfoHandler(m)(t.Context(), nil, tools.NoParams{})

			return o, e
		}},
		{"access_users_list", func(m *mockAPI) (map[string]any, error) {
			_, o, e := tools.NewAccessUsersListHandler(m)(t.Context(), nil, tools.NoParams{})

			return o, e
		}},
		{"access_info", func(m *mockAPI) (map[string]any, error) {
			_, o, e := tools.NewAccessInfoHandler(m)(t.Context(), nil, tools.NoParams{})

			return o, e
		}},
		{"announcements_list", func(m *mockAPI) (map[string]any, error) {
			_, o, e := tools.NewAnnouncementsListHandler(m)(t.Context(), nil, tools.AnnouncementsListParams{})

			return o, e
		}},
		{"announcements_feeds", func(m *mockAPI) (map[string]any, error) {
			_, o, e := tools.NewAnnouncementsFeedsHandler(m)(t.Context(), nil, tools.NoParams{})

			return o, e
		}},
		{"webcams_list", func(m *mockAPI) (map[string]any, error) {
			_, o, e := tools.NewWebcamsListHandler(m)(t.Context(), nil, tools.NoParams{})

			return o, e
		}},
		{"webcams_get", func(m *mockAPI) (map[string]any, error) {
			_, o, e := tools.NewWebcamsGetHandler(m)(t.Context(), nil, tools.WebcamNameParams{Name: "cam"})

			return o, e
		}},
		{"sensors_list", func(m *mockAPI) (map[string]any, error) {
			_, o, e := tools.NewSensorsListHandler(m)(t.Context(), nil, tools.SensorsListParams{})

			return o, e
		}},
		{"sensors_info", func(m *mockAPI) (map[string]any, error) {
			_, o, e := tools.NewSensorsInfoHandler(m)(t.Context(), nil, tools.SensorParams{Sensor: testSensor})

			return o, e
		}},
		{"sensors_measurements", func(m *mockAPI) (map[string]any, error) {
			_, o, e := tools.NewSensorsMeasurementsHandler(m)(t.Context(), nil, tools.SensorsMeasurementsParams{})

			return o, e
		}},
		{"spoolman_status", func(m *mockAPI) (map[string]any, error) {
			_, o, e := tools.NewSpoolmanStatusHandler(m)(t.Context(), nil, tools.NoParams{})

			return o, e
		}},
		{"spoolman_get_spool", func(m *mockAPI) (map[string]any, error) {
			_, o, e := tools.NewSpoolmanGetSpoolHandler(m)(t.Context(), nil, tools.NoParams{})

			return o, e
		}},
		{"wled_status", func(m *mockAPI) (map[string]any, error) {
			_, o, e := tools.NewWLEDStatusHandler(m)(t.Context(), nil, tools.StripParams{Strip: testStrip})

			return o, e
		}},
		{"extensions_list", func(m *mockAPI) (map[string]any, error) {
			_, o, e := tools.NewExtensionsListHandler(m)(t.Context(), nil, tools.NoParams{})

			return o, e
		}},
		{"analysis_status", func(m *mockAPI) (map[string]any, error) {
			_, o, e := tools.NewAnalysisStatusHandler(m)(t.Context(), nil, tools.NoParams{})

			return o, e
		}},
		{"files_directory", func(m *mockAPI) (map[string]any, error) {
			_, o, e := tools.NewFilesDirectoryHandler(m)(t.Context(), nil, tools.FilesDirectoryParams{})

			return o, e
		}},
		{"files_metadata", func(m *mockAPI) (map[string]any, error) {
			_, o, e := tools.NewFilesMetadataHandler(m)(t.Context(), nil, tools.FilenameParams{Filename: testGcodeFile})

			return o, e
		}},
	}

	for _, testCase := range tests {
		t.Run(testCase.name, func(t *testing.T) {
			t.Parallel()

			out, err := testCase.call(&mockAPI{result: fixture})
			if err != nil {
				t.Fatalf("handler: %v", err)
			}

			if _, wrapped := out["result"]; wrapped {
				t.Errorf("output carries a 'result' envelope: %v", out)
			}

			if out[probe] != float64(42) {
				t.Errorf("output = %v, want the payload's top-level %q key preserved", out, probe)
			}
		})
	}
}

// TestContract_DataReadStaysTopLevel verifies a data-read tool returns the
// Moonraker payload's fields at the top level, with no "result" envelope.
func TestContract_DataReadStaysTopLevel(t *testing.T) {
	t.Parallel()

	mock := &mockAPI{result: json.RawMessage(`{"system_info":{"cpu_count":4}}`)}

	_, out, err := tools.NewSystemInfoHandler(mock)(t.Context(), nil, tools.NoParams{})
	if err != nil {
		t.Fatalf("handler: %v", err)
	}

	if _, wrapped := out["result"]; wrapped {
		t.Errorf("output must not carry a 'result' envelope: %v", out)
	}

	if _, ok := out["system_info"]; !ok {
		t.Errorf("output = %v, want a top-level system_info key", out)
	}
}

// TestContract_ActionReturnsAck verifies an action tool whose Moonraker payload
// is the scalar "ok" normalizes to a uniform {"ok": true} acknowledgement.
func TestContract_ActionReturnsAck(t *testing.T) {
	t.Parallel()

	mock := &mockAPI{result: json.RawMessage(`"ok"`)}

	_, out, err := tools.NewPrintPauseHandler(mock)(t.Context(), nil, tools.NoParams{})
	if err != nil {
		t.Fatalf("handler: %v", err)
	}

	if out["ok"] != true {
		t.Errorf("output = %v, want {ok:true}", out)
	}
}

// TestContract_ActionPreservesObject verifies an action tool that returns a real
// object passes it through at the top level rather than discarding it.
func TestContract_ActionPreservesObject(t *testing.T) {
	t.Parallel()

	mock := &mockAPI{result: json.RawMessage(`{"item":{"path":"gcodes/b.gcode"}}`)}

	params := tools.SourceDestParams{Source: "gcodes/a.gcode", Dest: "gcodes/b.gcode"}

	_, out, err := tools.NewFilesMoveHandler(mock)(t.Context(), nil, params)
	if err != nil {
		t.Fatalf("handler: %v", err)
	}

	if _, ok := out["item"]; !ok {
		t.Errorf("output = %v, want the returned object preserved at top level", out)
	}
}

// TestContract_RootsWrapBareArray verifies the bare-array roots payload is
// wrapped under a "roots" key (MCP structured content must be an object).
func TestContract_RootsWrapBareArray(t *testing.T) {
	t.Parallel()

	mock := &mockAPI{result: json.RawMessage(`[{"name":"gcodes"},{"name":"config"}]`)}

	_, out, err := tools.NewFilesRootsHandler(mock)(t.Context(), nil, tools.NoParams{})
	if err != nil {
		t.Fatalf("handler: %v", err)
	}

	if len(out.Roots) != 2 {
		t.Errorf("roots = %v, want 2 entries", out.Roots)
	}
}

// TestContract_ThumbnailsWrapBareArray verifies the bare-array thumbnails payload
// is wrapped under a "thumbnails" key.
func TestContract_ThumbnailsWrapBareArray(t *testing.T) {
	t.Parallel()

	mock := &mockAPI{result: json.RawMessage(`[{"width":300},{"width":32}]`)}

	params := tools.FilenameParams{Filename: testGcodeFile}

	_, out, err := tools.NewFilesThumbnailsHandler(mock)(t.Context(), nil, params)
	if err != nil {
		t.Fatalf("handler: %v", err)
	}

	if len(out.Thumbnails) != 2 {
		t.Errorf("thumbnails = %v, want 2 entries", out.Thumbnails)
	}
}

// TestContract_APIKeyWrapsScalar verifies the bare-string API key is wrapped
// under an "api_key" key.
func TestContract_APIKeyWrapsScalar(t *testing.T) {
	t.Parallel()

	mock := &mockAPI{result: json.RawMessage(`"deadbeef"`)}

	_, out, err := tools.NewAccessAPIKeyHandler(mock)(t.Context(), nil, tools.NoParams{})
	if err != nil {
		t.Fatalf("handler: %v", err)
	}

	if out.APIKey != "deadbeef" {
		t.Errorf("api_key = %q, want deadbeef", out.APIKey)
	}
}

// TestContract_MQTTSubscribePassesScalarThrough verifies an MQTT payload that is
// a bare scalar is returned verbatim rather than collapsed to an acknowledgement.
func TestContract_MQTTSubscribePassesScalarThrough(t *testing.T) {
	t.Parallel()

	mock := &mockAPI{result: json.RawMessage(`"online"`)}

	params := tools.MQTTSubscribeParams{Topic: "klipper/state"}

	_, out, err := tools.NewMQTTSubscribeHandler(mock)(t.Context(), nil, params)
	if err != nil {
		t.Fatalf("handler: %v", err)
	}

	// The scalar payload is preserved under an object "payload" key (MCP requires
	// structured content to be an object).
	if out["payload"] != "online" {
		t.Errorf("out = %#v, want payload \"online\" preserved", out)
	}
}

// TestContract_ExtensionsRequestWrapsArray verifies a shape-variable agent
// response that is a bare array is returned under "response" so the structured
// content stays a JSON object.
func TestContract_ExtensionsRequestWrapsArray(t *testing.T) {
	t.Parallel()

	mock := &mockAPI{result: json.RawMessage(`[{"id":1},{"id":2}]`)}

	params := tools.ExtensionsRequestParams{Agent: testAgent, Method: "list"}

	_, out, err := tools.NewExtensionsRequestHandler(mock)(t.Context(), nil, params)
	if err != nil {
		t.Fatalf("handler: %v", err)
	}

	arr, ok := out["response"].([]any)
	if !ok {
		t.Fatalf("out[response] = %T, want []any passthrough", out["response"])
	}

	if len(arr) != 2 {
		t.Errorf("array len = %d, want 2", len(arr))
	}
}

// TestContract_SpoolmanProxyExposesV2Envelope verifies the proxy surfaces
// Moonraker's v2 {response, error} envelope at the top level rather than nesting
// it again under another "response" key (which would hide the error field).
func TestContract_SpoolmanProxyExposesV2Envelope(t *testing.T) {
	t.Parallel()

	mock := &mockAPI{result: json.RawMessage(`{"response":{"id":1},"error":null}`)}

	params := tools.SpoolmanProxyParams{RequestMethod: "GET", Path: "/v1/spool/1"}

	_, out, err := tools.NewSpoolmanProxyHandler(mock)(t.Context(), nil, params)
	if err != nil {
		t.Fatalf("handler: %v", err)
	}

	// The error field must be visible at the top level.
	if _, hasError := out["error"]; !hasError {
		t.Errorf("out = %v, want a top-level error field", out)
	}

	// response holds the actual payload, not the whole envelope again.
	data, ok := out["response"].(map[string]any)
	if !ok || data["id"] != float64(1) {
		t.Errorf("out[response] = %v, want the inner payload {id:1}, not a re-nested envelope", out["response"])
	}
}
