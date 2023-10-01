package handlers

/*
 * 	License: GPL-3.0-or-later
 * 	Authors:
 * 		Mateus Melchiades <matbme@duck.com>
 * 	Copyright: 2023
 */

import (
	"errors"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/vanilla-os/differ/core"
	"github.com/vanilla-os/differ/types"
	"gorm.io/gorm"
)

func HandleGetImages(c *gin.Context) {
	images, err := types.GetImages(core.DB)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"images": images})
}

func HandleFindImage(c *gin.Context) {
	image, err := types.GetImageByName(core.DB, c.Param("name"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"image": image})
}

func HandleAddImage(c *gin.Context) {
	var imageInput struct {
		Name     string `json:"name" binding:"required"`
		URL      string `json:"url" binding:"required"`
		Releases []types.Release
	}
	if err := c.ShouldBindJSON(&imageInput); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if len(imageInput.Releases) == 0 {
		imageInput.Releases = []types.Release{}
	}

	newImage := types.Image{
		Name:     imageInput.Name,
		URL:      imageInput.URL,
		Releases: imageInput.Releases,
	}
	if status := core.DB.Create(&newImage); status.Error != nil {
		errorCode := http.StatusInternalServerError
		if errors.Is(status.Error, gorm.ErrDuplicatedKey) {
			errorCode = http.StatusBadRequest
		}
		c.JSON(errorCode, gin.H{"error": status.Error.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"image": newImage})
}
