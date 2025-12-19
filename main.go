package main

import (
	"log"

	klaw "github.com/joaooliveirapro/klaw/src"
)

// TODO: (#10:open) implement issue body parsing
// TODO: (#12:open) implement ignore-todo keyword for todos that can be ignored
// TODO: (#13:open) implement closing todos with the CLI (and it reflecting on the Github issue)
// TODO: (#14:closed) implement optional skip closed todos logic when listing all
// TODO: (#16:closed) implement config file for extra features (keyword, defaults, skips, etc.)
// TODO: (#17:open) implement actions table: list all actions that can be commited (create or update certain tickets)
// TODO: (#18:closed) implement flag output (so users know what is set)
func main() {
	klaw.PrintLogo()

	args := klaw.ParseArgs()

	cfg, err := klaw.ParseConfigFile()
	if err != nil {
		log.Fatal(err)
	}

	m := klaw.NewManager(args, cfg)
	if !m.Args.Offline { // TODO: (#11:closed) implement allow offline (no-git) checks for comments
		repo, err := klaw.CheckGitRepo(m)
		if err != nil {
			log.Fatal(err)
		}
		m.AssignRepo(repo)
	}
	if m.Args.ListAll {
		klaw.HandleListAll(m)
	}
	if m.Args.Create {
		klaw.HandleCreate(m)
	}
}
