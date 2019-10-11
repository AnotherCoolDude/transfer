package actions

import (
	"github.com/AnotherCoolDude/transfer/models"
	"github.com/gobuffalo/buffalo"
	"net/http"
)

// TransferShow default implementation.
func TransferShow(c buffalo.Context) error {

	var bcprojects []models.BCProject
	err := BCClient.unmarshal("/projects.json", query{}, &bcprojects)
	if err != nil {
		c.Error(404, err)
	}
	respChan := make(chan asyncResponse)
	counter := 0
	for _, p := range bcprojects {
		if p.Projectno() == "" {
			continue
		}
		counter++
		go PAClient.async("GET", "projects", http.NoBody, map[string]string{"projectno": p.Projectno()}, respChan)
	}
	results := make([]asyncResponse, counter)
	var paprojects []models.PAProject
	for i := range results {
		results[i] = <-respChan
		if results[i].err != nil {
			return c.Error(404, results[i].err)
		}
		prj, err := unmarshalPAProjects(results[i].response)
		if err != nil {
			return c.Error(404, err)
		}
		paprojects = append(paprojects, prj...)
	}

	c.Set("basecamp", bcprojects)
	c.Set("proad", paprojects)
	return c.Render(200, r.HTML("transfer/show.html"))
}
