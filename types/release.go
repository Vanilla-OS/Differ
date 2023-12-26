package types

/*
 * 	License: GPL-3.0-or-later
 * 	Authors:
 * 		Mateus Melchiades <matbme@duck.com>
 * 	Copyright: 2023
 */

import (
	"time"

	"github.com/vanilla-os/differ/diff"
	"gorm.io/gorm"
)

type Package struct {
	gorm.Model `json:"-"`
	Name       string `json:"name"`
	Version    string `json:"version"`
}

type Release struct {
	gorm.Model `json:"-"`
	Digest     string    `json:"digest" gorm:"unique"`
	ImageID    uint      `json:"-"` // foreign key for Image
	Date       time.Time `json:"date"`
	Packages   []Package `json:"packages,omitempty" gorm:"many2many:release_packages;"`
}

func (re *Release) DiffPackages(other *Release) ([]diff.PackageDiff, []diff.PackageDiff, []diff.PackageDiff, []diff.PackageDiff) {
	thisPackagesMap := make(diff.Package, len(re.Packages))
	for _, pkg := range re.Packages {
		thisPackagesMap[pkg.Name] = pkg.Version
	}

	otherPackagesMap := make(diff.Package, len(other.Packages))
	for _, pkg := range other.Packages {
		otherPackagesMap[pkg.Name] = pkg.Version
	}

	return diff.DiffPackages(thisPackagesMap, otherPackagesMap)
}
