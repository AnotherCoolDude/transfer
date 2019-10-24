package models

import (
	"sort"
	"time"
)

// PAProject wraps the json respond from proad into a go struct
type PAProject struct {
	Urno         int       `json:"urno" db:"id"`
	Projectno    string    `json:"projectno" db:"projectno"`
	ProjectName  string    `json:"project_name" db:"project_name"`
	Manager      PAManager `json:"manager"`
	Type         string    `json:"type" db:"type"`
	Status       string    `json:"status" db:"status"`
	Orderno      string    `json:"orderno" db:"orderno"`
	OrderDate    string    `json:"order_date" db:"order_date"`
	DeliveryDate string    `json:"delivery_date" db:"delivery_date"`
	Description  string    `json:"description" db:"description"`
	Todos        []PATodo
}

// SortTodos sorts todos using the FromDatetime property
func (p *PAProject) SortTodos() {
	sort.Slice((*p).Todos, func(i, j int) bool {
		ti, err := time.Parse(time.RFC3339, ((*p).Todos)[i].FromDatetime)
		tj, err := time.Parse(time.RFC3339, ((*p).Todos)[j].FromDatetime)
		if err != nil {
			panic(0)
		}
		return ti.Before(tj)
	})
}

// OrderTodos modifies p's todos, so they match projects todos
func (p *PAProject) OrderTodos(project *BCProject) {
	tt := []PATodo{}
	usedIdx := []int{}
	existing := false

	for _, bct := range project.Todos {
		for i, pat := range p.Todos {
			if bct.Timestamp() == pat.Timestamp() {
				tt = append(tt, pat)
				usedIdx = append(usedIdx, i)
				existing = true
			}
		}
		if !existing {
			tt = append(tt, PATodo{})
		}
		existing = false
	}
	for j, t := range p.Todos {
		for _, i := range usedIdx {
			if i == j {
				existing = true
			}
		}
		if !existing {
			tt = append(tt, t)
		}
		existing = false
	}
	p.Todos = tt
}
