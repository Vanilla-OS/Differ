package types

/*
 * 	License: GPL-3.0-or-later
 * 	Authors:
 * 		Mateus Melchiades <matbme@duck.com>
 * 	Copyright: 2023
 */

import (
	"cmp"
	"slices"
	"time"

	"gorm.io/gorm"
)

type Package struct {
	gorm.Model `json:"-"`
	Name       string `json:"name"`
	Version    string `json:"version"`
}

type PackageDiff struct {
	Name                   string
	OldVersion, NewVersion string
}

type Release struct {
	gorm.Model `json:"-"`
	Digest     string    `json:"digest"`
	ImageID    uint      `json:"-"` // foreign key for Image
	Date       time.Time `json:"date"`
	Packages   []Package `json:"packages" gorm:"many2many:release_packages;"`
}

// PackageDiff returns the difference in packages between two images, organized into
// four slices: Added, Upgraded, Downgraded, and Removed packages, respectively.
func (re *Release) DiffPackages(other *Release) ([]PackageDiff, []PackageDiff, []PackageDiff, []PackageDiff) {
	added := []PackageDiff{}
	upgraded := []PackageDiff{}
	downgraded := []PackageDiff{}
	removed := []PackageDiff{}

	otherCopy := make([]Package, len(other.Packages))
	copy(otherCopy, other.Packages)

	for _, pkg := range re.Packages {
		pos := slices.IndexFunc(otherCopy, func(n Package) bool { return n.Name == pkg.Name })
		if pos != -1 {
			diff := PackageDiff{pkg.Name, pkg.Version, otherCopy[pos].Version}
			switch cmp.Compare(pkg.Version, otherCopy[pos].Version) {
			case -1:
				upgraded = append(upgraded, diff)
			case 1:
				downgraded = append(downgraded, diff)
			}

			// Clear package from copy so we can later check for removed packages
			otherCopy[pos] = Package{}
		} else {
			diff := PackageDiff{pkg.Name, pkg.Version, ""}
			added = append(removed, diff)
		}
	}

	for _, opkg := range otherCopy {
		dummy := Package{}
		if opkg != dummy {
			diff := PackageDiff{opkg.Name, "", opkg.Version}
			removed = append(removed, diff)
		}
	}

	return added, upgraded, downgraded, removed
}
