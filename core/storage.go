package core

/*
 * 	License: GPL-3.0-or-later
 * 	Authors:
 * 		Mateus Melchiades <matbme@duck.com>
 * 	Copyright: 2023
 */

import (
	"github.com/vanilla-os/differ/types"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

var DB *gorm.DB

func InitStorage(storagePath string) error {
	db, err := gorm.Open(sqlite.Open(storagePath), &gorm.Config{})
	if err != nil {
		return err
	}

	db.AutoMigrate(&types.Image{}, &types.Release{})

	DB = db

	return nil
}
