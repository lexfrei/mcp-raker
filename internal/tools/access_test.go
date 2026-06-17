package tools_test

import (
	"testing"

	"github.com/cockroachdb/errors"

	"github.com/lexfrei/mcp-raker/internal/tools"
)

const testUser = "bob"

func TestAccessReads(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name    string
		handler func(*mockAPI) error
		path    string
	}{
		{"user_info", func(m *mockAPI) error {
			_, _, err := tools.NewAccessUserInfoHandler(m)(t.Context(), nil, tools.NoParams{})

			return err
		}, "/access/user"},
		{"users_list", func(m *mockAPI) error {
			_, _, err := tools.NewAccessUsersListHandler(m)(t.Context(), nil, tools.NoParams{})

			return err
		}, "/access/users/list"},
		{"info", func(m *mockAPI) error {
			_, _, err := tools.NewAccessInfoHandler(m)(t.Context(), nil, tools.NoParams{})

			return err
		}, "/access/info"},
		{"api_key", func(m *mockAPI) error {
			_, _, err := tools.NewAccessAPIKeyHandler(m)(t.Context(), nil, tools.NoParams{})

			return err
		}, "/access/api_key"},
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

func TestAccessCreateUser(t *testing.T) {
	t.Parallel()

	missing := func() error {
		_, _, err := tools.NewAccessCreateUserHandler(&mockAPI{})(t.Context(), nil, tools.AccessCreateUserParams{Username: testUser})

		return err
	}()
	if !errors.Is(missing, tools.ErrValidation) {
		t.Errorf("missing password err = %v, want ErrValidation", missing)
	}

	mock := &mockAPI{result: okJSON}
	params := tools.AccessCreateUserParams{Username: testUser, Password: "pw"}

	_, _, err := tools.NewAccessCreateUserHandler(mock)(t.Context(), nil, params)
	if err != nil {
		t.Fatalf("handler: %v", err)
	}

	assertCall(t, mock, methodPost, "/access/user")
}

func TestAccessDeleteUser(t *testing.T) {
	t.Parallel()

	_, _, missing := tools.NewAccessDeleteUserHandler(&mockAPI{})(t.Context(), nil, tools.AccessDeleteUserParams{})
	if !errors.Is(missing, tools.ErrValidation) {
		t.Errorf("missing username err = %v, want ErrValidation", missing)
	}

	mock := &mockAPI{result: okJSON}

	_, _, err := tools.NewAccessDeleteUserHandler(mock)(t.Context(), nil, tools.AccessDeleteUserParams{Username: testUser})
	if err != nil {
		t.Fatalf("handler: %v", err)
	}

	assertCall(t, mock, methodDelete, "/access/user")
}

func TestAccessUserPassword_RequiresBoth(t *testing.T) {
	t.Parallel()

	_, _, err := tools.NewAccessUserPasswordHandler(&mockAPI{})(t.Context(), nil, tools.AccessUserPasswordParams{Password: "old"})
	if !errors.Is(err, tools.ErrValidation) {
		t.Fatalf("err = %v, want ErrValidation", err)
	}
}

func TestAccessCreateAPIKey(t *testing.T) {
	t.Parallel()

	mock := &mockAPI{result: okJSON}

	_, _, err := tools.NewAccessCreateAPIKeyHandler(mock)(t.Context(), nil, tools.NoParams{})
	if err != nil {
		t.Fatalf("handler: %v", err)
	}

	assertCall(t, mock, methodPost, "/access/api_key")
}
