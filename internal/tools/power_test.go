package tools_test

import (
	"testing"

	"github.com/cockroachdb/errors"

	"github.com/lexfrei/mcp-raker/internal/tools"
)

const devicePrinter = "printer"

func TestPowerDevices(t *testing.T) {
	t.Parallel()

	mock := &mockAPI{result: okJSON}

	_, _, err := tools.NewPowerDevicesHandler(mock)(t.Context(), nil, tools.NoParams{})
	if err != nil {
		t.Fatalf("handler: %v", err)
	}

	assertCall(t, mock, methodGet, "/machine/device_power/devices")
}

func TestPowerStatus(t *testing.T) {
	t.Parallel()

	_, _, missing := tools.NewPowerStatusHandler(&mockAPI{})(t.Context(), nil, tools.PowerDeviceParams{})
	if !errors.Is(missing, tools.ErrValidation) {
		t.Errorf("missing device err = %v, want ErrValidation", missing)
	}

	mock := &mockAPI{result: okJSON}

	_, _, err := tools.NewPowerStatusHandler(mock)(t.Context(), nil, tools.PowerDeviceParams{Device: devicePrinter})
	if err != nil {
		t.Fatalf("handler: %v", err)
	}

	// Status uses the singular endpoint with a device= query parameter.
	assertCall(t, mock, methodGet, "/machine/device_power/device")

	if mock.lastQuery.Get("device") != devicePrinter {
		t.Errorf("device = %q, want %s", mock.lastQuery.Get("device"), devicePrinter)
	}
}

func TestPowerActions(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name    string
		action  string
		handler func(*mockAPI, tools.PowerDeviceParams) error
	}{
		{actionOn, actionOn, func(m *mockAPI, p tools.PowerDeviceParams) error {
			_, _, err := tools.NewPowerOnHandler(m)(t.Context(), nil, p)

			return err
		}},
		{actionOff, actionOff, func(m *mockAPI, p tools.PowerDeviceParams) error {
			_, _, err := tools.NewPowerOffHandler(m)(t.Context(), nil, p)

			return err
		}},
		{actionToggle, actionToggle, func(m *mockAPI, p tools.PowerDeviceParams) error {
			_, _, err := tools.NewPowerToggleHandler(m)(t.Context(), nil, p)

			return err
		}},
	}

	for _, testCase := range tests {
		t.Run(testCase.name, func(t *testing.T) {
			t.Parallel()

			missing := testCase.handler(&mockAPI{}, tools.PowerDeviceParams{})
			if !errors.Is(missing, tools.ErrValidation) {
				t.Errorf("missing device err = %v, want ErrValidation", missing)
			}

			mock := &mockAPI{result: okJSON}

			err := testCase.handler(mock, tools.PowerDeviceParams{Device: devicePrinter})
			if err != nil {
				t.Fatalf("handler: %v", err)
			}

			// on/off/toggle use the singular endpoint with a {device, action} body.
			assertCall(t, mock, methodPost, "/machine/device_power/device")

			body, ok := mock.lastBody.(map[string]any)
			if !ok {
				t.Fatalf("body type = %T, want map[string]any", mock.lastBody)
			}

			if body["device"] != devicePrinter || body["action"] != testCase.action {
				t.Errorf("body = %v, want device=%s action=%s", body, devicePrinter, testCase.action)
			}
		})
	}
}
