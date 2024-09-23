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
// pas besoin de already seen vu que c est 1 slice donc 1 passage et pas de doublons
func isPref(alt1, alt2 Alternative, prefs []Alternative) bool {
	for _, a := range prefs {
		if a == alt1 {
			return true
		}
		if a == alt2 {
			return false
		}
	}
	return false
	//return rank(alt1, prefs) < rank(alt2, prefs) je laisse ca en cas de pb
}

// renvoie les meilleures alternatives pour un décomtpe donné
// un seul passage apr la boucle est suffisant on reinitialise bestAlts a chaque fois qu'on trouve un meilleur score
func maxCount(count Count) (bestAlts []Alternative) {
	bestAlts = make([]Alternative, 0)
	var max int
	for alt, c := range count {
		if c > max {
			max = c
			bestAlts = []Alternative{alt}
		} else if c == max {
			bestAlts = append(bestAlts, alt)
		}
	}
	return bestAlts
}

func checkProfile(prefs []Alternative, alts []Alternative) error {
	// Vérifier que prefs est complet
	if len(prefs) != len(alts) {
		return fmt.Errorf("le profil n'est pas complet")
	}
	// Vérifier que chaque alternative n'apparaît qu'une seule fois
	seen := make(map[Alternative]bool)
	for _, alt := range prefs {
		if seen[alt] {
			return fmt.Errorf("l'alternative %v apparaît plusieurs fois", alt)
		}
		seen[alt] = true
	}
	// Vérifier que chaque alternative de alts apparaît exactement une fois dans prefs
	for _, alt := range alts {
		if !seen[alt] {
			return fmt.Errorf("l'alternative %v n'apparaît pas dans le profil", alt)
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
func RecoverAlts(p Profile) []Alternative {
	alts := make([]Alternative, 0)
	for _, prefs := range p {
		for _, alt := range prefs {
			if !Contains(alts, alt) {
				alts = append(alts, alt)
			}
		}
	}
	return alts
}
