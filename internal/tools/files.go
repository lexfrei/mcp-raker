package tools

import (
	"context"
	"net/url"
	"strings"
	"unicode/utf8"

	"github.com/modelcontextprotocol/go-sdk/mcp"

	"github.com/lexfrei/mcp-raker/internal/moonraker"
)

// maxDownloadBytes caps how much of a downloaded file is returned inline.
const maxDownloadBytes = 262144

// rootOrDefault returns root, defaulting to the gcodes root when empty.
func rootOrDefault(root string) string {
	if root == "" {
		return rootGcodes
	}

	return root
}

// filePath builds the /server/files/<root>/<filename> path used by the download
// and delete tools, URL-escaping each path segment so characters like '#' and
// '?' in a filename cannot truncate or corrupt the request URL. The filename's
// own '/' separators are preserved so subdirectory paths still work.
func filePath(root, filename string) string {
	parts := strings.Split(filename, "/")

	escaped := make([]string, 0, len(parts)+1)
	escaped = append(escaped, url.PathEscape(rootOrDefault(root)))

	for _, part := range parts {
		escaped = append(escaped, url.PathEscape(part))
	}

	return "/server/files/" + strings.Join(escaped, "/")
}

// trimToValidUTF8 drops a trailing partial UTF-8 rune left by a byte-offset
// truncation, so the returned bytes always end on a rune boundary. A correctly
// encoded U+FFFD is preserved.
func trimToValidUTF8(data []byte) []byte {
	for len(data) > 0 {
		runeValue, size := utf8.DecodeLastRune(data)
		if runeValue == utf8.RuneError && size <= 1 {
			data = data[:len(data)-1]

			continue
		}

		break
	}

	return data
}

// FilesListResult is the output of moonraker_files_list.
type FilesListResult struct {
	Files []moonraker.FileEntry `json:"files"`
}

// FilesListParams defines the parameters for moonraker_files_list.
type FilesListParams struct {
	Root string `json:"root,omitempty" jsonschema:"File-manager root to list, e.g. 'gcodes' (default), 'config', or 'logs'"`
}

// FilesListTool returns the definition for moonraker_files_list.
func FilesListTool() *mcp.Tool {
	return &mcp.Tool{
		Name:        "moonraker_files_list",
		Description: "List every file under a file-manager root (GET /server/files/list).",
		Annotations: readOnly("List Files"),
	}
}

// NewFilesListHandler creates the handler for moonraker_files_list.
func NewFilesListHandler(api moonraker.API) mcp.ToolHandlerFor[FilesListParams, FilesListResult] {
	return func(ctx context.Context, _ *mcp.CallToolRequest, params FilesListParams) (*mcp.CallToolResult, FilesListResult, error) {
		query := url.Values{"root": {rootOrDefault(params.Root)}}

		files, err := decodeTyped[[]moonraker.FileEntry](api.Get(ctx, "/server/files/list", query))

		return nil, FilesListResult{Files: files}, err
	}
}

// FilesDirectoryParams defines the parameters for moonraker_files_directory.
type FilesDirectoryParams struct {
	Path     string `json:"path,omitempty"     jsonschema:"Directory path to list, e.g. 'gcodes' (default) or 'gcodes/subdir'"`
	Extended bool   `json:"extended,omitempty" jsonschema:"When true, include gcode metadata for each file"`
}

// FilesDirectoryTool returns the definition for moonraker_files_directory.
func FilesDirectoryTool() *mcp.Tool {
	return &mcp.Tool{
		Name:        "moonraker_files_directory",
		Description: "List a directory's files and subdirectories with disk usage (GET /server/files/directory).",
		Annotations: readOnly("List Directory"),
	}
}

// NewFilesDirectoryHandler creates the handler for moonraker_files_directory.
func NewFilesDirectoryHandler(api moonraker.API) mcp.ToolHandlerFor[FilesDirectoryParams, map[string]any] {
	return func(ctx context.Context, _ *mcp.CallToolRequest, params FilesDirectoryParams) (*mcp.CallToolResult, map[string]any, error) {
		query := url.Values{paramPath: {rootOrDefault(params.Path)}}
		if params.Extended {
			query.Set("extended", "true")
		}

		out, err := decodeResult(api.Get(ctx, "/server/files/directory", query))

		return nil, out, err
	}
}

// FilesRootsResult is the output of moonraker_files_roots. Moonraker returns a
// bare array, which MCP cannot expose as top-level structured content, so the
// roots are wrapped under a "roots" key.
type FilesRootsResult struct {
	Roots []map[string]any `json:"roots"`
}

// FilesRootsTool returns the definition for moonraker_files_roots.
func FilesRootsTool() *mcp.Tool {
	return &mcp.Tool{
		Name:        "moonraker_files_roots",
		Description: "List the registered file-manager root directories (GET /server/files/roots).",
		Annotations: readOnly("List Roots"),
	}
}

// NewFilesRootsHandler creates the handler for moonraker_files_roots.
func NewFilesRootsHandler(api moonraker.API) mcp.ToolHandlerFor[NoParams, FilesRootsResult] {
	return func(ctx context.Context, _ *mcp.CallToolRequest, _ NoParams) (*mcp.CallToolResult, FilesRootsResult, error) {
		roots, err := decodeTyped[[]map[string]any](api.Get(ctx, "/server/files/roots", nil))

		return nil, FilesRootsResult{Roots: roots}, err
	}
}

// FilenameParams names a gcode file relative to the gcodes root.
type FilenameParams struct {
	Filename string `json:"filename" jsonschema:"Path of the gcode file relative to the gcodes root, e.g. 'benchy.gcode'"`
}

// FilesMetadataTool returns the definition for moonraker_files_metadata.
func FilesMetadataTool() *mcp.Tool {
	return &mcp.Tool{
		Name:        "moonraker_files_metadata",
		Description: "Get parsed gcode metadata: print time, filament usage, slicer, and thumbnails (GET /server/files/metadata).",
		Annotations: readOnly("File Metadata"),
	}
}

// NewFilesMetadataHandler creates the handler for moonraker_files_metadata.
func NewFilesMetadataHandler(api moonraker.API) mcp.ToolHandlerFor[FilenameParams, map[string]any] {
	return func(ctx context.Context, _ *mcp.CallToolRequest, params FilenameParams) (*mcp.CallToolResult, map[string]any, error) {
		valErr := requireString(paramFilename, params.Filename)
		if valErr != nil {
			return nil, map[string]any{}, valErr
		}

		out, err := decodeResult(api.Get(ctx, "/server/files/metadata", url.Values{paramFilename: {params.Filename}}))

		return nil, out, err
	}
}

// FilesMetascanTool returns the definition for moonraker_files_metascan.
func FilesMetascanTool() *mcp.Tool {
	return &mcp.Tool{
		Name:        "moonraker_files_metascan",
		Description: "Force a fresh metadata scan of a gcode file (POST /server/files/metascan).",
		Annotations: write("Rescan Metadata"),
	}
}

// NewFilesMetascanHandler creates the handler for moonraker_files_metascan.
func NewFilesMetascanHandler(api moonraker.API) mcp.ToolHandlerFor[FilenameParams, map[string]any] {
	return func(ctx context.Context, _ *mcp.CallToolRequest, params FilenameParams) (*mcp.CallToolResult, map[string]any, error) {
		valErr := requireString(paramFilename, params.Filename)
		if valErr != nil {
			return nil, map[string]any{}, valErr
		}

		out, err := decodeResult(api.Post(ctx, "/server/files/metascan", url.Values{paramFilename: {params.Filename}}, nil))

		return nil, out, err
	}
}

// FileThumbnailsResult is the output of moonraker_files_thumbnails. Moonraker
// returns a bare array of thumbnails, wrapped here under a "thumbnails" key.
type FileThumbnailsResult struct {
	Thumbnails []map[string]any `json:"thumbnails"`
}

// FilesThumbnailsTool returns the definition for moonraker_files_thumbnails.
func FilesThumbnailsTool() *mcp.Tool {
	return &mcp.Tool{
		Name:        "moonraker_files_thumbnails",
		Description: "List the thumbnails embedded in a gcode file (GET /server/files/thumbnails).",
		Annotations: readOnly("File Thumbnails"),
	}
}

// NewFilesThumbnailsHandler creates the handler for moonraker_files_thumbnails.
func NewFilesThumbnailsHandler(api moonraker.API) mcp.ToolHandlerFor[FilenameParams, FileThumbnailsResult] {
	return func(ctx context.Context, _ *mcp.CallToolRequest, params FilenameParams) (*mcp.CallToolResult, FileThumbnailsResult, error) {
		valErr := requireString(paramFilename, params.Filename)
		if valErr != nil {
			return nil, FileThumbnailsResult{}, valErr
		}

		thumbs, err := decodeTyped[[]map[string]any](api.Get(ctx, "/server/files/thumbnails", url.Values{paramFilename: {params.Filename}}))

		return nil, FileThumbnailsResult{Thumbnails: thumbs}, err
	}
}

// PathParams names a directory path within a file-manager root.
type PathParams struct {
	Path  string `json:"path"            jsonschema:"Directory path including its root, e.g. 'gcodes/new_folder'"`
	Force bool   `json:"force,omitempty" jsonschema:"When true, delete the directory even if it is not empty"`
}

// FilesCreateDirectoryTool returns the definition for moonraker_files_create_directory.
func FilesCreateDirectoryTool() *mcp.Tool {
	return &mcp.Tool{
		Name:        "moonraker_files_create_directory",
		Description: "Create a directory within a file-manager root (POST /server/files/directory).",
		Annotations: write("Create Directory"),
	}
}

// NewFilesCreateDirectoryHandler creates the handler for moonraker_files_create_directory.
func NewFilesCreateDirectoryHandler(api moonraker.API) mcp.ToolHandlerFor[PathParams, map[string]any] {
	return func(ctx context.Context, _ *mcp.CallToolRequest, params PathParams) (*mcp.CallToolResult, map[string]any, error) {
		valErr := requireString(paramPath, params.Path)
		if valErr != nil {
			return nil, map[string]any{}, valErr
		}

		out, err := decodeResult(api.Post(ctx, "/server/files/directory", url.Values{paramPath: {params.Path}}, nil))

		return nil, out, err
	}
}

// FilesDeleteDirectoryTool returns the definition for moonraker_files_delete_directory.
func FilesDeleteDirectoryTool() *mcp.Tool {
	return &mcp.Tool{
		Name:        "moonraker_files_delete_directory",
		Description: "Delete a directory within a file-manager root (DELETE /server/files/directory).",
		Annotations: writeDestructive("Delete Directory"),
	}
}

// NewFilesDeleteDirectoryHandler creates the handler for moonraker_files_delete_directory.
func NewFilesDeleteDirectoryHandler(api moonraker.API) mcp.ToolHandlerFor[PathParams, map[string]any] {
	return func(ctx context.Context, _ *mcp.CallToolRequest, params PathParams) (*mcp.CallToolResult, map[string]any, error) {
		valErr := requireString(paramPath, params.Path)
		if valErr != nil {
			return nil, map[string]any{}, valErr
		}

		query := url.Values{paramPath: {params.Path}}
		if params.Force {
			query.Set("force", "true")
		}

		out, err := decodeResult(api.Delete(ctx, "/server/files/directory", query))

		return nil, out, err
	}
}

// SourceDestParams defines the parameters for the move and copy tools.
type SourceDestParams struct {
	Source string `json:"source" jsonschema:"Source path including its root, e.g. 'gcodes/a.gcode'"`
	Dest   string `json:"dest"   jsonschema:"Destination path including its root, e.g. 'gcodes/sub/a.gcode'"`
}

// validateSourceDest checks that both endpoints of a move/copy are present.
func validateSourceDest(params SourceDestParams) error {
	srcErr := requireString(paramSource, params.Source)
	if srcErr != nil {
		return srcErr
	}

	return requireString(paramDest, params.Dest)
}

// FilesMoveTool returns the definition for moonraker_files_move.
func FilesMoveTool() *mcp.Tool {
	return &mcp.Tool{
		Name:        "moonraker_files_move",
		Description: "Move or rename a file or directory (POST /server/files/move).",
		Annotations: write("Move File"),
	}
}

// NewFilesMoveHandler creates the handler for moonraker_files_move.
func NewFilesMoveHandler(api moonraker.API) mcp.ToolHandlerFor[SourceDestParams, map[string]any] {
	return func(ctx context.Context, _ *mcp.CallToolRequest, params SourceDestParams) (*mcp.CallToolResult, map[string]any, error) {
		valErr := validateSourceDest(params)
		if valErr != nil {
			return nil, map[string]any{}, valErr
		}

		query := url.Values{paramSource: {params.Source}, paramDest: {params.Dest}}

		out, err := decodeResult(api.Post(ctx, "/server/files/move", query, nil))

		return nil, out, err
	}
}

// FilesCopyTool returns the definition for moonraker_files_copy.
func FilesCopyTool() *mcp.Tool {
	return &mcp.Tool{
		Name:        "moonraker_files_copy",
		Description: "Copy a file or directory (POST /server/files/copy).",
		Annotations: write("Copy File"),
	}
}

// NewFilesCopyHandler creates the handler for moonraker_files_copy.
func NewFilesCopyHandler(api moonraker.API) mcp.ToolHandlerFor[SourceDestParams, map[string]any] {
	return func(ctx context.Context, _ *mcp.CallToolRequest, params SourceDestParams) (*mcp.CallToolResult, map[string]any, error) {
		valErr := validateSourceDest(params)
		if valErr != nil {
			return nil, map[string]any{}, valErr
		}

		query := url.Values{paramSource: {params.Source}, paramDest: {params.Dest}}

		out, err := decodeResult(api.Post(ctx, "/server/files/copy", query, nil))

		return nil, out, err
	}
}

// FilesZipParams defines the parameters for moonraker_files_zip.
type FilesZipParams struct {
	Items     []string `json:"items"                jsonschema:"Files or directories (including their root) to add to the archive"`
	Dest      string   `json:"dest,omitempty"       jsonschema:"Destination path for the zip file; omit to use a default location"`
	StoreOnly bool     `json:"store_only,omitempty" jsonschema:"When true, store without compression"`
}

// FilesZipTool returns the definition for moonraker_files_zip.
func FilesZipTool() *mcp.Tool {
	return &mcp.Tool{
		Name:        "moonraker_files_zip",
		Description: "Create a zip archive from files and directories (POST /server/files/zip).",
		Annotations: write("Create Zip"),
	}
}

// NewFilesZipHandler creates the handler for moonraker_files_zip.
func NewFilesZipHandler(api moonraker.API) mcp.ToolHandlerFor[FilesZipParams, map[string]any] {
	return func(ctx context.Context, _ *mcp.CallToolRequest, params FilesZipParams) (*mcp.CallToolResult, map[string]any, error) {
		valErr := requirePresent("items", len(params.Items))
		if valErr != nil {
			return nil, map[string]any{}, valErr
		}

		body := map[string]any{"items": params.Items, "store_only": params.StoreOnly}
		if params.Dest != "" {
			body[paramDest] = params.Dest
		}

		out, err := decodeResult(api.Post(ctx, "/server/files/zip", nil, body))

		return nil, out, err
	}
}

// FileDownloadParams defines the parameters for moonraker_files_download.
type FileDownloadParams struct {
	Root     string `json:"root,omitempty" jsonschema:"File-manager root, e.g. 'gcodes' (default), 'config', or 'logs'"`
	Filename string `json:"filename"       jsonschema:"Path of the file within the root, e.g. 'printer.cfg'"`
}

// FileDownload is the output of moonraker_files_download. Content has no
// omitempty: an empty file should still report content as "" rather than dropping
// the field, keeping the output shape uniform.
type FileDownload struct {
	Filename  string `json:"filename"`
	Size      int    `json:"size"`
	Truncated bool   `json:"truncated"`
	Content   string `json:"content"`
}

// FilesDownloadTool returns the definition for moonraker_files_download.
func FilesDownloadTool() *mcp.Tool {
	return &mcp.Tool{
		Name: "moonraker_files_download",
		Description: "Download a file's contents as text (GET /server/files/{root}/{filename}). " +
			"Large files are truncated.",
		Annotations: readOnly("Download File"),
	}
}

// NewFilesDownloadHandler creates the handler for moonraker_files_download.
func NewFilesDownloadHandler(api moonraker.API) mcp.ToolHandlerFor[FileDownloadParams, FileDownload] {
	return func(ctx context.Context, _ *mcp.CallToolRequest, params FileDownloadParams) (*mcp.CallToolResult, FileDownload, error) {
		valErr := requireString(paramFilename, params.Filename)
		if valErr != nil {
			return nil, FileDownload{}, valErr
		}

		data, err := api.GetRaw(ctx, filePath(params.Root, params.Filename), nil)
		if err != nil {
			return nil, FileDownload{}, moonrakerErr("download failed", err)
		}

		content := data
		truncated := false

		if len(content) > maxDownloadBytes {
			content = trimToValidUTF8(content[:maxDownloadBytes])
			truncated = true
		}

		result := FileDownload{
			Filename:  params.Filename,
			Size:      len(data),
			Truncated: truncated,
			Content:   string(content),
		}

		return nil, result, nil
	}
}

// FilesUploadParams defines the parameters for moonraker_files_upload.
type FilesUploadParams struct {
	Root       string `json:"root,omitempty"        jsonschema:"File-manager root to upload into, e.g. 'gcodes' (default)"`
	Path       string `json:"path,omitempty"        jsonschema:"Optional subdirectory within the root"`
	Filename   string `json:"filename"              jsonschema:"Name to store the uploaded file under, e.g. 'part.gcode'"`
	Content    string `json:"content,omitempty"     jsonschema:"The full text content of the file to upload"`
	StartPrint bool   `json:"start_print,omitempty" jsonschema:"When true, start printing the file immediately after upload"`
}

// FilesUploadTool returns the definition for moonraker_files_upload.
func FilesUploadTool() *mcp.Tool {
	return &mcp.Tool{
		Name:        "moonraker_files_upload",
		Description: "Upload a file to a file-manager root, optionally starting a print (POST /server/files/upload).",
		Annotations: write("Upload File"),
	}
}

// NewFilesUploadHandler creates the handler for moonraker_files_upload.
func NewFilesUploadHandler(api moonraker.API) mcp.ToolHandlerFor[FilesUploadParams, map[string]any] {
	return func(ctx context.Context, _ *mcp.CallToolRequest, params FilesUploadParams) (*mcp.CallToolResult, map[string]any, error) {
		valErr := requireString(paramFilename, params.Filename)
		if valErr != nil {
			return nil, map[string]any{}, valErr
		}

		out, err := decodeResult(api.Upload(ctx, &moonraker.UploadOptions{
			Root:       rootOrDefault(params.Root),
			Path:       params.Path,
			Filename:   params.Filename,
			Content:    []byte(params.Content),
			StartPrint: params.StartPrint,
		}))

		return nil, out, err
	}
}

// FilesDeleteTool returns the definition for moonraker_files_delete.
func FilesDeleteTool() *mcp.Tool {
	return &mcp.Tool{
		Name:        "moonraker_files_delete",
		Description: "Delete a file from a file-manager root (DELETE /server/files/{root}/{filename}).",
		Annotations: writeDestructive("Delete File"),
	}
}

// NewFilesDeleteHandler creates the handler for moonraker_files_delete.
func NewFilesDeleteHandler(api moonraker.API) mcp.ToolHandlerFor[FileDownloadParams, map[string]any] {
	return func(ctx context.Context, _ *mcp.CallToolRequest, params FileDownloadParams) (*mcp.CallToolResult, map[string]any, error) {
		valErr := requireString(paramFilename, params.Filename)
		if valErr != nil {
			return nil, map[string]any{}, valErr
		}

		out, err := decodeResult(api.Delete(ctx, filePath(params.Root, params.Filename), nil))

		return nil, out, err
	}
}
