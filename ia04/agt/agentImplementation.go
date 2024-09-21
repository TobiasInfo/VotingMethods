package agt

import (
	"fmt"
)

// Check if the ID of the agent is equal to the ID of the other agent.
func (ag *Agent) Equal(other AgentI) bool {
	// Check if the other object is a pointer to an Agent
	o, ok := other.(*Agent)
	if !ok {
		return false
	}
	return ag.ID == o.ID
}

// Check all fields of the agent are equal to the fields of the other agent.
func (ag *Agent) DeepEqual(other AgentI) bool {
	o, ok := other.(*Agent)
	if !ok {
		return false
	}
	if ag.ID != o.ID || ag.Name != o.Name {
		return false
	}
	// Compare preferences deeply
	if len(ag.Prefs) != len(o.Prefs) {
		return false
	}
	for i, alt := range ag.Prefs {
		if alt != o.Prefs[i] {
			return false
		}
	}
	return true
}

// Create a deep copy of the agent.
func (ag *Agent) Clone() AgentI {
	clone := *ag
	clone.Prefs = make([]Alternative, len(ag.Prefs))
	copy(clone.Prefs, ag.Prefs)
	return &clone
}

// Return a string representation of the agent.
func (ag *Agent) String() string {
	return fmt.Sprintf("Agent[ID=%d, Name=%s, Prefs=%v]", ag.ID, ag.Name, ag.Prefs)
}

// Check if the agent prefers alternative a over alternative b.
func (ag *Agent) Prefers(a, b Alternative) bool {
	for _, alt := range ag.Prefs {
		if alt == a {
			return true
		}
		if alt == b {
			return false
		}
	}
	return false
}

// Start the agent.
func (ag *Agent) Start() {
	fmt.Printf("Agent %s is ready to vote!\n", ag.Name)
}
