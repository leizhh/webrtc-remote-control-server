package main

import (
	"github.com/gin-gonic/gin"
	"io/ioutil"
	"net/http"
	"webrtc-remote-control-server/server"
)

func main() {
	router := gin.Default()
	router.Use(Cors())
	router.LoadHTMLGlob("templates/*")
	router.Static("/public/img/", "./public/img/")
	server.InitDeviceHub()

	router.GET("/ping", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"message": "pong",
		})
	})

	router.GET("/index", func(c *gin.Context) {
		c.HTML(http.StatusOK, "index.html", gin.H{
			"nav": "nav_home",
		})
	})

	router.GET("/doc", func(c *gin.Context) {
		doc, _ := ioutil.ReadFile("README.md")
		c.HTML(http.StatusOK, "doc.html", gin.H{
			"nav": "nav_doc",
			"doc": string(doc),
		})
	})

	router.GET("/answer", server.AnswerHandler)
	router.GET("/offer", server.OfferHandler)
	router.GET("/devices", server.GetDevices)

	router.Run("0.0.0.0:8080")
}

func Cors() gin.HandlerFunc {
	return func(c *gin.Context) {
		method := c.Request.Method

		c.Header("Access-Control-Allow-Origin", "*")
		c.Header("Access-Control-Allow-Headers", "Content-Type,AccessToken,X-CSRF-Token, Authorization, Token")
		c.Header("Access-Control-Allow-Methods", "POST, GET, OPTIONS")
		c.Header("Access-Control-Expose-Headers", "Content-Length, Access-Control-Allow-Origin, Access-Control-Allow-Headers, Content-Type")
		c.Header("Access-Control-Allow-Credentials", "true")

		//放行所有OPTIONS方法
		if method == "OPTIONS" {
			c.AbortWithStatus(http.StatusNoContent)
		}
		// 处理请求
		c.Next()
	}
}
