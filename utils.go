package main

import (
	"bufio"
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"io/fs"
	"net/http"
	"os"
	"path/filepath"
	"strings"
)

// ANSI colors
const (
	Reset        = "\033[0m"
	Bold         = "\033[1m"
	Dim          = "\033[2m"
	Red          = "\033[31m"
	Green        = "\033[32m"
	Yellow       = "\033[33m"
	Blue         = "\033[34m"
	Magenta      = "\033[35m"
	Cyan         = "\033[36m"
	Gray         = "\033[90m"
	GreenBGWhite = "\033[42;97m"
)

// color text using ANSI escape codes
func C(text string, color string) string {
	return fmt.Sprintf("%s%s%s", color, text, Reset)
}

// print error formatted
func IfErrPrint(err error) {
	if err != nil {
		fmt.Printf("%s %s", C("[error]", Red), err)
	}
}

// IO

// Find files in dir that match extension and exclude searchin in excluded dirs
func findFilesByExtension(dir, ext string, excluded map[string]int) []*File {
	filesToCheck := []*File{}
	filepath.Walk(dir, func(path string, info fs.FileInfo, err error) error {
		if err != nil {
			return nil
		}
		if info.IsDir() {
			if _, skip := excluded[info.Name()]; skip {
				return filepath.SkipDir // signals .Walk to skip this directory
			}
		}
		if strings.HasSuffix(info.Name(), ext) {
			filesToCheck = append(filesToCheck, &File{Path: path})
		}
		return nil
	})
	return filesToCheck
}

// Find the line number where substr occurs in the file
func findLineNumber(filepath, text string) (int, error) {
	f, err := os.Open(filepath)
	if err != nil {
		return -1, err
	}
	defer f.Close()
	scanner := bufio.NewScanner(f)
	lineNum := 1
	for scanner.Scan() {
		line := scanner.Text()
		if strings.Contains(line, strings.TrimSpace(text)) {
			return lineNum, nil
		}
		lineNum++
	}
	return -1, nil
}

// Read file and return its content
func ReadFile(filepath string) (string, error) {
	f, err := os.Open(filepath)
	if err != nil {
		return "", err
	}
	defer f.Close()
	contentB, err := io.ReadAll(f)
	if err != nil {
		return "", err
	}
	return string(contentB), nil
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
