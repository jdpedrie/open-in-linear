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

func main() {
	var issue string

	wsName := flag.String("workspace", "", "The linear workspace name")
	issueFlag := flag.String("issue", "", "The name of the linear issue")
	repoFlag := flag.String("repo", "",
		"The absolute path to a git repo. A linear issue name must be the "+
			"current working branch",
	)
	flag.Parse()

	if *wsName == "" {
		log.Fatal("no workspace name specified")
	}

	if *issueFlag != "" {
		issue = strings.ToUpper(*issueFlag)
	} else if *repoFlag != "" {
		var err error
		issue, err = getIssueFromRepo(*repoFlag)
		if err != nil {
			log.Fatal(err)
		}
	} else {
		dir, err := os.Getwd()
		if err != nil {
			log.Fatal(err)
		}

		issue, err = getIssueFromRepo(dir)
		if err != nil {
			log.Fatal(err)
		}
	}

	cmd := []*exec.Cmd{
		exec.Command("open", fmt.Sprintf("linear://%s/issue/%s", *wsName, issue)),
		exec.Command("open", fmt.Sprintf("https://linear.app/%s/issue/%s", *wsName, issue)),
	}

	var err error
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

func getIssueFromRepo(dir string) (string, error) {
	repo, err := git.PlainOpen(dir)
	if err != nil {
		return "", err
	}

	h, err := repo.Head()
	if err != nil {
		return "", err
	}

	log.Println(h.Name())

	m := refRegex.FindString(h.Name().String())
	if m == "" {
		return "", fmt.Errorf(
			"current working branch %s does not contain a linear issue name",
			h.Name().String(),
		)
	}

	return strings.ToUpper(m), nil
}
