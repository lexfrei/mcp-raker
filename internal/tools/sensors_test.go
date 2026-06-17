package tools_test

import (
	"testing"

	"github.com/cockroachdb/errors"

	"github.com/lexfrei/mcp-raker/internal/tools"
)

func TestSensorsList_Extended(t *testing.T) {
	t.Parallel()

	mock := &mockAPI{result: okJSON}

	_, _, err := tools.NewSensorsListHandler(mock)(t.Context(), nil, tools.SensorsListParams{Extended: true})
	if err != nil {
		t.Fatalf("handler: %v", err)
	}

	assertCall(t, mock, methodGet, "/machine/sensors/list")

	if mock.lastQuery.Get("extended") != queryTrue {
		t.Errorf("extended = %q, want true", mock.lastQuery.Get("extended"))
	}
}

func TestSensorsInfo_RequiresSensor(t *testing.T) {
	t.Parallel()

	_, _, err := tools.NewSensorsInfoHandler(&mockAPI{})(t.Context(), nil, tools.SensorParams{})
	if !errors.Is(err, tools.ErrValidation) {
		t.Fatalf("err = %v, want ErrValidation", err)
	}
}

func TestSensorsMeasurements(t *testing.T) {
	t.Parallel()

	mock := &mockAPI{result: okJSON}

	_, _, err := tools.NewSensorsMeasurementsHandler(mock)(t.Context(), nil, tools.SensorsMeasurementsParams{Sensor: testSensor})
	if err != nil {
		t.Fatalf("handler: %v", err)
	}

	assertCall(t, mock, methodGet, "/machine/sensors/measurements")

	if mock.lastQuery.Get("sensor") != testSensor {
		t.Errorf("sensor = %q, want chamber", mock.lastQuery.Get("sensor"))
	}
}
