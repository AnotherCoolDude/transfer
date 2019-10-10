package actions

import (
	"github.com/gobuffalo/buffalo"
	"net/http"
	"strconv"
)

var (
	// PAClient returns the proad client
	PAClient = defaultProadclient()
)

// ProadShow default implementation.
func ProadShow(c buffalo.Context) error {
	resp, err := PAClient.do("GET", "projects", http.NoBody, map[string]string{"projectno": "SEIN-0001-0190"})
	if err != nil {
		return c.Error(404, err)
	}
	// return c.Render(200, responseRenderer{response: resp})
	prj, err := unmarshalPAProjects(resp)
	if err != nil {
		c.Error(404, err)
	}
	todoResp, err := PAClient.do("GET", "tasks", http.NoBody, map[string]string{"project": strconv.Itoa(prj[0].Urno)})
	if err != nil {
		return c.Error(404, err)
	}
	tds, err := unmarshalPATodos(todoResp)
	if err != nil {
		return c.Error(404, err)
	}

	c.Set("proad", prj)
	c.Set("todos", tds)
	return c.Render(200, r.HTML("proad/show.html"))
}
