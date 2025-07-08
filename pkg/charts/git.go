// SPDX-FileCopyrightText: 2025 SAP SE or an SAP affiliate company
// SPDX-License-Identifier: Apache-2.0

package charts

import (
	"bytes"
	"errors"
	"os"
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
	remote    string
	directory string
}

func newGit(directory, remote string) (*git, error) {
	g := &git{
		directory: directory,
		remote:    remote,
	}

	err := g.testGitInstalled()
	if err != nil {
		return nil, err
	}

	err = g.testGitRepository()
	return g, err
}

func (g *git) testGitInstalled() error {
	if _, err := g.runGitCmd("--version"); err != nil {
		return errGitNotInstalled
	}
	return nil
}

func (g *git) testGitRepository() error {
	stdout, err := g.runGitCmd("rev-parse", "--is-inside-work-tree")
	if err != nil {
		return err
	}

	b, err := strconv.ParseBool(stdout)
	if err != nil {
		return err
	}

	if !b {
		return errNoGitRepository
	}

	return nil
}

func (g *git) fetch() error {
	stdout, err := g.runGitCmd("remote", "get-url", g.remote)
	if err != nil || stdout == "" {
		return errNoRemote
	}

	_, err = g.runGitCmd("fetch")
	return err
}

func (g *git) getChangedDirs(remote, commit string) ([]string, error) {
	stdOut, err := g.runGitCmd("diff", "--find-renames", "--name-only", remote, commit, "--", g.directory)
	if err != nil {
		return nil, err
	}

	var changedDirs []string
	lines := strings.Split(stdOut, "\n")
	for _, l := range lines {
		if p := l; p != "" {
			changedDirs = append(changedDirs, g.pathWithDirectory(p))
		}
	}

	return changedDirs, nil
}

func (g *git) getCommitHash(commit string) (string, error) {
	stdOut, err := g.runGitCmd("rev-parse", commit)
	return stdOut, err
}

func (g *git) getMergeBase(commit1, commit2 string) (string, error) {
	stdOut, err := g.runGitCmd("merge-base", commit1, commit2)
	if err != nil {
		return "", err
	}

	if stdOut == commit2 {
		stdOut = "HEAD"
	}

	return stdOut, err
}

func (g *git) runGitCmd(args ...string) (stdOutString string, err error) {
	var stdout bytes.Buffer

	cmd := exec.Command("git", append([]string{"-C", g.directory}, args...)...) //nolint:gosec // all arguments are used supplied
	cmd.Stdout = &stdout
	cmd.Stderr = os.Stderr
	err = cmd.Run()

	stdOutString = string(bytes.TrimSpace(stdout.Bytes()))

	return stdOutString, err
}

func (g *git) pathWithDirectory(path string) string {
	// Avoid duplicating folder name when joining path.
	base := filepath.Base(g.directory)
	path = strings.TrimPrefix(path, base)
	return filepath.Join(g.directory, path)
}
