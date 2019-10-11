package models

import (
	"time"
)

// BCTodo is a struct generated from basecamps json response
type BCTodo struct {
	ID                    int       `json:"id"`
	Status                string    `json:"status"`
	VisibleToClients      bool      `json:"visible_to_clients"`
	CreatedAt             time.Time `json:"created_at"`
	UpdatedAt             time.Time `json:"updated_at"`
	Title                 string    `json:"title"`
	InheritsStatus        bool      `json:"inherits_status"`
	Type                  string    `json:"type"`
	URL                   string    `json:"url"`
	AppURL                string    `json:"app_url"`
	BookmarkURL           string    `json:"bookmark_url"`
	SubscriptionURL       string    `json:"subscription_url"`
	CommentsCount         int       `json:"comments_count"`
	CommentsURL           string    `json:"comments_url"`
	Position              int       `json:"position"`
	Parent                BCParent  `json:"parent"`
	Bucket                BCBucket  `json:"bucket"`
	Creator               BCCreator `json:"creator"`
	Description           string    `json:"description"`
	Completed             bool      `json:"completed"`
	Content               string    `json:"content"`
	StartsOn              time.Time `json:"starts_on"`
	DueOn                 time.Time `json:"due_on"`
	Assignees             []string  `json:"assignees"`
	CompletionSubscribers []string  `json:"completion_subscribers"`
	CompletionURL         string    `json:"completion_url"`
}
