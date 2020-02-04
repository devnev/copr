package config

type Repository struct {
	Outputs []Output
}

type Output struct {
	Repository string
	Command    []string
	Directory string
}
