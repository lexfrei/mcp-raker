package tools_test

import (
	"testing"

	"github.com/cockroachdb/errors"

	"github.com/lexfrei/mcp-raker/internal/tools"
)

const testNS = "ns"

func TestDBList(t *testing.T) {
	t.Parallel()

	mock := &mockAPI{result: okJSON}

	_, _, err := tools.NewDBListHandler(mock)(t.Context(), nil, tools.NoParams{})
	if err != nil {
		t.Fatalf("handler: %v", err)
	}

	assertCall(t, mock, methodGet, "/server/database/list")
}

func TestDBGetItem_RequiresNamespace(t *testing.T) {
	t.Parallel()

	_, _, err := tools.NewDBGetItemHandler(&mockAPI{})(t.Context(), nil, tools.DBGetItemParams{})
	if !errors.Is(err, tools.ErrValidation) {
		t.Fatalf("err = %v, want ErrValidation", err)
	}
}

func TestDBPostItem(t *testing.T) {
	t.Parallel()

	missing := func() error {
		_, _, err := tools.NewDBPostItemHandler(&mockAPI{})(t.Context(), nil, tools.DBPostItemParams{Namespace: testNS})

		return err
	}()
	if !errors.Is(missing, tools.ErrValidation) {
		t.Errorf("missing key err = %v, want ErrValidation", missing)
	}

	mock := &mockAPI{result: okJSON}
	params := tools.DBPostItemParams{Namespace: testNS, Key: "k", Value: 42}

	_, _, err := tools.NewDBPostItemHandler(mock)(t.Context(), nil, params)
	if err != nil {
		t.Fatalf("handler: %v", err)
	}

	assertCall(t, mock, methodPost, "/server/database/item")
}

func TestDBDeleteItem_RequiresKeys(t *testing.T) {
	t.Parallel()

	_, _, err := tools.NewDBDeleteItemHandler(&mockAPI{})(t.Context(), nil, tools.DBDeleteItemParams{Namespace: testNS})
	if !errors.Is(err, tools.ErrValidation) {
		t.Fatalf("err = %v, want ErrValidation", err)
	}
}

func TestDBBackupAndCompact(t *testing.T) {
	t.Parallel()

	backup := &mockAPI{result: okJSON}

	_, _, backupErr := tools.NewDBBackupHandler(backup)(t.Context(), nil, tools.DBBackupParams{})
	if backupErr != nil {
		t.Fatalf("backup: %v", backupErr)
	}

	assertCall(t, backup, methodPost, "/server/database/backup")

	compact := &mockAPI{result: okJSON}

	_, _, compactErr := tools.NewDBCompactHandler(compact)(t.Context(), nil, tools.NoParams{})
	if compactErr != nil {
		t.Fatalf("compact: %v", compactErr)
	}

	assertCall(t, compact, methodPost, "/server/database/compact")
}
