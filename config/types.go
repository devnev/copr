package config

type Repository struct {
	Outputs []Output
}

type Output struct {
	Repository string
	BaseBranch string
	BranchCmd  []string
	Command    []string
	Directory  string
}
