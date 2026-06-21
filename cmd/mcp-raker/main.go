// Command mcp-raker is an MCP server exposing the Moonraker 3D-printer API
// (the web server that fronts the Klipper firmware) over stdio and, optionally,
// HTTP.
package main

import (
	"context"
	"crypto/subtle"
	"log/slog"
	"net/http"
	"os"
	"os/signal"
	"runtime"
	"syscall"
	"time"

	"github.com/cockroachdb/errors"
	"github.com/modelcontextprotocol/go-sdk/mcp"
	"golang.org/x/sync/errgroup"

	"github.com/lexfrei/mcp-raker/internal/config"
	"github.com/lexfrei/mcp-raker/internal/moonraker"
	"github.com/lexfrei/mcp-raker/internal/tools"
)

const (
	serverName        = "mcp-raker"
	readHeaderTimeout = 10 * time.Second
	shutdownTimeout   = 5 * time.Second
)

// version and revision are set via ldflags at build time.
var (
	version  = "dev"
	revision = "unknown"
)

func main() {
	logger := newLogger()

	err := run(logger)
	if err != nil {
		logger.Error("server failed", slog.Any("error", err))
		os.Exit(1)
	}
}

// newLogger builds the structured JSON logger. Logs go to stderr because stdout
// carries the JSON-RPC stream.
func newLogger() *slog.Logger {
	return slog.New(slog.NewJSONHandler(os.Stderr, &slog.HandlerOptions{
		Level: slog.LevelInfo,
	}))
}

func run(logger *slog.Logger) error {
	cfg, cfgErr := config.Load()
	if cfgErr != nil {
		return errors.Wrap(cfgErr, "invalid configuration")
	}

	httpErr := cfg.ValidateHTTP()
	if httpErr != nil {
		return errors.Wrap(httpErr, "invalid HTTP configuration")
	}

	transport, transportErr := cfg.ProxyTransport()
	if transportErr != nil {
		return errors.Wrap(transportErr, "invalid proxy configuration")
	}

	client, clientErr := moonraker.New(&moonraker.Options{
		BaseURL:   cfg.URL,
		APIKey:    cfg.APIKey,
		Token:     cfg.Token,
		Username:  cfg.Username,
		Password:  cfg.Password,
		TokenPath: cfg.TokenFile,
		UserAgent: cfg.UserAgent,
		Timeout:   cfg.Timeout,
		Transport: transport,
		Logger:    logger,
	})
	if clientErr != nil {
		return errors.Wrap(clientErr, "failed to create moonraker client")
	}

	server := mcp.NewServer(
		&mcp.Implementation{
			Name:    serverName,
			Version: version + "+" + revision,
		},
		newServerOptions(logger, cfg.EnableAdmin),
	)

	registerTools(server, client, cfg.EnableAdmin)

	logger.Info("starting server",
		slog.String("url", cfg.URL),
		slog.Bool("admin", cfg.EnableAdmin),
		slog.Bool("authenticated", cfg.HasAuth()))

	return serve(logger, server, cfg)
}

// newServerOptions wires the shared logger into the MCP server and describes the
// available tools and configuration surfaced to clients.
func newServerOptions(logger *slog.Logger, enableAdmin bool) *mcp.ServerOptions {
	instructions := "MCP server for Moonraker, the API server for Klipper 3D printers. Monitor and " +
		"control the printer: query status and temperatures, start/pause/resume/cancel prints, run " +
		"G-code, browse and manage gcode files, inspect print history and the job queue, control power " +
		"devices, and read sensors. Set MOONRAKER_URL to the printer address (default " +
		"http://localhost:7125) and, if the printer requires it, MOONRAKER_API_KEY or " +
		"MOONRAKER_USERNAME/MOONRAKER_PASSWORD."

	if enableAdmin {
		instructions += " Admin tools are enabled: OS shutdown/reboot, service control, the update " +
			"manager, and user management are available."
	} else {
		instructions += " Admin tools (OS, services, updates, user management) are disabled; set " +
			"MOONRAKER_ENABLE_ADMIN=true to enable them."
	}

	return &mcp.ServerOptions{
		Instructions: instructions,
		Logger:       logger,
	}
}

// registerTools registers the always-available tools, plus the destructive
// admin tools when MOONRAKER_ENABLE_ADMIN is set.
func registerTools(server *mcp.Server, api moonraker.API, enableAdmin bool) {
	registerCoreTools(server, api)

	if enableAdmin {
		registerAdminTools(server, api)
	}
}

// registerCoreTools registers the monitoring and everyday-control tools.
func registerCoreTools(server *mcp.Server, api moonraker.API) {
	mcp.AddTool(server, tools.MCPVersionTool(), tools.NewMCPVersionHandler(version, revision, runtime.Version()))

	// Server and printer status.
	mcp.AddTool(server, tools.ServerInfoTool(), tools.NewServerInfoHandler(api))
	mcp.AddTool(server, tools.ServerConfigTool(), tools.NewServerConfigHandler(api))
	mcp.AddTool(server, tools.TemperatureStoreTool(), tools.NewTemperatureStoreHandler(api))
	mcp.AddTool(server, tools.GcodeStoreTool(), tools.NewGcodeStoreHandler(api))
	mcp.AddTool(server, tools.PrinterInfoTool(), tools.NewPrinterInfoHandler(api))
	mcp.AddTool(server, tools.ObjectsListTool(), tools.NewObjectsListHandler(api))
	mcp.AddTool(server, tools.ObjectsQueryTool(), tools.NewObjectsQueryHandler(api))
	mcp.AddTool(server, tools.QueryEndstopsTool(), tools.NewQueryEndstopsHandler(api))
	mcp.AddTool(server, tools.GcodeScriptTool(), tools.NewGcodeScriptHandler(api))
	mcp.AddTool(server, tools.GcodeHelpTool(), tools.NewGcodeHelpHandler(api))
	mcp.AddTool(server, tools.EmergencyStopTool(), tools.NewEmergencyStopHandler(api))

	// Klipper restart is a printer endpoint, not an OS/service admin action, so
	// it stays outside the MOONRAKER_ENABLE_ADMIN gate alongside emergency stop.
	mcp.AddTool(server, tools.PrinterRestartTool(), tools.NewPrinterRestartHandler(api))
	mcp.AddTool(server, tools.FirmwareRestartTool(), tools.NewFirmwareRestartHandler(api))

	// Print job control.
	mcp.AddTool(server, tools.PrintStartTool(), tools.NewPrintStartHandler(api))
	mcp.AddTool(server, tools.PrintPauseTool(), tools.NewPrintPauseHandler(api))
	mcp.AddTool(server, tools.PrintResumeTool(), tools.NewPrintResumeHandler(api))
	mcp.AddTool(server, tools.PrintCancelTool(), tools.NewPrintCancelHandler(api))

	// Machine status.
	mcp.AddTool(server, tools.SystemInfoTool(), tools.NewSystemInfoHandler(api))
	mcp.AddTool(server, tools.ProcStatsTool(), tools.NewProcStatsHandler(api))
	mcp.AddTool(server, tools.SudoInfoTool(), tools.NewSudoInfoHandler(api))
	mcp.AddTool(server, tools.PeripheralsUSBTool(), tools.NewPeripheralsUSBHandler(api))
	mcp.AddTool(server, tools.PeripheralsSerialTool(), tools.NewPeripheralsSerialHandler(api))
	mcp.AddTool(server, tools.PeripheralsVideoTool(), tools.NewPeripheralsVideoHandler(api))
	mcp.AddTool(server, tools.PeripheralsCanbusTool(), tools.NewPeripheralsCanbusHandler(api))
	mcp.AddTool(server, tools.UpdateStatusTool(), tools.NewUpdateStatusHandler(api))

	// Power devices.
	mcp.AddTool(server, tools.PowerDevicesTool(), tools.NewPowerDevicesHandler(api))
	mcp.AddTool(server, tools.PowerStatusTool(), tools.NewPowerStatusHandler(api))
	mcp.AddTool(server, tools.PowerOnTool(), tools.NewPowerOnHandler(api))
	mcp.AddTool(server, tools.PowerOffTool(), tools.NewPowerOffHandler(api))
	mcp.AddTool(server, tools.PowerToggleTool(), tools.NewPowerToggleHandler(api))

	// File manager.
	mcp.AddTool(server, tools.FilesListTool(), tools.NewFilesListHandler(api))
	mcp.AddTool(server, tools.FilesDirectoryTool(), tools.NewFilesDirectoryHandler(api))
	mcp.AddTool(server, tools.FilesRootsTool(), tools.NewFilesRootsHandler(api))
	mcp.AddTool(server, tools.FilesMetadataTool(), tools.NewFilesMetadataHandler(api))
	mcp.AddTool(server, tools.FilesMetascanTool(), tools.NewFilesMetascanHandler(api))
	mcp.AddTool(server, tools.FilesThumbnailsTool(), tools.NewFilesThumbnailsHandler(api))
	mcp.AddTool(server, tools.FilesCreateDirectoryTool(), tools.NewFilesCreateDirectoryHandler(api))
	mcp.AddTool(server, tools.FilesDeleteDirectoryTool(), tools.NewFilesDeleteDirectoryHandler(api))
	mcp.AddTool(server, tools.FilesMoveTool(), tools.NewFilesMoveHandler(api))
	mcp.AddTool(server, tools.FilesCopyTool(), tools.NewFilesCopyHandler(api))
	mcp.AddTool(server, tools.FilesZipTool(), tools.NewFilesZipHandler(api))
	mcp.AddTool(server, tools.FilesDownloadTool(), tools.NewFilesDownloadHandler(api))
	mcp.AddTool(server, tools.FilesUploadTool(), tools.NewFilesUploadHandler(api))
	mcp.AddTool(server, tools.FilesDeleteTool(), tools.NewFilesDeleteHandler(api))

	// History and job queue.
	mcp.AddTool(server, tools.HistoryListTool(), tools.NewHistoryListHandler(api))
	mcp.AddTool(server, tools.HistoryTotalsTool(), tools.NewHistoryTotalsHandler(api))
	mcp.AddTool(server, tools.HistoryJobTool(), tools.NewHistoryJobHandler(api))
	mcp.AddTool(server, tools.HistoryResetTotalsTool(), tools.NewHistoryResetTotalsHandler(api))
	mcp.AddTool(server, tools.HistoryDeleteJobTool(), tools.NewHistoryDeleteJobHandler(api))
	mcp.AddTool(server, tools.JobQueueStatusTool(), tools.NewJobQueueStatusHandler(api))
	mcp.AddTool(server, tools.JobQueueEnqueueTool(), tools.NewJobQueueEnqueueHandler(api))
	mcp.AddTool(server, tools.JobQueueRemoveTool(), tools.NewJobQueueRemoveHandler(api))
	mcp.AddTool(server, tools.JobQueuePauseTool(), tools.NewJobQueuePauseHandler(api))
	mcp.AddTool(server, tools.JobQueueStartTool(), tools.NewJobQueueStartHandler(api))
	mcp.AddTool(server, tools.JobQueueJumpTool(), tools.NewJobQueueJumpHandler(api))

	// Database.
	mcp.AddTool(server, tools.DBListTool(), tools.NewDBListHandler(api))
	mcp.AddTool(server, tools.DBGetItemTool(), tools.NewDBGetItemHandler(api))
	mcp.AddTool(server, tools.DBPostItemTool(), tools.NewDBPostItemHandler(api))
	mcp.AddTool(server, tools.DBDeleteItemTool(), tools.NewDBDeleteItemHandler(api))
	mcp.AddTool(server, tools.DBBackupTool(), tools.NewDBBackupHandler(api))
	mcp.AddTool(server, tools.DBCompactTool(), tools.NewDBCompactHandler(api))

	// Access (read-only).
	mcp.AddTool(server, tools.AccessUserInfoTool(), tools.NewAccessUserInfoHandler(api))
	mcp.AddTool(server, tools.AccessUsersListTool(), tools.NewAccessUsersListHandler(api))
	mcp.AddTool(server, tools.AccessInfoTool(), tools.NewAccessInfoHandler(api))
	mcp.AddTool(server, tools.AccessAPIKeyTool(), tools.NewAccessAPIKeyHandler(api))

	// Announcements.
	mcp.AddTool(server, tools.AnnouncementsListTool(), tools.NewAnnouncementsListHandler(api))
	mcp.AddTool(server, tools.AnnouncementsUpdateTool(), tools.NewAnnouncementsUpdateHandler(api))
	mcp.AddTool(server, tools.AnnouncementsDismissTool(), tools.NewAnnouncementsDismissHandler(api))
	mcp.AddTool(server, tools.AnnouncementsFeedsTool(), tools.NewAnnouncementsFeedsHandler(api))
	mcp.AddTool(server, tools.AnnouncementsAddFeedTool(), tools.NewAnnouncementsAddFeedHandler(api))
	mcp.AddTool(server, tools.AnnouncementsRemoveFeedTool(), tools.NewAnnouncementsRemoveFeedHandler(api))

	// Webcams.
	mcp.AddTool(server, tools.WebcamsListTool(), tools.NewWebcamsListHandler(api))
	mcp.AddTool(server, tools.WebcamsGetTool(), tools.NewWebcamsGetHandler(api))
	mcp.AddTool(server, tools.WebcamsAddTool(), tools.NewWebcamsAddHandler(api))
	mcp.AddTool(server, tools.WebcamsDeleteTool(), tools.NewWebcamsDeleteHandler(api))
	mcp.AddTool(server, tools.WebcamsTestTool(), tools.NewWebcamsTestHandler(api))

	// Sensors.
	mcp.AddTool(server, tools.SensorsListTool(), tools.NewSensorsListHandler(api))
	mcp.AddTool(server, tools.SensorsInfoTool(), tools.NewSensorsInfoHandler(api))
	mcp.AddTool(server, tools.SensorsMeasurementsTool(), tools.NewSensorsMeasurementsHandler(api))

	// WLED.
	mcp.AddTool(server, tools.WLEDStripsTool(), tools.NewWLEDStripsHandler(api))
	mcp.AddTool(server, tools.WLEDStatusTool(), tools.NewWLEDStatusHandler(api))
	mcp.AddTool(server, tools.WLEDOnTool(), tools.NewWLEDOnHandler(api))
	mcp.AddTool(server, tools.WLEDOffTool(), tools.NewWLEDOffHandler(api))
	mcp.AddTool(server, tools.WLEDToggleTool(), tools.NewWLEDToggleHandler(api))
	mcp.AddTool(server, tools.WLEDSetTool(), tools.NewWLEDSetHandler(api))

	// Spoolman, analysis, MQTT, extensions, notifiers.
	mcp.AddTool(server, tools.SpoolmanStatusTool(), tools.NewSpoolmanStatusHandler(api))
	mcp.AddTool(server, tools.SpoolmanGetSpoolTool(), tools.NewSpoolmanGetSpoolHandler(api))
	mcp.AddTool(server, tools.SpoolmanSetSpoolTool(), tools.NewSpoolmanSetSpoolHandler(api))
	mcp.AddTool(server, tools.SpoolmanProxyTool(), tools.NewSpoolmanProxyHandler(api))
	mcp.AddTool(server, tools.AnalysisStatusTool(), tools.NewAnalysisStatusHandler(api))
	mcp.AddTool(server, tools.AnalysisEstimateTool(), tools.NewAnalysisEstimateHandler(api))
	mcp.AddTool(server, tools.AnalysisProcessTool(), tools.NewAnalysisProcessHandler(api))
	mcp.AddTool(server, tools.AnalysisDumpConfigTool(), tools.NewAnalysisDumpConfigHandler(api))
	mcp.AddTool(server, tools.MQTTPublishTool(), tools.NewMQTTPublishHandler(api))
	mcp.AddTool(server, tools.MQTTSubscribeTool(), tools.NewMQTTSubscribeHandler(api))
	mcp.AddTool(server, tools.ExtensionsListTool(), tools.NewExtensionsListHandler(api))
	mcp.AddTool(server, tools.ExtensionsRequestTool(), tools.NewExtensionsRequestHandler(api))
	mcp.AddTool(server, tools.NotifiersListTool(), tools.NewNotifiersListHandler(api))
}

// registerAdminTools registers the OS, service, update, and user-management
// tools, gated behind MOONRAKER_ENABLE_ADMIN. Every tool in this set carries a
// destructive annotation: the gate marks them as the dangerous, hard-to-undo
// surface, so a client should prompt before any of them. TestAdminToolsDestructive
// enforces this invariant over the whole set.
func registerAdminTools(server *mcp.Server, api moonraker.API) {
	mcp.AddTool(server, tools.LogsRolloverTool(), tools.NewLogsRolloverHandler(api))
	mcp.AddTool(server, tools.ServerRestartTool(), tools.NewServerRestartHandler(api))
	mcp.AddTool(server, tools.MachineShutdownTool(), tools.NewMachineShutdownHandler(api))
	mcp.AddTool(server, tools.MachineRebootTool(), tools.NewMachineRebootHandler(api))
	mcp.AddTool(server, tools.ServiceStartTool(), tools.NewServiceStartHandler(api))
	mcp.AddTool(server, tools.ServiceStopTool(), tools.NewServiceStopHandler(api))
	mcp.AddTool(server, tools.ServiceRestartTool(), tools.NewServiceRestartHandler(api))
	mcp.AddTool(server, tools.SudoPasswordTool(), tools.NewSudoPasswordHandler(api))
	mcp.AddTool(server, tools.UpdateRefreshTool(), tools.NewUpdateRefreshHandler(api))
	mcp.AddTool(server, tools.UpdateUpgradeTool(), tools.NewUpdateUpgradeHandler(api))
	mcp.AddTool(server, tools.UpdateRecoverTool(), tools.NewUpdateRecoverHandler(api))
	mcp.AddTool(server, tools.UpdateRollbackTool(), tools.NewUpdateRollbackHandler(api))
	mcp.AddTool(server, tools.AccessCreateUserTool(), tools.NewAccessCreateUserHandler(api))
	mcp.AddTool(server, tools.AccessDeleteUserTool(), tools.NewAccessDeleteUserHandler(api))
	mcp.AddTool(server, tools.AccessUserPasswordTool(), tools.NewAccessUserPasswordHandler(api))
	mcp.AddTool(server, tools.AccessCreateAPIKeyTool(), tools.NewAccessCreateAPIKeyHandler(api))
}

// serve runs the stdio transport and, when configured, an HTTP transport.
func serve(logger *slog.Logger, server *mcp.Server, cfg *config.Config) error {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM)

	go func() {
		select {
		case <-sigChan:
			cancel()
		case <-ctx.Done():
		}

		signal.Stop(sigChan)
	}()

	group, groupCtx := errgroup.WithContext(ctx)
	httpEnabled := cfg.HTTPEnabled()

	group.Go(func() error {
		runErr := server.Run(groupCtx, &mcp.StdioTransport{})
		if runErr != nil && groupCtx.Err() == nil {
			return errors.Wrap(runErr, "stdio server failed")
		}

		if !httpEnabled {
			cancel()
		}

		return nil
	})

	if httpEnabled {
		group.Go(func() error {
			return runHTTPServer(groupCtx, logger, server, cfg.HTTPAddr(), cfg.HTTPToken)
		})
	}

	//nolint:wrapcheck // errors are already wrapped inside the group goroutines.
	return group.Wait()
}

// runHTTPServer starts an HTTP transport for the MCP server. Sharing a single
// *mcp.Server across transports is safe: the SDK guards internal state with a
// mutex. When token is set, every request must carry a matching Bearer token;
// otherwise config validation has confined the transport to a loopback host.
func runHTTPServer(ctx context.Context, logger *slog.Logger, server *mcp.Server, addr, token string) error {
	var handler http.Handler = mcp.NewStreamableHTTPHandler(
		func(_ *http.Request) *mcp.Server { return server },
		nil,
	)

	handler = bearerAuth(handler, token)

	httpServer := &http.Server{
		Addr:              addr,
		Handler:           handler,
		ReadHeaderTimeout: readHeaderTimeout,
	}

	//nolint:gosec // G118: shutdown uses a fresh context because ctx is already cancelled.
	go func() {
		<-ctx.Done()

		shutdownCtx, shutdownCancel := context.WithTimeout(context.Background(), shutdownTimeout)
		defer shutdownCancel()

		shutdownErr := httpServer.Shutdown(shutdownCtx) //nolint:contextcheck // fresh context for graceful shutdown.
		if shutdownErr != nil {
			logger.Error("http server shutdown failed", slog.Any("error", shutdownErr))
		}
	}()

	logger.Info("http server listening", slog.String("addr", addr))

	listenErr := httpServer.ListenAndServe()
	if errors.Is(listenErr, http.ErrServerClosed) {
		return nil
	}

	return errors.Wrap(listenErr, "HTTP listen failed")
}

// bearerAuth wraps next so every request must present a matching
// "Authorization: Bearer <token>" header. An empty token disables the check;
// config validation has already confined that case to a loopback host.
func bearerAuth(next http.Handler, token string) http.Handler {
	if token == "" {
		return next
	}

	want := []byte("Bearer " + token)

	return http.HandlerFunc(func(writer http.ResponseWriter, request *http.Request) {
		got := []byte(request.Header.Get("Authorization"))
		if subtle.ConstantTimeCompare(got, want) != 1 {
			http.Error(writer, "unauthorized", http.StatusUnauthorized)

			return
		}

		next.ServeHTTP(writer, request)
	})
}
