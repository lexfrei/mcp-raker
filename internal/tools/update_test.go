package tools_test

import (
	"testing"

	"github.com/cockroachdb/errors"

	"github.com/lexfrei/mcp-raker/internal/tools"
)

func TestUpdateStatus(t *testing.T) {
	t.Parallel()

	mock := &mockAPI{result: okJSON}

	_, _, err := tools.NewUpdateStatusHandler(mock)(t.Context(), nil, tools.NoParams{})
	if err != nil {
		t.Fatalf("handler: %v", err)
	}

	assertCall(t, mock, methodGet, "/machine/update/status")
}

func TestUpdateUpgrade_OptionalName(t *testing.T) {
	t.Parallel()

	mock := &mockAPI{result: okJSON}

	_, _, err := tools.NewUpdateUpgradeHandler(mock)(t.Context(), nil, tools.UpdateNameParams{Name: "klipper"})
	if err != nil {
		t.Fatalf("handler: %v", err)
	}

	assertCall(t, mock, methodPost, "/machine/update/upgrade")

	if mock.lastQuery.Get("name") != svcKlipper {
		t.Errorf("name = %q, want %s", mock.lastQuery.Get("name"), svcKlipper)
	}
}

func TestUpdateRecover_RequiresName(t *testing.T) {
	t.Parallel()

	_, _, err := tools.NewUpdateRecoverHandler(&mockAPI{})(t.Context(), nil, tools.UpdateRecoverParams{})
	if !errors.Is(err, tools.ErrValidation) {
		t.Fatalf("err = %v, want ErrValidation", err)
	}
}

func TestUpdateRecover_Hard(t *testing.T) {
	t.Parallel()

	mock := &mockAPI{result: okJSON}
	params := tools.UpdateRecoverParams{Name: svcKlipper, Hard: true}

	_, _, err := tools.NewUpdateRecoverHandler(mock)(t.Context(), nil, params)
	if err != nil {
		t.Fatalf("handler: %v", err)
	}

	assertCall(t, mock, methodPost, "/machine/update/recover")

	if mock.lastQuery.Get("hard") != queryTrue {
		t.Errorf("hard = %q, want true", mock.lastQuery.Get("hard"))
	}
}

func TestUpdateRollback_RequiresName(t *testing.T) {
	t.Parallel()

	_, _, err := tools.NewUpdateRollbackHandler(&mockAPI{})(t.Context(), nil, tools.UpdateRollbackParams{})
	if !errors.Is(err, tools.ErrValidation) {
		t.Fatalf("err = %v, want ErrValidation", err)
	}
}

func TestUpdateRollback_SendsName(t *testing.T) {
	t.Parallel()

	mock := &mockAPI{result: okJSON}

	_, _, err := tools.NewUpdateRollbackHandler(mock)(t.Context(), nil, tools.UpdateRollbackParams{Name: svcKlipper})
	if err != nil {
		t.Fatalf("handler: %v", err)
	}

	assertCall(t, mock, methodPost, "/machine/update/rollback")

	if mock.lastQuery.Get("name") != svcKlipper {
		t.Errorf("name = %q, want %s", mock.lastQuery.Get("name"), svcKlipper)
	}
}

func TestUpdateRefresh(t *testing.T) {
	t.Parallel()

	mock := &mockAPI{result: okJSON}

	_, _, err := tools.NewUpdateRefreshHandler(mock)(t.Context(), nil, tools.UpdateNameParams{})
	if err != nil {
		t.Fatalf("handler: %v", err)
	}

	assertCall(t, mock, methodPost, "/machine/update/refresh")
}
