package server

type DBAgent struct {
	id        string
	createdAt string
	updatedAt string
}

type DBCommand struct {
	id        string
	agentID   string
	command   string
	result    string
	createdAt string
	updatedAt string
}
