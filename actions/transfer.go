package actions

import (
	"github.com/AnotherCoolDude/transfer/models"
	"github.com/gobuffalo/buffalo"
	"strconv"
)

// TransferShow default implementation.
func TransferShow(c buffalo.Context) error {

	// Projectpair wraps a pair of basecamp and Proad projects into a struct
	type Projectpair struct {
		Basecamp models.BCProject
		Proad    models.PAProject
	}
	projectsmap := map[string]Projectpair{}

	// channel for async requests
	unmarshalChan := make(chan asyncUnmarshal)

	// error to minimize allocation
	var err error

	// get projects from basecamp
	var bcprojects []models.BCProject
	err = BCClient.unmarshal("/projects.json", query{}, &bcprojects)
	if err != nil {
		c.Error(404, err)
	}

	// n := 0
	// for _, x := range a {
	// 	if keep(x) {
	// 		a[n] = x
	// 		n++
	// 	}
	// }
	// a = a[:n]

	// filter out projects without projectno
	idx := 0
	for _, p := range bcprojects {
		if p.Projectno() != "" {
			bcprojects[idx] = p
			idx++
		}
	}
	bcprojects = bcprojects[:idx]

	// get the todoset for each basecamp project
	counter := 0
	for _, p := range bcprojects {
		go BCClient.asyncUnmarshal(p.Dock[2].URL, query{}, unmarshalChan)
		counter++
	}
	result := make([]asyncUnmarshal, counter)
	bcsets := []models.BCTodoset{}

	for i := range result {
		var set models.BCTodoset
		result[i] = <-unmarshalChan
		err = result[i].decode(&set)
		if set == (models.BCTodoset{}) {
			continue
		}
		if err != nil {
			c.Error(404, err)
		}
		c.Logger().Debugf("Set: %s, %s\n", set.Name, set.Bucket.Name)
		bcsets = append(bcsets, set)
	}

	// get the todolists for the sets
	counter = 0
	for _, s := range bcsets {
		go BCClient.asyncUnmarshal(s.TodolistsURL, query{}, unmarshalChan)
		counter++
	}
	result = make([]asyncUnmarshal, counter)
	bclists := []models.BCTodolist{}

	for i := range result {
		var list []models.BCTodolist
		result[i] = <-unmarshalChan
		c.Logger().Debugf("%+v\n", result[i].model)
		err = result[i].decode(&list)
		c.Logger().Debug(list)
		if err != nil {
			c.Error(404, err)
		}
		if len(list) == 0 {
			continue
		}

		c.Logger().Debugf("List: %s, %s\n", list[0].Name, list[0].Bucket.Name)
		bclists = append(bclists, list...)
	}

	// get todos for the lists
	counter = 0
	for _, l := range bclists {
		go BCClient.asyncUnmarshal(l.TodosURL, query{}, unmarshalChan)
		counter++
	}
	result = make([]asyncUnmarshal, counter)
	bctodos := []models.BCTodo{}

	for i := range result {
		var todo []models.BCTodo
		result[i] = <-unmarshalChan
		err = result[i].decode(&todo)
		if len(todo) == 0 {
			continue
		}
		if err != nil {
			c.Error(404, err)
		}
		c.Logger().Debugf("Todo: %s, %s\n", todo[0].Title, todo[0].Bucket.Name)
		bctodos = append(bctodos, todo...)
	}

	// assign todos to basecamp projects
	for _, t := range bctodos {
		c.Logger().Debug(t)
		for i, p := range bcprojects {
			if t.Projectno() == p.Projectno() {
				bcprojects[i].Todos = append(bcprojects[i].Todos, t)
			}
		}
	}

	// for each basecamp project get the corresponding proad project
	counter = 0
	for _, p := range bcprojects {
		go PAClient.asyncUnmarshal("projects", map[string]string{"projectno": p.Projectno()}, unmarshalChan)
		counter++
	}
	result = make([]asyncUnmarshal, counter)
	paprojects := []models.PAProject{}

	for i := range result {
		result[i] = <-unmarshalChan
		var p []models.PAProject
		err = result[i].decode(&p)
		if err != nil {
			return c.Error(404, err)
		}
		paprojects = append(paprojects, p...)
	}

	// for each proad project get its todos
	counter = 0
	for _, p := range paprojects {
		go PAClient.asyncUnmarshal("tasks", query{"project": strconv.Itoa(p.Urno)}, unmarshalChan)
		counter++
	}
	result = make([]asyncUnmarshal, counter)
	for i := range result {
		result[i] = <-unmarshalChan

		var t []models.PATodo
		err = result[i].decode(&t)
		if len(t) == 0 {
			continue
		}
		if err != nil {
			return c.Error(404, err)
		}
		for n, p := range paprojects {
			if result[i].breadcrumb == strconv.Itoa(p.Urno) {
				paprojects[n].Todos = t
			}
		}
	}

	for _, bcp := range bcprojects {
		for _, pap := range paprojects {
			if pap.Projectno == bcp.Projectno() {
				projectsmap[pap.Projectno] = Projectpair{
					Basecamp: bcp,
					Proad:    pap,
				}
			}
		}
	}

	// return c.Render(200, render.JSON(paprojects))
	c.Set("pair", projectsmap)
	return c.Render(200, r.HTML("transfer/show.html"))
}
