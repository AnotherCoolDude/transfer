package models

// PAProject wraps the json respond from proad into a go struct
type PAProject struct {
	Urno         int    `json:"urno" db:"id"`
	Projectno    string `json:"projectno" db:"projectno"`
	ProjectName  string `json:"project_name" db:"project_name"`
	Type         string `json:"type" db:"type"`
	Status       string `json:"status" db:"status"`
	Orderno      string `json:"orderno" db:"orderno"`
	OrderDate    string `json:"order_date" db:"order_date"`
	DeliveryDate string `json:"delivery_date" db:"delivery_date"`
	Description  string `json:"description" db:"description"`
	Todos        []PATodo
}
