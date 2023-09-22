package handlers

import "github.com/gin-gonic/gin"

/*
 * 	License: GPL-3.0-or-later
 * 	Authors:
 * 		Mateus Melchiades <matbme@duck.com>
 * 	Copyright: 2023
 */

func HandleGetImages(c *gin.Context) {
	// TODO: Get list of all images being tracked
}

func HandleGetLatestRelease(c *gin.Context) {
	// TODO: Get latest release from image
}

func HandleGetRelease(c *gin.Context) {
	// TODO: Get specific release from image
}

func HandleGetReleaseDiff(c *gin.Context) {
	// TODO: Get diff betwenn two releases
}
