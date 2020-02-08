package gen

import (
	"flag"
	"fmt"
	"os"
	"os/exec"
	"reflect"
	"testing"

	"github.com/devnev/copr/config"
)

func handleCommands(t *testing.T, cb func([]string)) {
	const testCommandEnv = "GO_TEST_COMMAND"
	if env := os.Getenv(testCommandEnv); env == "" {
		// if there's no intended command, we're running in the test, setup command wrapper
		command = func(name string, arg ...string) *exec.Cmd {
			args := []string{"-test.run=" + t.Name(), "--", name}
			cmd := exec.Command(os.Args[0], append(args, arg...)...)
			cmd.Env = []string{testCommandEnv + "=" + t.Name()}
			return cmd
		}
	} else if env != t.Name() {
		// if we're running in a command test but it's not for our current test, something went wrong
		t.Fatalf("Test %s executed during %s test command execution", t.Name(), env)
	} else {
		cb(flag.Args())
	}
}

func TestFailsIfDefaultBranchCmdFails(t *testing.T) {
	handleCommands(t, func(args []string) {
		if reflect.DeepEqual(args, []string{"git", "rev-parse", "--abbrev-ref", "HEAD"}) {
			os.Exit(1)
		}
		fmt.Fprintf(os.Stderr, "got %q\n", args)
		os.Exit(2)
	})
	err := Do(".", config.Output{})
	if err == nil || err.Error() != "branch command failed with exit status 1" {
		t.Fatalf("expected branch command error, got %q", err)
	}
}
