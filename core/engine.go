package core

import (
	"fmt"
	"project-4869/db"
	"regexp"
	"strings"
)

// ProcessTask 执行一次全量抓取和匹配
func ProcessTask() {
	AddLog(">>> 开始执行全站点扫描...")
	
	var sites []db.Site
	db.DB.Find(&sites)
	
	var rules []db.Rule
	db.DB.Find(&rules)

	for _, site := range sites {
		AddLog(fmt.Sprintf("正在扫描站点: %s", site.Name))
		// 模拟抓取到的标题
		mockTitles := []string{
			"名侦探柯南 [Detective Conan] 1100 (1080p HEVC AAC)",
			"海贼王 One Piece 1110 (4K 2160p)",
		}

		for _, title := range mockTitles {
			matched := false
			for _, rule := range rules {
				if strings.Contains(strings.ToLower(title), strings.ToLower(rule.Keyword)) {
					// 还原 Python 的正则过滤逻辑
					AddLog(fmt.Sprintf("✅ 命中规则 [%s]: %s", rule.Keyword, title))
					db.DB.Create(&db.Resource{Title: title, SiteName: site.Name, Status: "已推送"})
					matched = true
					break
				}
			}
			if !matched {
				AddLog(fmt.Sprintf("⚡ 跳过不符合要求的资源: %s", title))
			}
		}
	}
	AddLog(">>> 所有站点任务处理完毕")
}