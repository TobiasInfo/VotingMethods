package comsoc

import "fmt"

// CondorcetWinner finds the Condorcet winner in the profile.
// A Condorcet winner is an alternative that wins against all others in pairwise comparisons.
func CondorcetWinner(p Profile) (bestAlts []Alternative, err error) {
	alts := RecoverAlts(p)
	if err := checkProfileAlternative(p, alts); err != nil {
		return nil, err
	}

	// Map to store pairwise wins: wins[a][b] means "number of voters preferring a over b"
	wins := make(map[Alternative]map[Alternative]int)
	for _, alt := range alts {
		wins[alt] = make(map[Alternative]int)
	}

	// Pairwise comparison: for each voter's preferences, increment win counts
	for _, prefs := range p {
		for i := 0; i < len(prefs)-1; i++ {
			for j := i + 1; j < len(prefs); j++ {
				wins[prefs[i]][prefs[j]]++ // prefs[i] is preferred over prefs[j]
			}
		}
	}

	// Check if any alternative is a Condorcet winner
	for _, a := range alts {
		// We assume a is the Condorcet winner until proven otherwise.
		isWinner := true
		for _, b := range alts {
			if a != b {
				// Check if more voters prefer 'a' over 'b' than 'b' over 'a'
				votesForA := wins[a][b]
				votesForB := wins[b][a]
				if votesForA <= votesForB {
					isWinner = false
					break
				}
			}
		}
		if isWinner {
			return []Alternative{a}, nil
		}
	}
	return nil, fmt.Errorf("no Condorcet winner")
}
