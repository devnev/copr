package gen

import (
	"io/ioutil"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"unicode"

	"github.com/devnev/copr/api"
	"gopkg.in/src-d/go-billy.v4/osfs"
	"gopkg.in/src-d/go-git.v4"
	"gopkg.in/src-d/go-git.v4/config"
	"gopkg.in/src-d/go-git.v4/plumbing"
	"gopkg.in/src-d/go-git.v4/storage/memory"
)

func Do(srcRoot, baseBranch, newBranch string, tracker api.Tracker) error {
	err := runCmd(srcRoot, tracker)
	if err != nil {
		return err
	}

	_, repo, worktree, err := makeClone(tracker.Repository, baseBranch)
	if err != nil {
		return err
	}

	err = copyOutput(filepath.Join(srcRoot, filepath.FromSlash(tracker.Output.Directory)), worktree.Filesystem.Root())
	if err != nil {
		return err
	}

	push, err := commit(worktree)
	if err != nil || !push {
		return err
	}

	return repo.Push(&git.PushOptions{
		RemoteName: git.DefaultRemoteName,
		RefSpecs: []config.RefSpec{
			config.RefSpec("+" + plumbing.NewBranchReferenceName("master") + ":" + plumbing.NewBranchReferenceName(newBranch)),
		},
	})
}

func runCmd(srcRoot string, tracker api.Tracker) error {
	cmd := exec.Command(tracker.Command[0], tracker.Command[1:]...)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	cmd.Dir = srcRoot
	return cmd.Run()
}

func makeClone(repositoryURL string, branch string) (store *memory.Storage, repo *git.Repository, worktree *git.Worktree, err error) {
	name := strings.Map(func(r rune) rune {
		if r > unicode.MaxASCII {
			return r
		}
		if unicode.IsControl(r) || unicode.IsPunct(r) {
			return '_'
		}
		return r
	}, repositoryURL)
	checkoutDir, err := ioutil.TempDir(".", name+"_")
	if err != nil {
		return
	}

	store = memory.NewStorage()
	repo, err = git.Clone(store, osfs.New(checkoutDir), &git.CloneOptions{
		URL:           repositoryURL,
		ReferenceName: plumbing.NewBranchReferenceName(branch),
		SingleBranch:  true,
		NoCheckout:    true,
		Depth:         1,
	})
	if err == nil {
		worktree, err = repo.Worktree()
	}

	return
}

func copyOutput(outputDir string, dstRoot string) error {
	return filepath.Walk(outputDir, func(path string, info os.FileInfo, err error) error {
		if (info.Mode() & os.ModeType) != 0 {
			return nil
		}
		rel, err := filepath.Rel(outputDir, path)
		if err != nil {
			return err
		}
		dst := filepath.Join(dstRoot, rel)
		return copyFile(path, dst)
	})
}

func copyFile(src, dst string) error {
	b, err := ioutil.ReadFile(src)
	if err != nil {
		return err
	}
	return ioutil.WriteFile(dst, b, 0644)
}

func commit(worktree *git.Worktree) (bool, error) {
	_, err := worktree.Add(".")
	if err != nil {
		return false, err
	}
	status, err := worktree.Status()
	if err != nil {
		return false, err
	}
	if status.IsClean() {
		return false, nil
	}
	_, err = worktree.Commit("Test", &git.CommitOptions{})
	if err != nil {
		return false, err
	}
	return true, nil
}
