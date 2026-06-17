package tools

// Shared query and body parameter keys, kept as constants so every tool builds
// requests in one uniform style and the goconst linter stays satisfied.
const (
	paramName      = "name"
	paramFilename  = "filename"
	paramDevice    = "device"
	paramService   = "service"
	paramStrip     = "strip"
	paramNamespace = "namespace"
	paramKey       = "key"
	paramSource    = "source"
	paramDest      = "dest"
	paramPath      = "path"
	paramSensor    = "sensor"
	paramUID       = "uid"
	paramUsername  = "username"
	paramPassword  = "password"
	paramScript    = "script"
	paramJobID     = "job_id"
	paramValue     = "value"
	paramEntryID   = "entry_id"
	paramAction    = "action"

	// rootGcodes is the default file-manager root.
	rootGcodes = "gcodes"
)
