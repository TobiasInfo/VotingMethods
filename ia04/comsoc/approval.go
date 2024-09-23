package comsoc

import "fmt"

func ApprovalSWF(p Profile, thresholds []int) (count Count, err error) {
	alts := RecoverAlts(p)
	if err := checkProfileAlternative(p, alts); err != nil {
		return nil, err
	}
	if len(p) != len(thresholds) {
		return nil, fmt.Errorf("le nombre de seuils ne correspond pas au nombre de préférences")
	}
	count = make(Count)
	for i := 0; i < len(p); i++ {
		for j := 0; j < thresholds[i]; j++ {
			count[p[i][j]]++
		}
	}
	return count, nil
}

func ApprovalSCF(p Profile, thresholds []int) (bestAlts []Alternative, err error) {
	count, err := ApprovalSWF(p, thresholds)
	if err != nil {
		return nil, err
	}
	return maxCount(count), nil
}
