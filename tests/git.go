package tests

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"net/url"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
)

// SetRepo changes the current test to the repository's working directory.  The
// given name refers to the subdirectory in the runner temp directory that the
// repository lives.
func (r *runner) SetRepo(name string) {
	r.Logf("$ cd %s", name)
	r.repoName = name
	if err := os.Chdir(r.repo().dir); err != nil {
		r.Fatal(err)
	}
}

// InitRepo initializes an empty git repository.  The given name is a
// subdirectory in the runner temp directory that the repository will live.
func (r *runner) InitRepo(name string) {
	dir := filepath.Join(r.dir, name)
	if err := os.MkdirAll(dir, 0777); err != nil {
		r.Fatal(err)
	}

	serverDir := dir + ".server"
	if err := os.MkdirAll(serverDir, 0777); err != nil {
		r.Fatal(err)
	}

	repo := &repo{
		dir:          dir,
		largeObjects: make(map[string][]byte),
		serverDir:    serverDir,
		configFile:   filepath.Join(r.dir, name+".gitconfig"),
	}
	repo.server = httptest.NewServer(r.httpHandler(repo))
	repo.serverURL = repo.server.URL + "/" + name + ".server"
	r.repos[name] = repo

	// write a config file for any clones of this repo
	cfg := fmt.Sprintf(`[filter "lfs"]
	required = true
	smudge = %s smudge %%f
	clean = %s clean %%f

[remote "origin"]
	url = %s
	fetch = +refs/heads/*:refs/remotes/origin/*
`, bin, bin, repo.serverURL)
	if err := ioutil.WriteFile(repo.configFile, []byte(cfg), 0755); err != nil {
		panic(err)
	}

	// set up the git server
	if err := os.Chdir(serverDir); err != nil {
		r.Fatal(err)
	}
	r.Git("init")
	r.Git("config", "http.receivepack", "true")
	r.Git("config", "receive.denyCurrentBranch", "ignore")
	r.Logf("git init server: %s/%s", repo.server.URL, name)

	// set up the local git clone
	r.SetRepo(name)
	r.Git("init")
	r.Git("remote", "add", "origin", repo.server.URL+"/"+name+".server")
	r.Git("config", "filter.lfs.smudge", fmt.Sprintf("%s smudge %%f", bin))
	r.Git("config", "filter.lfs.clean", fmt.Sprintf("%s clean %%f", bin))
	r.setupCredentials(repo.server.URL)
	r.WriteFile(".git/hooks/pre-push", []byte("#!/bin/sh\n"+bin+` pre-push "$@"`+"\n"))
	r.Logf("git init: %s", dir)
}

func (r *runner) CloneTo(name string) string {
	repository := r.repo()
	if err := os.Chdir(r.dir); err != nil {
		r.Fatal(err)
	}

	cmdName := "git"
	args := []string{"clone", repository.serverURL, name}
	r.logCmd(cmdName, args...)

	cmd := exec.Command(cmdName, args...)
	currEnv := os.Environ()
	cmd.Env = make([]string, len(currEnv)+1)
	for idx, e := range currEnv {
		cmd.Env[idx] = e
	}

	// ensures that filter.lfs.* config settings point to the test's git-lfs.
	cmd.Env[len(cmd.Env)-1] = "GIT_CONFIG=" + repository.configFile
	out := &bytes.Buffer{}
	cmd.Stdout = out
	cmd.Stderr = out

	if err := cmd.Run(); err != nil {
		r.Fatalf("%s\n\n%s", err, out.String())
	}

	r.repos[name] = &repo{
		runner:       r,
		dir:          filepath.Join(r.dir, name),
		largeObjects: repository.largeObjects,
		configFile:   repository.configFile,
		server:       repository.server,
		serverURL:    repository.serverURL,
	}

	r.SetRepo(name)
	r.Git("config", "filter.lfs.smudge", fmt.Sprintf("%s smudge %%f", bin))
	r.Git("config", "filter.lfs.clean", fmt.Sprintf("%s clean %%f", bin))
	return out.String()
}

func (r *runner) repo() *repo {
	repo := r.repos[r.repoName]
	if repo == nil {
		r.Fatalf("no repo found for %q", r.repoName)
	}
	return repo
}

func (r *runner) setupCredentials(rawurl string) {
	u, err := url.Parse(rawurl)
	if err != nil {
		r.Fatal(err)
	}

	input := fmt.Sprintf("protocol=http\nhost=%s\nusername=a\npassword=b", u.Host)
	out := &bytes.Buffer{}
	cmd := exec.Command("git", "credential", "approve")
	cmd.Stdin = strings.NewReader(input)
	cmd.Stdout = out
	cmd.Stderr = out
	if err := cmd.Run(); err != nil {
		r.Fatalf("%s\n\n%s", err, out.String())
	}
}

func (run *runner) httpHandler(repository *repo) http.Handler {
	mux := http.NewServeMux()
	mux.HandleFunc("/storage/", run.storageHandler(repository))
	mux.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		if strings.Contains(r.URL.Path, ".server.git/info/lfs") {
			run.Logf("git lfs %s %s", r.Method, r.URL)
			run.Logf("git lfs Accept: %s", r.Header.Get("Accept"))
			run.lfsHandler(repository, w, r)
			return
		}

		run.Logf("git http-backend %s %s", r.Method, r.URL)
		run.gitHandler(w, r)
	})

	return mux
}

type repo struct {
	runner       *runner
	dir          string
	largeObjects map[string][]byte
	configFile   string
	server       *httptest.Server
	serverDir    string
	serverURL    string
}

func (r *repo) Teardown() {
	r.server.Close()

	u, err := url.Parse(r.server.URL)
	if err != nil {
		r.runner.Fatal(err)
	}

	input := fmt.Sprintf("protocol=http\nhost=%s", u.Host)
	out := &bytes.Buffer{}
	cmd := exec.Command("git", "credential", "reject")
	cmd.Stdin = strings.NewReader(input)
	cmd.Stdout = out
	cmd.Stderr = out
	if err := cmd.Run(); err != nil {
		r.runner.Fatalf("%s\n\n%s", err, out.String())
	}
}