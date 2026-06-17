package tools_test

import (
	"testing"

	"github.com/lexfrei/mcp-raker/internal/tools"
)

// These tests cover the success path of handlers whose other tests only exercise
// parameter validation, so the request-building code (method, path, query, body)
// is actually verified.

func TestAnalysisProcess_Success(t *testing.T) {
	t.Parallel()

	mock := &mockAPI{result: okJSON}

	_, _, err := tools.NewAnalysisProcessHandler(mock)(t.Context(), nil, tools.FilenameParams{Filename: testGcodeFile})
	if err != nil {
		t.Fatalf("handler: %v", err)
	}

	assertCall(t, mock, methodPost, "/server/analysis/process")
}

func TestAnnouncementsDismiss_Success(t *testing.T) {
	t.Parallel()

	mock := &mockAPI{result: okJSON}

	_, _, err := tools.NewAnnouncementsDismissHandler(mock)(t.Context(), nil, tools.AnnouncementsDismissParams{EntryID: "e1"})
	if err != nil {
		t.Fatalf("handler: %v", err)
	}

	assertCall(t, mock, methodPost, "/server/announcements/dismiss")

	if mock.lastQuery.Get("entry_id") != "e1" {
		t.Errorf("entry_id = %q, want e1", mock.lastQuery.Get("entry_id"))
	}
}

func TestDBGetItem_Success(t *testing.T) {
	t.Parallel()

	mock := &mockAPI{result: okJSON}

	_, _, err := tools.NewDBGetItemHandler(mock)(t.Context(), nil, tools.DBGetItemParams{Namespace: testNS, Key: "k"})
	if err != nil {
		t.Fatalf("handler: %v", err)
	}

	assertCall(t, mock, methodGet, "/server/database/item")

	if mock.lastQuery.Get("namespace") != testNS {
		t.Errorf("namespace = %q, want %s", mock.lastQuery.Get("namespace"), testNS)
	}
}

func TestDBDeleteItem_Success(t *testing.T) {
	t.Parallel()

	mock := &mockAPI{result: okJSON}

	_, _, err := tools.NewDBDeleteItemHandler(mock)(t.Context(), nil, tools.DBDeleteItemParams{Namespace: testNS, Key: "k"})
	if err != nil {
		t.Fatalf("handler: %v", err)
	}

	assertCall(t, mock, methodDelete, "/server/database/item")
}

func TestFilesMetadata_Success(t *testing.T) {
	t.Parallel()

	mock := &mockAPI{result: okJSON}

	_, _, err := tools.NewFilesMetadataHandler(mock)(t.Context(), nil, tools.FilenameParams{Filename: testGcodeFile})
	if err != nil {
		t.Fatalf("handler: %v", err)
	}

	assertCall(t, mock, methodGet, "/server/files/metadata")

	if mock.lastQuery.Get("filename") != testGcodeFile {
		t.Errorf("filename = %q, want %s", mock.lastQuery.Get("filename"), testGcodeFile)
	}
}

func TestFilesMove_Success(t *testing.T) {
	t.Parallel()

	mock := &mockAPI{result: okJSON}

	_, _, err := tools.NewFilesMoveHandler(mock)(t.Context(), nil, tools.SourceDestParams{Source: "a", Dest: "b"})
	if err != nil {
		t.Fatalf("handler: %v", err)
	}

	assertCall(t, mock, methodPost, "/server/files/move")

	if mock.lastQuery.Get("source") != "a" || mock.lastQuery.Get("dest") != "b" {
		t.Errorf("query = %v, want source=a dest=b", mock.lastQuery)
	}
}

func TestFilesZip_Success(t *testing.T) {
	t.Parallel()

	mock := &mockAPI{result: okJSON}

	_, _, err := tools.NewFilesZipHandler(mock)(t.Context(), nil, tools.FilesZipParams{Items: []string{"gcodes/a.gcode"}})
	if err != nil {
		t.Fatalf("handler: %v", err)
	}

	assertCall(t, mock, methodPost, "/server/files/zip")
}

func TestHistoryJob_Success(t *testing.T) {
	t.Parallel()

	mock := &mockAPI{result: okJSON}

	_, _, err := tools.NewHistoryJobHandler(mock)(t.Context(), nil, tools.HistoryJobParams{UID: "5"})
	if err != nil {
		t.Fatalf("handler: %v", err)
	}

	assertCall(t, mock, methodGet, "/server/history/job")

	if mock.lastQuery.Get("uid") != "5" {
		t.Errorf("uid = %q, want 5", mock.lastQuery.Get("uid"))
	}
}

func TestJobQueueJump_Success(t *testing.T) {
	t.Parallel()

	mock := &mockAPI{result: okJSON}

	_, _, err := tools.NewJobQueueJumpHandler(mock)(t.Context(), nil, tools.JobQueueJumpParams{JobID: "j1"})
	if err != nil {
		t.Fatalf("handler: %v", err)
	}

	assertCall(t, mock, methodPost, "/server/job_queue/jump")

	if mock.lastQuery.Get("job_id") != "j1" {
		t.Errorf("job_id = %q, want j1", mock.lastQuery.Get("job_id"))
	}
}

func TestAccessUserPassword_Success(t *testing.T) {
	t.Parallel()

	mock := &mockAPI{result: okJSON}
	params := tools.AccessUserPasswordParams{Password: "old", NewPassword: "new"}

	_, _, err := tools.NewAccessUserPasswordHandler(mock)(t.Context(), nil, params)
	if err != nil {
		t.Fatalf("handler: %v", err)
	}

	assertCall(t, mock, methodPost, "/access/user/password")
}

func TestMQTTSubscribe_Success(t *testing.T) {
	t.Parallel()

	mock := &mockAPI{result: okJSON}

	_, _, err := tools.NewMQTTSubscribeHandler(mock)(t.Context(), nil, tools.MQTTSubscribeParams{Topic: "klipper/state"})
	if err != nil {
		t.Fatalf("handler: %v", err)
	}

	assertCall(t, mock, methodPost, "/server/mqtt/subscribe")
}

func TestExtensionsRequest_Success(t *testing.T) {
	t.Parallel()

	mock := &mockAPI{result: okJSON}
	params := tools.ExtensionsRequestParams{Agent: testAgent, Method: "ping"}

	_, _, err := tools.NewExtensionsRequestHandler(mock)(t.Context(), nil, params)
	if err != nil {
		t.Fatalf("handler: %v", err)
	}

	assertCall(t, mock, methodPost, "/server/extensions/request")
}

func TestSensorsInfo_Success(t *testing.T) {
	t.Parallel()

	mock := &mockAPI{result: okJSON}

	_, _, err := tools.NewSensorsInfoHandler(mock)(t.Context(), nil, tools.SensorParams{Sensor: testSensor})
	if err != nil {
		t.Fatalf("handler: %v", err)
	}

	assertCall(t, mock, methodGet, "/machine/sensors/info")

	if mock.lastQuery.Get("sensor") != testSensor {
		t.Errorf("sensor = %q, want chamber", mock.lastQuery.Get("sensor"))
	}
}

func TestWebcamsGet_Success(t *testing.T) {
	t.Parallel()

	mock := &mockAPI{result: okJSON}

	_, _, err := tools.NewWebcamsGetHandler(mock)(t.Context(), nil, tools.WebcamNameParams{Name: "cam"})
	if err != nil {
		t.Fatalf("handler: %v", err)
	}

	assertCall(t, mock, methodGet, "/server/webcams/item")

	if mock.lastQuery.Get("name") != "cam" {
		t.Errorf("name = %q, want cam", mock.lastQuery.Get("name"))
	}
}
