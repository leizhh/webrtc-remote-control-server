package main

import (
	"net/http"
	"github.com/gin-gonic/gin"
	"webrtc-server/server"
)

func main() {
	router := gin.Default()
	router.LoadHTMLGlob("templates/*")
	router.Static("/public/img/", "./public/img/")
	server.InitKChan()

	router.GET("/ping", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"message": "pong",
		})
	})

	router.GET("/index", func(c *gin.Context) {
		c.HTML(http.StatusOK, "index.html", gin.H{
			"title": "WebRTC",
		})
	})

	router.GET("/answer",server.AnswerHandler)
	router.GET("/offer",server.OfferHandler)
	router.GET("/devices",server.GetDevices)

	router.Run("0.0.0.0:8080") 
}