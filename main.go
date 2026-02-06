package main

import (
	"log"
	"project-4869/core"
	"project-4869/db"

	"github.com/gin-gonic/gin"
	"github.com/robfig/cron/v3"
)

func main() {
	// 1. 初始化数据库
	db.InitDB()

	// 2. 启动定时任务
	c := cron.New(cron.WithSeconds())
	c.AddFunc("0 0 * * * *", func() {
		log.Println("执行定时抓取任务...")
		core.RunScraper()
	})
	c.Start()

	// 3. 设置路由
	r := gin.Default()
	r.Static("/static", "./static")
	r.LoadHTMLGlob("static/*")

	r.GET("/", func(c *gin.Context) {
		c.HTML(200, "index.html", gin.H{"title": "Project 4869"})
	})

	// 接口示例
	r.POST("/api/run", func(c *gin.Context) {
		go core.RunScraper()
		c.JSON(200, gin.H{"message": "任务已启动"})
	})

	log.Println("服务启动在 http://localhost:4869")
	r.Run(":4869")
}