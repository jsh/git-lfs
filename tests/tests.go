package tests

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"testing"
)

var (
	rootDir string
	testDir string
	bin     string
)

type runner struct {
	dir      string
	repos    map[string]*repo
	repoName string
	*testing.T
}

// WriteFile writes the byte contents to the given path, which should be
// relative to this test's repository directory.
func (r *runner) WriteFile(path string, contents []byte) {
	err := ioutil.WriteFile(r.repoPath(path), contents, 0755)
	if err != nil {
		r.Fatal(err)
	}
}

// ReadFile reads the byte contents of the given path, which should be relative
// to this test's repository directory.
func (r *runner) ReadFile(path string) []byte {
	by, err := ioutil.ReadFile(r.repoPath(path))
	if err != nil {
		r.Fatal(err)
	}
	return by
}

func (r *runner) Cd(path string) {
	fullPath := r.repoPath(path)
	r.Logf("$ cd %s", fullPath)
	if err := os.Chdir(fullPath); err != nil {
		r.Fatal(err)
	}
}

func (r *runner) MkdirP(path string) {
	fullPath := r.repoPath(path)
	r.Logf("$ mkdir -p %s", fullPath)
	if err := os.MkdirAll(fullPath, 0755); err != nil {
		r.Fatal(err)
	}
}

// Git executes a Git command and returns the combined STDOUT and STDERR.
func (r *runner) Git(args ...string) string {
	return r.exec("git", args...)
}

// GitBlob gets the blob OID of the given path at the given commit.
func (r *runner) GitBlob(commitish, path string) string {
	out := r.Git("ls-tree", commitish)
	lines := strings.Split(out, "\n")
	for _, line := range lines {
		tabs := strings.Split(line, "\t")
		if len(tabs) < 2 {
			continue
		}

		attrs := strings.Split(tabs[0], " ")
		if len(attrs) < 3 {
			continue
		}

		if tabs[1] == path {
			return attrs[2]
		}
	}

	return ""
}

func (r *runner) exec(name string, args ...string) string {
	// replace standard "git lfs" commands with the compiled binary in ./bin
	if name == "git" && len(args) > 0 && args[0] == "lfs" {
		name = bin
		args = args[1:len(args)]
	}

	r.logCmd(name, args...)

	cmd := exec.Command(name, args...)
	out := &bytes.Buffer{}
	cmd.Stdout = out
	cmd.Stderr = out

	if err := cmd.Run(); err != nil {
		r.Fatalf("%s\n\n%s", err, out.String())
	}

	return out.String()
}

func (r *runner) logCmd(name string, args ...string) {
	loggedArgs := make([]string, len(args)+1)
	if strings.HasPrefix(name, "/") {
		rel, err := filepath.Rel(rootDir, name)
		if err == nil {
			loggedArgs[0] = rel
		} else {
			r.Errorf("Cannot make %q relative to %q", name, rootDir)
			loggedArgs[0] = name
		}
	} else {
		loggedArgs[0] = name
	}

	for idx, arg := range args {
		if strings.Contains(arg, " ") {
			arg = fmt.Sprintf(`"%s"`, arg)
		}
		loggedArgs[idx+1] = arg
	}

	r.Logf("$ %s", strings.Join(loggedArgs, " "))
}

func (r *runner) repoPath(path string) string {
	cleaned := filepath.Clean(path)
	repo := r.repo()
	if strings.HasPrefix(cleaned, "/") {
		r.Fatalf("%q is not relative to %q", path, repo.dir)
	}

	return filepath.Join(repo.dir, cleaned)
}

func Setup(t *testing.T) *runner {
	t.Parallel()
	t.Logf("working directory: %s", testDir)
	dir, err := ioutil.TempDir(testDir, "integration-test-")
	if err != nil {
		t.Fatal(err)
	}

	r := &runner{
		dir:   dir,
		repos: make(map[string]*repo),
		T:     t,
	}

	r.Logf("temp: %s", dir)
	r.InitRepo("repo")

	return r
}

func (r *runner) Teardown() {
	for _, repo := range r.repos {
		repo.Teardown()
	}

	if !r.T.Failed() {
		os.RemoveAll(r.dir)
	}
}

func init() {
	wd, err := os.Getwd()
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	rootDir = filepath.Join(wd, "..")
	testDir = filepath.Join(rootDir, "tmp")
	if err := os.MkdirAll(testDir, 0755); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	bin = filepath.Join(rootDir, "bin", "git-lfs")
	if _, err := os.Stat(bin); err != nil {
		fmt.Println("git-lfs is not compiled to " + bin)
		os.Exit(1)
	}
}