package tools

import (
	"context"
	"net/url"

	"github.com/modelcontextprotocol/go-sdk/mcp"

	"github.com/lexfrei/mcp-raker/internal/moonraker"
)

// AccessUserInfoTool returns the definition for moonraker_access_user_info.
func AccessUserInfoTool() *mcp.Tool {
	return &mcp.Tool{
		Name:        "moonraker_access_user_info",
		Description: "Get the currently authenticated user (GET /access/user).",
		Annotations: readOnly("Current User"),
	}
}

// NewAccessUserInfoHandler creates the handler for moonraker_access_user_info.
func NewAccessUserInfoHandler(api moonraker.API) mcp.ToolHandlerFor[NoParams, RawResult] {
	return func(ctx context.Context, _ *mcp.CallToolRequest, _ NoParams) (*mcp.CallToolResult, RawResult, error) {
		out, err := decodeRaw(api.Get(ctx, "/access/user", nil))

		return nil, out, err
	}
}

// AccessUsersListTool returns the definition for moonraker_access_users_list.
func AccessUsersListTool() *mcp.Tool {
	return &mcp.Tool{
		Name:        "moonraker_access_users_list",
		Description: "List all registered users (GET /access/users/list).",
		Annotations: readOnly("List Users"),
	}
}

// NewAccessUsersListHandler creates the handler for moonraker_access_users_list.
func NewAccessUsersListHandler(api moonraker.API) mcp.ToolHandlerFor[NoParams, RawResult] {
	return func(ctx context.Context, _ *mcp.CallToolRequest, _ NoParams) (*mcp.CallToolResult, RawResult, error) {
		out, err := decodeRaw(api.Get(ctx, "/access/users/list", nil))

		return nil, out, err
	}
}

// AccessInfoTool returns the definition for moonraker_access_info.
func AccessInfoTool() *mcp.Tool {
	return &mcp.Tool{
		Name:        "moonraker_access_info",
		Description: "Get authentication configuration and the caller's trust status (GET /access/info).",
		Annotations: readOnly("Auth Info"),
	}
}

// NewAccessInfoHandler creates the handler for moonraker_access_info.
func NewAccessInfoHandler(api moonraker.API) mcp.ToolHandlerFor[NoParams, RawResult] {
	return func(ctx context.Context, _ *mcp.CallToolRequest, _ NoParams) (*mcp.CallToolResult, RawResult, error) {
		out, err := decodeRaw(api.Get(ctx, "/access/info", nil))

		return nil, out, err
	}
}

// AccessAPIKeyTool returns the definition for moonraker_access_api_key.
func AccessAPIKeyTool() *mcp.Tool {
	return &mcp.Tool{
		Name:        "moonraker_access_api_key",
		Description: "Get the current API key (GET /access/api_key).",
		Annotations: readOnly("Get API Key"),
	}
}

// NewAccessAPIKeyHandler creates the handler for moonraker_access_api_key.
func NewAccessAPIKeyHandler(api moonraker.API) mcp.ToolHandlerFor[NoParams, RawResult] {
	return func(ctx context.Context, _ *mcp.CallToolRequest, _ NoParams) (*mcp.CallToolResult, RawResult, error) {
		out, err := decodeRaw(api.Get(ctx, "/access/api_key", nil))

		return nil, out, err
	}
}

// AccessCreateUserParams defines the parameters for moonraker_access_create_user (admin).
type AccessCreateUserParams struct {
	Username string `json:"username" jsonschema:"Username for the new local account"`
	Password string `json:"password" jsonschema:"Password for the new local account"`
}

// AccessCreateUserTool returns the definition for moonraker_access_create_user (admin).
func AccessCreateUserTool() *mcp.Tool {
	return &mcp.Tool{
		Name:        "moonraker_access_create_user",
		Description: "Create a new local user account (POST /access/user).",
		Annotations: writeDestructive("Create User"),
	}
}

// NewAccessCreateUserHandler creates the handler for moonraker_access_create_user.
func NewAccessCreateUserHandler(api moonraker.API) mcp.ToolHandlerFor[AccessCreateUserParams, RawResult] {
	return func(ctx context.Context, _ *mcp.CallToolRequest, params AccessCreateUserParams) (*mcp.CallToolResult, RawResult, error) {
		userErr := requireString(paramUsername, params.Username)
		if userErr != nil {
			return nil, RawResult{}, userErr
		}

		passErr := requireString(paramPassword, params.Password)
		if passErr != nil {
			return nil, RawResult{}, passErr
		}

		body := map[string]any{paramUsername: params.Username, paramPassword: params.Password}

		out, err := decodeRaw(api.Post(ctx, "/access/user", nil, body))

		return nil, out, err
	}
}

// AccessDeleteUserParams names a user to delete.
type AccessDeleteUserParams struct {
	Username string `json:"username" jsonschema:"Username of the account to delete"`
}

// AccessDeleteUserTool returns the definition for moonraker_access_delete_user (admin).
func AccessDeleteUserTool() *mcp.Tool {
	return &mcp.Tool{
		Name:        "moonraker_access_delete_user",
		Description: "Delete a local user account (DELETE /access/user).",
		Annotations: writeDestructive("Delete User"),
	}
}

// NewAccessDeleteUserHandler creates the handler for moonraker_access_delete_user.
func NewAccessDeleteUserHandler(api moonraker.API) mcp.ToolHandlerFor[AccessDeleteUserParams, RawResult] {
	return func(ctx context.Context, _ *mcp.CallToolRequest, params AccessDeleteUserParams) (*mcp.CallToolResult, RawResult, error) {
		valErr := requireString(paramUsername, params.Username)
		if valErr != nil {
			return nil, RawResult{}, valErr
		}

		out, err := decodeRaw(api.Delete(ctx, "/access/user", url.Values{paramUsername: {params.Username}}))

		return nil, out, err
	}
}

// AccessUserPasswordParams defines the parameters for moonraker_access_user_password (admin).
type AccessUserPasswordParams struct {
	Password    string `json:"password"     jsonschema:"The current password"`
	NewPassword string `json:"new_password" jsonschema:"The new password to set"`
}

// AccessUserPasswordTool returns the definition for moonraker_access_user_password (admin).
func AccessUserPasswordTool() *mcp.Tool {
	return &mcp.Tool{
		Name:        "moonraker_access_user_password",
		Description: "Change the current user's password (POST /access/user/password).",
		Annotations: writeDestructive("Change Password"),
	}
}

// NewAccessUserPasswordHandler creates the handler for moonraker_access_user_password.
func NewAccessUserPasswordHandler(api moonraker.API) mcp.ToolHandlerFor[AccessUserPasswordParams, RawResult] {
	return func(ctx context.Context, _ *mcp.CallToolRequest, params AccessUserPasswordParams) (*mcp.CallToolResult, RawResult, error) {
		passErr := requireString(paramPassword, params.Password)
		if passErr != nil {
			return nil, RawResult{}, passErr
		}

		newErr := requireString("new_password", params.NewPassword)
		if newErr != nil {
			return nil, RawResult{}, newErr
		}

		body := map[string]any{paramPassword: params.Password, "new_password": params.NewPassword}

		out, err := decodeRaw(api.Post(ctx, "/access/user/password", nil, body))

		return nil, out, err
	}
}

// AccessCreateAPIKeyTool returns the definition for moonraker_access_create_api_key (admin).
func AccessCreateAPIKeyTool() *mcp.Tool {
	return &mcp.Tool{
		Name:        "moonraker_access_create_api_key",
		Description: "Generate a new API key, replacing the previous one (POST /access/api_key).",
		Annotations: writeDestructive("Regenerate API Key"),
	}
}

// NewAccessCreateAPIKeyHandler creates the handler for moonraker_access_create_api_key.
func NewAccessCreateAPIKeyHandler(api moonraker.API) mcp.ToolHandlerFor[NoParams, RawResult] {
	return func(ctx context.Context, _ *mcp.CallToolRequest, _ NoParams) (*mcp.CallToolResult, RawResult, error) {
		out, err := decodeRaw(api.Post(ctx, "/access/api_key", nil, nil))

		return nil, out, err
	}
}
