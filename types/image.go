package types

import (
	"sort"

	"gorm.io/gorm"
)

/*
 * 	License: GPL-3.0-or-later
 * 	Authors:
 * 		Mateus Melchiades <matbme@duck.com>
 * 	Copyright: 2023
 */

type Image struct {
	gorm.Model
	Name     string
	URL      string
	Releases []Release
}

func GetImages(db *gorm.DB) ([]Image, error) {
	var images []Image
	err := db.Model(&Image{}).Preload("Releases").Find(&images).Error
	return images, err
}

func GetImageByName(db *gorm.DB, name string) (Image, error) {
	var image Image
	err := db.First(&image, "name = ?", name).Error
	return image, err
}

func (im *Image) GetLatestRelease() *Release {
	sort.Slice(im.Releases, func(i, j int) bool {
		return im.Releases[i].Date.After(im.Releases[j].Date)
	})

	return &im.Releases[0]
}

func (im *Image) NewRelease(db *gorm.DB, release *Release) error {
	return db.Create(release).Error
}
