package hooks

import (
	"bufio"
	"bytes"
	"context"
	"fmt"
	"net"
	"os"
	"strings"
	"time"

	api "github.com/shelmangroup/monohub/api"
	log "github.com/sirupsen/logrus"
	"google.golang.org/grpc"
	"gopkg.in/alecthomas/kingpin.v2"
)

const (
	EmptySHA     = "0000000000000000000000000000000000000000"
	BranchPrefix = "refs/heads/"
)

var (
	command = kingpin.Command("hook", "Hook")

	pre        = command.Command("pre-receive", "pre receive hook")
	post       = command.Command("post-receive", "post receive hook.")
	grpcSocket = command.Flag("grpc-socket", "grpc hook server endpoint").Short('g').Required().String()
)

func grpcClient() api.GitHooksClient {
	conn, err := grpc.Dial(
		*grpcSocket,
		grpc.WithInsecure(),
		grpc.WithDialer(func(addr string, timeout time.Duration) (net.Conn, error) {
			return net.DialTimeout("unix", addr, timeout)
		}))
	if err != nil {
		log.Fatalf("did not connect: %v", err)
	}
	return api.NewGitHooksClient(conn)
}

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

		log.Debugf("branch: %s  oldCommit: %s  newCommit %s   refFullName %s\n", branchName, oldCommitID, newCommitID, refFullName)

	}

	req := api.PreReceiveRequest{}

	c := grpcClient()
	c.PreReceive(context.Background(), &req)
	return nil
}

func RunHookPostReceive() error {
	return nil
}

func fail(userMessage, logMessage string, args ...interface{}) {
	fmt.Fprintln(os.Stderr, "Monohub:", userMessage)

	if len(logMessage) > 0 {
		fmt.Fprintf(os.Stderr, logMessage+"\n", args...)
	}
	return
}
