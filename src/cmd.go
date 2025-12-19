package klaw

import (
	"flag"
	"fmt"

	"github.com/jedib0t/go-pretty/v6/table"
)

func PrintLogo() {
	fmt.Println(C(`
██ ▄█▀ ██     ▄████▄ ██     ██ 
████   ██     ██▄▄██ ██ ▄█▄ ██ 
██ ▀█▄ ██████ ██  ██  ▀██▀██▀`+"\n", Green))
}

func ParseArgs() *Flags {
	listall := flag.Bool("l", true, "list all todos (tracked and untracked) [requires gittoken]")
	skipclosed := flag.Bool("s", false, "Issues with state closed will not be included in the table")
	create := flag.Bool("c", false, "create issues: convert all new todos in source files into GitHub issues [requires gittoken]")
	update := flag.Bool("u", false, "update issues: updates state of tracked todos based on remote issues [requires gittoken]")
	offline := flag.Bool("o", false, "offline mode: find TODOs in source files")
	printflags := flag.Bool("f", false, "print flags set")
	flag.Usage = func() { // Override usage string
		fmt.Printf("Usage of KLAW:\n")
		flag.PrintDefaults()
	}
	flag.Parse()
	f := Flags{
		ListAll:    *listall,
		Create:     *create,
		Update:     *update,
		Offline:    *offline,
		SkipClosed: *skipclosed,
		PrintFlags: *printflags,
	}
	if f.PrintFlags {
		fmt.Println(C("Flags table", Green))
		t := createTable("#", "Flag", "State")
		t.AppendRows([]table.Row{
			{"1", "-l listall", *listall},
			{"2", "-c create", *create},
			{"3", "-u update", *update},
			{"4", "-o offline", *offline},
			{"5", "-s skipclosed", *skipclosed},
			{"6", "-f printflags", *printflags},
		})
		t.Render()
		fmt.Println() // spacing
	}
	return &f
}
