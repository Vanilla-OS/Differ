package handlers

/*
 * 	License: GPL-3.0-or-later
 * 	Authors:
 * 		Mateus Melchiades <matbme@duck.com>
 * 	Copyright: 2023
 */

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"time"

	"github.com/bytedance/sonic"
	"github.com/gin-gonic/gin"
	"github.com/vanilla-os/differ/core"
	"github.com/vanilla-os/differ/diff"
	"github.com/vanilla-os/differ/types"
	"gorm.io/gorm"
)

func HandleGetLatestRelease(c *gin.Context) {
	imageName := c.Param("name")
	image, err := types.GetImageByName(core.DB, imageName)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"release": image.GetLatestRelease()})
}

func HandleFindRelease(c *gin.Context) {
	imageName := c.Param("name")
	image, err := types.GetImageByName(core.DB, imageName)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	releaseDigest := c.Param("digest")
	release, err := image.GetReleaseByDigest(core.DB, releaseDigest)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"release": release})
}

func HandleAddRelease(c *gin.Context) {
	imageName := c.Param("name")
	image, err := types.GetImageByName(core.DB, imageName)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	var releaseInput struct {
		Digest   string          `json:"digest" binding:"required"`
		Date     time.Time       `json:"date"`
		Packages []types.Package `json:"packages" binding:"required"`
	}
	if err := c.ShouldBindJSON(&releaseInput); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if releaseInput.Date.IsZero() {
		releaseInput.Date = time.Now()
	}

	newRelease, err := image.NewRelease(core.DB, &types.Release{
		Digest:   releaseInput.Digest,
		ImageID:  image.ID,
		Date:     releaseInput.Date,
		Packages: releaseInput.Packages,
	})
	if err != nil {
		errorCode := http.StatusInternalServerError
		if errors.Is(err, gorm.ErrDuplicatedKey) {
			errorCode = http.StatusBadRequest
		}
		c.JSON(errorCode, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"release": newRelease})
}

func HandleGetReleaseDiff(c *gin.Context) {
	var diffInput struct {
		OldDigest string `json:"old_digest" binding:"required"`
		NewDigest string `json:"new_digest" binding:"required"`
	}
	if err := c.ShouldBindJSON(&diffInput); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	cacheKey := fmt.Sprintf("%s-%s", diffInput.OldDigest, diffInput.NewDigest)
	cacheDiff, _ := core.CacheManager.Get(context.Background(), cacheKey)

	// Cache hit
	if cacheDiff != nil {
		var diff struct {
			OldDigest                            string `json:"_old_digest"`
			NewDigest                            string `json:"_new_digest"`
			Added, Upgraded, Downgraded, Removed []diff.PackageDiff
		}
		err := sonic.Unmarshal(cacheDiff, &diff)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}

		diff.OldDigest = diffInput.OldDigest
		diff.NewDigest = diffInput.NewDigest

		c.JSON(http.StatusOK, diff)
		return
	}

	// Diff not in cache. Generate diff and store for future queries.
	imageName := c.Param("name")
	image, err := types.GetImageByName(core.DB, imageName)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	oldRelease, err := image.GetReleaseByDigest(core.DB, diffInput.OldDigest)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	newRelease, err := image.GetReleaseByDigest(core.DB, diffInput.NewDigest)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	added, upgraded, downgraded, removed := oldRelease.DiffPackages(newRelease)

	cacheDiffEntry := struct {
		Added, Upgraded, Downgraded, Removed []diff.PackageDiff
	}{
		Added:      added,
		Upgraded:   upgraded,
		Downgraded: downgraded,
		Removed:    removed,
	}
	cacheBytes, err := sonic.Marshal(cacheDiffEntry)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	err = core.CacheManager.Set(context.Background(), cacheKey, cacheBytes)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"_old_digest": diffInput.OldDigest,
		"_new_digest": diffInput.NewDigest,
		"added":       added,
		"upgraded":    upgraded,
		"downgraded":  downgraded,
		"removed":     removed,
	})
}
