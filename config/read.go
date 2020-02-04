package config

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"

	"gopkg.in/yaml.v3"
)

var Names []string = names()

func FindPath(startDir string) (string, error) {
	dir, err := filepath.Abs(startDir)
	if err != nil {
		return "", err
	}
	for {
		var found []string
		for _, name := range Names {
			path := filepath.Join(dir, name)
			if fi, err := os.Stat(path); err == nil && !fi.IsDir() {
				found = append(found, path)
			} else if err != nil && !os.IsNotExist(err) {
				return "", err
			}
		}
		if len(found) == 1 {
			return found[0], nil
		}
		if len(found) > 1 {
			return "", fmt.Errorf("multiple candidates: %q", found)
		}
		if newDir := filepath.Dir(dir); newDir == dir {
			return "", fmt.Errorf("no config found")
		} else {
			dir = newDir
		}
	}
}

func ReadPath(path string) (conf Repository, err error) {
	buf, err := ioutil.ReadFile(path)
	if err != nil {
		return conf, err
	}

	switch {
	case strings.HasSuffix(path, ".json"):
		conf, err = ReadJSON(bytes.NewReader(buf))
	case strings.HasSuffix(path, ".yml"):
		conf, err = ReadYAML(bytes.NewReader(buf))
	case strings.HasSuffix(path, ".yaml"):
		conf, err = ReadYAML(bytes.NewReader(buf))
	default:
		conf, err = ReadDetected(bytes.NewReader(buf))
	}
	return conf, err
}

func ReadYAML(r io.Reader) (config Repository, err error) {
	dec := yaml.NewDecoder(r)
	dec.KnownFields(true)
	return config, dec.Decode(&config)
}

func ReadJSON(r io.Reader) (config Repository, err error) {
	dec := json.NewDecoder(r)
	dec.DisallowUnknownFields()
	return config, dec.Decode(&config)
}

func ReadDetected(r io.ReadSeeker) (config Repository, err error) {
	firstByte := make([]byte, 1)
	n, err := r.Read(firstByte)
	if err != nil {
		return config, err
	} else if n == 0 {
		return config, fmt.Errorf("unexpected empty read")
	}
	r.Seek(-1, io.SeekCurrent)
	if firstByte[0] == '{' {
		return ReadJSON(r)
	} else {
		return ReadYAML(r)
	}
}

func names() []string {
	var names []string
	parts := [][]string{{"", "."}, {"copr"}, {"", ".yml", ".yaml", ".json"}}
	var f func([]string)
	f = func(iter []string) {
		if len(iter) == len(parts) {
			names = append(names, strings.Join(iter, ""))
			return
		}
		for _, part := range parts[len(iter)] {
			f(append(iter, part))
		}
	}
	f(make([]string, 0, len(parts)))
	return names
}
