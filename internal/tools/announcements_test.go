package tools_test

import (
	"testing"

	"github.com/cockroachdb/errors"

	"github.com/lexfrei/mcp-raker/internal/tools"
)

func TestAnnouncementsList_IncludeDismissed(t *testing.T) {
	t.Parallel()

	mock := &mockAPI{result: okJSON}

	_, _, err := tools.NewAnnouncementsListHandler(mock)(t.Context(), nil, tools.AnnouncementsListParams{IncludeDismissed: true})
	if err != nil {
		t.Fatalf("handler: %v", err)
	}

	assertCall(t, mock, methodGet, "/server/announcements/list")

	if mock.lastQuery.Get("include_dismissed") != queryTrue {
		t.Errorf("include_dismissed = %q, want true", mock.lastQuery.Get("include_dismissed"))
	}
}

func TestAnnouncementsUpdate(t *testing.T) {
	t.Parallel()

	mock := &mockAPI{result: okJSON}

	_, _, err := tools.NewAnnouncementsUpdateHandler(mock)(t.Context(), nil, tools.NoParams{})
	if err != nil {
		t.Fatalf("handler: %v", err)
	}

	assertCall(t, mock, methodPost, "/server/announcements/update")
}

func TestAnnouncementsDismiss_RequiresEntryID(t *testing.T) {
	t.Parallel()

	_, _, err := tools.NewAnnouncementsDismissHandler(&mockAPI{})(t.Context(), nil, tools.AnnouncementsDismissParams{})
	if !errors.Is(err, tools.ErrValidation) {
		t.Fatalf("err = %v, want ErrValidation", err)
	}
}

func TestAnnouncementsFeeds(t *testing.T) {
	t.Parallel()

	mock := &mockAPI{result: okJSON}

	_, _, err := tools.NewAnnouncementsFeedsHandler(mock)(t.Context(), nil, tools.NoParams{})
	if err != nil {
		t.Fatalf("handler: %v", err)
	}

	assertCall(t, mock, methodGet, "/server/announcements/feeds")
}

func TestAnnouncementsFeedMutations(t *testing.T) {
	t.Parallel()

	_, _, missing := tools.NewAnnouncementsAddFeedHandler(&mockAPI{})(t.Context(), nil, tools.FeedParams{})
	if !errors.Is(missing, tools.ErrValidation) {
		t.Errorf("missing name err = %v, want ErrValidation", missing)
	}

	add := &mockAPI{result: okJSON}

	_, _, addErr := tools.NewAnnouncementsAddFeedHandler(add)(t.Context(), nil, tools.FeedParams{Name: "mainsail"})
	if addErr != nil {
		t.Fatalf("add: %v", addErr)
	}

	assertCall(t, add, methodPost, "/server/announcements/feed")

	remove := &mockAPI{result: okJSON}

	_, _, removeErr := tools.NewAnnouncementsRemoveFeedHandler(remove)(t.Context(), nil, tools.FeedParams{Name: "mainsail"})
	if removeErr != nil {
		t.Fatalf("remove: %v", removeErr)
	}

	assertCall(t, remove, methodDelete, "/server/announcements/feed")
}
