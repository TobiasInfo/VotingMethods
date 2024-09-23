package comsoc

import "fmt"

func CondorcetWinner(p Profile) (bestAlts []Alternative, err error) {
	alts := RecoverAlts(p)
	if err := checkProfileAlternative(p, alts); err != nil {
		return nil, err
	}

	wins := make(map[Alternative]map[Alternative]int)
	for _, prefs := range p {
		for i := 0; i < len(prefs)-1; i++ {
			for j := i + 1; j < len(prefs); j++ {
				wins[prefs[i]][prefs[j]]++
			}
		}
	}

	return nil, fmt.Errorf("no Condorcet winner")
}
