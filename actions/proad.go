package actions

import (
	"fmt"
	"github.com/AnotherCoolDude/transfer/models"
	"github.com/gobuffalo/buffalo"
	"net/http"
	"sync"
)

var (
	// PAClient returns the proad client
	PAClient = defaultProadclient()
)

// ProadShow default implementation.
func ProadShow(c buffalo.Context) error {

	var projects []models.PAProject
	projectsresp, err := PAClient.do("GET", "projects", http.NoBody, query{"projectno": "SEIN-0001-0190"})
	if err != nil {
		return c.Error(404, err)
	}
	err = unmarshalProad(projectsresp, &projects)
	if err != nil {
		return c.Error(404, err)
	}
	if len(projects) == 0 {
		return c.Error(400, fmt.Errorf("zero projects found: %v", projects))
	}

	sem := make(chan int, 4)
	errChan := make(chan error, 1)
	var wg sync.WaitGroup
	wg.Add(len(projects))

	for i := range projects {
		go PAClient.fetchTodosAsync(&projects[i], sem, &wg, errChan)
	}
	wg.Wait()
	close(errChan)
	if err = <-errChan; err != nil {
		return c.Error(404, err)
	}

	return c.Render(404, r.JSON(projects))
}
