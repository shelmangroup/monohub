package hooks

import (
	"bufio"
	"bytes"
	"fmt"
	"os"
	"os/exec"
	"strings"

	// log "github.com/sirupsen/logrus"
	"gopkg.in/alecthomas/kingpin.v2"
)

const (
	EmptySHA     = "0000000000000000000000000000000000000000"
	BranchPrefix = "refs/heads/"
)

var (
	command = kingpin.Command("hook", "Hook")

	pre      = command.Command("pre-receive", "pre receive hook")
	post     = command.Command("post-receive", "post receive hook.")
	repoPath = command.Flag("repo-path", "repo directory").Short('d').Required().String()
)

func PreFullCommand() string {
	return pre.FullCommand()
}

func PostFullCommand() string {
	return post.FullCommand()
}

func RunHookPreReceive() error {
	buf := bytes.NewBuffer(nil)
	scanner := bufio.NewScanner(os.Stdin)
	for scanner.Scan() {
		buf.Write(scanner.Bytes())
		buf.WriteByte('\n')

		fields := bytes.Fields(scanner.Bytes())
		if len(fields) != 3 {
			continue
		}

		oldCommitID := string(fields[0])
		newCommitID := string(fields[1])
		refFullName := string(fields[2])

		branchName := strings.TrimPrefix(refFullName, BranchPrefix)

		if branchName == "master" {
			// detect force push
			if EmptySHA != oldCommitID {
				cmd := exec.Command("git", "rev-list", "--max-count=1", oldCommitID, "^"+newCommitID)
				cmd.Dir = *repoPath
				output, err := cmd.Output()
				if err != nil {
					fail("Internal error", "Fail to detect force push: %v", err)
				} else if len(output) > 0 {
					fail(fmt.Sprintf("branch %s is protected from force push", branchName), "")
				}
			}
			// check and deletion
			if newCommitID == EmptySHA {
				fail(fmt.Sprintf("branch %s is protected from deletion", branchName), "")
			}
		}
	}
	return nil
}

func RunHookPostReceive() error {
	return nil
}

func fail(userMessage, logMessage string, args ...interface{}) {
	fmt.Fprintln(os.Stderr, "Monohub:", userMessage)

	if len(logMessage) > 0 {
		fmt.Fprintf(os.Stderr, logMessage+"\n", args...)
		return
	}
	os.Exit(1)
}
