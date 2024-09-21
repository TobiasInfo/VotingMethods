package comsoc

// TODO : Implements these fonctions
func BordaSWF(p Profile) (count Count, err error) {
	count = make(Count)

	return count, nil
}
func BordaSCF(p Profile) (bestAlts []Alternative, err error) {
	count, err := BordaSWF(p)
	if err != nil {
		return nil, err
	}
	return maxCount(count), nil
}
