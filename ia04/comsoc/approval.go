package comsoc

// TODO : Implements these fonctions
func ApprovalSWF(p Profile, thresholds []int) (count Count, err error) {
	count = make(Count)

	return count, nil
}
func ApprovalSCF(p Profile, thresholds []int) (bestAlts []Alternative, err error) {
	count, err := ApprovalSWF(p, thresholds)
	if err != nil {
		return nil, err
	}
	return maxCount(count), nil
}
