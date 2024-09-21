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
