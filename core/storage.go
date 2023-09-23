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

func InitStorage(storagePath string) (*gorm.DB, error) {
	db, err := gorm.Open(sqlite.Open(storagePath), &gorm.Config{})
	if err != nil {
		return nil, err
	}

	db.AutoMigrate(&types.Image{}, &types.Release{})

	return db, nil
}
