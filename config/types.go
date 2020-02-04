package config

type Repository struct {
	Outputs []Output
}

type Output struct {
	Repository string
	Base       string
	Branch     []string
	Generate   []string
	Directory  string
}
