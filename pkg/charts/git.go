package charts

import (
	"bytes"
	"errors"
	"os/exec"
	"path/filepath"
	"strconv"
	"strings"
)

var (
	errGitNotInstalled = errors.New("git is not installed")
	errNoGitRepository = errors.New("folder is not a git repository")
	errNoRemote        = errors.New("no remote configured in git repository")
)

type git struct {
	remote string

	cmd       string
	directory string
	baseArgs  []string
}

func newGit(directory, remote string) (*git, error) {
	g := &git{
		cmd:       "git",
		directory: directory,
		remote:    remote,
	}

	if err := g.testGitInstalled(); err != nil {
		return nil, err
	}

	return g, g.testGitRepository()
}

func (g *git) testGitInstalled() error {
	if _, _, err := g.runGitCmd("--version"); err != nil {
		return errGitNotInstalled
	}
	return nil
}

func (g *git) testGitRepository() error {
	stdout, _, err := g.runGitCmd("rev-parse", "--is-inside-work-tree")
	if err != nil {
		return errNoGitRepository
	}

	if b, err := strconv.ParseBool(stdout); err == nil && b == true {
		return nil
	}

	return errNoGitRepository
}

func (g *git) fetch() error {
	stdout, _, err := g.runGitCmd("remote", "get-url", g.remote)
	if err != nil || stdout == "" {
		return errNoRemote
	}

	_, _, err = g.runGitCmd("fetch")
	return err
}

func (g *git) getChangedDirs(remote, commit string) ([]string, error) {
	stdOut, _, err := g.runGitCmd("diff", "--find-renames", "--name-only", remote, commit, "--", g.directory)
	if err != nil {
		return nil, err
	}

	var changedDirs []string
	lines := strings.Split(stdOut, "\n")
	for _, l := range lines {
		if p := string(l); p != "" {
			changedDirs = append(changedDirs, g.pathWithDirectory(p))
		}
	}

	return changedDirs, nil
}

func (g *git) getCommitHash(commit string) (string, error) {
	stdOut, _, err := g.runGitCmd("rev-parse", commit)
	return stdOut, err
}

func (g *git) getMergeBase(commit1, commit2 string) (string, error) {
	stdOut, _, err := g.runGitCmd("merge-base", commit1, commit2)
	if err != nil {
		return "", err
	}

	if stdOut == commit2 {
		stdOut = "HEAD"
	}

	return stdOut, err
}

func (g *git) runGitCmd(args ...string) (string, string, error) {
	var (
		stdout,
		stderr bytes.Buffer
	)
	cmd := exec.Command(g.cmd, append([]string{"-C", g.directory}, args...)...)
	cmd.Stdout = &stdout
	cmd.Stderr = &stderr
	err := cmd.Run()

	stdOutString := string(bytes.TrimSpace(stdout.Bytes()))
	stdErrString := string(bytes.TrimSpace(stderr.Bytes()))

	return stdOutString, stdErrString, err
}

func (g *git) pathWithDirectory(path string) string {
	// Avoid duplicating folder name when joining path.
	base := filepath.Base(g.directory)
	if strings.HasPrefix(path, base) {
		path = strings.TrimPrefix(path, base)
	}
	return filepath.Join(g.directory, path)
}
