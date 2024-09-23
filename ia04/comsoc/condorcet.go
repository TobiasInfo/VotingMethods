package comsoc

import "fmt"

func CondorcetWinner(p Profile) (bestAlts []Alternative, err error) {
	alts := RecoverAlts(p)
	if err := checkProfileAlternative(p, alts); err != nil {
		return nil, err
	}

	// TODO : implement the CondorcetWinner algorithm
	return nil, fmt.Errorf("CondorcetWinner not implemented yet")

}
