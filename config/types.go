package config

const (
	Filename = ".copr"
)

type Repository struct {
	Trackers []Tracker
}

type Tracker struct {
	Repository string
	Command    []string
	Output     OutputFormat
}

type OutputFormat struct {
	Directory string
}
