package main

import (
	"io"
	"net/http"
	"project-4869/core"
	"project-4869/db"
	"github.com/gin-gonic/gin"
)

func main() {
	db.InitDB()

	r := gin.Default()
	r.Static("/static", "./static")
	r.LoadHTMLGlob("static/*")

	r.GET("/", func(c *gin.Context) {
		c.HTML(http.StatusOK, "index.html", nil)
	})

	// 实时日志
	r.GET("/api/logs", func(c *gin.Context) {
		ch := core.GetLogChan()
		c.Stream(func(w io.Writer) bool {
			if msg, ok := <-ch; ok {
				c.SSEvent("message", msg)
				return true
			}
			return false
		})
	})

	// 站点管理 API
	r.POST("/api/sites", func(c *gin.Context) {
		var site db.Site
		c.BindJSON(&site)
		db.DB.Create(&site)
		c.JSON(200, gin.H{"status": "ok"})
	})

	// 规则管理 API
	r.POST("/api/rules", func(c *gin.Context) {
		var rule db.Rule
		c.BindJSON(&rule)
		db.DB.Create(&rule)
		c.JSON(200, gin.H{"status": "ok"})
	})

	r.POST("/api/run", func(c *gin.Context) {
		go core.ProcessTask()
		c.JSON(200, gin.H{"message": "任务已启动"})
	})

	r.Run(":4869")
}