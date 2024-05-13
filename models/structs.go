package models

type BranchCheck string

const (
	CurrentBranch  BranchCheck = "current"
	SpecificBranch BranchCheck = "specific"
)

type Response struct {
	Message string `json:"message"`
}

type ServerConfig struct {
	Port string
}

type ProjectEntry struct {
	Name         string
	AbsolutePath string
	SharedSecret string
	SyncBranch   BranchCheck // Default: current. [ current > , <branch_name>]
	BranchName   string      //
}

type TriggerObject struct {
	ProjectName  string `json:"project_name"`
	SharedSecret string `json:"shared_secret"`
}
