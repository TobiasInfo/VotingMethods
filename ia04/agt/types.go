package agt

type Alternative int

// Agent represents a voting agent with preferences.
type Agent struct {
	ID    AgentID
	Name  string
	Prefs []Alternative
}

// AgentID is a type to represent a unique agent ID.
type AgentID int