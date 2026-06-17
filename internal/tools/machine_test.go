package tools_test

import (
	"testing"

	"github.com/cockroachdb/errors"

	"github.com/lexfrei/mcp-raker/internal/tools"
)

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
			_, _, err := tools.NewProcStatsHandler(m)(t.Context(), nil, tools.NoParams{})

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
