package main

import (
	"log"
	"project-4869/core"
	"project-4869/db"
	"github.com/gin-gonic/gin"
	"github.com/robfig/cron/v3"
)

func main() {
	db.InitDB()

	// 启动定时任务
	c := cron.New(cron.WithSeconds())
	c.AddFunc("0 0 * * * *", func() {
		core.AddLog("系统提示: 触发定时抓取任务")
		core.RunScraper()
	})
	c.Start()

	r := gin.Default()
	r.Static("/static", "./static")
	r.LoadHTMLGlob("static/*")

	r.GET("/", func(c *gin.Context) {
		c.HTML(200, "index.html", nil)
	})

	// 实时日志流 (SSE)
	r.GET("/api/logs", func(c *gin.Context) {
		c.Header("Content-Type", "text/event-stream")
		c.Header("Cache-Control", "no-cache")
		c.Header("Connection", "keep-alive")
		
		msgChan := core.GetLogChan()
		c.Stream(func(w interface{}) bool {
			if msg, ok := <-msgChan; ok {
				c.SSEvent("message", msg)
				return true
			}
			return false
		})
	})

	// 保存配置
	r.POST("/api/config", func(c *gin.Context) {
		var config db.SystemConfig
		if err := c.ShouldBindJSON(&config); err != nil {
			c.JSON(400, gin.H{"error": err.Error()})
			return
		}
		db.SaveConfig(config)
		core.AddLog("系统提示: 配置已更新")
		c.JSON(200, gin.H{"message": "保存成功"})
	})

	r.POST("/api/run", func(c *gin.Context) {
		go core.RunScraper()
		c.JSON(200, gin.H{"message": "任务已启动"})
	})

	log.Println("Project 4869 运行在 http://localhost:4869")
	r.Run(":4869")
}