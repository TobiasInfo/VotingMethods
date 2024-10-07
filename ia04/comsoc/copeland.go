package comsoc

func CopelandSWF(p Profile) (Count, error) {
	alts := RecoverAlts(p)
	if err := checkProfileAlternative(p, alts); err != nil {
		return nil, err
	}

	// Initialize a count map to store Copeland scores for each alternative
	count := make(Count)
	for _, alt := range alts {
		count[alt] = 0
	}

	// Pairwise comparison: for each pair of alternatives (a, b), calculate scores
	for i := 0; i < len(alts); i++ {
		for j := i + 1; j < len(alts); j++ {
			a, b := alts[i], alts[j]
			votesForA, votesForB := 0, 0

			// Count how many voters prefer 'a' over 'b' and vice versa
			for _, prefs := range p {
				if rank(a, prefs) < rank(b, prefs) {
					votesForA++
				} else {
					votesForB++
				}
			}

			// Update scores based on the majority preference
			if votesForA > votesForB {
				count[a]++
				count[b]--
			} else if votesForA < votesForB {
				count[a]--
				count[b]++
			}
			// In case of a tie, no changes to the score
		}
	}

	return count, nil
}

func CopelandSCF(p Profile) (bestAlts []Alternative, err error) {
	// Get the Copeland scores for each alternative using CopelandSWF
	count, err := CopelandSWF(p)
	if err != nil {
		return nil, err
	}

	return maxCount(count), nil
}
