package gitreport

import (
	"encoding/json"
	"errors"
	"fmt"
	"os/exec"
	"path/filepath"
	"strings"
	"time"
)

var (
	JSONOutFormat             = `{%n  "commit": "%H",%n  "refs": "%D",%n  "subject": "%s",%n  "body": "%b",%n  "author": {%n    "name": "%aN",%n    "email": "%aE",%n    "date": "%aI"%n  },%n  "commiter": {%n    "name": "%cN",%n    "email": "%cE",%n    "date": "%cI"%n  }%n}`
	JSONMinOutFormat          = `{"commit": "%H","refs": "%D","subject": "%s","body": "%b","author": {"name": "%aN","email": "%aE","date": "%aI"},"commiter": {"name": "%cN","email": "%cE","date": "%cI"}}`
	JSONOutFormatWithComma    = fmt.Sprintf("%s,", JSONOutFormat)
	JSONMinOutFormatWithComma = fmt.Sprintf("%s,", JSONMinOutFormat)
)

type LogOption struct {
	Authors           []string
	Since             *time.Time
	Until             *time.Time
	FetchAllBranch    bool
	Limit             int
	FilterMergeCommit bool
}

type User struct {
	Name  string    `json:"name"`
	Email string    `json:"email"`
	Date  time.Time `json:"date"`
}

type GitCommit struct {
	Hash      string `json:"commit"`
	Refs      string `json:"refs"`
	Subject   string `json:"subject"`
	Body      string `json:"body"`
	Author    *User  `json:"author"`
	Committer *User  `json:"committer"`
}

func (c *GitCommit) Message() string {
	return fmt.Sprintf("%s %s", c.Subject, c.Body)
}

type Result struct {
	Repo    string
	Commits []*GitCommit
}

type GitClientOptions struct {
	Repo string
}

type IGitClient interface {
	Log(*LogOption) (*Result, error)
}

type nativeGitWrapper struct {
	repo string
}

func NewGitClient(repoPath string) IGitClient {
	return &nativeGitWrapper{
		repo: repoPath,
	}
}

func (c *nativeGitWrapper) Log(options *LogOption) (*Result, error) {
	gitPath, err := ResolvePath(c.repo)
	if err != nil {
		return nil, err
	}

	gitOptions := []string{"log"}
	gitOptions = append(gitOptions, convertLogOptions(options)...)
	cmd := exec.Command("git", gitOptions...)
	cmd.Dir = gitPath

	out, err := cmd.Output()
	if err != nil {
		return nil, err
	}

	jsonStr := string(out)
	jsonStr = strings.TrimSpace(jsonStr)
	if jsonStr == "" {
		return &Result{
			Repo:    filepath.Base(c.repo),
			Commits: []*GitCommit{},
		}, nil
	}
	if jsonStr[len(jsonStr)-1] != ',' {
		return nil, errors.New("Invalid return from git log")
	}
	jsonStr = jsonStr[:len(jsonStr)-1]
	jsonStr = fmt.Sprintf(`[%s]`, jsonStr)

	var commits []*GitCommit
	err = json.Unmarshal([]byte(jsonStr), &commits)
	if err != nil {
		return nil, err
	}

	return &Result{
		Repo:    filepath.Base(c.repo),
		Commits: commits,
	}, nil
}

func convertLogOptions(option *LogOption) []string {
	outOptions := []string{
		fmt.Sprintf("--pretty=format:%s", JSONMinOutFormatWithComma),
	}
	if option.FetchAllBranch {
		outOptions = append(outOptions, "--all")
	}
	if option.FilterMergeCommit {
		outOptions = append(outOptions, "--no-merges")
	}
	if option.Limit > 0 {
		outOptions = append(outOptions, fmt.Sprintf("-n %d", option.Limit))
	}
	for _, author := range option.Authors {
		outOptions = append(outOptions, fmt.Sprintf(`--author=%s`, author))
	}
	if option.Since != nil {
		isoFormat := option.Since.Format(time.RFC3339)
		outOptions = append(outOptions, fmt.Sprintf(`--since="%s"`, isoFormat))
	}
	if option.Until != nil {
		isoFormat := option.Until.Format(time.RFC3339)
		outOptions = append(outOptions, fmt.Sprintf(`--until="%s"`, isoFormat))
	}
	return outOptions
}
