package tools

import (
	"encoding/json"
	"testing"

	"github.com/cockroachdb/errors"
)

// errBoom is a static error for exercising the helper error paths.
var errBoom = errors.New("boom")

func TestDecodeResult_ObjectStaysTopLevel(t *testing.T) {
	t.Parallel()

	out, err := decodeResult(json.RawMessage(`{"klippy_state":"ready","count":3}`), nil)
	if err != nil {
		t.Fatalf("decodeResult: %v", err)
	}

	if _, wrapped := out["result"]; wrapped {
		t.Errorf("output must not carry a 'result' envelope: %v", out)
	}

	if out["klippy_state"] != "ready" {
		t.Errorf("klippy_state = %v, want ready", out["klippy_state"])
	}
}

func TestDecodeResult_OkScalarBecomesAck(t *testing.T) {
	t.Parallel()

	out, err := decodeResult(json.RawMessage(`"ok"`), nil)
	if err != nil {
		t.Fatalf("decodeResult: %v", err)
	}

	if out["ok"] != true {
		t.Errorf("output = %v, want {ok:true}", out)
	}
}

func TestDecodeResult_EmptyBecomesAck(t *testing.T) {
	t.Parallel()

	out, err := decodeResult(nil, nil)
	if err != nil {
		t.Fatalf("decodeResult: %v", err)
	}

	if out["ok"] != true {
		t.Errorf("output = %v, want {ok:true}", out)
	}
}

func TestDecodeResult_NonObjectBecomesAck(t *testing.T) {
	t.Parallel()

	// A null or any non-object success payload normalizes to a uniform ack.
	for _, raw := range []string{`null`, `42`, `[1,2,3]`} {
		out, err := decodeResult(json.RawMessage(raw), nil)
		if err != nil {
			t.Fatalf("decodeResult(%s): %v", raw, err)
		}

		if out["ok"] != true {
			t.Errorf("decodeResult(%s) = %v, want {ok:true}", raw, out)
		}
	}
}

func TestDecodeResult_PropagatesError(t *testing.T) {
	t.Parallel()

	_, err := decodeResult(nil, errBoom)
	if !errors.Is(err, ErrMoonraker) {
		t.Fatalf("err = %v, want ErrMoonraker", err)
	}
}

func TestDecodeResult_DecodeFailure(t *testing.T) {
	t.Parallel()

	_, err := decodeResult(json.RawMessage(`{bad`), nil)
	if !errors.Is(err, ErrMoonraker) {
		t.Fatalf("err = %v, want ErrMoonraker on malformed JSON", err)
	}
}

func TestDecodePassthrough_PreservesShape(t *testing.T) {
	t.Parallel()

	obj, err := decodePassthrough(json.RawMessage(`{"a":1}`), nil)
	if err != nil {
		t.Fatalf("decodePassthrough object: %v", err)
	}

	if _, ok := obj.(map[string]any); !ok {
		t.Errorf("object payload = %T, want map[string]any", obj)
	}

	arr, err := decodePassthrough(json.RawMessage(`[1,2]`), nil)
	if err != nil {
		t.Fatalf("decodePassthrough array: %v", err)
	}

	if _, ok := arr.([]any); !ok {
		t.Errorf("array payload = %T, want []any", arr)
	}
}

func TestDecodePassthrough_PropagatesError(t *testing.T) {
	t.Parallel()

	_, err := decodePassthrough(nil, errBoom)
	if !errors.Is(err, ErrMoonraker) {
		t.Fatalf("err = %v, want ErrMoonraker", err)
	}
}

func TestDecodePassthrough_PreservesNull(t *testing.T) {
	t.Parallel()

	// A literal JSON null from an upstream proxy must survive verbatim. Returning
	// a nil interface would make the SDK drop the result as "no content".
	out, err := decodePassthrough(json.RawMessage(`null`), nil)
	if err != nil {
		t.Fatalf("decodePassthrough: %v", err)
	}

	if out == nil {
		t.Fatal("JSON null was dropped; want it preserved")
	}

	if raw, ok := out.(json.RawMessage); !ok || string(raw) != "null" {
		t.Errorf("out = %#v, want json.RawMessage(\"null\")", out)
	}
}

func TestDecodePassthrough_EmptyBodyIsNil(t *testing.T) {
	t.Parallel()

	// An empty body (no response at all) legitimately yields no value.
	out, err := decodePassthrough(nil, nil)
	if err != nil {
		t.Fatalf("decodePassthrough: %v", err)
	}

	if out != nil {
		t.Errorf("out = %#v, want nil for an empty body", out)
	}
}
