package main

import (
	"awesomeProject1/scraper"
	"awesomeProject1/status"
	"encoding/json"
	"fmt"
	"log/slog"
	"net/http"
	"os"
	"sync"

	"github.com/gin-gonic/gin"
)

var (
	totalURLs     int
	processedURLs int
	statusMutex   sync.Mutex
)

type ScrapeRequest struct {
	URLs []string `json:"urls" binding:"required"`
}

func RunServer() {
	r := gin.Default()

	r.POST("/scrape", func(c *gin.Context) {
		var req ScrapeRequest
		if err := c.ShouldBindJSON(&req); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}
		statusMutex.Lock()
		totalURLs += len(req.URLs)
		statusMutex.Unlock()

		results := scraper.Run(req.URLs, func() {
			statusMutex.Lock()
			processedURLs++
			statusMutex.Unlock()
		})

		file, _ := os.Create("results.json")
		defer file.Close()
		encoder := json.NewEncoder(file)
		encoder.SetIndent("", "	")
		encoder.Encode(results)

		c.JSON(http.StatusOK, results)
	})

	r.GET("/status", func(c *gin.Context) {
		statusMutex.Lock()
		defer statusMutex.Unlock()

		url := c.Query("url")
		slog.Info(url)
		if url == "" {
			c.JSON(http.StatusOK, gin.H{
				"totalURLs":     totalURLs,
				"processedURLs": processedURLs,
				"pending":       totalURLs - processedURLs,
				"URLS":          status.GetAll(),
			})
		}
		c.JSON(http.StatusOK, gin.H{
			"status":        status.GetStatus(url),
			"time of start": status.GetTimeToAction(url),
			"duration":      status.GetDuration(url),
		})
	})

	r.GET("/version", func(c *gin.Context) {
		v := os.Getenv("VERSION")
		fmt.Println("VERSIoooon", v)
		c.JSON(http.StatusOK, gin.H{
			"version": v,
		})
	})

	r.GET("/health", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{
			"status": "ok",
		})
	})

	r.Run(":8080")
}
