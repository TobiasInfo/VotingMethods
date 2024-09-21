package comsoc

import "fmt"

// renvoie l'indice ou se trouve alt dans prefs
func rank(alt Alternative, prefs []Alternative) int {
	for i, a := range prefs {
		if a == alt {
			return i
		}
	}
	// Par convention, on renvoie -1 si l'alternative n'est pas trouvée
	return -1
}

// renvoie vrai ssi alt1 est préférée à alt2
func isPref(alt1, alt2 Alternative, prefs []Alternative) bool {
	// Cette solution est naïve, on pourrait faire mieux
	// Elle améliore la lisibilité du code
	// Mais impose de parcourir prefs deux fois
	return rank(alt1, prefs) < rank(alt2, prefs)

	// Solution plus efficace, mais moins lisible (A vérifier)
	/*
		ever_seen_alt1 := false
		ever_seen_alt2 := false
		for _, a := range prefs {
			if a == alt1 {
				if ever_seen_alt2 {
					return false
				}
				ever_seen_alt1 = true
			}
			if a == alt2 {
				if ever_seen_alt1 {
					return true
				}
				ever_seen_alt2 = true
			}
		}
		// Si alt1 et alt2 ne sont pas dans prefs, on considère qu'elles sont équivalentes
		return false
	*/

}

// type Count map[Alternative]int
// renvoie les meilleures alternatives pour un décomtpe donné
func maxCount(count Count) (bestAlts []Alternative) {
	bestAlts = make([]Alternative, 0)
	// Chercher le max dans la map
	var max int
	for _, c := range count {
		if c > max {
			max = c
		}
	}
	// Parcourir la map pour trouver les alternatives qui ont le même score
	for alt, c := range count {
		if c == max {
			bestAlts = append(bestAlts, alt)
		}
	}
	return bestAlts
}

// vérifie les préférences d'un agent, par ex. qu'ils sont tous complets
// et que chaque alternative n'apparaît qu'une seule fois
func checkProfile(prefs []Alternative, alts []Alternative) error {
	// Vérifier que prefs est complet
	if len(prefs) != len(alts) {
		return fmt.Errorf("Le profil n'est pas complet")
	}
	// Vérifier que chaque alternative n'apparaît qu'une seule fois
	seen := make(map[Alternative]bool)
	for _, alt := range prefs {
		if seen[alt] {
			return fmt.Errorf("L'alternative %v apparaît plusieurs fois", alt)
		}
		seen[alt] = true
	}
	// Vérifier que chaque alternative de alts apparaît exactement une fois dans prefs
	for _, alt := range alts {
		if !seen[alt] {
			return fmt.Errorf("L'alternative %v n'apparaît pas dans le profil", alt)
		}
	}
	return nil
}

// vérifie le profil donné, par ex. qu'ils sont tous complets et
// que chaque alternative de alts apparaît exactement une fois par préférences
func checkProfileAlternative(prefs Profile, alts []Alternative) error {
	for _, p := range prefs {
		if err := checkProfile(p, alts); err != nil {
			return err
		}
	}
	return nil
}

func Contains(s []Alternative, e Alternative) bool {
	for _, a := range s {
		if a == e {
			return true
		}
	}
	return false
}

func SWF(p Profile) (count Count, err error) {
	// Recover all the alternatives
	alts := make([]Alternative, 0)
	for _, prefs := range p {
		for _, alt := range prefs {
			if !Contains(alts, alt) {
				alts = append(alts, alt)
			}
		}
	}

	// Check if the profile is valid
	err = checkProfileAlternative(p, alts)
	if err != nil {
		return nil, err
	}

	// Initialize the count for each alternative
	count = make(Count)
	for _, alt := range alts {
		count[alt] = 0
	}

	// Score each alternative based on its rank in each agent's preferences
	for _, prefs := range p {
		numPrefs := len(prefs) // Total number of preferences for the agent
		for i, alt := range prefs {
			// Give a score of len(prefs) - 1 for the first choice, len(prefs) - 2 for the second choice, etc.
			count[alt] += numPrefs - 1 - i
		}
	}

	return count, nil
}


func SCF(p Profile) (bestAlts []Alternative, err error) {
	// Get the score (count) of each alternative using SWF
	count, err := SWF(p)
	if err != nil {
		return nil, err
	}

	// Use maxCount to find the alternatives with the highest score
	bestAlts = maxCount(count)

	// Return the best alternatives
	return bestAlts, nil
}
