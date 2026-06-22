package tools_test

import (
	"testing"

	"github.com/lexfrei/mcp-raker/internal/tools"
)

func TestMQTTSubscribe_OmitsZeroTimeout(t *testing.T) {
	t.Parallel()

	mock := &mockAPI{result: okJSON}

	_, _, err := tools.NewMQTTSubscribeHandler(mock)(t.Context(), nil, tools.MQTTSubscribeParams{Topic: "klipper/state"})
	if err != nil {
		t.Fatalf("handler: %v", err)
	}

	body, ok := mock.lastBody.(map[string]any)
	if !ok {
		t.Fatalf("body = %T, want map", mock.lastBody)
	}

	// A zero timeout must not be sent: it would make Moonraker time out at once.
	if _, has := body["timeout"]; has {
		t.Errorf("body = %v, want no timeout field when unset", body)
	}
}

func TestMQTTSubscribe_SendsTimeoutWhenSet(t *testing.T) {
	t.Parallel()

	mock := &mockAPI{result: okJSON}

	params := tools.MQTTSubscribeParams{Topic: "klipper/state", Timeout: 5}

	_, _, err := tools.NewMQTTSubscribeHandler(mock)(t.Context(), nil, params)
	if err != nil {
		t.Fatalf("handler: %v", err)
	}

	body, ok := mock.lastBody.(map[string]any)
	if !ok {
		t.Fatalf("body = %T, want map", mock.lastBody)
	}

	if body["timeout"] != float64(5) {
		t.Errorf("timeout = %v, want 5", body["timeout"])
	}
}
