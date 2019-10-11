package models

import (
	"time"
)

// BCTodolist is a struct generated from basecamps json response
type BCTodolist struct {
	ID               int       `json:"id"`
	Status           string    `json:"status"`
	VisibleToClients bool      `json:"visible_to_clients"`
	CreatedAt        time.Time `json:"created_at"`
	UpdatedAt        time.Time `json:"updated_at"`
	Title            string    `json:"title"`
	InheritsStatus   bool      `json:"inherits_status"`
	Type             string    `json:"type"`
	URL              string    `json:"url"`
	AppURL           string    `json:"app_url"`
	BookmarkURL      string    `json:"bookmark_url"`
	SubscriptionURL  string    `json:"subscription_url"`
	CommentsCount    int       `json:"comments_count"`
	CommentsURL      string    `json:"comments_url"`
	Position         int       `json:"position"`
	Parent           BCParent  `json:"parent"`
	Bucket           BCBucket  `json:"bucket"`
	Creator          BCCreator `json:"creator"`
	Description      string    `json:"description"`
	Completed        bool      `json:"completed"`
	CompletedRatio   string    `json:"completed_ratio"`
	Name             string    `json:"name"`
	TodosURL         string    `json:"todos_url"`
	GroupsURL        string    `json:"groups_url"`
	AppTodosURL      string    `json:"app_todos_url"`
}

// BCParent is a struct generated from basecamps json response
type BCParent struct {
	ID     int    `json:"id"`
	Title  string `json:"title"`
	Type   string `json:"type"`
	URL    string `json:"url"`
	AppURL string `json:"app_url"`
}
