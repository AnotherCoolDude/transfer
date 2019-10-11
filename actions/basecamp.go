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

	var wg sync.WaitGroup

	var projects []models.BCProject
	err := BCClient.unmarshal("/projects.json", query{}, &projects)
	if err != nil {
		return c.Error(404, err)
	}

	sets := []models.BCTodoset{}
	for _, p := range projects {
		var pset models.BCTodoset
		err := BCClient.unmarshal(p.Dock[2].URL, query{}, &pset)
		if err != nil {
			wg.Done()
			return err
		}
		sets = append(sets, pset)
		return nil
	}
	wg.Wait()

	lists := []models.BCTodolist{}
	for _, s := range sets {
		var slists []models.BCTodolist
		err := BCClient.unmarshal(s.TodolistsURL, query{}, &slists)
		if err != nil {
			return c.Error(404, err)
		}
		lists = append(lists, slists...)
	}

	todos := []models.BCTodo{}
	for _, l := range lists {
		var ltodos []models.BCTodo
		err := BCClient.unmarshal(l.TodosURL, query{}, &ltodos)
		if err != nil {
			return c.Error(404, err)
		}
		todos = append(todos, ltodos...)
	}

	return c.Render(200, r.JSON(todos))

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
