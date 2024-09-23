package comsoc

func MajoritySWF(p Profile) (count Count, err error) {
	alts := RecoverAlts(p)
	if err := checkProfileAlternative(p, alts); err != nil {
		return nil, err
	}
	count = make(Count)
	for _, prefs := range p {
		firstChoice := prefs[0]
		count[firstChoice] += 1
	}
	return count, nil
}

func MajoritySCF(p Profile) (bestAlts []Alternative, err error) {
	count, err := MajoritySWF(p)
	if err != nil {
		return nil, err
	}
	return maxCount(count), nil
}
