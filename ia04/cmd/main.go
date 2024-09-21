package main

import (
	"fmt"
	"ia04/comsoc"
	"ia04/agt"
)

func main() {
	// Create some agents
	agent1 := &agt.Agent{ID: 1, Name: "Alice", Prefs: []agt.Alternative{1, 2, 3}}
	agent2 := &agt.Agent{ID: 2, Name: "Bob", Prefs: []agt.Alternative{2, 3, 1}}
	agent3 := &agt.Agent{ID: 3, Name: "Charlie", Prefs: []agt.Alternative{3, 1, 2}}

	// Group the agents into a profile
	profile := comsoc.Profile{
		agent1.Prefs,
		agent2.Prefs,
		agent3.Prefs,
	}

	// Run the Majority SCF procedure
	bestAlts, err := comsoc.MajoritySCF(profile)
	if err != nil {
		fmt.Println("Error:", err)
		return
	}
	fmt.Println("Best alternatives according to Majority SCF:", bestAlts)
}
