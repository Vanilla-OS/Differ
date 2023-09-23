package types

/*
 * 	License: GPL-3.0-or-later
 * 	Authors:
 * 		Mateus Melchiades <matbme@duck.com>
 * 	Copyright: 2023
 */

import (
	"slices"
	"testing"
)

func TestDiffPackages(t *testing.T) {
	sampleNew := Release{
		Packages: []Package{
			{"pkg1", "1.0"}, // Added
			{"pkg3", "1.0"}, // Upgraded
			{"pkg4", "2.0"}, // Downgraded
			{"pkg5", "1.0"}, // Not changed
		},
	}

	sampleOld := Release{
		Packages: []Package{
			{"pkg2", "1.0"}, // Removed
			{"pkg3", "2.0"},
			{"pkg4", "1.0"},
			{"pkg5", "1.0"},
		},
	}

	added, upgraded, downgraded, removed := sampleNew.DiffPackages(&sampleOld)

	if !slices.Equal(added, []PackageDiff{{"pkg1", "1.0", ""}}) {
		t.Fatalf("DiffPackages added = %v, expected {\"pkg1\", \"1.0\", \"\"}", added)
	}
	if !slices.Equal(upgraded, []PackageDiff{{"pkg3", "1.0", "2.0"}}) {
		t.Fatalf("DiffPackages upgraded = %v, expected {\"pkg3\", \"1.0\", \"2.0\"}", upgraded)
	}
	if !slices.Equal(downgraded, []PackageDiff{{"pkg4", "2.0", "1.0"}}) {
		t.Fatalf("DiffPackages downgraded = %v, expected {\"pkg4\", \"2.0\", \"1.0\"}", downgraded)
	}
	if !slices.Equal(removed, []PackageDiff{{"pkg2", "", "1.0"}}) {
		t.Fatalf("DiffPackages rmeoved = %v, expected {\"pkg2\", \"\", \"1.0\"}", removed)
	}
}
