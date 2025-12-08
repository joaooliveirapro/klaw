package main

func main() {
	printLogo()
	// Parse args
	args, err := ParseArgs()
	if err != nil {
		printErr(err)
		return
	}
	// Check .git folder exists in directory
	repo, err := CheckGitfolder(args.Directory)
	if err != nil {
		printErr(err)
		return
	}
	// Check .gittoken file exists
	token, err := LoadGitToken()
	if err != nil {
		printErr(err)
		return
	}
	// Get all local todos
	localTodos, err := FindAllTodosInDirectory(args.Directory, args.Extension)
	if err != nil {
		printErr(err)
		return
	}
	if args.ListAll {
		handleListAll(args, repo, localTodos, token)
	}
	if args.Create {
		handleCreate(repo, localTodos, token)
	}
}
