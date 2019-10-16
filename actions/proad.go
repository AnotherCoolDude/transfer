package actions

import (
	"net/http"
	"strconv"

	"github.com/AnotherCoolDude/transfer/models"
	"github.com/gobuffalo/buffalo"
)

var (
	// PAClient returns the proad client
	PAClient = defaultProadclient()
)

// ProadShow default implementation.
func ProadShow(c buffalo.Context) error {
	projectresp, err := PAClient.do("GET", "projects", http.NoBody, map[string]string{"projectno": "SEIN-0001-0190"})
	if err != nil {
		return c.Error(404, err)
	}
	var projects []models.PAProject
	err = unmarshalProad(projectresp, &projects)
	if err != nil {
		c.Error(404, err)
	}

	todoResp, err := PAClient.do("GET", "tasks", http.NoBody, map[string]string{"project": strconv.Itoa(projects[0].Urno)})
	if err != nil {
		return c.Error(404, err)
	}
	var tds []models.PATodo
	err = unmarshalProad(todoResp, &tds)
	if err != nil {
		return c.Error(404, err)
	}

	c.Set("proad", projects)
	c.Set("todos", tds)
	return c.Render(200, r.HTML("proad/show.html"))
}
