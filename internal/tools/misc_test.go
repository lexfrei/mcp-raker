package tools_test

import (
	"testing"

	"github.com/cockroachdb/errors"

	"github.com/lexfrei/mcp-raker/internal/tools"
)

func TestMQTTPublish_RequiresTopic(t *testing.T) {
	t.Parallel()

	_, _, err := tools.NewMQTTPublishHandler(&mockAPI{})(t.Context(), nil, tools.MQTTPublishParams{})
	if !errors.Is(err, tools.ErrValidation) {
		t.Fatalf("err = %v, want ErrValidation", err)
	}
}

func TestMQTTPublish(t *testing.T) {
	t.Parallel()

	mock := &mockAPI{result: okJSON}

	_, _, err := tools.NewMQTTPublishHandler(mock)(t.Context(), nil, tools.MQTTPublishParams{Topic: "klipper/status"})
	if err != nil {
		t.Fatalf("handler: %v", err)
	}

	assertCall(t, mock, methodPost, "/server/mqtt/publish")
}

func TestMQTTPublish_RejectsBadQOS(t *testing.T) {
	t.Parallel()

	params := tools.MQTTPublishParams{Topic: "t", QOS: 5}

	_, _, err := tools.NewMQTTPublishHandler(&mockAPI{})(t.Context(), nil, params)
	if !errors.Is(err, tools.ErrValidation) {
		t.Fatalf("err = %v, want ErrValidation for qos out of range", err)
	}
}

func TestMQTTSubscribe_RequiresTopic(t *testing.T) {
	t.Parallel()

	_, _, err := tools.NewMQTTSubscribeHandler(&mockAPI{})(t.Context(), nil, tools.MQTTSubscribeParams{})
	if !errors.Is(err, tools.ErrValidation) {
		t.Fatalf("err = %v, want ErrValidation", err)
	}
}

func TestExtensionsList(t *testing.T) {
	t.Parallel()

	mock := &mockAPI{result: okJSON}

	_, _, err := tools.NewExtensionsListHandler(mock)(t.Context(), nil, tools.NoParams{})
	if err != nil {
		t.Fatalf("handler: %v", err)
	}

	assertCall(t, mock, methodGet, "/server/extensions/list")
}

func TestExtensionsRequest_RequiresAgentAndMethod(t *testing.T) {
	t.Parallel()

	_, _, err := tools.NewExtensionsRequestHandler(&mockAPI{})(t.Context(), nil, tools.ExtensionsRequestParams{Agent: testAgent})
	if !errors.Is(err, tools.ErrValidation) {
		t.Fatalf("err = %v, want ErrValidation", err)
	}
}

func TestExtensionsRequest_OmitsNilArguments(t *testing.T) {
	t.Parallel()

	mock := &mockAPI{result: okJSON}
	params := tools.ExtensionsRequestParams{Agent: testAgent, Method: "ping"}

	_, _, err := tools.NewExtensionsRequestHandler(mock)(t.Context(), nil, params)
	if err != nil {
		t.Fatalf("handler: %v", err)
	}

	body, ok := mock.lastBody.(map[string]any)
	if !ok {
		t.Fatalf("body type = %T, want map[string]any", mock.lastBody)
	}

	if _, present := body["arguments"]; present {
		t.Errorf("body has arguments=%v, want the key omitted when nil", body["arguments"])
	}
}

func TestNotifiersList(t *testing.T) {
	t.Parallel()

	mock := &mockAPI{result: okJSON}

	_, _, err := tools.NewNotifiersListHandler(mock)(t.Context(), nil, tools.NoParams{})
	if err != nil {
		t.Fatalf("handler: %v", err)
	}

	assertCall(t, mock, methodGet, "/server/notifiers/list")
}
