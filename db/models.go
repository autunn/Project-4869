package db

import (
	"log"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

var DB *gorm.DB

// Site 站点配置
type Site struct {
	ID       uint   `gorm:"primaryKey" json:"id"`
	Name     string `json:"name"`     // 站点名称
	URL      string `json:"url"`      // RSS或搜索地址
	Cookie   string `json:"cookie"`   // 登录凭证
}

// Rule 订阅规则
type Rule struct {
	ID        uint   `gorm:"primaryKey" json:"id"`
	Keyword   string `json:"keyword"`    // 必须包含的词
	Exclude   string `json:"exclude"`    // 不能包含的词
	Quality   string `json:"quality"`    // 1080p, 4K 等
}

// Resource 抓取历史
type Resource struct {
	ID        uint   `gorm:"primaryKey" json:"id"`
	Title     string `json:"title"`
	SiteName  string `json:"site_name"`
	Size      string `json:"size"`
	Status    string `json:"status"` // 已过滤, 已推送
}

func InitDB() {
	var err error
	DB, err = gorm.Open(sqlite.Open("data/project4869.db"), &gorm.Config{})
	if err != nil {
		log.Fatal("数据库连接失败")
	}
	DB.AutoMigrate(&Site{}, &Rule{}, &Resource{})
}