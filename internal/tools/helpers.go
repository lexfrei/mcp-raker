package tools

import (
	"encoding/json"

	"github.com/modelcontextprotocol/go-sdk/mcp"
)

// NoParams is the input type for tools that take no parameters.
type NoParams struct{}

// RawResult carries an arbitrary Moonraker JSON payload for tools that do not
// need a typed output shape. It mirrors Moonraker's own {"result": ...}
// envelope so the model sees the response verbatim.
type RawResult struct {
	Result any `json:"result"`
}

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

// decodeRaw turns a client result and error into a passthrough RawResult.
func decodeRaw(raw json.RawMessage, err error) (RawResult, error) {
	if err != nil {
		return RawResult{}, moonrakerErr("request failed", err)
	}

	if len(raw) == 0 {
		return RawResult{}, nil
	}

	var value any

	unErr := json.Unmarshal(raw, &value)
	if unErr != nil {
		return RawResult{}, moonrakerErr("decode response", unErr)
	}

	return RawResult{Result: value}, nil
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
