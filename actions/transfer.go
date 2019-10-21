package actions

import (
	"github.com/AnotherCoolDude/transfer/models"
	"github.com/gobuffalo/buffalo"
	"sync"
)

// TransferShow default implementation.
func TransferShow(c buffalo.Context) error {

	// Projectpair wraps a pair of basecamp and Proad projects into a struct
	type Projectpair struct {
		Basecamp models.BCProject
		Proad    models.PAProject
	}
	projectsmap := map[string]Projectpair{}

	// error to minimize allocation
	var err error

	// get projects from basecamp
	var bcprojects []models.BCProject

	// fetch basecampprojects
	err = BCClient.unmarshal("/projects.json", query{}, &bcprojects)
	if err != nil {
		c.Error(404, err)
	}

	// filter out projects without projectno and add paprojects
	idx := 0
	for _, p := range bcprojects {
		if p.Projectno() != "" {
			bcprojects[idx] = p
			idx++
		}
	}
	bcprojects = bcprojects[:idx]

	// fetch basecamptodos concurrent
	sem := make(chan int, 4) // 4 jobs at a time
	errChan := make(chan error, 1)
	var wg sync.WaitGroup
	wg.Add(len(bcprojects))

	for i := range bcprojects {
		go BCClient.fetchTodosAsync(&bcprojects[i], sem, &wg, errChan)
	}
	wg.Wait()
	close(errChan)
	if err = <-errChan; err != nil {
		return c.Error(404, err)
	}

	// fetch proadprojects concurrent
	paprojects := make([]models.PAProject, len(bcprojects))
	errChan = make(chan error, 1)
	wg.Add(len(bcprojects))
	for i := range bcprojects {
		go PAClient.fetchProjectAsync(bcprojects[i].Projectno(), &paprojects[i], sem, &wg, errChan)
	}
	wg.Wait()
	close(errChan)
	if err = <-errChan; err != nil {
		return c.Error(404, err)
	}

	// fetch proadtodos concurrent
	errChan = make(chan error, 1)
	wg.Add(len(paprojects))
	for i := range paprojects {
		go PAClient.fetchTodosAsync(&paprojects[i], sem, &wg, errChan)
	}
	wg.Wait()
	close(errChan)
	if err = <-errChan; err != nil {
		return c.Error(404, err)
	}

	//sort todos
	for i := range bcprojects {
		bcprojects[i].SortTodos()
		paprojects[i].SortTodos()
	}

	tt := []models.Todo{}

	// fill struct to better access data in template
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
