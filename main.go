package main

/*
 * 	License: GPL-3.0-or-later
 * 	Authors:
 * 		Mateus Melchiades <matbme@duck.com>
 * 	Copyright: 2023
 */

import (
	"errors"
	"fmt"
	"os"
	"os/exec"

	"github.com/gin-gonic/gin"
	"github.com/vanilla-os/differ/core"
	"github.com/vanilla-os/differ/core/handlers"
)

func setupRouter(dbPath string) (*gin.Engine, error) {
	// Create database if it doesn't exist
	if _, err := os.Stat(dbPath); os.IsNotExist(err) {
		dbUser, ok := os.LookupEnv("admin_user")
		if !ok {
			return nil, errors.New("admin_user environment variable not found")
		}

		dbPass, ok := os.LookupEnv("admin_password")
		if !ok {
			return nil, errors.New("admin_user environment variable not found")
		}

		err = exec.Command(
			"sh",
			"-c",
			fmt.Sprintf("sqlite3 %s 'create table auth(\"ID\" INTEGER, name, pass TEXT, PRIMARY KEY(\"ID\")); insert into auth values(1, \"%s\", \"%s\");'", dbPath, dbUser, dbPass),
		).Run()
		if err != nil {
			return nil, err
		}
	}

	// Initialize storage database
	err := core.InitStorage(dbPath)
	if err != nil {
		return nil, errors.New("Failed to init storage: " + err.Error())
	}

	// Initialize cache
	err = core.InitCache()
	if err != nil {
		return nil, errors.New("Failed to init cache: " + err.Error())
	}

	// Fetches authentications from storage
	auths, err := core.FetchAuthorizations()
	if err != nil {
		return nil, errors.New("Failed to fetch authorizations from storage: " + err.Error())
	}

	// If auths is empty, we run the API in "read-only" mode.
	// In other words, we won't be able to add any images or releases via the API.
	readOnly := len(auths) == 0

	var authRequired gin.HandlerFunc
	if !readOnly {
		authRequired = gin.BasicAuth(auths)
	}

	r := gin.Default()
	r.SetTrustedProxies(nil)

	// Endpoint to check if API is running
	r.GET("/status", handlers.HandleStatus)

	// Manipulate images
	images := r.Group("/images")
	{
		// List all images
		images.GET("/", handlers.HandleGetImages)
		// List specific image
		images.GET("/:name", handlers.HandleFindImage)
		// Creates new image (Auth required)
		if !readOnly {
			images.POST("/new", authRequired, handlers.HandleAddImage)
		}

		// Release-related endpoints
		// Diffs two releases
		images.GET("/:name/diff", handlers.HandleGetReleaseDiff)
		// Gets latest release
		images.GET("/:name/latest", handlers.HandleGetLatestRelease)
		// Gets specific release with digest
		images.GET("/:name/:digest", handlers.HandleFindRelease)
		// Creates new release (Auth required)
		if !readOnly {
			images.POST("/:name/new", authRequired, handlers.HandleAddRelease)
		}
	}

	return r, nil
}

func main() {
	var dbPath string
	if len(os.Args) > 1 {
		dbPath = os.Args[1]
	} else if dbPathArg, ok := os.LookupEnv("db_path"); ok {
		dbPath = dbPathArg
	} else {
		panic("No path to DB was provided. You must either pass it as a positional argument or by setting the 'db_path' environment variable")
	}

	router, err := setupRouter(dbPath)
	if err != nil {
		panic(err)
	}

	router.Run()
}
