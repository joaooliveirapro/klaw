package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"io/fs"
	"net/http"
	"path/filepath"
	"regexp"
	"strings"
)

type GitIssue struct {
	Title    string `json:"title"`
	Body     string `json:"body"`
	Id       uint64 `json:"id,omitempty"`
	Number   uint64 `json:"number,omitempty"`
	State    string `json:"state,omitempty"`
	Assignee string `json:"-,omitempty"`
}

type RepoInfo struct {
	Owner  string
	Repo   string
	Issues []*GitIssue
}

// Create a GitHub issue with the title, body provided
func (r *RepoInfo) CreateIssue(title, body, token string) (*GitIssue, error) {
	issuesURL := fmt.Sprintf("https://api.github.com/repos/%s/%s/issues", r.Owner, r.Repo)
	payload := GitIssue{
		Title:    title,
		Body:     body,
		Assignee: r.Owner, // TODO: (#9:closed) not working; check token permissions include push access
	}
	var newIssue *GitIssue
	resp, err := Post(issuesURL, token, &payload)
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
func (r *RepoInfo) LoadIssues(token string) error {
	issuesURL := fmt.Sprintf("https://api.github.com/repos/%s/%s/issues?state=all", r.Owner, r.Repo)
	resp, err := Get(issuesURL, token)
	if err != nil {
		return err
	}
	err = json.Unmarshal(resp, &r.Issues)
	return err
}

/* Utils */

// Validates .git folder exists and parses Repo and Owner name
func CheckGitfolder(dir string) (*RepoInfo, error) {
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

// Reads and return token from .gittoken file
func LoadGitToken() (string, error) {
	token, err := ReadFile("./.gittoken")
	if err != nil {
		return "", err
	}
	return token, err
}

/* HTTP utils */
// Makes a POST request using Authorization: Bearer <token> and
// JSON payload
func Post(url string, token string, payload *GitIssue) ([]byte, error) {
	payloadBytes, err := json.Marshal(payload)
	if err != nil {
		return nil, fmt.Errorf("could not marshal payload: %v", err)
	}
	// Create a request object
	req, err := http.NewRequest("POST", url, bytes.NewReader(payloadBytes))
	if err != nil {
		return nil, fmt.Errorf("get: could not create request %s %v", url, err)
	}
	// Add headers
	req.Header.Add("Accept", "application/vnd.github+json")
	req.Header.Add("X-GitHub-Api-Version", "2022-11-28")
	req.Header.Add("Authorization", fmt.Sprintf("Bearer %s", token))
	// Create a new HTTP Client
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("get: could not make request %s %v", url, err)
	}
	defer resp.Body.Close()
	// Read the response body
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("get: could not read response body %v", err)
	}
	// Response is OK
	if resp.StatusCode != http.StatusCreated {
		return nil, fmt.Errorf("get: HTTP code error %d at %s", resp.StatusCode, url)
	}
	return body, nil
}

// Makes a GET request using Authorization: Bearer <token>
func Get(url string, token string) ([]byte, error) {
	// Create a request object
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return nil, fmt.Errorf("get: could not create request %s %v", url, err)
	}
	// Add headers
	req.Header.Add("Accept", "application/vnd.github+json")
	req.Header.Add("X-GitHub-Api-Version", "2022-11-28")
	req.Header.Add("Authorization", fmt.Sprintf("Bearer %s", token))
	// Create a new HTTP Client
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("get: could not make request %s %v", url, err)
	}
	defer resp.Body.Close()
	// Read the response body
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("get: could not read response body %v", err)
	}
	// Response is OK
	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("get: HTTP code error %d at %s", resp.StatusCode, url)
	}
	return body, nil
}
