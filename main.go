package main

import (
	"fmt"
	"github.com/devnev/copr/gen"
	"log"
	"os"
	"path/filepath"

	"github.com/devnev/copr/config"
)

func main() {
	err := run()
	if err != nil {
		fmt.Println(err.Error())
		os.Exit(1)
	}
}

func run() error {
	startDir, err := os.Getwd()
	if err != nil {
		return err
	}

	path, err := config.FindPath(startDir)
	if err != nil {
		return err
	}
	log.Printf("Using config at %q", path)

	repoConf, err := config.ReadPath(path)
	if err != nil {
		return err
	}

	log.Printf("Processing with config %+v", repoConf)

	for _, out := range repoConf.Outputs {
		err = gen.Do(filepath.Dir(path), out)
		if err != nil {
			return err
		}
	}

	return err
}
