package main

/*
 * 	License: GPL-3.0-or-later
 * 	Authors:
 * 		Mateus Melchiades <matbme@duck.com>
 * 	Copyright: 2023
 */

import (
	"github.com/gin-gonic/gin"
	"github.com/vanilla-os/differ/core"
	"github.com/vanilla-os/differ/core/handlers"
)

func main() {
	err := core.InitStorage("test.db")
	if err != nil {
		panic("Failed to init storage: " + err.Error())
	}

	r := gin.Default()

	// Endpoint to check if API is running
	r.GET("/status", handlers.HandleStatus)

	// Manipulate images
	r.GET("/images", handlers.HandleGetImages)
	r.POST("/images/new", handlers.HandleAddImage)
	r.GET("/images/:name", handlers.HandleFindImage)

	// Manipulate releases

	r.Run()
}
