package types

/*
 * 	License: GPL-3.0-or-later
 * 	Authors:
 * 		Mateus Melchiades <matbme@duck.com>
 * 	Copyright: 2023
 */

import (
	"fmt"
	"sort"

	"gorm.io/gorm"
)

type Image struct {
	gorm.Model `json:"-"`
	Name       string    `json:"name" gorm:"unique"`
	URL        string    `json:"url" gorm:"unique"`
	Releases   []Release `json:"releases"`
}

func GetImages(db *gorm.DB) ([]Image, error) {
	var images []Image
	err := db.Model(&Image{}).Preload("Releases").Preload("Releases.Packages").Find(&images).Error
	return images, err
}

func GetImageByName(db *gorm.DB, name string) (Image, error) {
	var image Image
	err := db.Preload("Releases").Preload("Releases.Packages").First(&image, "name = ?", name).Error
	if err == gorm.ErrRecordNotFound {
		return image, fmt.Errorf("no image found with name %s", name)
	}

	return image, err
}

func (im *Image) GetLatestRelease() *Release {
	if len(im.Releases) == 0 {
		return nil
	}

	sort.Slice(im.Releases, func(i, j int) bool {
		return im.Releases[i].Date.After(im.Releases[j].Date)
	})

	return &im.Releases[0]
}

func (im *Image) GetReleaseByDigest(db *gorm.DB, digest string) (*Release, error) {
	if len(im.Releases) == 0 {
		return nil, fmt.Errorf("no release found with digest %s", digest)
	}

	var release Release
	err := db.Preload("Packages").First(&release, "digest = ?", digest).Error
	if err == gorm.ErrRecordNotFound {
		return nil, fmt.Errorf("no release found with digest %s", digest)
	}

	return &release, err
}

func (im *Image) NewRelease(db *gorm.DB, release *Release) (*Release, error) {
	status := db.Create(release)
	if status.Error != nil {
		return nil, status.Error
	}

	var newRelease Release
	err := db.Preload("Packages").First(&newRelease, release.ID).Error
	if status.Error != nil {
		return nil, err
	}

	return &newRelease, nil
}
