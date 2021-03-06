package models

import (
	"time"
	"unicode"
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

// Projectno returns the projectno of the todo
func (t *BCTodo) Projectno() string {
	if len(t.Bucket.Name) < 14 {
		return ""
	}
	nr := t.Bucket.Name[:14]
	for i := 0; i < 4; i++ {
		r := rune(nr[i])
		if !unicode.IsUpper(r) {
			return ""
		}
	}
	return nr
}

// Timestamp is a identifier for comparing with other todos
func (t BCTodo) Timestamp() string {
	return t.CreatedAt.Format(time.RFC3339)
}

// Identifier returns a unique identifier
func (t BCTodo) Identifier() int {
	return t.ID
}

// ClientType returns the type of Todo
func (t BCTodo) ClientType() string {
	return "basecamp"
}
