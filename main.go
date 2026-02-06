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

	c := cron.New(cron.WithSeconds())
	_, err := c.AddFunc("0 0 * * * *", func() {
		log.Println("执行定时抓取任务...")
		core.RunScraper()
	})
	if err != nil {
		log.Fatalf("Cron 任务添加失败: %v", err)
	}
	c.Start()

	r := gin.Default()
	r.Static("/static", "./static")
	r.LoadHTMLGlob("static/*")

	r.GET("/", func(c *gin.Context) {
		c.HTML(200, "index.html", gin.H{"title": "Project 4869"})
	})

	r.POST("/api/run", func(c *gin.Context) {
		go core.RunScraper()
		c.JSON(200, gin.H{"message": "任务已启动"})
	})

	log.Println("服务启动在 http://localhost:4869")
	if err := r.Run(":4869"); err != nil {
		log.Fatal("服务运行失败: ", err)
	}
}