package types

/*
 * 	License: GPL-3.0-or-later
 * 	Authors:
 * 		Mateus Melchiades <matbme@duck.com>
 * 	Copyright: 2023
 */

import (
	"cmp"
	"regexp"
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
	Name       string `json:"name"`
	OldVersion string `json:"old_version,omitempty"`
	NewVersion string `json:"new_version,omitempty"`
}

type Release struct {
	gorm.Model `json:"-"`
	Digest     string    `json:"digest" gorm:"unique"`
	ImageID    uint      `json:"-"` // foreign key for Image
	Date       time.Time `json:"date"`
	Packages   []Package `json:"packages" gorm:"many2many:release_packages;"`
}

// This monstruosity is an adaptation of the regex for semver (available in https://semver.org/).
// It SHOULD be able to capture every type of exoteric versioning scheme out there.
var versionRegex = regexp.MustCompile(`^(?:(?P<prefix>\d+):)?(?P<major>\d+[a-zA-Z]?)(?:\.(?P<minor>\d+))?(?:\.(?P<patch>\d+))?(?:[-~](?P<prerelease>(?:\d+|\d*[a-zA-Z-][0-9a-zA-Z-]*)(?:\.(?:\d+|\d*[a-zA-Z-][0-9a-zA-Z-]*))*))?(?:[+.](?P<buildmetadata>[0-9a-zA-Z-+.]+(?:\.[0-9a-zA-Z-]+)*))?$`)

// compareVersions has the same behavior as cmp.Compare, but for package versions. It parses
// both version strings and checks for differences in major, minor, patch, pre-release, etc.
func compareVersions(a, b string) int {
	aMatchStr := versionRegex.FindStringSubmatch(a)
	aMatch := make(map[string]string)
	for i, name := range versionRegex.SubexpNames() {
		if i != 0 && name != "" && aMatchStr[i] != "" {
			aMatch[name] = aMatchStr[i]
		}
	}

	bMatchStr := versionRegex.FindStringSubmatch(b)
	bMatch := make(map[string]string)
	for i, name := range versionRegex.SubexpNames() {
		if i != 0 && name != "" && bMatchStr[i] != "" {
			bMatch[name] = bMatchStr[i]
		}
	}

	compResult := 0

	compOrder := []string{"prefix", "major", "minor", "patch", "prerelease", "buildmetadata"}
	for _, comp := range compOrder {
		aValue, aOk := aMatch[comp]
		bValue, bOk := bMatch[comp]
		// If neither version has component or if they equal
		if !aOk && !bOk {
			continue
		}
		// If a has component but b doesn't, package was upgraded, unless it's prerelease
		if aOk && !bOk {
			if comp == "prerelease" {
				compResult = -1
			} else {
				compResult = 1
			}
			break
		}
		// If b has component but a doesn't, package was downgraded
		if !aOk && bOk {
			compResult = -1
			break
		}

		// If both have, do regular compare
		abComp := cmp.Compare(aValue, bValue)
		if abComp == 0 {
			continue
		}
		compResult = abComp
		break
	}

	return compResult
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
			diff := PackageDiff{pkg.Name, otherCopy[pos].Version, pkg.Version}
			switch compareVersions(pkg.Version, otherCopy[pos].Version) {
			case -1:
				downgraded = append(downgraded, diff)
			case 1:
				upgraded = append(upgraded, diff)
			}

			// Clear package from copy so we can later check for removed packages
			otherCopy[pos] = Package{}
		} else {
			diff := PackageDiff{pkg.Name, "", pkg.Version}
			added = append(removed, diff)
		}
	}

	for _, opkg := range otherCopy {
		dummy := Package{}
		if opkg != dummy {
			diff := PackageDiff{opkg.Name, opkg.Version, ""}
			removed = append(removed, diff)
		}
	}

	return added, upgraded, downgraded, removed
}
