package main

import (
	"errors"
	"flag"
	"fmt"
	"slices"
	"strings"
)

type Flags struct {
	ListAll   bool
	Create    bool
	Update    bool
	Directory string
	Extension string
}

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

func printLogo() {
	fmt.Println(C(`
██ ▄█▀ ██     ▄████▄ ██     ██ 
████   ██     ██▄▄██ ██ ▄█▄ ██ 
██ ▀█▄ ██████ ██  ██  ▀██▀██▀`+"\n", Green))

	fmt.Println(C("usage: klaw -<flags> <dir> <ext>", Gray))
	fmt.Println(C("\t-h - print help\n\t-l - list all\n\t-c - convert TODO comments to issues\n\t-u - update TODO state in comments\n", Gray))

}

// color text using ANSI escape codes
func C(text string, color string) string {
	return fmt.Sprintf("%s%s%s", color, text, Reset)
}

// print error formatted
func printErr(err error) {
	fmt.Printf("%s %s", C("[error]", Red), err)
}

func ParseArgs() (*Flags, error) {
	/*
		-h help
		-l list all todos
			// -lp sort by urgency // TODO: (#7:open) implement this logic
			// -lt sort by timestamp
			// 	-r to reverse order
		-c create todos
		-r resolve todos
	*/
	listall := flag.Bool("l", true, "list all todos (tracked and untracked)")
	create := flag.Bool("c", false, "create issues: convert all new todos in source fileds to GitHub issues")
	update := flag.Bool("u", false, "update issues: updates state of tracked todos based on remote issues")
	flag.Parse()
	directory := "." // defaults to cwd
	extension := ""  // no default
	if len(flag.Args()) > 1 {
		directory = flag.Arg(0)
		extension = flag.Arg(1)
	} else {
		if extension == "" {
			return nil, errors.New("missing file extension arg")
		}
	}
	return &Flags{
		ListAll:   *listall,
		Create:    *create,
		Update:    *update,
		Directory: directory,
		Extension: extension,
	}, nil
}

func handleListAll(args *Flags, repo *RepoInfo, localTodos []*Todo, token string) {
	fmt.Println(C("\n TODOS ", Green))

	// List all remote todos
	for _, todo := range localTodos {
		if !todo.Tracked {
			padded := fmt.Sprintf("%-17s", fmt.Sprintf("%s:%d", todo.File, todo.LineNumber)) // pad based on visible chars
			fmt.Printf("%-17s (%s) %s\n", C(padded, Red), C("UNTRACKED", Red), todo.Text)
		}
	}
	err := repo.LoadIssues(token)
	if err != nil {
		printErr(err)
	}
	for _, issue := range repo.Issues {
		// Check if remote issue has a twin todo
		idx := slices.IndexFunc(localTodos, func(todo *Todo) bool {
			return todo.IssueNumber == fmt.Sprint(issue.Number)
		})
		ok := idx > -1
		if !ok {
			// Issue opened on GitHub but not on src files
			// fmt.Printf("%-17s (%s:%s) %s\n", C("[not-in-src]", Gray), C(fmt.Sprintf("#%d", issue.Number), Yellow), issue.State, issue.Title)
			padded := fmt.Sprintf("%-17s", "not-in-src") // pad based on visible chars
			fmt.Printf("%s (%s:%s) %s\n", C(padded, Gray), C(fmt.Sprintf("#%d", issue.Number), Yellow), issue.State, issue.Title)
		} else {
			todo := localTodos[idx]
			// Issue both on GitHub and on src files
			// Check state is updated
			if todo.IssueState == issue.State {
				// Issue and TODO up to date
				s := fmt.Sprintf("%s:%s", todo.File, fmt.Sprint(todo.LineNumber))
				padded := fmt.Sprintf("%-17s", s) // pad based on visible chars
				fmt.Printf("%s (%s:%s) %s\n", C(padded, Cyan), C(fmt.Sprintf("#%d", issue.Number), Yellow), issue.State, issue.Title)
			} else {
				// Update if args flag -u was set
				if args.Update {
					// Replace todo state with the new one
					err := RewriteTodoState(todo, issue)
					if err != nil {
						fmt.Printf("err: %v\n", err)
					}
				} else {
					s := fmt.Sprintf("%s:%s", todo.File, fmt.Sprint(todo.LineNumber))
					padded := fmt.Sprintf("%-17s", s) // pad based on visible chars
					fmt.Printf("%s (%s:<actual:%s!%s>) %s\n", C(padded, Yellow), C(fmt.Sprintf("#%d", issue.Number), Yellow), C(issue.State, Green), C(todo.IssueState, Red), issue.Title)
				}
			}
		}
	}
}

func handleCreate(repo *RepoInfo, localTodos []*Todo, token string) {
	for _, todo := range localTodos {
		if todo.Tracked {
			continue
		}
		var proceed string
		fmt.Printf("\ncreate issue\n\ttitle: %s\n\tbody: %s\n\n? [Y/n]", todo.Text, todo.Text)
		fmt.Scanln(&proceed)
		if strings.ToLower(proceed) == "y" {
			newIssue, err := repo.CreateIssue(todo.Text, todo.Text, token)
			if err != nil {
				fmt.Println(err)
			}
			err = RewriteUntrackedTodo(todo, newIssue)
			if err != nil {
				fmt.Println(err)
			}
		} else {
			fmt.Println("skipping.")
		}
	}
}
