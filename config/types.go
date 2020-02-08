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

func (r *Repository) SetDefaults() {
	for _, o := range r.Outputs {
		o.SetDefaults()
	}
}

func (o *Output) SetDefaults() {
	o.Base = "master"
	o.Branch = []string{"git", "rev-parse", "--abbrev-ref", "HEAD"}
}
