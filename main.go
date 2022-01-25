package main

import (
	"flag"
	"fmt"
	"log"
	"os"
	"os/exec"
	"regexp"
	"strings"

	git "github.com/go-git/go-git/v5"
)

var refRegex = regexp.MustCompile(`([a-zA-Z0-9]{1,5}-\d{1,})`)

const WS_ENV = "LINEAR_WORKSPACE"

func main() {
	var (
		issue             string
		workspace         string
		issueFromRepo     string
		workspaceFromRepo string
	)

	workspace = *flag.String("workspace", "", "The linear workspace name")
	issue = strings.ToUpper(*flag.String("issue", "", "The name of the linear issue"))
	repoFlag := flag.String("repo", "",
		"The absolute path to a git repo. A linear issue name must be the "+
			"current working branch",
	)
	flag.Parse()

	var repo *git.Repository
	if *repoFlag != "" {
		var err error
		repo, err = getRepo(*repoFlag)
		if err != nil {
			log.Fatal(err)
		}
	} else {
		dir, err := os.Getwd()
		if err != nil {
			log.Fatal(err)
		}

		repo, err = getRepo(dir)
		if err != nil {
			log.Fatal(err)
		}
	}

	if repo != nil {
		var err error
		workspaceFromRepo, err = getWorkspaceFromRepo(repo)
		if err != nil {
			log.Fatal(err)
		}

		issueFromRepo, err = getIssueFromRepo(repo)
		if err != nil {
			log.Fatal(err)
		}
	}

	var err error
	if workspace == "" {
		if workspaceFromRepo != "" {
			workspace = workspaceFromRepo
		} else if ws := os.Getenv(WS_ENV); ws != "" {
			workspace = ws
		} else {
			log.Fatal("no workspace name provided")
		}
	}

	if issue == "" && issueFromRepo != "" {
		issue = issueFromRepo
	}

	if issue == "" {
		log.Fatal("no issue name provided")
	}

	cmd := []*exec.Cmd{
		exec.Command("open", fmt.Sprintf("linear://%s/issue/%s", workspace, issue)),
		exec.Command("open", fmt.Sprintf("https://linear.app/%s/issue/%s", workspace, issue)),
	}

	for _, c := range cmd {
		err = c.Run()
		if err == nil {
			break
		}
	}

	if err != nil {
		log.Fatal(err)
	}
}

func getRepo(dir string) (*git.Repository, error) {
	return git.PlainOpen(dir)
}

func getIssueFromRepo(repo *git.Repository) (string, error) {
	h, err := repo.Head()
	if err != nil {
		return "", err
	}

	return strings.ToUpper(refRegex.FindString(h.Name().String())), nil
}

func getWorkspaceFromRepo(repo *git.Repository) (string, error) {
	c, err := repo.Config()
	if err != nil {
		return "", err
	}

	var ws string
	for _, section := range c.Raw.Sections {
		if section.Name != "linear" {
			continue
		}

		for _, opt := range section.Options {
			if opt.Key != "workspace" {
				continue
			}

			ws = opt.Value
		}
	}

	return ws, nil
}
