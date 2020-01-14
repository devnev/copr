package gen

import (
	"github.com/devnev/copr/api"
	"gopkg.in/src-d/go-billy.v4/osfs"
	"gopkg.in/src-d/go-git.v4"
	"gopkg.in/src-d/go-git.v4/config"
	"gopkg.in/src-d/go-git.v4/plumbing"
	"gopkg.in/src-d/go-git.v4/storage/memory"
	"io/ioutil"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"unicode"
)

func Do(srcRoot, baseBranch, newBranch string, tracker api.Tracker) error {
	cmd := exec.Command(tracker.Command[0], tracker.Command[1:]...)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	cmd.Dir = srcRoot
	err := cmd.Run()
	if err != nil {
		return err
	}

	name := strings.Map(func(r rune) rune {
		if r > unicode.MaxASCII {
			return r
		}
		if unicode.IsControl(r) || unicode.IsPunct(r) {
			return '_'
		}
		return r
	}, tracker.Repository)
	checkoutDir, err := ioutil.TempDir(".", name + "_")
	if err != nil {
		return err
	}

	store := memory.NewStorage()
	repo, err := git.Clone(store, osfs.New(checkoutDir), &git.CloneOptions{
		URL: tracker.Repository,
		ReferenceName: plumbing.NewBranchReferenceName(baseBranch),
		SingleBranch: true,
		NoCheckout: true,
		Depth: 1,
	})
	if err != nil {
		return err
	}
	worktree, err := repo.Worktree()
	if err != nil {
		return err
	}

	outputDir := filepath.Join(srcRoot, filepath.FromSlash(tracker.Output.Directory))
	err = filepath.Walk(outputDir, func(path string, info os.FileInfo, err error) error {
		if (info.Mode() & os.ModeType) != 0 {
			return nil
		}
		rel, err := filepath.Rel(outputDir, path)
		if err != nil {
			return err
		}
		dst := filepath.Join(worktree.Filesystem.Root(), rel)
		return copyFile(path, dst)
	})
	if err != nil {
		return err
	}

	_, err = worktree.Add(".")
	status, err := worktree.Status()
	if status.IsClean() {
		return nil
	}
	_, err = worktree.Commit("Test", &git.CommitOptions{})
	err = repo.Push(&git.PushOptions{
		RemoteName: git.DefaultRemoteName,
		RefSpecs: []config.RefSpec{
			config.RefSpec(plumbing.NewBranchReferenceName("master") + ":" + plumbing.NewBranchReferenceName(branch)),
		},
	})
	return nil
}

func copyFile(src, dst string) error {
	b, err := ioutil.ReadFile(src)
	if err != nil {
		return err
	}
	return ioutil.WriteFile(dst, b, 0644)
}
