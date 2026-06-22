package tools_test

import (
	"testing"

	"github.com/cockroachdb/errors"

	"github.com/lexfrei/mcp-raker/internal/moonraker"
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

func TestSensorsList_GracefulOnNotFound(t *testing.T) {
	t.Parallel()

	// Moonraker returns 404 when no [sensor] section is configured. The tool
	// should degrade to an empty result rather than surfacing the error.
	mock := &mockAPI{err: moonraker.ErrNotFound}

	_, out, err := tools.NewSensorsListHandler(mock)(t.Context(), nil, tools.SensorsListParams{})
	if err != nil {
		t.Fatalf("handler should not error on 404: %v", err)
	}

	sensors, ok := out["sensors"].(map[string]any)
	if !ok {
		t.Fatalf("out = %v, want a top-level empty sensors object", out)
	}

	if len(sensors) != 0 {
		t.Errorf("sensors = %v, want empty", sensors)
	}
}

func TestSensorsList_PropagatesOtherErrors(t *testing.T) {
	t.Parallel()

	mock := &mockAPI{err: errStub}

	_, _, err := tools.NewSensorsListHandler(mock)(t.Context(), nil, tools.SensorsListParams{})
	if !errors.Is(err, tools.ErrMoonraker) {
		t.Fatalf("err = %v, want ErrMoonraker for non-404 failures", err)
	}
}

func TestSensorsMeasurements_GracefulOnNotFoundForAll(t *testing.T) {
	t.Parallel()

	// No specific sensor requested: a 404 means no [sensor] section is configured.
	mock := &mockAPI{err: moonraker.ErrNotFound}

	_, out, err := tools.NewSensorsMeasurementsHandler(mock)(t.Context(), nil, tools.SensorsMeasurementsParams{})
	if err != nil {
		t.Fatalf("handler should not error on 404 for all sensors: %v", err)
	}

	if len(out) != 0 {
		t.Errorf("out = %v, want an empty measurements object", out)
	}
}

func TestSensorsMeasurements_PropagatesNotFoundForNamedSensor(t *testing.T) {
	t.Parallel()

	// A 404 for a specifically requested sensor is a real "no such sensor" error.
	mock := &mockAPI{err: moonraker.ErrNotFound}

	_, _, err := tools.NewSensorsMeasurementsHandler(mock)(t.Context(), nil, tools.SensorsMeasurementsParams{Sensor: testSensor})
	if !errors.Is(err, tools.ErrMoonraker) {
		t.Fatalf("err = %v, want the 404 to propagate for a named sensor", err)
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
