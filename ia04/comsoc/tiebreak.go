package comsoc

import "fmt"

func TieBreakFactory(orderedAlts []Alternative) func([]Alternative) (Alternative, error) {
	return func(alts []Alternative) (Alternative, error) {
		altMap := make(map[Alternative]bool)
		for _, alt := range alts {
			altMap[alt] = true
		}

		for _, alt := range orderedAlts {
			if altMap[alt] {
				return alt, nil
			}
		}

		return -1, fmt.Errorf("aucune alternative ne correspond")
	}
}

func SWFFactory(swf func(p Profile) (Count, error), tb func([]Alternative) (Alternative, error)) func(Profile) ([]Alternative, error) {
	return func(p Profile) ([]Alternative, error) {
		count, err := swf(p)
		if err != nil {
			return nil, fmt.Errorf("failed to get count from SWF: %w", err)
		}
		bestAlts := make([]Alternative, 0)
		for {
			bestAlt := maxCount(count)
			if len(bestAlt) == 0 {
				break
			}
			if len(bestAlt) == 1 {
				bestAlts = append(bestAlts, bestAlt[0])
				delete(count, bestAlt[0])
			} else {
				alt, err := tb(bestAlt)
				if err != nil {
					return nil, fmt.Errorf("failed to get best alternative from tie-breaking function: %w", err)
				}
				bestAlts = append(bestAlts, alt)
				delete(count, alt)
			}
		}
		return bestAlts, nil
	}
}

func SCFFactory(scf func(p Profile) ([]Alternative, error), tb func([]Alternative) (Alternative, error)) func(Profile) (Alternative, error) {
	return func(p Profile) (Alternative, error) {
		alts, err := scf(p)
		if err != nil {
			return -1, fmt.Errorf("failed to get alternatives from SCF: %w", err)
		}
		if len(alts) == 1 {
			return alts[0], nil
		}
		alt, err := tb(alts)
		if err != nil {
			return -1, fmt.Errorf("failed to get best alternative from tie-breaking function: %w", err)
		}
		return alt, nil
	}
}
