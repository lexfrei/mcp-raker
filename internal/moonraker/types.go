package moonraker

// ServerInfo is the result of GET /server/info.
type ServerInfo struct {
	KlippyConnected  bool     `json:"klippy_connected"`
	KlippyState      string   `json:"klippy_state"`
	Components       []string `json:"components"`
	FailedComponents []string `json:"failed_components"`
	Warnings         []string `json:"warnings"`
	WebsocketCount   int      `json:"websocket_count"`
	MoonrakerVersion string   `json:"moonraker_version"`
	APIVersion       []int    `json:"api_version"`
	APIVersionString string   `json:"api_version_string"`
}

// PrinterInfo is the result of GET /printer/info.
type PrinterInfo struct {
	State           string `json:"state"`
	StateMessage    string `json:"state_message"`
	Hostname        string `json:"hostname"`
	SoftwareVersion string `json:"software_version"`
	CPUInfo         string `json:"cpu_info"`
	KlipperPath     string `json:"klipper_path"`
	PythonPath      string `json:"python_path"`
	LogFile         string `json:"log_file"`
	ConfigFile      string `json:"config_file"`
}

// ObjectsList is the result of GET /printer/objects/list: the names of every
// printer object that can be queried.
type ObjectsList struct {
	Objects []string `json:"objects"`
}

// ObjectsQuery is the result of POST /printer/objects/query. Status maps each
// requested object name to its current field values; the shape varies by
// object, so it is kept dynamic.
type ObjectsQuery struct {
	Eventtime float64        `json:"eventtime"`
	Status    map[string]any `json:"status"`
}

// FileEntry is one entry in a flat file listing (GET /server/files/list).
type FileEntry struct {
	Path        string  `json:"path"`
	Modified    float64 `json:"modified"`
	Size        int64   `json:"size"`
	Permissions string  `json:"permissions"`
}
