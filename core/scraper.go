package core

import (
	"log"
	"project-4869/db"
	"strings"
	"time"

	"github.com/playwright-community/playwright-go"
)

func RunFullScraper() {
	log.Println("[Scraper] Starting Full Scrape via Playwright...")

	// 启动 Playwright
	pw, err := playwright.Run()
	if err != nil {
		log.Printf("[Scraper] Could not start playwright: %v", err)
		return
	}
	defer pw.Stop()

	// 启动浏览器 (Headless)
	// 注意：Docker 容器内通过 ENV 设置了 PLAYWRIGHT_BROWSERS_PATH，驱动会自动找到浏览器
	browser, err := pw.Chromium.Launch(playwright.BrowserTypeLaunchOptions{
		Headless: playwright.Bool(true),
	})
	if err != nil {
		log.Printf("[Scraper] Could not launch browser: %v", err)
		return
	}
	defer browser.Close()

	page, err := browser.NewPage()
	if err != nil {
		log.Printf("[Scraper] Could not create page: %v", err)
		return
	}

	log.Println("[Scraper] Navigating to sbsub.com...")
	if _, err = page.Goto("https://www.sbsub.com/data/", playwright.PageGotoOptions{Timeout: playwright.Float(60000)}); err != nil {
		log.Printf("[Scraper] Navigation failed: %v", err)
		return
	}

	// 1. 处理版权同意按钮
	if agreeBtn := page.Locator("#agree"); agreeBtn != nil {
		if vis, _ := agreeBtn.IsVisible(); vis {
			agreeBtn.Click()
			log.Println("[Scraper] Clicked agree button")
		}
	}

	// 2. 等待内容加载
	page.WaitForSelector("div.resdiv-l", playwright.PageWaitForSelectorOptions{Timeout: playwright.Float(10000)})

	// 获取所有行
	rows, err := page.Locator("div.resdiv-l").All()
	if err != nil {
		log.Printf("[Scraper] Could not get rows: %v", err)
		return
	}

	log.Printf("[Scraper] Found %d rows to process...", len(rows))
	count := 0

	for _, row := range rows {
		text, _ := row.InnerText()
		
		// 查找该行内的所有 Modal 触发按钮
		buttons, _ := row.Locator("a[data-toggle='modal']").All()
		
		for _, btn := range buttons {
			// 点击打开 Modal
			btn.Click()
			
			// 等待 Input 出现
			inputLoc := page.Locator("input.reslink:visible")
			// 增加容错，最多等待2秒
			err := inputLoc.WaitFor(playwright.LocatorWaitForOptions{Timeout: playwright.Float(2000)})
			if err != nil {
				// 如果没出现，可能是还没加载或出错了，关闭模态框尝试下一个
				page.Keyboard().Press("Escape")
				continue
			}
			
			magnet, _ := inputLoc.InputValue()
			
			// 关闭 Modal
			page.Keyboard().Press("Escape")
			// 简单等待 Modal 动画消失
			time.Sleep(200 * time.Millisecond)

			if magnet != "" && strings.HasPrefix(magnet, "magnet:") {
				meta := ParseTitle(text)
				if meta.Episode == "" { meta.Episode = "0" }
				
				record := db.Magnet{
					MagnetLink:  magnet,
					Episode:     meta.Episode,
					RawTitle:    text,
					Resolution:  meta.Resolution,
					Container:   meta.Container,
					Subtitle:    meta.Subtitle,
					SourceType:  meta.SourceType,
					PublishDate: time.Now().Format("2006-01-02"),
					CreatedAt:   time.Now(),
				}
				
				res := db.DB.Where(db.Magnet{MagnetLink: magnet, Episode: meta.Episode}).FirstOrCreate(&record)
				if res.RowsAffected > 0 {
					count++
				}
			}
		}
	}

	log.Printf("[Scraper] Finished. Added %d new items.", count)
}