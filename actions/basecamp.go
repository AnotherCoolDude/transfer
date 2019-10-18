package actions

import (
	"encoding/json"
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
		wg.Add(1)
		set, _ := BCClient.receiveFromURLAsync(&wg, p.Dock[2].URL, "Todoset", query{})
		sets = append(sets, set.(models.BCTodoset))
	}
	wg.Wait()

	for _, p := range projects {
		var pset models.BCTodoset
		err := BCClient.unmarshal(p.Dock[2].URL, query{}, &pset)
		if err != nil {
			wg.Done()
			return c.Error(404, err)
		}
		sets = append(sets, pset)
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

func (c *basecampclient) receiveFromURLAsync(wg *sync.WaitGroup, url string, modelType string, q query) {
	defer wg.Done()
	resp, err := c.do("GET", url, http.NoBody, q)
	if err != nil {
		return nil, err
	}
	bb, err := responseBytes(resp)
	if err != nil {
		return nil, err
	}

	var set models.BCTodoset
	var lists []models.BCTodolist
	var todos []models.BCTodo

	switch modelType {
	case "Todoset":
		err = json.Unmarshal(bb, &set)
		if err != nil {
			return nil, err
		}
		return set, nil
	case "Todolist":
		err = json.Unmarshal(bb, &lists)
		if err != nil {
			return nil, err
		}
		return lists, nil
	case "Todo":
		err = json.Unmarshal(bb, &todos)
		if err != nil {
			return nil, err
		}
		return todos, nil
	default:
		return nil, err
	}
}
