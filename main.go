package main

import (
	"log"
)

// TODO: (#10:open) implement issue body parsing
// TODO: (#12:open) implement ignore-todo keyword for todos that can be ignored
// TODO: (#13:open) implement closing todos with the CLI (and it reflecting on the Github issue)
// TODO: (#14:open) implement optional skip closed todos logic when listing all
// TODO: (#16:open) implement config file for extra features (keyword, defaults, skips, etc.)
// TODO: (#17:open) implement actions table: list all actions that can be commited (create or update certain tickets)
// TODO: (#18:open) implement flag output (so users know what is set)
func main() {
	printLogo()

	args := ParseArgs()

	cfg, err := ParseConfigFile()
	if err != nil {
		log.Fatal(err)
	}

	m := NewManager(args, cfg)
	if !m.Args.Offline { // TODO: (#11:open) implement allow offline (no-git) checks for comments
		repo, err := CheckGitRepo(m)
		if err != nil {
			log.Fatal(err)
		}
		m.AssignRepo(repo)
	}
	if m.Args.ListAll {
		HandleListAll(m)
	}
	if m.Args.Create {
		HandleCreate(m)
	}
}
