package main

/*
 * 	License: GPL-3.0-or-later
 * 	Authors:
 * 		Mateus Melchiades <matbme@duck.com>
 * 	Copyright: 2023
 */

import (
	"github.com/gin-gonic/gin"
	"github.com/vanilla-os/differ/core/handlers"
)

func main() {
	r := gin.Default()

	r.GET("/status", handlers.HandleStatus)

	r.Run()
}
