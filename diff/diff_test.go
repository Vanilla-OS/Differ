package diff

import (
	"fmt"
	"slices"
	"testing"
)

func TestDiffPackages(t *testing.T) {
	newPackages := Package{
		"libei1":                   "1.2.1-1",
		"xdg-desktop-portal-gnome": "46.2-1",
		"libdisplay-info1":         "0.1.1-1",
	}

	oldPackages := Package{
		"libwinpr2-2t64":           "2.11.5+dfsg1-1",
		"xdg-desktop-portal-gnome": "44.2-4+b1",
		"libdisplay-info1":         "0.1.1-2+b1",
	}

	expected := struct {
		added, upgraded, downgraded, removed []PackageDiff
	}{
		[]PackageDiff{{Name: "libei1", NewVersion: "1.2.1-1"}},
		[]PackageDiff{{Name: "xdg-desktop-portal-gnome", NewVersion: "46.2-1", PreviousVersion: "44.2-4+b1"}},
		[]PackageDiff{{Name: "libdisplay-info1", NewVersion: "0.1.1-1", PreviousVersion: "0.1.1-2+b1"}},
		[]PackageDiff{{Name: "libwinpr2-2t64", PreviousVersion: "2.11.5+dfsg1-1"}},
	}

	added, upgraded, downgraded, removed := DiffPackages(oldPackages, newPackages)

	if !slices.Equal(expected.added, added) {
		fmt.Printf("Incorrect parsing of added. Expected %v, got %v\n", expected.added, added)
		t.Fail()
	}
	if !slices.Equal(expected.upgraded, upgraded) {
		fmt.Printf("Incorrect parsing of upgraded. Expected %v, got %v\n", expected.upgraded, upgraded)
		t.Fail()
	}
	if !slices.Equal(expected.downgraded, downgraded) {
		fmt.Printf("Incorrect parsing of downgraded. Expected %v, got %v\n", expected.downgraded, downgraded)
		t.Fail()
	}
	if !slices.Equal(expected.removed, removed) {
		fmt.Printf("Incorrect parsing of removed. Expected %v, got %v\n", expected.removed, removed)
		t.Fail()
	}
}
