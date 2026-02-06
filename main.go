package main

import (
	"io"
	"log"
	"net/http"
	"project-4869/core"
	"project-4869/db"

	"github.com/gin-gonic/gin"
	"github.com/robfig/cron/v3"
)

func main() {
	db.InitDB()

	c := cron.New(cron.WithSeconds())
	c.AddFunc("0 0 * * * *", func() {
		core.AddLog("系统消息: 定时任务触发")
		core.RunScraper()
	})
	c.Start()

	r := gin.Default()
	r.Static("/static", "./static")
	r.LoadHTMLGlob("static/*")

	r.GET("/", func(c *gin.Context) {
		c.HTML(http.StatusOK, "index.html", nil)
	})

	// 实时日志通道 (SSE) - 已修正 io.Writer 类型
	r.GET("/api/logs", func(c *gin.Context) {
		c.Header("Content-Type", "text/event-stream")
		c.Header("Cache-Control", "no-cache")
		c.Header("Connection", "keep-alive")
		
		ch := core.GetLogChan()
		c.Stream(func(w io.Writer) bool {
			if msg, ok := <-ch; ok {
				c.SSEvent("message", msg)
				return true
			}
			return false
		})
	})

	r.POST("/api/config", func(c *gin.Context) {
		var cfg db.SystemConfig
		if err := c.ShouldBindJSON(&cfg); err != nil {
			c.JSON(400, gin.H{"status": "error"})
			return
		}
		db.SaveConfig(cfg)
		core.AddLog("系统消息: 配置已成功保存")
		c.JSON(200, gin.H{"status": "ok"})
	})

	r.POST("/api/run", func(c *gin.Context) {
		go core.RunScraper()
		c.JSON(200, gin.H{"message": "任务已启动"})
	})

	r.Run(":4869")
}