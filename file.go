package main

import (
	"fmt"
	"log"
	"os"
	"regexp"
	"strings"
)

// Finds all todos in the src files that match extension in the directory
func FindAllTodosInDirectory(m *Manager) []*Todo {
	// TODO: (#8:open) implement logic to cover multiple extensions
	localTodos := []*Todo{}
	// Directories to exclude from search
	excluded := make(map[string]int)
	for _, f := range m.Config.ExcludeFolders {
		excluded[f] = 1
	}
	for _, ext := range m.Config.Extensions {
		files := findFilesByExtension(m.Config.Directory, ext, excluded)
		for _, f := range files {
			err := f.findTodosInFile(m.Config.TodoCommentSymbol, m.Config.TodoKeyword)
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

// Checks file and loads all todos found within
func (f *File) findTodosInFile(todoCommentSymbol, todoKeyword string) error {
	fileContent, err := ReadFile(f.Path)
	if err != nil {
		return err
	}
	todoRgx := `(?i)` + regexp.QuoteMeta(todoCommentSymbol) +
		`\s*` + regexp.QuoteMeta(todoKeyword) +
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

// Rewrites the state (#id:state) of a todo to the new state found in the
// twin GitHub issue; Example: when an issue is closed it updated the local
// todos state to closed;
func RewriteTodoState(todo *Todo, issue *GitIssue) error {
	content, err := ReadFile(todo.File)
	if err != nil {
		return err
	}
	// KEYWORD: (#123:open) todotext
	// KEYWORD: (#123:closed) todotext
	oldtodostr := fmt.Sprintf("(#%s:%s) %s", todo.IssueNumber, todo.IssueState, todo.Text) // todo contains old values
	newtodostr := fmt.Sprintf("(#%d:%s) %s", issue.Number, issue.State, todo.Text)         // issue contains updated values
	// to ensure replaceAll works properly, the newtodostr text has to be todo.Text
	newcontent := strings.ReplaceAll(content, oldtodostr, newtodostr)
	err = os.WriteFile(todo.File, []byte(newcontent), 0644)
	if err != nil {
		return err
	}
	s := fmt.Sprintf("%s:%s", todo.File, fmt.Sprint(todo.LineNumber))
	padded := fmt.Sprintf("%-17s", s) // pad based on visible chars
	fmt.Printf("%s (%s:%s) %s %s\n",
		C(padded, Cyan), C(fmt.Sprintf("#%d", issue.Number), Yellow), issue.State,
		todo.Text, C("updated", Green))
	return nil
}

// Rewrites the untracekd todo with the corresponding repr of its twin
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
