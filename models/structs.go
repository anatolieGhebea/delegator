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
	HeaderBitBucket string      = "X-Event-Key"
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
	// Event      string                 `json:"event"`
	Ref        string                 `json:"ref"`
	Repository map[string]interface{} `json:"repository"`
}

// manage bitbucket event
type BitBucketEventSource struct {
	Push       Push                   `json:"push"`
	Repository map[string]interface{} `json:"repository"`
}

type Push struct {
	Changes []Change `json:"changes"`
}

type Change struct {
	Change NewBranch `json:"new"`
}

type NewBranch struct {
	Name string `json:"name"`
}

var Configuration Config = Config{}
