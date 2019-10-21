package models

import (
	"time"
)

// PATodo wraps the proad response for todos in a struct
type PATodo struct {
	Urno             int    `json:"urno" db:"id"`
	Shortinfo        string `json:"shortinfo" db:"shortinfo"`
	FromDatetime     string `json:"from_datetime" db:"from_datetime"`
	UntilDatetime    string `json:"until_datetime" db:"until_datetime"`
	ReminderDatetime string `json:"reminder_datetime" db:"reminder_datetime"`
	Status           string `json:"status" db:"status"`
	Priority         string `json:"priority" db:"priority"`
	Description      string `json:"description" db:"description"`
}

// Timestamp is a identifier for comparing with other todos
func (t PATodo) Timestamp() string {
	date, err := time.Parse(time.RFC3339, t.FromDatetime)
	if err != nil {
		panic(0)
	}
	return date.Format(time.RFC3339)
}

// Identifier returns a unique identifier
func (t PATodo) Identifier() int {
	return t.Urno
}

// ClientType returns the type of Todo
func (t PATodo) ClientType() string {
	return "proad"
}
