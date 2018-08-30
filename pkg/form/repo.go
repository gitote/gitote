package form

import (
	"net/url"
	"strings"

	"github.com/go-macaron/binding"
	"gitote.com/yoginth/com"
	"gopkg.in/macaron.v1"

	"gitote.com/gitote/gitote/models"
)

type CreateRepo struct {
	UserID      int64  `binding:"Required"`
	RepoName    string `binding:"Required;AlphaDashDot;MaxSize(100)"`
	Private     bool
	Description string `binding:"MaxSize(512)"`
	AutoInit    bool
	Gitignores  string
	License     string
	Readme      string
}

func (f *CreateRepo) Validate(ctx *macaron.Context, errs binding.Errors) binding.Errors {
	return validate(errs, ctx.Data, f, ctx.Locale)
}

type MigrateRepo struct {
	CloneAddr    string `json:"clone_addr" binding:"Required"`
	AuthUsername string `json:"auth_username"`
	AuthPassword string `json:"auth_password"`
	Uid          int64  `json:"uid" binding:"Required"`
	RepoName     string `json:"repo_name" binding:"Required;AlphaDashDot;MaxSize(100)"`
	Mirror       bool   `json:"mirror"`
	Private      bool   `json:"private"`
	Description  string `json:"description" binding:"MaxSize(512)"`
}

func (f *MigrateRepo) Validate(ctx *macaron.Context, errs binding.Errors) binding.Errors {
	return validate(errs, ctx.Data, f, ctx.Locale)
}

// ParseRemoteAddr checks if given remote address is valid,
// and returns composed URL with needed username and password.
// It also checks if given user has permission when remote address
// is actually a local path.
func (f MigrateRepo) ParseRemoteAddr(user *models.User) (string, error) {
	remoteAddr := strings.TrimSpace(f.CloneAddr)

	// Remote address can be HTTP/HTTPS/Git URL or local path.
	if strings.HasPrefix(remoteAddr, "http://") ||
		strings.HasPrefix(remoteAddr, "https://") ||
		strings.HasPrefix(remoteAddr, "git://") {
		u, err := url.Parse(remoteAddr)
		if err != nil {
			return "", models.ErrInvalidCloneAddr{IsURLError: true}
		}
		if len(f.AuthUsername)+len(f.AuthPassword) > 0 {
			u.User = url.UserPassword(f.AuthUsername, f.AuthPassword)
		}
		remoteAddr = u.String()
	} else if !user.CanImportLocal() {
		return "", models.ErrInvalidCloneAddr{IsPermissionDenied: true}
	} else if !com.IsDir(remoteAddr) {
		return "", models.ErrInvalidCloneAddr{IsInvalidPath: true}
	}

	return remoteAddr, nil
}

type RepoSetting struct {
	RepoName      string `binding:"Required;AlphaDashDot;MaxSize(100)"`
	Description   string `binding:"MaxSize(512)"`
	Website       string `binding:"Url;MaxSize(100)"`
	Branch        string
	Interval      int
	MirrorAddress string
	Private       bool
	EnablePrune   bool

	// Advanced settings
	EnableWiki            bool
	AllowPublicWiki       bool
	EnableExternalWiki    bool
	ExternalWikiURL       string
	EnableIssues          bool
	AllowPublicIssues     bool
	EnableExternalTracker bool
	ExternalTrackerURL    string
	TrackerURLFormat      string
	TrackerIssueStyle     string
	EnablePulls           bool
	PullsIgnoreWhitespace bool
	PullsAllowRebase      bool
}

func (f *RepoSetting) Validate(ctx *macaron.Context, errs binding.Errors) binding.Errors {
	return validate(errs, ctx.Data, f, ctx.Locale)
}

type ProtectBranch struct {
	Protected          bool
	RequirePullRequest bool
	EnableWhitelist    bool
	WhitelistUsers     string
	WhitelistTeams     string
}

func (f *ProtectBranch) Validate(ctx *macaron.Context, errs binding.Errors) binding.Errors {
	return validate(errs, ctx.Data, f, ctx.Locale)
}

type Webhook struct {
	Events       string
	Create       bool
	Delete       bool
	Fork         bool
	Push         bool
	Issues       bool
	IssueComment bool
	PullRequest  bool
	Release      bool
	Active       bool
}

func (f Webhook) PushOnly() bool {
	return f.Events == "push_only"
}

func (f Webhook) SendEverything() bool {
	return f.Events == "send_everything"
}

func (f Webhook) ChooseEvents() bool {
	return f.Events == "choose_events"
}

type NewWebhook struct {
	PayloadURL  string `binding:"Required;Url"`
	ContentType int    `binding:"Required"`
	Secret      string
	Webhook
}

func (f *NewWebhook) Validate(ctx *macaron.Context, errs binding.Errors) binding.Errors {
	return validate(errs, ctx.Data, f, ctx.Locale)
}

type NewSlackHook struct {
	PayloadURL string `binding:"Required;Url"`
	Channel    string `binding:"Required"`
	Username   string
	IconURL    string
	Color      string
	Webhook
}

func (f *NewSlackHook) Validate(ctx *macaron.Context, errs binding.Errors) binding.Errors {
	return validate(errs, ctx.Data, f, ctx.Locale)
}

type NewDiscordHook struct {
	PayloadURL string `binding:"Required;Url"`
	Username   string
	IconURL    string
	Color      string
	Webhook
}

func (f *NewDiscordHook) Validate(ctx *macaron.Context, errs binding.Errors) binding.Errors {
	return validate(errs, ctx.Data, f, ctx.Locale)
}

type NewDingtalkHook struct {
	PayloadURL string `binding:"Required;Url"`
	Webhook
}

func (f *NewDingtalkHook) Validate(ctx *macaron.Context, errs binding.Errors) binding.Errors {
	return validate(errs, ctx.Data, f, ctx.Locale)
}

type NewIssue struct {
	Title       string `binding:"Required;MaxSize(255)"`
	LabelIDs    string `form:"label_ids"`
	MilestoneID int64
	AssigneeID  int64
	Content     string
	Files       []string
}

func (f *NewIssue) Validate(ctx *macaron.Context, errs binding.Errors) binding.Errors {
	return validate(errs, ctx.Data, f, ctx.Locale)
}

type CreateComment struct {
	Content string
	Status  string `binding:"OmitEmpty;In(reopen,close)"`
	Files   []string
}

func (f *CreateComment) Validate(ctx *macaron.Context, errs binding.Errors) binding.Errors {
	return validate(errs, ctx.Data, f, ctx.Locale)
}

type CreateMilestone struct {
	Title    string `binding:"Required;MaxSize(50)"`
	Content  string
	Deadline string
}

func (f *CreateMilestone) Validate(ctx *macaron.Context, errs binding.Errors) binding.Errors {
	return validate(errs, ctx.Data, f, ctx.Locale)
}

type CreateLabel struct {
	ID    int64
	Title string `binding:"Required;MaxSize(50)" locale:"repo.issues.label_title"`
	Color string `binding:"Required;Size(7)" locale:"repo.issues.label_color"`
}

func (f *CreateLabel) Validate(ctx *macaron.Context, errs binding.Errors) binding.Errors {
	return validate(errs, ctx.Data, f, ctx.Locale)
}

type InitializeLabels struct {
	TemplateName string `binding:"Required"`
}

func (f *InitializeLabels) Validate(ctx *macaron.Context, errs binding.Errors) binding.Errors {
	return validate(errs, ctx.Data, f, ctx.Locale)
}

type NewRelease struct {
	TagName    string `binding:"Required"`
	Target     string `form:"tag_target" binding:"Required"`
	Title      string `binding:"Required"`
	Content    string
	Draft      string
	Prerelease bool
	Files      []string
}

func (f *NewRelease) Validate(ctx *macaron.Context, errs binding.Errors) binding.Errors {
	return validate(errs, ctx.Data, f, ctx.Locale)
}

type EditRelease struct {
	Title      string `binding:"Required"`
	Content    string
	Draft      string
	Prerelease bool
	Files      []string
}

func (f *EditRelease) Validate(ctx *macaron.Context, errs binding.Errors) binding.Errors {
	return validate(errs, ctx.Data, f, ctx.Locale)
}

type NewWiki struct {
	OldTitle string
	Title    string `binding:"Required"`
	Content  string `binding:"Required"`
	Message  string
}

// FIXME: use code generation to generate this method.
func (f *NewWiki) Validate(ctx *macaron.Context, errs binding.Errors) binding.Errors {
	return validate(errs, ctx.Data, f, ctx.Locale)
}

type EditRepoFile struct {
	TreePath      string `binding:"Required;MaxSize(500)"`
	Content       string `binding:"Required"`
	CommitSummary string `binding:"MaxSize(100)`
	CommitMessage string
	CommitChoice  string `binding:"Required;MaxSize(50)"`
	NewBranchName string `binding:"AlphaDashDotSlash;MaxSize(100)"`
	LastCommit    string
}

func (f *EditRepoFile) Validate(ctx *macaron.Context, errs binding.Errors) binding.Errors {
	return validate(errs, ctx.Data, f, ctx.Locale)
}

func (f *EditRepoFile) IsNewBrnach() bool {
	return f.CommitChoice == "commit-to-new-branch"
}

type EditPreviewDiff struct {
	Content string
}

func (f *EditPreviewDiff) Validate(ctx *macaron.Context, errs binding.Errors) binding.Errors {
	return validate(errs, ctx.Data, f, ctx.Locale)
}

type UploadRepoFile struct {
	TreePath      string `binding:MaxSize(500)"`
	CommitSummary string `binding:"MaxSize(100)`
	CommitMessage string
	CommitChoice  string `binding:"Required;MaxSize(50)"`
	NewBranchName string `binding:"AlphaDashDot;MaxSize(100)"`
	Files         []string
}

func (f *UploadRepoFile) Validate(ctx *macaron.Context, errs binding.Errors) binding.Errors {
	return validate(errs, ctx.Data, f, ctx.Locale)
}

func (f *UploadRepoFile) IsNewBrnach() bool {
	return f.CommitChoice == "commit-to-new-branch"
}

type RemoveUploadFile struct {
	File string `binding:"Required;MaxSize(50)"`
}

func (f *RemoveUploadFile) Validate(ctx *macaron.Context, errs binding.Errors) binding.Errors {
	return validate(errs, ctx.Data, f, ctx.Locale)
}

type DeleteRepoFile struct {
	CommitSummary string `binding:"MaxSize(100)`
	CommitMessage string
	CommitChoice  string `binding:"Required;MaxSize(50)"`
	NewBranchName string `binding:"AlphaDashDot;MaxSize(100)"`
}

func (f *DeleteRepoFile) Validate(ctx *macaron.Context, errs binding.Errors) binding.Errors {
	return validate(errs, ctx.Data, f, ctx.Locale)
}

func (f *DeleteRepoFile) IsNewBrnach() bool {
	return f.CommitChoice == "commit-to-new-branch"
}
