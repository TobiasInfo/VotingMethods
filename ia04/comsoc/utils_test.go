package comsoc

import "testing"

func TestRank(t *testing.T) {
	alt := Alternative(1)
	prefs := []Alternative{1, 2, 3}
	if rank(alt, prefs) != 0 {
		t.Error("rank(1, [1, 2, 3]) should return 0")
	}
	alt2 := Alternative(2)
	if rank(alt2, prefs) != 1 {
		t.Error("rank(2, [1, 2, 3]) should return 1")
	}
	alt3 := Alternative(5)
	if rank(alt3, prefs) != -1 {
		t.Error("rank(5, [1, 2, 3]) should return -1")
	}
}

func TestIsPref(t *testing.T) {
	alt1 := Alternative(1)
	alt2 := Alternative(2)
	prefs := []Alternative{1, 2, 3}
	if !isPref(alt1, alt2, prefs) {
		t.Error("1 should be preferred to 2 in [1, 2, 3]")
	}
	if isPref(alt2, alt1, prefs) {
		t.Error("2 should not be preferred to 1 in [1, 2, 3]")
	}
}

func TestMaxCount(t *testing.T) {
	count := Count{1: 1, 2: 1, 3: 2}
	bestAlts := maxCount(count)
	if len(bestAlts) != 1 || bestAlts[0] != 3 {
		t.Error("The best alternative should be 3")
	}

	count2 := Count{1: 1, 2: 2, 3: 1, 4: 2}
	bestAlts2 := maxCount(count2)
	if len(bestAlts2) != 2 || !Contains(bestAlts2, 2) || !Contains(bestAlts2, 4) {
		t.Error("The best alternatives should be 2 and 4")
	}
}

func TestCheckProfile(t *testing.T) {
	alts := []Alternative{1, 2, 3}
	prefs := []Alternative{1, 2, 3}
	if err := checkProfile(prefs, alts, true); err != nil {
		t.Error("The profile should be valid")
	}

	prefs = []Alternative{1, 2, 4}
	if err := checkProfile(prefs, alts, true); err == nil {
		t.Error("The profile should be invalid")
	}

	prefs = []Alternative{1, 2, 3, 4}
	if err := checkProfile(prefs, alts, true); err == nil {
		t.Error("The profile should be invalid")
	}
}

func TestCheckProfileAlternative(t *testing.T) {
	alts := []Alternative{1, 2, 3}

	p := Profile{{1, 2, 3}, {1, 2, 3}}
	if err := checkProfileAlternative(p, alts, true); err != nil {
		t.Error("The alternative should be valid")
	}

	p = Profile{{1, 2}, {1, 2, 3}}
	if err := checkProfileAlternative(p, alts, true); err == nil {
		t.Error("The alternative should be invalid")
	}

	p = Profile{{1, 2, 4}, {1, 2, 3}}
	if err := checkProfileAlternative(p, alts, true); err == nil {
		t.Error("The alternative should be invalid")
	}
}

// TODD Add tests for the SWF and SCF functions
