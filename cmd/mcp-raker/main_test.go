package main

import (
	"encoding/json"
	"log/slog"
	"net/http"
	"net/http/httptest"
	"slices"
	"strings"
	"testing"

	"github.com/modelcontextprotocol/go-sdk/mcp"

	"github.com/lexfrei/mcp-raker/internal/moonraker"
)

// adminToolCount is the number of tools gated behind MOONRAKER_ENABLE_ADMIN.
const adminToolCount = 16

func testLogger() *slog.Logger {
	return slog.New(slog.DiscardHandler)
}

// listTools registers tools for the given admin mode and returns them by name,
// as advertised over an in-memory transport (annotations included).
func listTools(t *testing.T, enableAdmin bool) map[string]*mcp.Tool {
	t.Helper()

	client, err := moonraker.New(&moonraker.Options{})
	if err != nil {
		t.Fatalf("moonraker.New: %v", err)
	}

	server := mcp.NewServer(&mcp.Implementation{Name: serverName, Version: "test"}, newServerOptions(testLogger(), enableAdmin))
	registerTools(server, client, enableAdmin)

	serverTransport, clientTransport := mcp.NewInMemoryTransports()

	serverSession, err := server.Connect(t.Context(), serverTransport, nil)
	if err != nil {
		t.Fatalf("server connect: %v", err)
	}
	defer func() { _ = serverSession.Close() }()

	mcpClient := mcp.NewClient(&mcp.Implementation{Name: "test-client", Version: "test"}, nil)

	clientSession, err := mcpClient.Connect(t.Context(), clientTransport, nil)
	if err != nil {
		t.Fatalf("client connect: %v", err)
	}
	defer func() { _ = clientSession.Close() }()

	result, err := clientSession.ListTools(t.Context(), nil)
	if err != nil {
		t.Fatalf("ListTools: %v", err)
	}

	tools := make(map[string]*mcp.Tool, len(result.Tools))
	for _, tool := range result.Tools {
		tools[tool.Name] = tool
	}

	return tools
}

// listToolNames returns the set of registered tool names for the given mode.
func listToolNames(t *testing.T, enableAdmin bool) map[string]bool {
	t.Helper()

	names := make(map[string]bool)
	for name := range listTools(t, enableAdmin) {
		names[name] = true
	}

	return names
}

// TestAdminToolsDestructive enforces that every tool gated behind
// MOONRAKER_ENABLE_ADMIN advertises a destructive hint. It derives the admin set
// from the actual registration (admin-with minus core-only), so a new admin tool
// with the wrong annotation fails here rather than slipping through.
func TestAdminToolsDestructive(t *testing.T) {
	t.Parallel()

	core := listToolNames(t, false)
	admin := listTools(t, true)

	adminOnly := 0

	for name, tool := range admin {
		if core[name] {
			continue
		}

		adminOnly++

		if tool.Annotations == nil || tool.Annotations.DestructiveHint == nil || !*tool.Annotations.DestructiveHint {
			t.Errorf("admin tool %q must carry a destructive hint", name)
		}
	}

	if adminOnly != adminToolCount {
		t.Errorf("found %d admin-only tools, want %d", adminOnly, adminToolCount)
	}
}

func TestRegisterTools_CoreOnly(t *testing.T) {
	t.Parallel()

	names := listToolNames(t, false)

	if !names["moonraker_printer_info"] {
		t.Error("core tool moonraker_printer_info not registered")
	}

	if !names["moonraker_emergency_stop"] {
		t.Error("core tool moonraker_emergency_stop not registered")
	}

	// Klipper restart is a printer action, not an OS/service admin action, so it
	// is available without MOONRAKER_ENABLE_ADMIN.
	if !names["moonraker_printer_restart"] {
		t.Error("core tool moonraker_printer_restart not registered")
	}

	if !names["moonraker_firmware_restart"] {
		t.Error("core tool moonraker_firmware_restart not registered")
	}

	// The Moonraker server restart (a service action) stays behind the admin gate.
	if names["moonraker_server_restart"] {
		t.Error("admin tool moonraker_server_restart registered without MOONRAKER_ENABLE_ADMIN")
	}

	if names["moonraker_machine_shutdown"] {
		t.Error("admin tool moonraker_machine_shutdown registered without MOONRAKER_ENABLE_ADMIN")
	}

	if len(names) < 90 {
		t.Errorf("registered %d core tools, want at least 90", len(names))
	}
}

func TestRegisterTools_WithAdmin(t *testing.T) {
	t.Parallel()

	core := listToolNames(t, false)
	admin := listToolNames(t, true)

	if !admin["moonraker_machine_shutdown"] {
		t.Error("admin tool moonraker_machine_shutdown not registered with admin enabled")
	}

	if !admin["moonraker_update_rollback"] {
		t.Error("admin tool moonraker_update_rollback not registered with admin enabled")
	}

	if len(admin) != len(core)+adminToolCount {
		t.Errorf("admin mode registered %d tools, want %d (core %d + %d admin)",
			len(admin), len(core)+adminToolCount, len(core), adminToolCount)
	}
}

// TestToolOutputSchemasNoBooleanProperty guards against the SDK reflecting an
// any-typed output field into a bare boolean property schema
// ("properties":{"x":true}), which strict MCP clients reject and which fails the
// whole tools/list. It scans every property of every tool's output schema, so a
// future passthrough field is covered too, not just the current "result" key.
func TestToolOutputSchemasNoBooleanProperty(t *testing.T) {
	t.Parallel()

	for name, tool := range listTools(t, true) {
		if tool.OutputSchema == nil {
			continue
		}

		raw, err := json.Marshal(tool.OutputSchema)
		if err != nil {
			t.Fatalf("%s: marshal output schema: %v", name, err)
		}

		var schema struct {
			Properties map[string]json.RawMessage `json:"properties"`
		}
		if err := json.Unmarshal(raw, &schema); err != nil {
			t.Fatalf("%s: unmarshal output schema: %v", name, err)
		}

		for prop, propSchema := range schema.Properties {
			if value := strings.TrimSpace(string(propSchema)); value == "true" || value == "false" {
				t.Errorf("%s: property %q has a bare boolean schema %s (rejected by strict clients)", name, prop, value)
			}
		}
	}
}

// TestOptionalParamsNotRequired verifies that filter and paging parameters with
// sensible defaults are not marked required in the generated input schema, while
// genuinely required parameters still are. The SDK marks every struct field
// required unless its json tag carries omitempty, so a regression here means a
// caller is forced to send placeholder values.
func TestOptionalParamsNotRequired(t *testing.T) {
	t.Parallel()

	tools := listTools(t, true)

	required := func(name string) []string {
		tool, ok := tools[name]
		if !ok {
			t.Fatalf("tool %q not registered", name)
		}

		schema, ok := tool.InputSchema.(map[string]any)
		if !ok {
			t.Fatalf("tool %q input schema is %T, want a JSON object", name, tool.InputSchema)
		}

		raw, _ := schema["required"].([]any)

		out := make([]string, 0, len(raw))
		for _, value := range raw {
			if field, ok := value.(string); ok {
				out = append(out, field)
			}
		}

		return out
	}

	// These tools have only optional filter/paging parameters.
	allOptional := []string{
		"moonraker_history_list",
		"moonraker_files_list",
		"moonraker_temperature_store",
		"moonraker_gcode_store",
		"moonraker_sensors_list",
		"moonraker_sensors_measurements",
		"moonraker_files_directory",
		"moonraker_update_refresh",
		"moonraker_db_backup",
	}
	for _, name := range allOptional {
		if got := required(name); len(got) != 0 {
			t.Errorf("%s: required = %v, want none", name, got)
		}
	}

	// These tools keep their genuinely required parameters.
	mustRequire := map[string]string{
		"moonraker_print_start":  "filename",
		"moonraker_history_job":  "uid",
		"moonraker_gcode_script": "script",
		"moonraker_db_post_item": "key",
		"moonraker_wled_set":     "strip",
		"moonraker_sensors_info": "sensor",
	}
	for name, field := range mustRequire {
		if !slices.Contains(required(name), field) {
			t.Errorf("%s: required = %v, want it to include %q", name, required(name), field)
		}
	}
}

func TestBearerAuth(t *testing.T) {
	t.Parallel()

	okHandler := http.HandlerFunc(func(writer http.ResponseWriter, _ *http.Request) {
		writer.WriteHeader(http.StatusOK)
	})

	const token = "secret"

	tests := []struct {
		name   string
		token  string
		header string
		want   int
	}{
		{name: "no token disables the check", token: "", header: "", want: http.StatusOK},
		{name: "valid token", token: token, header: "Bearer " + token, want: http.StatusOK},
		{name: "missing header", token: token, header: "", want: http.StatusUnauthorized},
		{name: "wrong token", token: token, header: "Bearer nope", want: http.StatusUnauthorized},
	}

	for _, testCase := range tests {
		t.Run(testCase.name, func(t *testing.T) {
			t.Parallel()

			handler := bearerAuth(okHandler, testCase.token)
			req := httptest.NewRequestWithContext(t.Context(), http.MethodGet, "/", nil)

			if testCase.header != "" {
				req.Header.Set("Authorization", testCase.header)
			}

			rec := httptest.NewRecorder()
			handler.ServeHTTP(rec, req)

			if rec.Code != testCase.want {
				t.Errorf("status = %d, want %d", rec.Code, testCase.want)
			}
		})
	}
}
