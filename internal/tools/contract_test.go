package tools_test

import (
	"encoding/json"
	"testing"

	"github.com/lexfrei/mcp-raker/internal/tools"
)

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

// TestContract_ProxyPassesArrayThrough verifies the shape-variable proxy returns
// an array payload unchanged at the top level.
func TestContract_ProxyPassesArrayThrough(t *testing.T) {
	t.Parallel()

	mock := &mockAPI{result: json.RawMessage(`[{"id":1},{"id":2}]`)}

	params := tools.SpoolmanProxyParams{RequestMethod: "GET", Path: "/v1/spool"}

	_, out, err := tools.NewSpoolmanProxyHandler(mock)(t.Context(), nil, params)
	if err != nil {
		t.Fatalf("handler: %v", err)
	}

	arr, ok := out.([]any)
	if !ok {
		t.Fatalf("output = %T, want []any passthrough", out)
	}

	if len(arr) != 2 {
		t.Errorf("array len = %d, want 2", len(arr))
	}
}
