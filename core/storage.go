package core

/*
 * 	License: GPL-3.0-or-later
 * 	Authors:
 * 		Mateus Melchiades <matbme@duck.com>
 * 	Copyright: 2023
 */

import (
	"errors"
	"fmt"

	"github.com/gin-gonic/gin"
	"github.com/vanilla-os/differ/types"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

var DB *gorm.DB

// InitStorage initializes the databse from a given path, creating it if necessary.
func InitStorage(storagePath string) error {
	db, err := gorm.Open(sqlite.Open(storagePath), &gorm.Config{TranslateError: true})
	if err != nil {
		return err
	}

	db.AutoMigrate(&types.Image{}, &types.Release{})

	DB = db

	return nil
}

// FetchAuthorizations fetches from the database a Gin Accounts map for using in BasicAuth requests.
// Authorizations must be added to the database beforehand in a table called `auth` with the following rows:
//   - ID: int PK
//   - name: text
//   - pass: text
//
// The following statement can be used to create the auth table:
//
//	CREATE TABLE "auth" (
//		"ID"	INTEGER,
//		"name"	TEXT,
//		"pass"	TEXT,
//		PRIMARY KEY("ID")
//	)
func FetchAuthorizations() (gin.Accounts, error) {
	// Assert database has been initialized
	if DB == nil {
		return nil, errors.New("db has not been initialized yet. You must call InitStorage first")
	}

	type AuthResult struct {
		Name, Pass string
	}

	var result []AuthResult
	err := DB.Raw("SELECT name, pass FROM auth").Scan(&result).Error
	if err != nil {
		return nil, err
	}

	auths := gin.Accounts{}
	for _, auth := range result {
		auths[auth.Name] = auth.Pass
	}

	// Warn if auth database is empty.
	// This isn't necessarily an error but we won't be able to add any images or releases via API.
	if len(auths) == 0 {
		fmt.Println("\033[1;33mWARN:\033[0m No authorized users found in database, starting API in read-only mode. If you intend to use the API for manipulating images/releases, please make sure to have at least one entry in the `auth` table.")
	}

	return auths, nil
}
