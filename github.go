package main

import (
	"encoding/json"
	"fmt"
	"io/fs"
	"path/filepath"
	"regexp"
	"strings"
)

// Create a GitHub issue with the title, body provided
func (r *RepoInfo) CreateIssue(title, body string) (*GitIssue, error) {
	issuesURL := fmt.Sprintf("https://api.github.com/repos/%s/%s/issues", r.Owner, r.Repo)
	payload := GitIssue{
		Title:    title,
		Body:     body,
		Assignee: r.Owner, // TODO: (#9:closed) not working; check token permissions include push access
	}
	var newIssue *GitIssue
	resp, err := Post(issuesURL, r.Token, &payload)
	if err != nil {
		return nil, err
	}
	err = json.Unmarshal(resp, &newIssue)
	if err != nil {
		return nil, err
	}
	return newIssue, nil
}

// Load all issues (open and closed) into r.Issues
func (r *RepoInfo) LoadIssues() error {
	issuesURL := fmt.Sprintf("https://api.github.com/repos/%s/%s/issues?state=all", r.Owner, r.Repo)
	resp, err := Get(issuesURL, r.Token)
	if err != nil {
		return err
	}
	err = json.Unmarshal(resp, &r.Issues)
	return err
}

/* Utils */
func CheckGitRepo(m *Manager) (*RepoInfo, error) {
	repo, err := checkGitfolder(m) // Check .git folder exists in directory
	if err != nil {
		return nil, err
	}
	repo.Token = m.Config.Gittoken
	return repo, nil
}

// Validates .git folder exists and parses Repo and Owner name
func checkGitfolder(m *Manager) (*RepoInfo, error) {
	dir := m.Config.Directory
	var gitfilepath string
	filepath.Walk(dir, func(path string, info fs.FileInfo, err error) error {
		if strings.HasSuffix(info.Name(), ".git") {
			gitfilepath = path
		}
		return nil
	})
	if gitfilepath == "" {
		return nil, fmt.Errorf("no .git folder found")
	}
	content, err := ReadFile(fmt.Sprintf("%s/config", gitfilepath))
	if err != nil {
		return nil, err
	}
	rgx := regexp.MustCompile(`(?i)git@github\.com:([^/]+)/([^\.]+)\.git`)
	matches := rgx.FindStringSubmatch(content)
	owner, repo := matches[1], matches[2]
	fmt.Printf("Repo:  %s\nOwner: %s\n", C(repo, Green), C(owner, Green))
	return &RepoInfo{Owner: owner, Repo: repo}, nil
}
