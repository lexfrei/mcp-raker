package tools

import (
	"context"
	"net/url"
	"strconv"

	"github.com/modelcontextprotocol/go-sdk/mcp"

	"github.com/lexfrei/mcp-raker/internal/moonraker"
)

// AnnouncementsListParams defines the parameters for moonraker_announcements_list.
type AnnouncementsListParams struct {
	IncludeDismissed bool `json:"include_dismissed,omitempty" jsonschema:"When true, also include dismissed announcements"`
}

// AnnouncementsListTool returns the definition for moonraker_announcements_list.
func AnnouncementsListTool() *mcp.Tool {
	return &mcp.Tool{
		Name:        "moonraker_announcements_list",
		Description: "List Moonraker announcements and alerts (GET /server/announcements/list).",
		Annotations: readOnly("List Announcements"),
	}
}

// NewAnnouncementsListHandler creates the handler for moonraker_announcements_list.
func NewAnnouncementsListHandler(api moonraker.API) mcp.ToolHandlerFor[AnnouncementsListParams, map[string]any] {
	return func(ctx context.Context, _ *mcp.CallToolRequest, params AnnouncementsListParams) (*mcp.CallToolResult, map[string]any, error) {
		query := url.Values{}
		if params.IncludeDismissed {
			query.Set("include_dismissed", "true")
		}

		out, err := decodeResult(api.Get(ctx, "/server/announcements/list", query))

		return nil, out, err
	}
}

// AnnouncementsUpdateTool returns the definition for moonraker_announcements_update.
func AnnouncementsUpdateTool() *mcp.Tool {
	return &mcp.Tool{
		Name:        "moonraker_announcements_update",
		Description: "Check the configured feeds for new announcements (POST /server/announcements/update).",
		Annotations: write("Check Announcements"),
	}
}

// NewAnnouncementsUpdateHandler creates the handler for moonraker_announcements_update.
func NewAnnouncementsUpdateHandler(api moonraker.API) mcp.ToolHandlerFor[NoParams, map[string]any] {
	return func(ctx context.Context, _ *mcp.CallToolRequest, _ NoParams) (*mcp.CallToolResult, map[string]any, error) {
		out, err := decodeResult(api.Post(ctx, "/server/announcements/update", nil, nil))

		return nil, out, err
	}
}

// AnnouncementsDismissParams defines the parameters for moonraker_announcements_dismiss.
type AnnouncementsDismissParams struct {
	EntryID  string `json:"entry_id"            jsonschema:"Identifier of the announcement to dismiss"`
	WakeTime int    `json:"wake_time,omitempty" jsonschema:"Optional seconds after which the announcement reappears"`
}

// AnnouncementsDismissTool returns the definition for moonraker_announcements_dismiss.
func AnnouncementsDismissTool() *mcp.Tool {
	return &mcp.Tool{
		Name:        "moonraker_announcements_dismiss",
		Description: "Dismiss an announcement (POST /server/announcements/dismiss).",
		Annotations: write("Dismiss Announcement"),
	}
}

// NewAnnouncementsDismissHandler creates the handler for moonraker_announcements_dismiss.
func NewAnnouncementsDismissHandler(api moonraker.API) mcp.ToolHandlerFor[AnnouncementsDismissParams, map[string]any] {
	return func(ctx context.Context, _ *mcp.CallToolRequest, params AnnouncementsDismissParams) (*mcp.CallToolResult, map[string]any, error) {
		valErr := requireString(paramEntryID, params.EntryID)
		if valErr != nil {
			return nil, map[string]any{}, valErr
		}

		query := url.Values{paramEntryID: {params.EntryID}}
		if params.WakeTime > 0 {
			query.Set("wake_time", strconv.Itoa(params.WakeTime))
		}

		out, err := decodeResult(api.Post(ctx, "/server/announcements/dismiss", query, nil))

		return nil, out, err
	}
}

// AnnouncementsFeedsTool returns the definition for moonraker_announcements_feeds.
func AnnouncementsFeedsTool() *mcp.Tool {
	return &mcp.Tool{
		Name:        "moonraker_announcements_feeds",
		Description: "List the configured announcement feeds (GET /server/announcements/feeds).",
		Annotations: readOnly("List Feeds"),
	}
}

// NewAnnouncementsFeedsHandler creates the handler for moonraker_announcements_feeds.
func NewAnnouncementsFeedsHandler(api moonraker.API) mcp.ToolHandlerFor[NoParams, map[string]any] {
	return func(ctx context.Context, _ *mcp.CallToolRequest, _ NoParams) (*mcp.CallToolResult, map[string]any, error) {
		out, err := decodeResult(api.Get(ctx, "/server/announcements/feeds", nil))

		return nil, out, err
	}
}

// FeedParams names an announcement feed.
type FeedParams struct {
	Name string `json:"name" jsonschema:"Name of the announcement feed"`
}

// AnnouncementsAddFeedTool returns the definition for moonraker_announcements_add_feed.
func AnnouncementsAddFeedTool() *mcp.Tool {
	return &mcp.Tool{
		Name:        "moonraker_announcements_add_feed",
		Description: "Subscribe to an announcement feed (POST /server/announcements/feed).",
		Annotations: write("Add Feed"),
	}
}

// NewAnnouncementsAddFeedHandler creates the handler for moonraker_announcements_add_feed.
func NewAnnouncementsAddFeedHandler(api moonraker.API) mcp.ToolHandlerFor[FeedParams, map[string]any] {
	return func(ctx context.Context, _ *mcp.CallToolRequest, params FeedParams) (*mcp.CallToolResult, map[string]any, error) {
		valErr := requireString(paramName, params.Name)
		if valErr != nil {
			return nil, map[string]any{}, valErr
		}

		out, err := decodeResult(api.Post(ctx, "/server/announcements/feed", url.Values{paramName: {params.Name}}, nil))

		return nil, out, err
	}
}

// AnnouncementsRemoveFeedTool returns the definition for moonraker_announcements_remove_feed.
func AnnouncementsRemoveFeedTool() *mcp.Tool {
	return &mcp.Tool{
		Name:        "moonraker_announcements_remove_feed",
		Description: "Unsubscribe from an announcement feed (DELETE /server/announcements/feed).",
		Annotations: writeDestructive("Remove Feed"),
	}
}

// NewAnnouncementsRemoveFeedHandler creates the handler for moonraker_announcements_remove_feed.
func NewAnnouncementsRemoveFeedHandler(api moonraker.API) mcp.ToolHandlerFor[FeedParams, map[string]any] {
	return func(ctx context.Context, _ *mcp.CallToolRequest, params FeedParams) (*mcp.CallToolResult, map[string]any, error) {
		valErr := requireString(paramName, params.Name)
		if valErr != nil {
			return nil, map[string]any{}, valErr
		}

		out, err := decodeResult(api.Delete(ctx, "/server/announcements/feed", url.Values{paramName: {params.Name}}))

		return nil, out, err
	}
}
