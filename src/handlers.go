package klaw

import (
	"fmt"
	"strconv"
	"strings"

	"github.com/jedib0t/go-pretty/v6/table"
)

func HandleListAll(m *Manager) {
	localTodos := FindAllTodosInDirectory(m)
	m.AssignLocalTodos(localTodos)

	ott := createTable("#", "File:Line #", "Issue #:Status", "Title")

	if m.Args.Offline {
		handleListOffline(m, ott)
	} else {
		handleListOnline(m, ott)
	}
	ott.Render()
}

func handleListOffline(m *Manager, ott table.Writer) {
	for todo_idx, todo := range m.LocalTodos { // Check local (always)
		if m.Args.SkipClosed && todo.IssueState == "closed" { // Skip closed TODOS
			continue
		}
		fileAndLinenum := fmt.Sprintf("%s:%d", C(todo.File, Cyan), todo.LineNumber)
		issueAndStatus := C("untracked", Red)
		if todo.Tracked {
			issueAndStatus = fmt.Sprintf("%s:%s", C(fmt.Sprintf("#%s", todo.IssueNumber), Yellow), todo.IssueState)
		}
		row := table.Row{todo_idx + 1, fileAndLinenum, issueAndStatus, todo.Text}
		ott.AppendRow(row)
	}
	fmt.Printf("%s %s\n", C("OFFLINE mode", Yellow), C("- status of tracked TODOs may not be up to date", Gray))
}

func handleListOnline(m *Manager, ott table.Writer) {
	err := m.Repo.LoadIssues()
	IfErrPrint(err)

	trackedTodos := map[string]*Todo{}
	untrackedTodos := []*Todo{}
	for _, todo := range m.LocalTodos {
		if todo.Tracked {
			trackedTodos[todo.IssueNumber] = todo
		} else {
			untrackedTodos = append(untrackedTodos, todo)
		}
	}

	i := 1 // cumulative count for local todos and remote issues
	for _, issue := range m.Repo.Issues {
		if m.Args.SkipClosed && issue.State == "closed" { // Skip closed issues
			continue
		}
		fileAndLinenum, issueAndStatus := "", ""

		var row table.Row
		todo, ok := trackedTodos[fmt.Sprint(issue.Number)]
		if !ok { // Issue opened on GitHub but not on src files
			fileAndLinenum = C("not-in-src", Gray)
			issueAndStatus = C(fmt.Sprintf("#%d:%s", issue.Number, issue.State), Yellow)
		} else { // Issue open on Github and has todo in source file
			// Check if issue state == todo state
			if todo.IssueState == issue.State { // Remote and local todo state match
				fileAndLinenum = fmt.Sprintf("%s:%s", C(todo.File, Cyan), C(strconv.Itoa(todo.LineNumber), Yellow))
				issueAndStatus = C(fmt.Sprintf("#%d:%s", issue.Number, issue.State), Yellow)
			} else { // Remote and local todo state differ
				if m.Args.Update {
					// -u flag set: update the state of todos to match remote issues
					err := RewriteTodoState(todo, issue) // Replace todo state with the new one
					IfErrPrint(err)
				} else {
					// -u not set: leave todos it as is
					fileAndLinenum = fmt.Sprintf("%s:%s", C(todo.File, Cyan), C(strconv.Itoa(todo.LineNumber), Yellow))
					issueAndStatus = C(fmt.Sprintf("#%d:%s (!%s)", issue.Number, issue.State, todo.IssueState), Yellow)
				}
			}
		}
		row = append(row, i, fileAndLinenum, issueAndStatus, issue.Title)
		ott.AppendRow(row)
		i += 1
	}

	for _, todo := range untrackedTodos {
		fileAndLinenum := fmt.Sprintf("%s:%d", C(todo.File, Cyan), todo.LineNumber)
		issueAndStatus := C("untracked", Red)
		ott.AppendRow(table.Row{i, fileAndLinenum, issueAndStatus, todo.Text})
	}
}

func HandleCreate(m *Manager) {
	for _, todo := range m.LocalTodos {
		if todo.Tracked { // skip todos already tracked in Github
			continue
		}
		var proceed string
		fmt.Printf("\ncreate issue\n\ttitle: %s\n\tbody: %s\n\n? [Y/n]", todo.Text, todo.Text)
		fmt.Scanln(&proceed)
		if strings.ToLower(proceed) == "y" {
			newIssue, err := m.Repo.CreateIssue(todo.Text, todo.Text)
			if err != nil {
				IfErrPrint(err)
				return
			}
			err = RewriteUntrackedTodo(todo, newIssue)
			IfErrPrint(err)
		} else {
			fmt.Println("skipping.")
		}
	}
}
