package gen

import (
	"io/ioutil"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"unicode"

	"gopkg.in/src-d/go-billy.v4/osfs"
	"gopkg.in/src-d/go-git.v4"
	gitconfig "gopkg.in/src-d/go-git.v4/config"
	"gopkg.in/src-d/go-git.v4/plumbing"
	"gopkg.in/src-d/go-git.v4/storage/memory"

	"github.com/devnev/copr/config"
)


func Do(srcRoot string, output config.Output) error {
	newBranch, err := readCmd(srcRoot, output.Branch[0], output.Branch[1:]...)
	if err != nil {
		return err
	}

	err = runCmd(srcRoot, output.Generate[0], output.Generate[1:]...)
	if err != nil {
		return err
	}

	_, repo, worktree, err := makeClone(output.Repository, output.Base)
	if err != nil {
		return err
	}

	err = copyOutput(filepath.Join(srcRoot, filepath.FromSlash(output.Directory)), worktree.Filesystem.Root())
	if err != nil {
		return err
	}

	push, err := commit(worktree)
	if err != nil || !push {
		return err
	}

	return repo.Push(&git.PushOptions{
		RemoteName: git.DefaultRemoteName,
		RefSpecs: []gitconfig.RefSpec{
			gitconfig.RefSpec("+" + plumbing.NewBranchReferenceName("master") + ":" + plumbing.NewBranchReferenceName(newBranch)),
		},
	})
}

func readCmd(dir string, name string, args ...string) (string, error) {
	cmd := exec.Command(name, args...)
	cmd.Stderr = os.Stderr
	cmd.Dir = dir
	out, err := cmd.Output()
	return strings.TrimSpace(string(out)), err
}

func runCmd(srcRoot string, name string, args ...string) error {
	cmd := exec.Command(name, args...)
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
