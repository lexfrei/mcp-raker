package moonraker_test

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"log/slog"
	"net/http"
	"net/http/httptest"
	"os"
	"path/filepath"
	"strings"
	"sync"
	"testing"

	"github.com/lexfrei/mcp-raker/internal/moonraker"
)

const (
	loginPath = "/access/login"
	testUser  = "user"
	testPass  = "pass"
)

// newClient builds a client pointed at srv with the given options applied.
func newClient(t *testing.T, srv *httptest.Server, opts *moonraker.Options) *moonraker.Client {
	t.Helper()

	if opts == nil {
		opts = &moonraker.Options{}
	}

	opts.BaseURL = srv.URL

	client, err := moonraker.New(opts)
	if err != nil {
		t.Fatalf("New: %v", err)
	}

	return client
}

func TestGet_UnwrapsResultEnvelope(t *testing.T) {
	t.Parallel()

	srv := httptest.NewServer(http.HandlerFunc(func(writer http.ResponseWriter, request *http.Request) {
		if request.URL.Path != "/server/info" {
			t.Errorf("path = %q, want /server/info", request.URL.Path)
		}

		_, _ = writer.Write([]byte(`{"result":{"klippy_state":"ready","moonraker_version":"v0.9.0"}}`))
	}))
	defer srv.Close()

	client := newClient(t, srv, nil)

	raw, err := client.Get(t.Context(), "/server/info", nil)
	if err != nil {
		t.Fatalf("Get: %v", err)
	}

	var info moonraker.ServerInfo

	unErr := json.Unmarshal(raw, &info)
	if unErr != nil {
		t.Fatalf("unmarshal: %v", unErr)
	}

	if info.KlippyState != "ready" || info.MoonrakerVersion != "v0.9.0" {
		t.Errorf("info = %+v, want klippy_state=ready moonraker_version=v0.9.0", info)
	}
}

func TestGet_SendsAPIKeyHeader(t *testing.T) {
	t.Parallel()

	var gotKey string

	srv := httptest.NewServer(http.HandlerFunc(func(writer http.ResponseWriter, request *http.Request) {
		gotKey = request.Header.Get("X-Api-Key")

		_, _ = writer.Write([]byte(`{"result":"ok"}`))
	}))
	defer srv.Close()

	client := newClient(t, srv, &moonraker.Options{APIKey: "secret-key"})

	_, err := client.Get(t.Context(), "/server/info", nil)
	if err != nil {
		t.Fatalf("Get: %v", err)
	}

	if gotKey != "secret-key" {
		t.Errorf("X-Api-Key = %q, want secret-key", gotKey)
	}
}

func TestGet_StatusErrorSurfacesMessage(t *testing.T) {
	t.Parallel()

	srv := httptest.NewServer(http.HandlerFunc(func(writer http.ResponseWriter, _ *http.Request) {
		writer.WriteHeader(http.StatusBadRequest)
		_, _ = writer.Write([]byte(`{"error":{"code":400,"message":"bad request"}}`))
	}))
	defer srv.Close()

	client := newClient(t, srv, nil)

	_, err := client.Get(t.Context(), "/printer/print/start", nil)
	if !errors.Is(err, moonraker.ErrAPI) {
		t.Fatalf("err = %v, want ErrAPI", err)
	}
}

func TestSend_JWTLoginAndRefreshOn401(t *testing.T) {
	t.Parallel()

	var (
		logins    int
		refreshes int
		authHdrs  []string
	)

	srv := httptest.NewServer(http.HandlerFunc(func(writer http.ResponseWriter, request *http.Request) {
		switch request.URL.Path {
		case loginPath:
			logins++

			_, _ = writer.Write([]byte(`{"result":{"token":"jwt-1","refresh_token":"refresh-1"}}`))
		case "/access/refresh_jwt":
			refreshes++

			_, _ = writer.Write([]byte(`{"result":{"token":"jwt-2"}}`))
		default:
			authHdrs = append(authHdrs, request.Header.Get("Authorization"))
			// Reject the first authed call so the client refreshes the token.
			if len(authHdrs) == 1 {
				writer.WriteHeader(http.StatusUnauthorized)

				return
			}

			_, _ = writer.Write([]byte(`{"result":"ok"}`))
		}
	}))
	defer srv.Close()

	client := newClient(t, srv, &moonraker.Options{Username: testUser, Password: testPass})

	_, err := client.Get(t.Context(), "/printer/info", nil)
	if err != nil {
		t.Fatalf("Get: %v", err)
	}

	if logins != 1 {
		t.Errorf("logins = %d, want 1", logins)
	}

	if refreshes != 1 {
		t.Errorf("refreshes = %d, want 1", refreshes)
	}

	if len(authHdrs) != 2 || authHdrs[0] != "Bearer jwt-1" || authHdrs[1] != "Bearer jwt-2" {
		t.Errorf("auth headers = %v, want [Bearer jwt-1, Bearer jwt-2]", authHdrs)
	}
}

func TestGet_APIKeyTakesPrecedenceNoReauth(t *testing.T) {
	t.Parallel()

	var logins int

	srv := httptest.NewServer(http.HandlerFunc(func(writer http.ResponseWriter, request *http.Request) {
		if request.URL.Path == loginPath {
			logins++
		}
		// Always reject so a refresh would be attempted if reauth were enabled.
		writer.WriteHeader(http.StatusUnauthorized)
	}))
	defer srv.Close()

	// Both an API key and credentials are set; the API key must win and a 401
	// must be terminal rather than triggering a username/password login.
	client := newClient(t, srv, &moonraker.Options{APIKey: "key", Username: testUser, Password: testPass})

	_, err := client.Get(t.Context(), "/printer/info", nil)
	if !errors.Is(err, moonraker.ErrNotAuthenticated) {
		t.Fatalf("err = %v, want ErrNotAuthenticated", err)
	}

	if logins != 0 {
		t.Errorf("logins = %d, want 0 (API key must not trigger a login)", logins)
	}
}

func TestSend_RawTokenWithCredentialsLogsInOn401(t *testing.T) {
	t.Parallel()

	var (
		logins   int
		authHdrs []string
	)

	srv := httptest.NewServer(http.HandlerFunc(func(writer http.ResponseWriter, request *http.Request) {
		if request.URL.Path == loginPath {
			logins++

			_, _ = writer.Write([]byte(`{"result":{"token":"fresh","refresh_token":"r"}}`))

			return
		}

		authHdrs = append(authHdrs, request.Header.Get("Authorization"))
		// Reject the supplied raw token once; the client has credentials, so it
		// logs in and retries with the freshly minted token.
		if len(authHdrs) == 1 {
			writer.WriteHeader(http.StatusUnauthorized)

			return
		}

		_, _ = writer.Write([]byte(`{"result":"ok"}`))
	}))
	defer srv.Close()

	client := newClient(t, srv, &moonraker.Options{Token: "stale", Username: testUser, Password: testPass})

	_, err := client.Get(t.Context(), "/printer/info", nil)
	if err != nil {
		t.Fatalf("Get: %v", err)
	}

	if logins != 1 {
		t.Errorf("logins = %d, want 1", logins)
	}

	if len(authHdrs) != 2 || authHdrs[0] != "Bearer stale" || authHdrs[1] != "Bearer fresh" {
		t.Errorf("auth headers = %v, want [Bearer stale, Bearer fresh]", authHdrs)
	}
}

func TestRawTokenNotPersisted(t *testing.T) {
	t.Parallel()

	srv := httptest.NewServer(http.HandlerFunc(func(writer http.ResponseWriter, _ *http.Request) {
		_, _ = writer.Write([]byte(`{"result":"ok"}`))
	}))
	defer srv.Close()

	tokenPath := filepath.Join(t.TempDir(), "token.json")
	client := newClient(t, srv, &moonraker.Options{Token: "raw", TokenPath: tokenPath})

	_, err := client.Get(t.Context(), "/printer/info", nil)
	if err != nil {
		t.Fatalf("Get: %v", err)
	}

	_, statErr := os.Stat(tokenPath)
	if !os.IsNotExist(statErr) {
		t.Errorf("token file %q exists; a raw MOONRAKER_TOKEN must not be persisted", tokenPath)
	}
}

func TestReauth_ConcurrentSingleFlight(t *testing.T) {
	t.Parallel()

	var (
		mu     sync.Mutex
		logins int
	)

	srv := httptest.NewServer(http.HandlerFunc(func(writer http.ResponseWriter, request *http.Request) {
		if request.URL.Path == loginPath {
			mu.Lock()
			logins++
			mu.Unlock()

			_, _ = writer.Write([]byte(`{"result":{"token":"fresh","refresh_token":"r"}}`))

			return
		}

		if request.Header.Get("Authorization") == "Bearer stale" {
			writer.WriteHeader(http.StatusUnauthorized)

			return
		}

		_, _ = writer.Write([]byte(`{"result":"ok"}`))
	}))
	defer srv.Close()

	client := newClient(t, srv, &moonraker.Options{Token: "stale", Username: testUser, Password: testPass})

	const workers = 8

	var wg sync.WaitGroup

	errs := make([]error, workers)

	for i := range workers {
		wg.Add(1)

		go func(idx int) {
			defer wg.Done()

			_, errs[idx] = client.Get(context.Background(), "/printer/info", nil)
		}(i)
	}

	wg.Wait()

	for _, callErr := range errs {
		if callErr != nil {
			t.Fatalf("Get: %v", callErr)
		}
	}

	mu.Lock()
	got := logins
	mu.Unlock()

	if got != 1 {
		t.Errorf("logins = %d, want 1 (concurrent 401s must single-flight one login)", got)
	}
}

func TestSaveToken_LogsWriteFailure(t *testing.T) {
	t.Parallel()

	var buf bytes.Buffer

	logger := slog.New(slog.NewTextHandler(&buf, &slog.HandlerOptions{Level: slog.LevelWarn}))

	srv := httptest.NewServer(http.HandlerFunc(func(writer http.ResponseWriter, request *http.Request) {
		if request.URL.Path == loginPath {
			_, _ = writer.Write([]byte(`{"result":{"token":"t","refresh_token":"r"}}`))

			return
		}

		_, _ = writer.Write([]byte(`{"result":"ok"}`))
	}))
	defer srv.Close()

	// A token path nested under an existing regular file makes the cache write
	// fail, which must be logged rather than silently dropped.
	blocker := filepath.Join(t.TempDir(), "file")

	writeErr := os.WriteFile(blocker, []byte("x"), 0o600)
	if writeErr != nil {
		t.Fatalf("setup: %v", writeErr)
	}

	tokenPath := filepath.Join(blocker, "token.json")
	client := newClient(t, srv, &moonraker.Options{
		Username: testUser, Password: testPass, TokenPath: tokenPath, Logger: logger,
	})

	_, err := client.Get(t.Context(), "/printer/info", nil)
	if err != nil {
		t.Fatalf("Get: %v", err)
	}

	if !strings.Contains(buf.String(), "token cache") {
		t.Errorf("logger output = %q, want a token-cache warning", buf.String())
	}
}

func TestGet_UnauthenticatedWithoutCredentials(t *testing.T) {
	t.Parallel()

	srv := httptest.NewServer(http.HandlerFunc(func(writer http.ResponseWriter, _ *http.Request) {
		writer.WriteHeader(http.StatusUnauthorized)
	}))
	defer srv.Close()

	client := newClient(t, srv, nil)

	_, err := client.Get(context.Background(), "/printer/info", nil)
	if !errors.Is(err, moonraker.ErrNotAuthenticated) {
		t.Fatalf("err = %v, want ErrNotAuthenticated", err)
	}
}
