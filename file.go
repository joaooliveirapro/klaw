package main

import (
	"bufio"
	"fmt"
	"io"
	"io/fs"
	"log"
	"os"
	"path/filepath"
	"regexp"
	"strings"
)

// Checks file and loads all todos found within
func (f *File) FindTodosInFile(m *Manager) error {
	fileContent, err := ReadFile(f.Path)
	if err != nil {
		return err
	}

	todoRgx := `(?i)` + regexp.QuoteMeta(m.Config.TodoCommentSymbol) +
		`\s*` + regexp.QuoteMeta(m.Config.TodoKeyword) +
		`:\s*(\(\#(\d+):(\w+)\))?(.*)`

	rgx := regexp.MustCompile(todoRgx)
	matches := rgx.FindAllStringSubmatch(fileContent, -1) // -1: no limit
	for _, m := range matches {
		issueNumber := m[2]
		issueState := m[3]
		body := m[4]
		tracked := false
		if issueNumber != "" { // if there's an issue # it's a tracked todo
			tracked = true
		}
		txt := strings.TrimSpace(body)
		linenumber, err := findLineNumber(f.Path, body)
		if err != nil {
			return err
		}
		todo := Todo{
			Text:        txt,
			IssueNumber: issueNumber,
			IssueState:  issueState,
			Tracked:     tracked,
			File:        f.Path,
			LineNumber:  linenumber,
		}
		f.Todos = append(f.Todos, &todo)
	}
	return nil
}

// Finds all todos in the src files that match extension in the directory
func FindAllTodosInDirectory(m *Manager) []*Todo {
	// TODO: (#8:open) implement logic to cover multiple extensions
	localTodos := []*Todo{}
	for _, ext := range m.Config.Extensions {
		files := findFilesByExtension(ext, m)
		for _, f := range files {
			err := f.FindTodosInFile(m)
			if err != nil {
				log.Println(err) // Not a critical error if line number isn't found
				return nil
			}
			if len(f.Todos) == 0 {
				continue
			}
			localTodos = append(localTodos, f.Todos...)
		}
	}
	return localTodos
}

// Rewrites the state (#id:state) of a todo to the new state found in the
// twin GitHub issue
func RewriteTodoState(todo *Todo, issue *GitIssue) error {
	content, err := ReadFile(todo.File)
	if err != nil {
		return err
	}
	// KEYWORD: (#123:open) todotext
	// KEYWORD: (#123:closed) todotext
	oldtodostr := fmt.Sprintf("(#%s:%s) %s", todo.IssueNumber, todo.IssueState, todo.Text) // todo contains old values
	newtodostr := fmt.Sprintf("(#%d:%s) %s", issue.Number, issue.State, todo.Text)         // issue contains updated values
	newcontent := strings.ReplaceAll(content, oldtodostr, newtodostr)
	err = os.WriteFile(todo.File, []byte(newcontent), 0644)
	if err != nil {
		return err
	}
	s := fmt.Sprintf("%s:%s", todo.File, fmt.Sprint(todo.LineNumber))
	padded := fmt.Sprintf("%-17s", s) // pad based on visible chars
	fmt.Printf("%s (%s:%s) %s %s\n",
		C(padded, Cyan), C(fmt.Sprintf("#%d", issue.Number), Yellow), issue.State,
		issue.Title, C("updated", Green))
	return nil
}

// Rewrutes the untracekd todo with the corresponding repr of its twin
// GitHub issue, including (#issuenumber:state) original text
func RewriteUntrackedTodo(todo *Todo, issue *GitIssue) error {
	// KEYWORD: todotext
	// KEYWORD: (#123:open) todotext
	content, err := ReadFile(todo.File)
	if err != nil {
		return err
	}
	newcontent := strings.ReplaceAll(content, todo.Text, fmt.Sprintf("(#%d:%s) %s", issue.Number, issue.State, todo.Text))
	err = os.WriteFile(todo.File, []byte(newcontent), 0644)
	if err != nil {
		return err
	}

	s := fmt.Sprintf("%s:%s", todo.File, fmt.Sprint(todo.LineNumber))
	padded := fmt.Sprintf("%-17s", s) // pad based on visible chars
	fmt.Printf("%s (%s:%s) %s %s\n",
		C(padded, Cyan), C(fmt.Sprintf("#%d", issue.Number), Yellow), issue.State,
		issue.Title, C("updated", Green))
	return nil
}

/* Utils */

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

// Find files that match extension
func findFilesByExtension(extenstion string, m *Manager) []*File {
	filesToCheck := []*File{}
	excluded := make(map[string]int)
	for _, f := range m.Config.ExcludeFolders {
		excluded[f] = 1
	}
	filepath.Walk(m.Config.Directory, func(path string, info fs.FileInfo, err error) error {
		if err != nil {
			return nil
		}
		if info.IsDir() {
			if _, skip := excluded[info.Name()]; skip {
				return filepath.SkipDir // signals .Walk to skip this directory
			}
		}
		if strings.HasSuffix(info.Name(), extenstion) {
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
