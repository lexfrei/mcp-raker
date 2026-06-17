package tools_test

import (
	"encoding/json"
	"strings"
	"testing"
	"unicode/utf8"

	"github.com/cockroachdb/errors"

	"github.com/lexfrei/mcp-raker/internal/tools"
)

const testCfgFile = "printer.cfg"

func TestFilesList_DefaultRootAndDecode(t *testing.T) {
	t.Parallel()

	mock := &mockAPI{result: json.RawMessage(`[{"path":"a.gcode","size":10,"modified":1.0,"permissions":"rw"}]`)}

	_, out, err := tools.NewFilesListHandler(mock)(t.Context(), nil, tools.FilesListParams{})
	if err != nil {
		t.Fatalf("handler: %v", err)
	}

	assertCall(t, mock, methodGet, "/server/files/list")

	if mock.lastQuery.Get("root") != "gcodes" {
		t.Errorf("root = %q, want gcodes", mock.lastQuery.Get("root"))
	}

	if len(out.Files) != 1 || out.Files[0].Path != testGcodeFile {
		t.Errorf("files = %+v, want one entry a.gcode", out.Files)
	}
}

func TestFilesMetadata_RequiresFilename(t *testing.T) {
	t.Parallel()

	_, _, err := tools.NewFilesMetadataHandler(&mockAPI{})(t.Context(), nil, tools.FilenameParams{})
	if !errors.Is(err, tools.ErrValidation) {
		t.Fatalf("err = %v, want ErrValidation", err)
	}
}

func TestFilesMove_RequiresBoth(t *testing.T) {
	t.Parallel()

	_, _, noSrc := tools.NewFilesMoveHandler(&mockAPI{})(t.Context(), nil, tools.SourceDestParams{Dest: "b"})
	if !errors.Is(noSrc, tools.ErrValidation) {
		t.Errorf("missing source err = %v, want ErrValidation", noSrc)
	}

	_, _, noDest := tools.NewFilesMoveHandler(&mockAPI{})(t.Context(), nil, tools.SourceDestParams{Source: "a"})
	if !errors.Is(noDest, tools.ErrValidation) {
		t.Errorf("missing dest err = %v, want ErrValidation", noDest)
	}
}

func TestFilesZip_RequiresItems(t *testing.T) {
	t.Parallel()

	_, _, err := tools.NewFilesZipHandler(&mockAPI{})(t.Context(), nil, tools.FilesZipParams{})
	if !errors.Is(err, tools.ErrValidation) {
		t.Fatalf("err = %v, want ErrValidation", err)
	}
}

func TestFilesDownload_Truncates(t *testing.T) {
	t.Parallel()

	big := strings.Repeat("x", 300000)
	mock := &mockAPI{rawResult: []byte(big)}
	params := tools.FileDownloadParams{Filename: testCfgFile}

	_, out, err := tools.NewFilesDownloadHandler(mock)(t.Context(), nil, params)
	if err != nil {
		t.Fatalf("handler: %v", err)
	}

	if mock.lastMethod != methodGetRaw || mock.lastPath != "/server/files/gcodes/printer.cfg" {
		t.Errorf("call = %s %s, want GETRAW /server/files/gcodes/printer.cfg", mock.lastMethod, mock.lastPath)
	}

	if !out.Truncated || out.Size != len(big) || len(out.Content) != 262144 {
		t.Errorf("download = {trunc:%v size:%d len:%d}, want trunc size=%d len=262144", out.Truncated, out.Size, len(out.Content), len(big))
	}
}

func TestFilesDownload_RequiresFilename(t *testing.T) {
	t.Parallel()

	_, _, err := tools.NewFilesDownloadHandler(&mockAPI{})(t.Context(), nil, tools.FileDownloadParams{})
	if !errors.Is(err, tools.ErrValidation) {
		t.Fatalf("err = %v, want ErrValidation", err)
	}
}

func TestFilesUpload(t *testing.T) {
	t.Parallel()

	mock := &mockAPI{result: okJSON}
	params := tools.FilesUploadParams{Filename: "part.gcode", Content: gcodeG28, StartPrint: true}

	_, _, err := tools.NewFilesUploadHandler(mock)(t.Context(), nil, params)
	if err != nil {
		t.Fatalf("handler: %v", err)
	}

	if mock.lastMethod != methodUpload || mock.lastUpload == nil {
		t.Fatalf("upload not invoked: method=%s", mock.lastMethod)
	}

	if mock.lastUpload.Root != "gcodes" || mock.lastUpload.Filename != "part.gcode" || !mock.lastUpload.StartPrint {
		t.Errorf("upload opts = %+v, want gcodes/part.gcode start_print", mock.lastUpload)
	}
}

func TestFilesDelete(t *testing.T) {
	t.Parallel()

	mock := &mockAPI{result: okJSON}
	params := tools.FileDownloadParams{Root: "config", Filename: testCfgFile}

	_, _, err := tools.NewFilesDeleteHandler(mock)(t.Context(), nil, params)
	if err != nil {
		t.Fatalf("handler: %v", err)
	}

	assertCall(t, mock, methodDelete, "/server/files/config/printer.cfg")
}

func TestFilesDelete_EscapesFilename(t *testing.T) {
	t.Parallel()

	mock := &mockAPI{result: okJSON}
	params := tools.FileDownloadParams{Filename: "a#b?c.gcode"}

	_, _, err := tools.NewFilesDeleteHandler(mock)(t.Context(), nil, params)
	if err != nil {
		t.Fatalf("handler: %v", err)
	}

	if mock.lastPath != "/server/files/gcodes/a%23b%3Fc.gcode" {
		t.Errorf("path = %q, want escaped /server/files/gcodes/a%%23b%%3Fc.gcode", mock.lastPath)
	}
}

func TestFilesDownload_EscapesFilename(t *testing.T) {
	t.Parallel()

	mock := &mockAPI{rawResult: []byte("data")}
	params := tools.FileDownloadParams{Filename: "sub/a#b.gcode"}

	_, _, err := tools.NewFilesDownloadHandler(mock)(t.Context(), nil, params)
	if err != nil {
		t.Fatalf("handler: %v", err)
	}

	if mock.lastPath != "/server/files/gcodes/sub/a%23b.gcode" {
		t.Errorf("path = %q, want escaped subdir path with %%23", mock.lastPath)
	}
}

func TestFilesDownload_ValidUTF8AtBoundary(t *testing.T) {
	t.Parallel()

	// A two-byte rune straddles the 262144-byte cap: byte 262144 is the first
	// byte of 'é', so a naive cut would leave a partial rune.
	data := strings.Repeat("a", 262143) + "é" + "tail"
	mock := &mockAPI{rawResult: []byte(data)}

	_, out, err := tools.NewFilesDownloadHandler(mock)(t.Context(), nil, tools.FileDownloadParams{Filename: testCfgFile})
	if err != nil {
		t.Fatalf("handler: %v", err)
	}

	if !out.Truncated {
		t.Error("expected truncated download")
	}

	if !utf8.ValidString(out.Content) {
		t.Error("truncated content is not valid UTF-8")
	}

	if len(out.Content) != 262143 {
		t.Errorf("content length = %d, want 262143 (partial rune trimmed)", len(out.Content))
	}
}

func TestFilesDownload_WrapsError(t *testing.T) {
	t.Parallel()

	mock := &mockAPI{rawErr: errStub}

	_, _, err := tools.NewFilesDownloadHandler(mock)(t.Context(), nil, tools.FileDownloadParams{Filename: "x.gcode"})
	if !errors.Is(err, tools.ErrMoonraker) {
		t.Fatalf("err = %v, want ErrMoonraker", err)
	}
}

func TestFilesRoots_WrapsError(t *testing.T) {
	t.Parallel()

	mock := &mockAPI{err: errStub}

	_, _, err := tools.NewFilesRootsHandler(mock)(t.Context(), nil, tools.NoParams{})
	if !errors.Is(err, tools.ErrMoonraker) {
		t.Fatalf("err = %v, want ErrMoonraker", err)
	}
}
