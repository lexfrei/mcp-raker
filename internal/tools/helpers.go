package tools

import (
	"encoding/json"

	"github.com/modelcontextprotocol/go-sdk/mcp"
)

// NoParams is the input type for tools that take no parameters.
type NoParams struct{}

// ptrBool returns a pointer to value, for the *bool annotation hint fields.
func ptrBool(value bool) *bool { return &value }

// readOnly builds annotations for a tool that only reads printer state.
func readOnly(title string) *mcp.ToolAnnotations {
	return &mcp.ToolAnnotations{
		Title:         title,
		ReadOnlyHint:  true,
		OpenWorldHint: ptrBool(true),
	}
}

// write builds annotations for a tool that changes printer or server state but
// is not destructive.
func write(title string) *mcp.ToolAnnotations {
	return &mcp.ToolAnnotations{
		Title:           title,
		ReadOnlyHint:    false,
		DestructiveHint: ptrBool(false),
		IdempotentHint:  false,
		OpenWorldHint:   ptrBool(true),
	}
}

// writeDestructive builds annotations for a tool whose effect is disruptive or
// hard to undo: stopping a print, deleting a file, rebooting the host.
func writeDestructive(title string) *mcp.ToolAnnotations {
	return &mcp.ToolAnnotations{
		Title:           title,
		ReadOnlyHint:    false,
		DestructiveHint: ptrBool(true),
		IdempotentHint:  false,
		OpenWorldHint:   ptrBool(true),
	}
}

// decodeResult normalizes a Moonraker result into a top-level JSON object.
//
// The client already strips Moonraker's {"result": ...} envelope, so an
// object payload is returned verbatim with its fields at the top level — a
// consumer always reads the data directly, never via a ".result" key. A scalar
// "ok", an empty body, or any other non-object success payload normalizes to a
// uniform {"ok": true} acknowledgement so action tools share one shape.
//
// Tools whose endpoint returns a bare array or a meaningful scalar (e.g. the
// file list or an API key) must use a typed named-key wrapper instead, because
// MCP structured content must be a JSON object.
func decodeResult(raw json.RawMessage, err error) (map[string]any, error) {
	if err != nil {
		return nil, moonrakerErr("request failed", err)
	}

	if len(raw) == 0 {
		return map[string]any{"ok": true}, nil
	}

	var value any

	unErr := json.Unmarshal(raw, &value)
	if unErr != nil {
		return nil, moonrakerErr("decode response", unErr)
	}

	if obj, ok := value.(map[string]any); ok {
		return obj, nil
	}

	return map[string]any{"ok": true}, nil
}

// decodePassthrough returns the Moonraker payload unchanged, for tools whose
// response shape varies by request (the Spoolman and extension proxies). Such
// tools register no output schema, so any JSON value — object, array, or
// scalar — passes through at the top level.
func decodePassthrough(raw json.RawMessage, err error) (any, error) {
	if err != nil {
		return nil, moonrakerErr("request failed", err)
	}

	if len(raw) == 0 {
		//nolint:nilnil // An empty proxy response legitimately carries no value and no error.
		return nil, nil
	}

	var value any

	unErr := json.Unmarshal(raw, &value)
	if unErr != nil {
		return nil, moonrakerErr("decode response", unErr)
	}

	return value, nil
}

// decodeTyped unmarshals a client result and error into T.
func decodeTyped[T any](raw json.RawMessage, err error) (T, error) {
	var out T

	if err != nil {
		return out, moonrakerErr("request failed", err)
	}

	if len(raw) == 0 {
		return out, nil
	}

	unErr := json.Unmarshal(raw, &out)
	if unErr != nil {
		return out, moonrakerErr("decode response", unErr)
	}

	return out, nil
}
