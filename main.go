package main

import (
	"bufio"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"project-4869/core"
	"project-4869/db"
	"strconv"
	"sync"

	"github.com/gin-gonic/gin"
	"github.com/robfig/cron/v3"
)

// 全局变量管理 Cron
var (
	c           *cron.Cron
	rssEntryID  cron.EntryID
	rssEntryMux sync.Mutex
)

// RSS 配置结构
type RSSConfig struct {
	CronExpression string `json:"cron_expression"`
	Enabled        bool   `json:"enabled"`
}

// 默认配置
var currentRSSConfig = RSSConfig{
	CronExpression: "0 */1 * * *", // 默认每小时
	Enabled:        false,
}

const (
	RSS_CONFIG_FILE = "data/rss_config.json"
	LOG_FILE        = "logs/app.log"
)

func main() {
	// 1. 初始化目录
	os.MkdirAll("data", 0755)
	os.MkdirAll("logs", 0755)

	// 2. 配置双向日志 (同时输出到控制台和文件)
	logFile, err := os.OpenFile(LOG_FILE, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0666)
	if err != nil {
		fmt.Println("Failed to open log file:", err)
	} else {
		// 同时写文件和标准输出
		mw := io.MultiWriter(os.Stdout, logFile)
		log.SetOutput(mw)
		// 也要让 Gin 的日志输出到这里 (可选)
		// gin.DefaultWriter = mw
	}

	db.InitDB()

	// 3. 启动定时任务系统
	c = cron.New()
	c.Start()
	loadAndApplyRSSConfig() // 加载并应用配置

	// 4. Web Server
	r := gin.Default()
	r.Static("/static", "./static")
	r.StaticFile("/", "./static/index.html")

	api := r.Group("/api")
	{
		// ... (原有的 /magnets, /scrape/full, /emby/missing, /cd2/add 保持不变) ...
		// 为了节省篇幅，这里只列出由于功能增强而修改/新增的 API，
		// 请保留原有的 CD2 和 Magnets 接口逻辑！

		// --- 复用原有的接口逻辑 (请确保你也复制了 CD2/Emby/Magnets 的代码) ---
		api.GET("/magnets", func(c *gin.Context) {
			var magnets []db.Magnet
			db.DB.Order("id desc").Find(&magnets)
			maxEp := 0
			grouped := make(map[int][]db.Magnet)
			for _, m := range magnets {
				epNum, _ := strconv.Atoi(m.Episode)
				if epNum > maxEp { maxEp = epNum }
				if epNum > 0 { grouped[epNum] = append(grouped[epNum], m) }
			}
			c.JSON(200, gin.H{"data": magnets, "max_episode": maxEp, "grouped_by_episode": grouped})
		})
		
		api.POST("/scrape/full", func(c *gin.Context) {
			go core.RunFullScraper()
			c.JSON(200, gin.H{"message": "Scraper started"})
		})

		api.POST("/emby/missing", func(c *gin.Context) {
			var conf core.EmbyConfig
			if err := c.BindJSON(&conf); err != nil { c.JSON(400, gin.H{"error": err.Error()}); return }
			missing, total, err := core.CheckEmbyMissing(conf)
			if err != nil { c.JSON(500, gin.H{"error": err.Error()}); return }
			c.JSON(200, gin.H{"missing_episodes": missing, "total_count": total})
		})

		api.POST("/cd2/add", func(c *gin.Context) {
			type Req struct { Magnet string `json:"magnet"` }
			var req Req
			if err := c.BindJSON(&req); err != nil { c.JSON(400, gin.H{"error": "Invalid JSON"}); return }
			confBytes, _ := os.ReadFile("data/cd2_config.json")
			var conf core.CD2Config
			json.Unmarshal(confBytes, &conf)
			msg, err := core.AddToCD2(conf, req.Magnet)
			if err != nil { c.JSON(500, gin.H{"error": err.Error()}); return }
			c.JSON(200, gin.H{"message": msg})
		})

		api.POST("/cd2/config", func(c *gin.Context) {
			var conf core.CD2Config
			c.BindJSON(&conf)
			bytes, _ := json.Marshal(conf)
			os.WriteFile("data/cd2_config.json", bytes, 0644)
			c.JSON(200, gin.H{"message": "Saved"})
		})
		
		api.DELETE("/database", func(c *gin.Context) {
			db.DB.Exec("DELETE FROM magnets")
			c.JSON(200, gin.H{"message": "Database cleared"})
		})

		// --- 新增/增强的接口 ---

		// 1. 获取日志 (读取最后 50 行)
		api.GET("/system/logs", func(c *gin.Context) {
			lines, err := readLastLines(LOG_FILE, 50)
			if err != nil {
				c.JSON(200, gin.H{"logs": []string{"No logs yet or error reading logs."}})
				return
			}
			c.JSON(200, gin.H{"logs": lines})
		})

		// 2. RSS 配置修改
		api.POST("/rss/config", func(c *gin.Context) {
			var newConfig RSSConfig
			if err := c.BindJSON(&newConfig); err != nil {
				c.JSON(400, gin.H{"error": "无效的配置格式"})
				return
			}

			// 验证 Cron 表达式有效性
			if newConfig.Enabled {
				parser := cron.NewParser(cron.Minute | cron.Hour | cron.Dom | cron.Month | cron.Dow)
				if _, err := parser.Parse(newConfig.CronExpression); err != nil {
					c.JSON(400, gin.H{"error": "Cron 表达式无效: " + err.Error()})
					return
				}
			}

			// 保存并应用
			rssEntryMux.Lock()
			currentRSSConfig = newConfig
			saveRSSConfig()
			applyRSSConfigLocked() // 重新调度
			rssEntryMux.Unlock()

			status := "RSS已关闭"
			if newConfig.Enabled {
				status = "RSS已启用: " + newConfig.CronExpression
			}
			c.JSON(200, gin.H{"message": status})
		})
	}

	log.Println("Server running on :4869")
	r.Run(":4869")
}

// --- 辅助函数 ---

// 读取文件最后 N 行
func readLastLines(filename string, n int) ([]string, error) {
	file, err := os.Open(filename)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	var lines []string
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		lines = append(lines, scanner.Text())
	}
	
	if len(lines) > n {
		return lines[len(lines)-n:], nil
	}
	return lines, nil
}

func loadAndApplyRSSConfig() {
	data, err := os.ReadFile(RSS_CONFIG_FILE)
	if err == nil {
		json.Unmarshal(data, &currentRSSConfig)
	}
	rssEntryMux.Lock()
	defer rssEntryMux.Unlock()
	applyRSSConfigLocked()
}

func saveRSSConfig() {
	data, _ := json.Marshal(currentRSSConfig)
	os.WriteFile(RSS_CONFIG_FILE, data, 0644)
}

// 必须在锁内调用
func