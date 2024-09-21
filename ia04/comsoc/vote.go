package comsoc

func SWF(p Profile) (count Count, err error) {
	count = make(Count)
	for _, prefs := range p {
		for i, alt := range prefs {
			count[alt] += len(prefs) - i - 1
		}
	}
	return count, nil
}
func SCF(p Profile) (bestAlts []Alternative, err error) {
	count, err := SWF(p)
	if err != nil {
		return nil, err
	}
	return maxCount(count), nil
}
