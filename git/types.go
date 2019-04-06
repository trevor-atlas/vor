package git

import (
	"encoding/json"
	"fmt"
	"time"

	"trevoratlas.com/vor/utils"
)

// PullRequestBody The POST body to create a pull request on github
type PullRequestBody struct {
	Title string `json:"title,omitempty"`
	Body  string `json:"body,omitempty"`
	Head  string `json:"head,omitempty"`
	Base  string `json:"base,omitempty"`
}

type JSONSerializer interface {
	Print()
	Marshal() error
	Unmarshal([]byte) error
}

func (prs *PullRequestResponse) PrintJSON() {
	jsonString, parseErr := prs.Marshal()
	if parseErr != nil {
		utils.Exit("error marshalling json")
	}
	fmt.Printf("%s\n", jsonString)
}

func (prs *PullRequestResponse) Marshal() (string, error) {
	data, err := json.MarshalIndent(prs, "", "    ")
	if err != nil { return "", err }
	return string(data), nil
}

func (prs * PullRequestResponse) Unmarshal(b []byte) error {
	return json.Unmarshal(b, prs)
}

type PullRequestResponse struct {
	Errors            []*struct{
		Resource string `json:"resource,omitempty"`
		Field    string `json:"field,omitempty"`
		Code     string `json:"code,omitempty"`
	} `json:"errors,omitempty"`
	Documentation_url string `json:"documentation_url,omitempty"`
	Message           string `json:"message,omitempty"`
	ID                int    `json:"id,omitempty"`
	NodeID            string `json:"node_id,omitempty"`
	URL               string `json:"url,omitempty"`
	HTMLURL           string `json:"html_url,omitempty"`
	DiffURL           string `json:"diff_url,omitempty"`
	PatchURL          string `json:"patch_url,omitempty"`
	IssueURL          string `json:"issue_url,omitempty"`
	CommitsURL        string `json:"commits_url,omitempty"`
	ReviewCommentsURL string `json:"review_comments_url,omitempty"`
	ReviewCommentURL  string `json:"review_comment_url,omitempty"`
	CommentsURL       string `json:"comments_url,omitempty"`
	StatusesURL       string `json:"statuses_url,omitempty"`
	Number            int    `json:"number,omitempty"`
	State             string `json:"state,omitempty"`
	Title             string `json:"title,omitempty"`
	Body              string `json:"body,omitempty"`
	Assignee          *struct {
		Login             string `json:"login,omitempty"`
		ID                int    `json:"id,omitempty"`
		NodeID            string `json:"node_id,omitempty"`
		AvatarURL         string `json:"avatar_url,omitempty"`
		GravatarID        string `json:"gravatar_id,omitempty"`
		URL               string `json:"url,omitempty"`
		HTMLURL           string `json:"html_url,omitempty"`
		FollowersURL      string `json:"followers_url,omitempty"`
		FollowingURL      string `json:"following_url,omitempty"`
		GistsURL          string `json:"gists_url,omitempty"`
		StarredURL        string `json:"starred_url,omitempty"`
		SubscriptionsURL  string `json:"subscriptions_url,omitempty"`
		OrganizationsURL  string `json:"organizations_url,omitempty"`
		ReposURL          string `json:"repos_url,omitempty"`
		EventsURL         string `json:"events_url,omitempty"`
		ReceivedEventsURL string `json:"received_events_url,omitempty"`
		Type              string `json:"type,omitempty"`
		SiteAdmin         bool   `json:"site_admin,omitempty"`
	} `json:"assignee,omitempty"`
	Labels []*struct {
		ID          int    `json:"id,omitempty"`
		NodeID      string `json:"node_id,omitempty"`
		URL         string `json:"url,omitempty"`
		Name        string `json:"name,omitempty"`
		Description string `json:"description,omitempty"`
		Color       string `json:"color,omitempty"`
		Default     bool   `json:"default,omitempty"`
	} `json:"labels,omitempty"`
	Milestone *struct {
		URL         string `json:"url,omitempty"`
		HTMLURL     string `json:"html_url,omitempty"`
		LabelsURL   string `json:"labels_url,omitempty"`
		ID          int    `json:"id,omitempty"`
		NodeID      string `json:"node_id,omitempty"`
		Number      int    `json:"number,omitempty"`
		State       string `json:"state,omitempty"`
		Title       string `json:"title,omitempty"`
		Description string `json:"description,omitempty"`
		Creator     *struct {
			Login             string `json:"login,omitempty"`
			ID                int    `json:"id,omitempty"`
			NodeID            string `json:"node_id,omitempty"`
			AvatarURL         string `json:"avatar_url,omitempty"`
			GravatarID        string `json:"gravatar_id,omitempty"`
			URL               string `json:"url,omitempty"`
			HTMLURL           string `json:"html_url,omitempty"`
			FollowersURL      string `json:"followers_url,omitempty"`
			FollowingURL      string `json:"following_url,omitempty"`
			GistsURL          string `json:"gists_url,omitempty"`
			StarredURL        string `json:"starred_url,omitempty"`
			SubscriptionsURL  string `json:"subscriptions_url,omitempty"`
			OrganizationsURL  string `json:"organizations_url,omitempty"`
			ReposURL          string `json:"repos_url,omitempty"`
			EventsURL         string `json:"events_url,omitempty"`
			ReceivedEventsURL string `json:"received_events_url,omitempty"`
			Type              string `json:"type,omitempty"`
			SiteAdmin         bool   `json:"site_admin,omitempty"`
		} `json:"creator,omitempty"`
		OpenIssues   int       `json:"open_issues,omitempty"`
		ClosedIssues int       `json:"closed_issues,omitempty"`
		CreatedAt    *time.Time `json:"created_at,omitempty"`
		UpdatedAt    *time.Time `json:"updated_at,omitempty"`
		ClosedAt     *time.Time `json:"closed_at,omitempty"`
		DueOn        *time.Time `json:"due_on,omitempty"`
	} `json:"milestone,omitempty"`
	Locked           bool      `json:"locked,omitempty"`
	ActiveLockReason string    `json:"active_lock_reason,omitempty"`
	CreatedAt        *time.Time `json:"created_at,omitempty"`
	UpdatedAt        *time.Time `json:"updated_at,omitempty"`
	ClosedAt         *time.Time `json:"closed_at,omitempty"`
	MergedAt         *time.Time `json:"merged_at,omitempty"`
	Head             *struct {
		Label string `json:"label,omitempty"`
		Ref   string `json:"ref,omitempty"`
		Sha   string `json:"sha,omitempty"`
		User  *struct {
			Login             string `json:"login,omitempty"`
			ID                int    `json:"id,omitempty"`
			NodeID            string `json:"node_id,omitempty"`
			AvatarURL         string `json:"avatar_url,omitempty"`
			GravatarID        string `json:"gravatar_id,omitempty"`
			URL               string `json:"url,omitempty"`
			HTMLURL           string `json:"html_url,omitempty"`
			FollowersURL      string `json:"followers_url,omitempty"`
			FollowingURL      string `json:"following_url,omitempty"`
			GistsURL          string `json:"gists_url,omitempty"`
			StarredURL        string `json:"starred_url,omitempty"`
			SubscriptionsURL  string `json:"subscriptions_url,omitempty"`
			OrganizationsURL  string `json:"organizations_url,omitempty"`
			ReposURL          string `json:"repos_url,omitempty"`
			EventsURL         string `json:"events_url,omitempty"`
			ReceivedEventsURL string `json:"received_events_url,omitempty"`
			Type              string `json:"type,omitempty"`
			SiteAdmin         bool   `json:"site_admin,omitempty"`
		} `json:"user,omitempty"`
		Repo *struct {
			ID       int    `json:"id,omitempty"`
			NodeID   string `json:"node_id,omitempty"`
			Name     string `json:"name,omitempty"`
			FullName string `json:"full_name,omitempty"`
			Owner    *struct {
				Login             string `json:"login,omitempty"`
				ID                int    `json:"id,omitempty"`
				NodeID            string `json:"node_id,omitempty"`
				AvatarURL         string `json:"avatar_url,omitempty"`
				GravatarID        string `json:"gravatar_id,omitempty"`
				URL               string `json:"url,omitempty"`
				HTMLURL           string `json:"html_url,omitempty"`
				FollowersURL      string `json:"followers_url,omitempty"`
				FollowingURL      string `json:"following_url,omitempty"`
				GistsURL          string `json:"gists_url,omitempty"`
				StarredURL        string `json:"starred_url,omitempty"`
				SubscriptionsURL  string `json:"subscriptions_url,omitempty"`
				OrganizationsURL  string `json:"organizations_url,omitempty"`
				ReposURL          string `json:"repos_url,omitempty"`
				EventsURL         string `json:"events_url,omitempty"`
				ReceivedEventsURL string `json:"received_events_url,omitempty"`
				Type              string `json:"type,omitempty"`
				SiteAdmin         bool   `json:"site_admin,omitempty"`
			} `json:"owner,omitempty"`
			Private          bool        `json:"private,omitempty"`
			HTMLURL          string      `json:"html_url,omitempty"`
			Description      string      `json:"description,omitempty"`
			Fork             bool        `json:"fork,omitempty"`
			URL              string      `json:"url,omitempty"`
			ArchiveURL       string      `json:"archive_url,omitempty"`
			AssigneesURL     string      `json:"assignees_url,omitempty"`
			BlobsURL         string      `json:"blobs_url,omitempty"`
			BranchesURL      string      `json:"branches_url,omitempty"`
			CollaboratorsURL string      `json:"collaborators_url,omitempty"`
			CommentsURL      string      `json:"comments_url,omitempty"`
			CommitsURL       string      `json:"commits_url,omitempty"`
			CompareURL       string      `json:"compare_url,omitempty"`
			ContentsURL      string      `json:"contents_url,omitempty"`
			ContributorsURL  string      `json:"contributors_url,omitempty"`
			DeploymentsURL   string      `json:"deployments_url,omitempty"`
			DownloadsURL     string      `json:"downloads_url,omitempty"`
			EventsURL        string      `json:"events_url,omitempty"`
			ForksURL         string      `json:"forks_url,omitempty"`
			GitCommitsURL    string      `json:"git_commits_url,omitempty"`
			GitRefsURL       string      `json:"git_refs_url,omitempty"`
			GitTagsURL       string      `json:"git_tags_url,omitempty"`
			GitURL           string      `json:"git_url,omitempty"`
			IssueCommentURL  string      `json:"issue_comment_url,omitempty"`
			IssueEventsURL   string      `json:"issue_events_url,omitempty"`
			IssuesURL        string      `json:"issues_url,omitempty"`
			KeysURL          string      `json:"keys_url,omitempty"`
			LabelsURL        string      `json:"labels_url,omitempty"`
			LanguagesURL     string      `json:"languages_url,omitempty"`
			MergesURL        string      `json:"merges_url,omitempty"`
			MilestonesURL    string      `json:"milestones_url,omitempty"`
			NotificationsURL string      `json:"notifications_url,omitempty"`
			PullsURL         string      `json:"pulls_url,omitempty"`
			ReleasesURL      string      `json:"releases_url,omitempty"`
			SSHURL           string      `json:"ssh_url,omitempty"`
			StargazersURL    string      `json:"stargazers_url,omitempty"`
			StatusesURL      string      `json:"statuses_url,omitempty"`
			SubscribersURL   string      `json:"subscribers_url,omitempty"`
			SubscriptionURL  string      `json:"subscription_url,omitempty"`
			TagsURL          string      `json:"tags_url,omitempty"`
			TeamsURL         string      `json:"teams_url,omitempty"`
			TreesURL         string      `json:"trees_url,omitempty"`
			CloneURL         string      `json:"clone_url,omitempty"`
			MirrorURL        string      `json:"mirror_url,omitempty"`
			HooksURL         string      `json:"hooks_url,omitempty"`
			SvnURL           string      `json:"svn_url,omitempty"`
			Homepage         string      `json:"homepage,omitempty"`
			Language         string       `json:"language,omitempty"`
			ForksCount       int         `json:"forks_count,omitempty"`
			StargazersCount  int         `json:"stargazers_count,omitempty"`
			WatchersCount    int         `json:"watchers_count,omitempty"`
			Size             int         `json:"size,omitempty"`
			DefaultBranch    string      `json:"default_branch,omitempty"`
			OpenIssuesCount  int         `json:"open_issues_count,omitempty"`
			Topics           []string    `json:"topics,omitempty"`
			HasIssues        bool        `json:"has_issues,omitempty"`
			HasProjects      bool        `json:"has_projects,omitempty"`
			HasWiki          bool        `json:"has_wiki,omitempty"`
			HasPages         bool        `json:"has_pages,omitempty"`
			HasDownloads     bool        `json:"has_downloads,omitempty"`
			Archived         bool        `json:"archived,omitempty"`
			PushedAt         *time.Time   `json:"pushed_at,omitempty"`
			CreatedAt        *time.Time   `json:"created_at,omitempty"`
			UpdatedAt        *time.Time   `json:"updated_at,omitempty"`
			Permissions      *struct {
				Admin bool `json:"admin,omitempty"`
				Push  bool `json:"push,omitempty"`
				Pull  bool `json:"pull,omitempty"`
			} `json:"permissions,omitempty"`
			AllowRebaseMerge bool `json:"allow_rebase_merge,omitempty"`
			AllowSquashMerge bool `json:"allow_squash_merge,omitempty"`
			AllowMergeCommit bool `json:"allow_merge_commit,omitempty"`
			SubscribersCount int  `json:"subscribers_count,omitempty"`
			NetworkCount     int  `json:"network_count,omitempty"`
		} `json:"repo,omitempty"`
	} `json:"head,omitempty"`
	Base *struct {
		Label string `json:"label,omitempty"`
		Ref   string `json:"ref,omitempty"`
		Sha   string `json:"sha,omitempty"`
		User  *struct {
			Login             string `json:"login,omitempty"`
			ID                int    `json:"id,omitempty"`
			NodeID            string `json:"node_id,omitempty"`
			AvatarURL         string `json:"avatar_url,omitempty"`
			GravatarID        string `json:"gravatar_id,omitempty"`
			URL               string `json:"url,omitempty"`
			HTMLURL           string `json:"html_url,omitempty"`
			FollowersURL      string `json:"followers_url,omitempty"`
			FollowingURL      string `json:"following_url,omitempty"`
			GistsURL          string `json:"gists_url,omitempty"`
			StarredURL        string `json:"starred_url,omitempty"`
			SubscriptionsURL  string `json:"subscriptions_url,omitempty"`
			OrganizationsURL  string `json:"organizations_url,omitempty"`
			ReposURL          string `json:"repos_url,omitempty"`
			EventsURL         string `json:"events_url,omitempty"`
			ReceivedEventsURL string `json:"received_events_url,omitempty"`
			Type              string `json:"type,omitempty"`
			SiteAdmin         bool   `json:"site_admin,omitempty"`
		} `json:"user,omitempty"`
		Repo *struct {
			ID       int    `json:"id,omitempty"`
			NodeID   string `json:"node_id,omitempty"`
			Name     string `json:"name,omitempty"`
			FullName string `json:"full_name,omitempty"`
			Owner    *struct {
				Login             string `json:"login,omitempty"`
				ID                int    `json:"id,omitempty"`
				NodeID            string `json:"node_id,omitempty"`
				AvatarURL         string `json:"avatar_url,omitempty"`
				GravatarID        string `json:"gravatar_id,omitempty"`
				URL               string `json:"url,omitempty"`
				HTMLURL           string `json:"html_url,omitempty"`
				FollowersURL      string `json:"followers_url,omitempty"`
				FollowingURL      string `json:"following_url,omitempty"`
				GistsURL          string `json:"gists_url,omitempty"`
				StarredURL        string `json:"starred_url,omitempty"`
				SubscriptionsURL  string `json:"subscriptions_url,omitempty"`
				OrganizationsURL  string `json:"organizations_url,omitempty"`
				ReposURL          string `json:"repos_url,omitempty"`
				EventsURL         string `json:"events_url,omitempty"`
				ReceivedEventsURL string `json:"received_events_url,omitempty"`
				Type              string `json:"type,omitempty"`
				SiteAdmin         bool   `json:"site_admin,omitempty"`
			} `json:"owner,omitempty"`
			Private          bool        `json:"private,omitempty"`
			HTMLURL          string      `json:"html_url,omitempty"`
			Description      string      `json:"description,omitempty"`
			Fork             bool        `json:"fork,omitempty"`
			URL              string      `json:"url,omitempty"`
			ArchiveURL       string      `json:"archive_url,omitempty"`
			AssigneesURL     string      `json:"assignees_url,omitempty"`
			BlobsURL         string      `json:"blobs_url,omitempty"`
			BranchesURL      string      `json:"branches_url,omitempty"`
			CollaboratorsURL string      `json:"collaborators_url,omitempty"`
			CommentsURL      string      `json:"comments_url,omitempty"`
			CommitsURL       string      `json:"commits_url,omitempty"`
			CompareURL       string      `json:"compare_url,omitempty"`
			ContentsURL      string      `json:"contents_url,omitempty"`
			ContributorsURL  string      `json:"contributors_url,omitempty"`
			DeploymentsURL   string      `json:"deployments_url,omitempty"`
			DownloadsURL     string      `json:"downloads_url,omitempty"`
			EventsURL        string      `json:"events_url,omitempty"`
			ForksURL         string      `json:"forks_url,omitempty"`
			GitCommitsURL    string      `json:"git_commits_url,omitempty"`
			GitRefsURL       string      `json:"git_refs_url,omitempty"`
			GitTagsURL       string      `json:"git_tags_url,omitempty"`
			GitURL           string      `json:"git_url,omitempty"`
			IssueCommentURL  string      `json:"issue_comment_url,omitempty"`
			IssueEventsURL   string      `json:"issue_events_url,omitempty"`
			IssuesURL        string      `json:"issues_url,omitempty"`
			KeysURL          string      `json:"keys_url,omitempty"`
			LabelsURL        string      `json:"labels_url,omitempty"`
			LanguagesURL     string      `json:"languages_url,omitempty"`
			MergesURL        string      `json:"merges_url,omitempty"`
			MilestonesURL    string      `json:"milestones_url,omitempty"`
			NotificationsURL string      `json:"notifications_url,omitempty"`
			PullsURL         string      `json:"pulls_url,omitempty"`
			ReleasesURL      string      `json:"releases_url,omitempty"`
			SSHURL           string      `json:"ssh_url,omitempty"`
			StargazersURL    string      `json:"stargazers_url,omitempty"`
			StatusesURL      string      `json:"statuses_url,omitempty"`
			SubscribersURL   string      `json:"subscribers_url,omitempty"`
			SubscriptionURL  string      `json:"subscription_url,omitempty"`
			TagsURL          string      `json:"tags_url,omitempty"`
			TeamsURL         string      `json:"teams_url,omitempty"`
			TreesURL         string      `json:"trees_url,omitempty"`
			CloneURL         string      `json:"clone_url,omitempty"`
			MirrorURL        string      `json:"mirror_url,omitempty"`
			HooksURL         string      `json:"hooks_url,omitempty"`
			SvnURL           string      `json:"svn_url,omitempty"`
			Homepage         string      `json:"homepage,omitempty"`
			Language         interface{} `json:"language,omitempty"`
			ForksCount       int         `json:"forks_count,omitempty"`
			StargazersCount  int         `json:"stargazers_count,omitempty"`
			WatchersCount    int         `json:"watchers_count,omitempty"`
			Size             int         `json:"size,omitempty"`
			DefaultBranch    string      `json:"default_branch,omitempty"`
			OpenIssuesCount  int         `json:"open_issues_count,omitempty"`
			Topics           []string    `json:"topics,omitempty"`
			HasIssues        bool        `json:"has_issues,omitempty"`
			HasProjects      bool        `json:"has_projects,omitempty"`
			HasWiki          bool        `json:"has_wiki,omitempty"`
			HasPages         bool        `json:"has_pages,omitempty"`
			HasDownloads     bool        `json:"has_downloads,omitempty"`
			Archived         bool        `json:"archived,omitempty"`
			PushedAt         *time.Time   `json:"pushed_at,omitempty"`
			CreatedAt        *time.Time   `json:"created_at,omitempty"`
			UpdatedAt        *time.Time   `json:"updated_at,omitempty"`
			Permissions      *struct {
				Admin bool `json:"admin,omitempty"`
				Push  bool `json:"push,omitempty"`
				Pull  bool `json:"pull,omitempty"`
			} `json:"permissions,omitempty"`
			AllowRebaseMerge bool `json:"allow_rebase_merge,omitempty"`
			AllowSquashMerge bool `json:"allow_squash_merge,omitempty"`
			AllowMergeCommit bool `json:"allow_merge_commit,omitempty"`
			SubscribersCount int  `json:"subscribers_count,omitempty"`
			NetworkCount     int  `json:"network_count,omitempty"`
		} `json:"repo,omitempty"`
	} `json:"base,omitempty"`
	Links *struct {
		Self *struct {
			Href string `json:"href,omitempty"`
		} `json:"self,omitempty"`
		HTML *struct {
			Href string `json:"href,omitempty"`
		} `json:"html,omitempty"`
		Issue *struct {
			Href string `json:"href,omitempty"`
		} `json:"issue,omitempty"`
		Comments *struct {
			Href string `json:"href,omitempty"`
		} `json:"comments,omitempty"`
		ReviewComments *struct {
			Href string `json:"href,omitempty"`
		} `json:"review_comments,omitempty"`
		ReviewComment *struct {
			Href string `json:"href,omitempty"`
		} `json:"review_comment,omitempty"`
		Commits *struct {
			Href string `json:"href,omitempty"`
		} `json:"commits,omitempty"`
		Statuses *struct {
			Href string `json:"href,omitempty"`
		} `json:"statuses,omitempty"`
	} `json:"_links,omitempty"`
	User *struct {
		Login             string `json:"login,omitempty"`
		ID                int    `json:"id,omitempty"`
		NodeID            string `json:"node_id,omitempty"`
		AvatarURL         string `json:"avatar_url,omitempty"`
		GravatarID        string `json:"gravatar_id,omitempty"`
		URL               string `json:"url,omitempty"`
		HTMLURL           string `json:"html_url,omitempty"`
		FollowersURL      string `json:"followers_url,omitempty"`
		FollowingURL      string `json:"following_url,omitempty"`
		GistsURL          string `json:"gists_url,omitempty"`
		StarredURL        string `json:"starred_url,omitempty"`
		SubscriptionsURL  string `json:"subscriptions_url,omitempty"`
		OrganizationsURL  string `json:"organizations_url,omitempty"`
		ReposURL          string `json:"repos_url,omitempty"`
		EventsURL         string `json:"events_url,omitempty"`
		ReceivedEventsURL string `json:"received_events_url,omitempty"`
		Type              string `json:"type,omitempty"`
		SiteAdmin         bool   `json:"site_admin,omitempty"`
	} `json:"user,omitempty"`
}
