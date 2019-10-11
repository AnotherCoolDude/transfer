package models

import (
	"time"
)

// BCTodoset is a struct generated from basecamps json response
type BCTodoset struct {
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
	Position         int       `json:"position"`
	Bucket           BCBucket  `json:"bucket"`
	Creator          BCCreator `json:"creator"`
	Completed        bool      `json:"completed"`
	CompletedRatio   string    `json:"completed_ratio"`
	Name             string    `json:"name"`
	TodolistsCount   int       `json:"todolists_count"`
	TodolistsURL     string    `json:"todolists_url"`
	AppTodoslistsURL string    `json:"app_todoslists_url"`
}

// BCBucket is a struct generated from basecamps json response
type BCBucket struct {
	ID   int    `json:"id"`
	Name string `json:"name"`
	Type string `json:"type"`
}

// BCCreator is a struct generated from basecamps json response
type BCCreator struct {
	ID             int       `json:"id"`
	AttachableSgid string    `json:"attachable_sgid"`
	Name           string    `json:"name"`
	EmailAddress   string    `json:"email_address"`
	PersonableType string    `json:"personable_type"`
	Title          string    `json:"title"`
	Bio            string    `json:"bio"`
	CreatedAt      time.Time `json:"created_at"`
	UpdatedAt      time.Time `json:"updated_at"`
	Admin          bool      `json:"admin"`
	Owner          bool      `json:"owner"`
	Client         bool      `json:"client"`
	TimeZone       string    `json:"time_zone"`
	AvatarURL      string    `json:"avatar_url"`
	Company        BCCompany `json:"company"`
}

// BCCompany is a struct generated from basecamps json response
type BCCompany struct {
	ID   int    `json:"id"`
	Name string `json:"name"`
}
