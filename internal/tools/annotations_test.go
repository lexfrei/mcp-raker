package tools_test

import (
	"testing"

	"github.com/modelcontextprotocol/go-sdk/mcp"

	"github.com/lexfrei/mcp-raker/internal/tools"
)

// isDestructive reports whether a tool's annotations mark it destructive.
func isDestructive(tool *mcp.Tool) bool {
	return tool.Annotations != nil &&
		tool.Annotations.DestructiveHint != nil &&
		*tool.Annotations.DestructiveHint
}

func TestDestructiveAnnotations(t *testing.T) {
	t.Parallel()

	destructive := map[string]*mcp.Tool{
		"print_start":           tools.PrintStartTool(),
		"print_cancel":          tools.PrintCancelTool(),
		"gcode_script":          tools.GcodeScriptTool(),
		"emergency_stop":        tools.EmergencyStopTool(),
		"files_delete":          tools.FilesDeleteTool(),
		"files_delete_dir":      tools.FilesDeleteDirectoryTool(),
		"power_off":             tools.PowerOffTool(),
		"machine_shutdown":      tools.MachineShutdownTool(),
		"machine_reboot":        tools.MachineRebootTool(),
		"service_stop":          tools.ServiceStopTool(),
		"service_restart":       tools.ServiceRestartTool(),
		"sudo_password":         tools.SudoPasswordTool(),
		"update_upgrade":        tools.UpdateUpgradeTool(),
		"update_rollback":       tools.UpdateRollbackTool(),
		"access_delete_user":    tools.AccessDeleteUserTool(),
		"access_create_api_key": tools.AccessCreateAPIKeyTool(),
	}

	for name, tool := range destructive {
		if !isDestructive(tool) {
			t.Errorf("%s should carry a destructive hint", name)
		}
	}
}

func TestNonDestructiveAnnotations(t *testing.T) {
	t.Parallel()

	writes := map[string]*mcp.Tool{
		"print_pause": tools.PrintPauseTool(),
		"power_on":    tools.PowerOnTool(),
	}

	for name, tool := range writes {
		if isDestructive(tool) {
			t.Errorf("%s should not carry a destructive hint", name)
		}

		if tool.Annotations == nil || tool.Annotations.ReadOnlyHint {
			t.Errorf("%s should be a non-read-only write tool", name)
		}
	}
}

func TestReadOnlyAnnotations(t *testing.T) {
	t.Parallel()

	reads := map[string]*mcp.Tool{
		"printer_info":  tools.PrinterInfoTool(),
		"server_info":   tools.ServerInfoTool(),
		"files_list":    tools.FilesListTool(),
		"history_list":  tools.HistoryListTool(),
		"power_devices": tools.PowerDevicesTool(),
	}

	for name, tool := range reads {
		if tool.Annotations == nil || !tool.Annotations.ReadOnlyHint {
			t.Errorf("%s should carry a read-only hint", name)
		}
	}
}
