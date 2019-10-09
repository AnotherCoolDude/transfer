package models

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
