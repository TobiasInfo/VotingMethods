package comsoc

import "fmt"

// TODO : implement CondorcetWinner function
func CondorcetWinner(p Profile) (bestAlts []Alternative, err error) {
	// check if the profile is valid
	if err = checkProfileAlternative(p, RecoverAlts(p)); err != nil {
		return nil, fmt.Errorf("invalid profile: %w", err)
	}

	return nil, fmt.Errorf("not implemented yet")
}
