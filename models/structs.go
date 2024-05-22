package models

type BranchCheck string
type EventSource string

const (
	CurrentBranch   BranchCheck = "current"
	SpecificBranch  BranchCheck = "specific"
	GitHubHook      EventSource = "github"
	BitBucketHook   EventSource = "bitbucket"
	GenericHook     EventSource = "generic"
	HeaderGitHub    string      = "X-Github-Event"
	HeaderBitBucket string      = "X-Bitbucket-Event"
)

type Response struct {
	Message string `json:"message"`
}

type Server struct {
	Port             string `json:"port"`
	LogRetentionDays int    `json:"log_retention_days"`
}

type EventHook struct {
	EventSource    EventSource // Default: generic. [ github, bitbucket, generic]
	Name           string
	RepositoryName string
	AbsolutePath   string
	SharedSecret   string
	SyncBranch     BranchCheck // Default: current. [ current > , <branch_name>]
	BranchName     string      //
}

type Config struct {
	Server   Server      `json:"server"`
	Triggers []EventHook `json:"triggers"`
}

type GenericEventSource struct {
	Name         string `json:"Name"`
	SharedSecret string `json:"SharedSecret"`
}

type GitHubEventSource struct {
	Event      string                 `json:"event"`
	Repository map[string]interface{} `json:"repository"`
	Ref        string                 `json:"ref"`
}
type BitBucketEventSource struct {
	Event      string                 `json:"event"`
	Repository map[string]interface{} `json:"repository"`
}

var Configuration Config = Config{}
