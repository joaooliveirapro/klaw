package main

type Manager struct {
	Args       *Flags
	Repo       *RepoInfo
	Config     *Config
	LocalTodos []*Todo
}

func NewManager(args *Flags, cfg *Config) *Manager {
	return &Manager{
		Args:   args,
		Config: cfg,
	}
}

func (m *Manager) AssignRepo(r *RepoInfo) {
	m.Repo = r
}

func (m *Manager) AssignLocalTodos(todos []*Todo) {
	m.LocalTodos = todos
}

type Flags struct {
	ListAll    bool
	Create     bool
	Update     bool
	Offline    bool
	SkipClosed bool
	PrintFlags bool
}

type Todo struct {
	Text        string
	IssueNumber string
	IssueState  string
	Tracked     bool
	File        string
	LineNumber  int
}

type File struct {
	Path  string  // full path to file
	Todos []*Todo // list of todos found in the file
}

type GitIssue struct {
	Title    string `json:"title"`
	Body     string `json:"body"`
	Id       uint64 `json:"id,omitempty"`
	Number   uint64 `json:"number,omitempty"`
	State    string `json:"state,omitempty"`
	Assignee string `json:"-,omitempty"`
}

type RepoInfo struct {
	Owner  string
	Repo   string
	Token  string
	Issues []*GitIssue
}
