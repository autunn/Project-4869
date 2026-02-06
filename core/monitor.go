package core

import (
	"log"
	"project-4869/db"
	"strings"
	"time"

	"github.com/mmcdole/gofeed"
)

const RSS_URL = "https://www.sbsub.com/data/rss/"

func RunRSSMonitor() {
	fp := gofeed.NewParser()
	fp.UserAgent = "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/120.0.0.0 Safari/537.36"
	
	feed, err := fp.ParseURL(RSS_URL)
	if err != nil {
		log.Printf("[RSS] Error parsing RSS: %v", err)
		return
	}

	log.Printf("[RSS] Check: Found %d entries", len(feed.Items))
	newCount := 0

	for _, item := range feed.Items {
		magnet := ""
		if strings.HasPrefix(item.Link, "magnet:") {
			magnet = item.Link
		} else {
			for _, enc := range item.Enclosures {
				if strings.HasPrefix(enc.URL, "magnet:") {
					magnet = enc.URL
					break
				}
			}
		}

		if magnet == "" {
			continue
		}

		meta := ParseTitle(item.Title)
		if meta.Episode == "" {
			meta.Episode = "0"
		}

		// 使用 FirstOrCreate 避免重复
		record := db.Magnet{
			MagnetLink:  magnet,
			Episode:     meta.Episode,
			RawTitle:    item.Title,
			Resolution:  meta.Resolution,
			Container:   meta.Container,
			Subtitle:    meta.Subtitle,
			SourceType:  meta.SourceType,
			PublishDate: item.Published,
			CreatedAt:   time.Now(),
		}

		result := db.DB.Where(db.Magnet{MagnetLink: magnet, Episode: meta.Episode}).FirstOrCreate(&record)
		if result.RowsAffected > 0 {
			newCount++
			log.Printf("[RSS] Added new: %s", item.Title)
		}
	}
	if newCount > 0 {
		log.Printf("[RSS] Finished. Added %d new items.", newCount)
	}
}