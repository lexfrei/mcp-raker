package tools_test

import (
	"testing"

	"github.com/cockroachdb/errors"

	"github.com/lexfrei/mcp-raker/internal/tools"
)

func TestWLEDStrips(t *testing.T) {
	t.Parallel()

	mock := &mockAPI{result: okJSON}

	_, _, err := tools.NewWLEDStripsHandler(mock)(t.Context(), nil, tools.NoParams{})
	if err != nil {
		t.Fatalf("handler: %v", err)
	}

	assertCall(t, mock, methodGet, "/machine/wled/strips")
}

func TestWLEDStripControls(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name    string
		method  string
		path    string
		handler func(*mockAPI, tools.StripParams) error
	}{
		{"status", methodGet, "/machine/wled/status", func(m *mockAPI, p tools.StripParams) error {
			_, _, err := tools.NewWLEDStatusHandler(m)(t.Context(), nil, p)

			return err
		}},
		{actionOn, methodPost, "/machine/wled/on", func(m *mockAPI, p tools.StripParams) error {
			_, _, err := tools.NewWLEDOnHandler(m)(t.Context(), nil, p)

			return err
		}},
		{actionOff, methodPost, "/machine/wled/off", func(m *mockAPI, p tools.StripParams) error {
			_, _, err := tools.NewWLEDOffHandler(m)(t.Context(), nil, p)

			return err
		}},
		{actionToggle, methodPost, "/machine/wled/toggle", func(m *mockAPI, p tools.StripParams) error {
			_, _, err := tools.NewWLEDToggleHandler(m)(t.Context(), nil, p)

			return err
		}},
	}

	for _, testCase := range tests {
		t.Run(testCase.name, func(t *testing.T) {
			t.Parallel()

			missing := testCase.handler(&mockAPI{}, tools.StripParams{})
			if !errors.Is(missing, tools.ErrValidation) {
				t.Errorf("missing strip err = %v, want ErrValidation", missing)
			}

			mock := &mockAPI{result: okJSON}

			err := testCase.handler(mock, tools.StripParams{Strip: testStrip})
			if err != nil {
				t.Fatalf("handler: %v", err)
			}

			assertCall(t, mock, testCase.method, testCase.path)
		})
	}
}

func TestWLEDSet(t *testing.T) {
	t.Parallel()

	zero := 0
	mock := &mockAPI{result: okJSON}
	params := tools.WLEDSetParams{Strip: testStrip, Preset: &zero, Brightness: 200, Intensity: &zero, Speed: &zero}

	_, _, err := tools.NewWLEDSetHandler(mock)(t.Context(), nil, params)
	if err != nil {
		t.Fatalf("handler: %v", err)
	}

	assertCall(t, mock, methodPost, "/machine/wled/strip")

	if mock.lastQuery.Get("brightness") != "200" {
		t.Errorf("brightness = %q, want 200", mock.lastQuery.Get("brightness"))
	}

	// Preset, intensity, and speed of 0 are valid values and must be sent.
	for _, key := range []string{"preset", "intensity", "speed"} {
		if !mock.lastQuery.Has(key) || mock.lastQuery.Get(key) != "0" {
			t.Errorf("%s = %q, want 0 to be sent", key, mock.lastQuery.Get(key))
		}
	}
}

func TestWLEDSet_RejectsOutOfRange(t *testing.T) {
	t.Parallel()

	high := 999
	low := -1

	cases := []tools.WLEDSetParams{
		{Strip: testStrip, Brightness: 999},
		{Strip: testStrip, Intensity: &high},
		{Strip: testStrip, Speed: &low},
	}

	for _, params := range cases {
		_, _, err := tools.NewWLEDSetHandler(&mockAPI{})(t.Context(), nil, params)
		if !errors.Is(err, tools.ErrValidation) {
			t.Errorf("params %+v: err = %v, want ErrValidation", params, err)
		}
	}
}

func TestWLEDSet_OmitsUnsetPointers(t *testing.T) {
	t.Parallel()

	mock := &mockAPI{result: okJSON}

	_, _, err := tools.NewWLEDSetHandler(mock)(t.Context(), nil, tools.WLEDSetParams{Strip: testStrip})
	if err != nil {
		t.Fatalf("handler: %v", err)
	}

	for _, key := range []string{"preset", "intensity", "speed"} {
		if mock.lastQuery.Has(key) {
			t.Errorf("query has %s, want it omitted when nil", key)
		}
	}
}
