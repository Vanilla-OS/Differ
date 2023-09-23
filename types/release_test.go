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
			{Name: "pkg1", Version: "1.0"}, // Added
			{Name: "pkg3", Version: "2.0"}, // Upgraded
			{Name: "pkg4", Version: "1.0"}, // Downgraded
			{Name: "pkg5", Version: "1.0"}, // Not changed
		},
	}

	sampleOld := Release{
		Packages: []Package{
			{Name: "pkg2", Version: "1.0"}, // Removed
			{Name: "pkg3", Version: "1.0"},
			{Name: "pkg4", Version: "2.0"},
			{Name: "pkg5", Version: "1.0"},
		},
	}

	added, upgraded, downgraded, removed := sampleNew.DiffPackages(&sampleOld)

	if !slices.Equal(added, []PackageDiff{{"pkg1", "", "1.0"}}) {
		t.Fatalf("DiffPackages added = %v, expected {\"pkg1\", \"1.0\", \"\"}", added)
	}
	if !slices.Equal(upgraded, []PackageDiff{{"pkg3", "1.0", "2.0"}}) {
		t.Fatalf("DiffPackages upgraded = %v, expected {\"pkg3\", \"1.0\", \"2.0\"}", upgraded)
	}
	if !slices.Equal(downgraded, []PackageDiff{{"pkg4", "2.0", "1.0"}}) {
		t.Fatalf("DiffPackages downgraded = %v, expected {\"pkg4\", \"2.0\", \"1.0\"}", downgraded)
	}
	if !slices.Equal(removed, []PackageDiff{{"pkg2", "1.0", ""}}) {
		t.Fatalf("DiffPackages rmeoved = %v, expected {\"pkg2\", \"\", \"1.0\"}", removed)
	}
}
