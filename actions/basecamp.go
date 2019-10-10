package actions

import (
	"github.com/AnotherCoolDude/transfer/models"
	"github.com/gobuffalo/buffalo"
	"net/http"
)

var (
	// BCClient returns the basecamp client
	BCClient = defaultBasecampclient()
)

// BasecampShow default implementation.
func BasecampShow(c buffalo.Context) error {

	resp, err := BCClient.do("GET", "/projects.json", http.NoBody)
	if err != nil {
		return c.Error(404, err)
	}

	bcprj, err := unmarshalBCProjects(resp)

	if err != nil {
		return c.Error(404, err)
	}

	paprj := []models.PAProject{}

	for _, p := range bcprj {
		if p.Projectno() == "" {
			continue
		}
		paresp, err := PAClient.do("GET", "projects", http.NoBody, map[string]string{"projectno": p.Projectno()})
		if err != nil {
			return c.Error(404, err)
		}
		prjs, err := unmarshalPAProjects(paresp)
		if err != nil {
			return c.Error(404, err)
		}
		paprj = append(paprj, prjs...)
	}

	c.Set("basecamp", bcprj)
	c.Set("proad", paprj)
	return c.Render(200, r.HTML("basecamp/show.html"))
	// return c.Render(200, responseRenderer{unmarshalledBytes: bytes})
}

// BasecampCallback handles the callback from basecamp upon authentication
func BasecampCallback(c buffalo.Context) error {
	err := BCClient.handleCallback(c.Request())

	if err != nil {
		return c.Error(404, err)
	}

	if err = BCClient.receiveID(); err != nil {
		return c.Error(404, err)
	}

	return c.Redirect(http.StatusTemporaryRedirect, "rootPath()")
}
