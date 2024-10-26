package comsoc

func BordaSWF(p Profile) (count Count, err error) {
	alts := RecoverAlts(p)
	if err := checkProfileAlternative(p, alts); err != nil {
		return nil, err
	}
	count = make(Count)
	for _, alt := range alts {
		count[alt] = 0
	}
	for _, prefs := range p {
		for i, alt := range prefs {
			count[alt] += len(prefs) - 1 - i
		}
	}
	return count, nil
}
func BordaSCF(p Profile) (bestAlts []Alternative, err error) {
	count, err := BordaSWF(p)
	if err != nil {
		return nil, err
	}
	return MaxCount(count), nil
}
