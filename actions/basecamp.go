package actions

import (
	"encoding/json"
	"github.com/AnotherCoolDude/transfer/models"
	"github.com/gobuffalo/buffalo"
	"io/ioutil"
	"net/http"
)

var (
	bcClient = defaultBasecampclient()
)

// BasecampShow default implementation.
func BasecampShow(c buffalo.Context) error {
	if !bcClient.isValid() {
		authURL := bcClient.authCodeURL()
		return c.Redirect(http.StatusTemporaryRedirect, authURL)
	}

	resp, err := bcClient.do("GET", "/projects.json", http.NoBody)
	if err != nil {
		return c.Error(404, err)
	}

	if err != nil {
		return c.Error(404, err)
	}

	bytes, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return c.Error(404, err)
	}
	defer resp.Body.Close()

	var projects []models.BCProject
	err = json.Unmarshal(bytes, &projects)
	if err != nil {
		c.Error(404, err)
	}
	for _, p := range projects {
		c.Logger().Debug(p.Projectno())
	}

	// Todo: check if proadclient can fetch projects based on projectno()

	c.Set("projects", projects)
	// return c.Render(200, r.HTML("basecamp/show.html"))
	return c.Render(200, responseRenderer{unmarshalledBytes: bytes})
}

// BasecampCallback handles the callback from basecamp upon authentication
func BasecampCallback(c buffalo.Context) error {
	err := bcClient.handleCallback(c.Request())

	if err != nil {
		return c.Error(404, err)
	}

	if err = bcClient.receiveID(); err != nil {
		return c.Error(404, err)
	}

	return c.Redirect(http.StatusTemporaryRedirect, "basecampShowPath()")
}
