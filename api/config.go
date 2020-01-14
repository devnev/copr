package api

import (
	"gopkg.in/yaml.v3"
	"io"
)

const (
	Filename = ".copr"
)

type Config struct {
	Trackers []Tracker
}

type Tracker struct {
	Repository string
	Command []string
	Output OutputFormat
}

type OutputFormat struct {
	Directory string
}

func ReadConfig(r io.Reader) (config Config, err error) {
	dec := yaml.NewDecoder(r)
	dec.KnownFields(true)
	err = dec.Decode(&config)
	return config, err
}
