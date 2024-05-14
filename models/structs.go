package models

type BranchCheck string

const (
	CurrentBranch  BranchCheck = "current"
	SpecificBranch BranchCheck = "specific"
)

type Response struct {
	Message string `json:"message"`
}

type Server struct {
	Port             string `json:"port"`
	LogRetentionDays int    `json:"log_retention_days"`
}

type TriggerEntry struct {
	Name         string
	AbsolutePath string
	SharedSecret string
	SyncBranch   BranchCheck // Default: current. [ current > , <branch_name>]
	BranchName   string      //
}

type Config struct {
	Server   Server         `json:"server"`
	Triggers []TriggerEntry `json:"triggers"`
}

type TriggerRequest struct {
	Name         string `json:"name"`
	SharedSecret string `json:"shared_secret"`
}

var Configuration Config = Config{}
