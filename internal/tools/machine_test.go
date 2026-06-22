package tools_test

import (
	"encoding/json"
	"strconv"
	"strings"
	"testing"

	"github.com/cockroachdb/errors"

	"github.com/lexfrei/mcp-raker/internal/tools"
)

// procStatsFixture builds a proc_stats result whose CPU history has n points,
// each tagged with its index so a test can tell which points survived trimming.
func procStatsFixture(n int) json.RawMessage {
	points := make([]string, n)
	for i := range points {
		points[i] = `{"cpu_usage":` + strconv.Itoa(i) + `}`
	}

	return json.RawMessage(`{"cpu_temp":50,"moonraker_stats":[` + strings.Join(points, ",") + `]}`)
}

// procStatsCount returns the number of CPU history points in a handler result.
func procStatsCount(t *testing.T, out map[string]any) int {
	t.Helper()

	stats, ok := out["moonraker_stats"].([]any)
	if !ok {
		t.Fatalf("moonraker_stats = %v, want a slice", out["moonraker_stats"])
	}

	return len(stats)
}

func TestProcStats_TrimsToFiveByDefault(t *testing.T) {
	t.Parallel()

	mock := &mockAPI{result: procStatsFixture(10)}

	_, out, err := tools.NewProcStatsHandler(mock)(t.Context(), nil, tools.ProcStatsParams{})
	if err != nil {
		t.Fatalf("handler: %v", err)
	}

	if got := procStatsCount(t, out); got != 5 {
		t.Fatalf("kept %d points, want 5", got)
	}

	// The most recent points are the ones that survive: indices 5..9.
	first := out["moonraker_stats"].([]any)[0].(map[string]any)
	if first["cpu_usage"] != float64(5) {
		t.Errorf("first kept point cpu_usage = %v, want 5 (the last 5 points)", first["cpu_usage"])
	}
}

func TestProcStats_SamplesZeroKeepsAll(t *testing.T) {
	t.Parallel()

	all := 0
	mock := &mockAPI{result: procStatsFixture(10)}

	_, out, err := tools.NewProcStatsHandler(mock)(t.Context(), nil, tools.ProcStatsParams{Samples: &all})
	if err != nil {
		t.Fatalf("handler: %v", err)
	}

	if got := procStatsCount(t, out); got != 10 {
		t.Errorf("kept %d points, want all 10", got)
	}
}

func TestProcStats_SamplesN(t *testing.T) {
	t.Parallel()

	two := 2
	mock := &mockAPI{result: procStatsFixture(10)}

	_, out, err := tools.NewProcStatsHandler(mock)(t.Context(), nil, tools.ProcStatsParams{Samples: &two})
	if err != nil {
		t.Fatalf("handler: %v", err)
	}

	if got := procStatsCount(t, out); got != 2 {
		t.Errorf("kept %d points, want 2", got)
	}
}

func TestProcStats_NegativeSamplesKeepsAll(t *testing.T) {
	t.Parallel()

	negative := -1
	mock := &mockAPI{result: procStatsFixture(10)}

	_, out, err := tools.NewProcStatsHandler(mock)(t.Context(), nil, tools.ProcStatsParams{Samples: &negative})
	if err != nil {
		t.Fatalf("handler: %v", err)
	}

	if got := procStatsCount(t, out); got != 10 {
		t.Errorf("kept %d points, want all 10 (a negative count means full history)", got)
	}
}

func TestProcStats_FewerThanLimitUnchanged(t *testing.T) {
	t.Parallel()

	mock := &mockAPI{result: procStatsFixture(3)}

	_, out, err := tools.NewProcStatsHandler(mock)(t.Context(), nil, tools.ProcStatsParams{})
	if err != nil {
		t.Fatalf("handler: %v", err)
	}

	if got := procStatsCount(t, out); got != 3 {
		t.Errorf("kept %d points, want 3", got)
	}
}

func TestMachineReads(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name    string
		handler func(*mockAPI) error
		path    string
	}{
		{"system_info", func(m *mockAPI) error {
			_, _, err := tools.NewSystemInfoHandler(m)(t.Context(), nil, tools.NoParams{})

			return err
		}, "/machine/system_info"},
		{"proc_stats", func(m *mockAPI) error {
			_, _, err := tools.NewProcStatsHandler(m)(t.Context(), nil, tools.ProcStatsParams{})

			return err
		}, "/machine/proc_stats"},
		{"sudo_info", func(m *mockAPI) error {
			_, _, err := tools.NewSudoInfoHandler(m)(t.Context(), nil, tools.NoParams{})

			return err
		}, "/machine/sudo/info"},
		{"usb", func(m *mockAPI) error {
			_, _, err := tools.NewPeripheralsUSBHandler(m)(t.Context(), nil, tools.NoParams{})

			return err
		}, "/machine/peripherals/usb"},
		{"serial", func(m *mockAPI) error {
			_, _, err := tools.NewPeripheralsSerialHandler(m)(t.Context(), nil, tools.NoParams{})

			return err
		}, "/machine/peripherals/serial"},
		{"video", func(m *mockAPI) error {
			_, _, err := tools.NewPeripheralsVideoHandler(m)(t.Context(), nil, tools.NoParams{})

			return err
		}, "/machine/peripherals/video"},
		{"canbus", func(m *mockAPI) error {
			_, _, err := tools.NewPeripheralsCanbusHandler(m)(t.Context(), nil, tools.PeripheralsCanbusParams{})

			return err
		}, "/machine/peripherals/canbus"},
	}

	for _, testCase := range tests {
		t.Run(testCase.name, func(t *testing.T) {
			t.Parallel()

			mock := &mockAPI{result: okJSON}

			err := testCase.handler(mock)
			if err != nil {
				t.Fatalf("handler: %v", err)
			}

			assertCall(t, mock, methodGet, testCase.path)
		})
	}
}

func TestMachineDestructive(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name    string
		handler func(*mockAPI) error
		path    string
	}{
		{"shutdown", func(m *mockAPI) error {
			_, _, err := tools.NewMachineShutdownHandler(m)(t.Context(), nil, tools.NoParams{})

			return err
		}, "/machine/shutdown"},
		{"reboot", func(m *mockAPI) error {
			_, _, err := tools.NewMachineRebootHandler(m)(t.Context(), nil, tools.NoParams{})

			return err
		}, "/machine/reboot"},
	}

	for _, testCase := range tests {
		t.Run(testCase.name, func(t *testing.T) {
			t.Parallel()

			mock := &mockAPI{result: okJSON}

			err := testCase.handler(mock)
			if err != nil {
				t.Fatalf("handler: %v", err)
			}

			assertCall(t, mock, methodPost, testCase.path)
		})
	}
}

func TestServiceControls_RequireService(t *testing.T) {
	t.Parallel()

	_, _, startErr := tools.NewServiceStartHandler(&mockAPI{})(t.Context(), nil, tools.ServiceParams{})
	if !errors.Is(startErr, tools.ErrValidation) {
		t.Errorf("start err = %v, want ErrValidation", startErr)
	}

	_, _, stopErr := tools.NewServiceStopHandler(&mockAPI{})(t.Context(), nil, tools.ServiceParams{})
	if !errors.Is(stopErr, tools.ErrValidation) {
		t.Errorf("stop err = %v, want ErrValidation", stopErr)
	}

	_, _, restartErr := tools.NewServiceRestartHandler(&mockAPI{})(t.Context(), nil, tools.ServiceParams{})
	if !errors.Is(restartErr, tools.ErrValidation) {
		t.Errorf("restart err = %v, want ErrValidation", restartErr)
	}
}

func TestServiceStart_SendsService(t *testing.T) {
	t.Parallel()

	mock := &mockAPI{result: okJSON}

	_, _, err := tools.NewServiceStartHandler(mock)(t.Context(), nil, tools.ServiceParams{Service: svcKlipper})
	if err != nil {
		t.Fatalf("handler: %v", err)
	}

	assertCall(t, mock, methodPost, "/machine/services/start")

	if mock.lastQuery.Get("service") != svcKlipper {
		t.Errorf("service = %q, want %s", mock.lastQuery.Get("service"), svcKlipper)
	}
}

func TestSudoPassword(t *testing.T) {
	t.Parallel()

	_, _, valErr := tools.NewSudoPasswordHandler(&mockAPI{})(t.Context(), nil, tools.SudoPasswordParams{})
	if !errors.Is(valErr, tools.ErrValidation) {
		t.Errorf("err = %v, want ErrValidation", valErr)
	}

	mock := &mockAPI{result: okJSON}

	_, _, err := tools.NewSudoPasswordHandler(mock)(t.Context(), nil, tools.SudoPasswordParams{Password: "secret"})
	if err != nil {
		t.Fatalf("handler: %v", err)
	}

	assertCall(t, mock, methodPost, "/machine/sudo/password")
}
