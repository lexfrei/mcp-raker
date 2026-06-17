package moonraker

import (
	"bytes"
	"context"
	"encoding/json"
	"io"
	"log/slog"
	"net/http"
	"os"
	"path/filepath"

	"github.com/cockroachdb/errors"
)

// Permissions for the persisted token cache: private directory, private file.
const (
	tokenDirPerm  = 0o700
	tokenFilePerm = 0o600
)

// loginResult is the payload returned by /access/login and /access/refresh_jwt.
type loginResult struct {
	Token        string `json:"token"`
	RefreshToken string `json:"refresh_token"`
}

// persistedToken is the on-disk token cache format.
type persistedToken struct {
	Token        string `json:"token"`
	RefreshToken string `json:"refresh_token"`
}

// ensureAuth guarantees a Bearer token is available, logging in on first use.
func (c *Client) ensureAuth(ctx context.Context) error {
	if c.currentToken() != "" {
		return nil
	}

	c.authMu.Lock()
	defer c.authMu.Unlock()

	// Another goroutine may have logged in while we waited for the lock.
	if c.currentToken() != "" {
		return nil
	}

	return c.login(ctx)
}

// reauth refreshes a rejected token. stale is the token that just failed; if
// another goroutine already refreshed it while we waited for the lock, reauth is
// a no-op. It tries the refresh token first and falls back to a full login.
func (c *Client) reauth(ctx context.Context, stale string) error {
	c.authMu.Lock()
	defer c.authMu.Unlock()

	if c.currentToken() != stale {
		return nil
	}

	if c.username == "" || c.password == "" {
		return ErrNotAuthenticated
	}

	if c.refreshTokenValue() != "" {
		refreshErr := c.refresh(ctx)
		if refreshErr == nil {
			return nil
		}
	}

	return c.login(ctx)
}

// login exchanges username/password for a JWT via /access/login and persists it.
// The caller must hold authMu.
func (c *Client) login(ctx context.Context) error {
	if c.username == "" || c.password == "" {
		return ErrNoCredentials
	}

	body := map[string]string{
		"username": c.username,
		"password": c.password,
		"source":   "moonraker",
	}

	var result loginResult

	err := c.authRequest(ctx, "/access/login", body, &result)
	if err != nil {
		return err
	}

	if result.Token == "" {
		return ErrLoginFailed
	}

	c.storeTokens(result.Token, result.RefreshToken)

	return nil
}

// refresh exchanges the stored refresh token for a new access token via
// /access/refresh_jwt. The caller must hold authMu.
func (c *Client) refresh(ctx context.Context) error {
	refreshToken := c.refreshTokenValue()
	if refreshToken == "" {
		return ErrLoginFailed
	}

	body := map[string]string{"refresh_token": refreshToken}

	var result loginResult

	err := c.authRequest(ctx, "/access/refresh_jwt", body, &result)
	if err != nil {
		return err
	}

	if result.Token == "" {
		return ErrLoginFailed
	}

	// refresh_jwt does not issue a new refresh token, so keep the existing one.
	newRefresh := result.RefreshToken
	if newRefresh == "" {
		newRefresh = refreshToken
	}

	c.storeTokens(result.Token, newRefresh)

	return nil
}

// authRequest posts an unauthenticated JSON request to an /access endpoint and
// decodes the result envelope into out. It never attaches a Bearer header so a
// stale token cannot interfere with login or refresh.
func (c *Client) authRequest(ctx context.Context, path string, body, out any) error {
	payload, err := json.Marshal(body)
	if err != nil {
		return errors.Wrap(err, "marshal auth request")
	}

	req, reqErr := http.NewRequestWithContext(ctx, http.MethodPost, c.baseURL+path, bytes.NewReader(payload))
	if reqErr != nil {
		return errors.Wrap(reqErr, "build auth request")
	}

	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Accept", "application/json")
	req.Header.Set("User-Agent", c.userAgent)

	if c.apiKey != "" {
		req.Header.Set("X-Api-Key", c.apiKey)
	}

	resp, doErr := c.http.Do(req)
	if doErr != nil {
		return loginErr(doErr, "auth request")
	}
	defer func() { _ = resp.Body.Close() }()

	return decodeAuthResponse(resp, out)
}

// decodeAuthResponse classifies the auth status and unwraps the result.
func decodeAuthResponse(resp *http.Response, out any) error {
	data, readErr := io.ReadAll(resp.Body)
	if readErr != nil {
		return loginErr(readErr, "read auth response")
	}

	if resp.StatusCode < http.StatusOK || resp.StatusCode >= http.StatusMultipleChoices {
		return errors.Wrapf(ErrLoginFailed, "status %d", resp.StatusCode)
	}

	var envelope resultEnvelope

	envErr := json.Unmarshal(data, &envelope)
	if envErr != nil {
		return loginErr(envErr, "decode auth response")
	}

	unErr := json.Unmarshal(envelope.Result, out)
	if unErr != nil {
		return loginErr(unErr, "decode auth result")
	}

	return nil
}

// loadToken seeds the token cache from disk, best-effort.
func (c *Client) loadToken() {
	if c.tokenPath == "" {
		return
	}

	data, err := os.ReadFile(c.tokenPath)
	if err != nil {
		return
	}

	var stored persistedToken

	unErr := json.Unmarshal(data, &stored)
	if unErr != nil {
		return
	}

	c.mu.Lock()
	if stored.Token != "" {
		c.token = stored.Token
	}

	if stored.RefreshToken != "" {
		c.refreshToken = stored.RefreshToken
	}
	c.mu.Unlock()
}

// saveToken persists the token cache. It is best-effort but logs a warning on
// failure, since a cache that cannot be written forces a fresh login on every
// run. A disabled cache (empty token path) is skipped silently.
func (c *Client) saveToken(token, refresh string) {
	if c.tokenPath == "" {
		return
	}

	//nolint:gosec // G117: persisting the refresh token to a private 0600 cache file is the point of this method.
	data, err := json.Marshal(persistedToken{Token: token, RefreshToken: refresh})
	if err != nil {
		c.logger.Warn("failed to marshal token cache", slog.Any("error", err))

		return
	}

	mkErr := os.MkdirAll(filepath.Dir(c.tokenPath), tokenDirPerm)
	if mkErr != nil {
		c.logger.Warn("failed to create token cache directory",
			slog.String("path", c.tokenPath), slog.Any("error", mkErr))

		return
	}

	writeErr := os.WriteFile(c.tokenPath, data, tokenFilePerm)
	if writeErr != nil {
		c.logger.Warn("failed to write token cache",
			slog.String("path", c.tokenPath), slog.Any("error", writeErr))
	}
}
