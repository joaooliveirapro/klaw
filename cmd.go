package main

import (
	"errors"
	"flag"
)

type Flags struct {
	ListAll   bool
	Create    bool
	Update    bool
	Directory string
	Extension string
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
