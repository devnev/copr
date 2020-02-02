package main

import (
	"fmt"
	"os"

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

	repoConf, err := config.ReadPath(path)
	if err != nil {
		return err
	}

	fmt.Println(repoConf)
	return err
}
