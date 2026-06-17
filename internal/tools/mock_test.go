package tools_test

import (
	"context"
	"encoding/json"
	"net/url"
	"testing"

	"github.com/cockroachdb/errors"

	"github.com/lexfrei/mcp-raker/internal/moonraker"
)

// assertCall fails the test unless the mock recorded a call with the given verb
// and path.
func assertCall(t *testing.T, mock *mockAPI, method, path string) {
	t.Helper()

	if mock.lastMethod != method || mock.lastPath != path {
		t.Errorf("call = %s %s, want %s %s", mock.lastMethod, mock.lastPath, method, path)
	}
}

// okJSON is a minimal successful Moonraker result body for tests.
var okJSON = json.RawMessage(`"ok"`)

// svcKlipper is a sample service name reused across tests.
const svcKlipper = "klipper"

// gcodeG28 is a sample G-code command reused across tests.
const gcodeG28 = "G28"

// queryTrue is the string value boolean query flags are set to.
const queryTrue = "true"

// testGcodeFile is a sample gcode filename reused across tests.
const testGcodeFile = "a.gcode"

// testSensor is a sample sensor name reused across tests.
const testSensor = "chamber"

// testStrip is a sample WLED strip name reused across tests.
const testStrip = "lights"

// testAgent is a sample extension agent name reused across tests.
const testAgent = "moonagent"

// Power/WLED action labels reused across tests.
const (
	actionOn     = "on"
	actionOff    = "off"
	actionToggle = "toggle"
)

// HTTP verb labels the mock records, shared across the tool tests.
const (
	methodGet    = "GET"
	methodPost   = "POST"
	methodDelete = "DELETE"
	methodUpload = "UPLOAD"
	methodGetRaw = "GETRAW"
)

// errStub is a static error for exercising the error path of handlers.
var errStub = errors.New("stub failure")

// mockAPI is a configurable moonraker.API for handler tests. It returns the
// configured result/error for each verb and records the last call's arguments.
type mockAPI struct {
	result json.RawMessage
	err    error

	rawResult []byte
	rawErr    error

	lastMethod string
	lastPath   string
	lastQuery  url.Values
	lastBody   any
	lastUpload *moonraker.UploadOptions
}

func (m *mockAPI) Get(_ context.Context, path string, query url.Values) (json.RawMessage, error) {
	m.lastMethod = methodGet
	m.lastPath = path
	m.lastQuery = query

	return m.result, m.err
}

func (m *mockAPI) Post(_ context.Context, path string, query url.Values, body any) (json.RawMessage, error) {
	m.lastMethod = methodPost
	m.lastPath = path
	m.lastQuery = query
	m.lastBody = body

	return m.result, m.err
}

func (m *mockAPI) Delete(_ context.Context, path string, query url.Values) (json.RawMessage, error) {
	m.lastMethod = methodDelete
	m.lastPath = path
	m.lastQuery = query

	return m.result, m.err
}

func (m *mockAPI) Upload(_ context.Context, opts *moonraker.UploadOptions) (json.RawMessage, error) {
	m.lastMethod = methodUpload
	m.lastUpload = opts

	return m.result, m.err
}

func (m *mockAPI) GetRaw(_ context.Context, path string, query url.Values) ([]byte, error) {
	m.lastMethod = methodGetRaw
	m.lastPath = path
	m.lastQuery = query

	return m.rawResult, m.rawErr
}
