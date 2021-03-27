package entities

import "time"

type Article struct {
	GUID          string    `json:"guid"`
	Title         string    `json:"title"`
	Description   string    `json:"description"`
	Link          string    `json:"link"`
	PublishedTime time.Time `json:"published_date"`
	Provider      string    `json:"provider"`
	Category      string    `json:"category"`
}

type Articles []Article
