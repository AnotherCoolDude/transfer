package actions

import (
	"github.com/AnotherCoolDude/transfer/models"
	"github.com/gobuffalo/buffalo"
	"net/http"
	"sync"
)

var (
	// BCClient returns the basecamp client
	BCClient = defaultBasecampclient()
)

// BasecampShow default implementation.
func BasecampShow(c buffalo.Context) error {

	var projects []models.BCProject
	err := BCClient.unmarshal("/projects.json", query{}, &projects)
	if err != nil {
		return c.Error(404, err)
	}

	sem := make(chan int, 4)
	errChan := make(chan error, 1)
	var wg sync.WaitGroup
	wg.Add(len(projects))
	for i := range projects {
		go BCClient.fetchTodosAsync(&projects[i], sem, &wg, errChan)
	}
	wg.Wait()
	close(errChan)
	if <-errChan != nil {
		return c.Error(404, err)
	}
	return c.Render(200, r.JSON(projects))

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
