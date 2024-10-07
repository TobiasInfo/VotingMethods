package comsoc

// fct qui prend en entrée un profile et retourne la préférence du groupe sur l'ensemble des alternatives (ordre de préférence)
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

// fct qui prend en entrée l'ensemble des profiles et retourne une alternative ou un ensemble d'alternatives si on a un ex-eco
func MajoritySCF(p Profile) (bestAlts []Alternative, err error) {
	count, err := MajoritySWF(p)
	if err != nil {
		return nil, err
	}
	return maxCount(count), nil
}
