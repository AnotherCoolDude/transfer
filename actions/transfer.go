package actions

import (
	"github.com/AnotherCoolDude/transfer/models"
	"github.com/gobuffalo/buffalo"
	"sync"
)

// TransferShow default implementation.
func TransferShow(c buffalo.Context) error {

	// // Projectpair wraps a pair of basecamp and Proad projects into a struct
	// type Projectpair struct {
	// 	Basecamp models.BCProject
	// 	Proad    models.PAProject
	// }
	// projectsmap := map[string]Projectpair{}

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

	c.Logger().Debug("fetched all projects and todos")

	// sort todos
	sortedProjects := map[string]*models.ClienttypeSorting{}
	clienttypes := []string{}
	tt := models.Todos{}
	for i, prj := range bcprojects {
		if len(prj.Todos) == 0 {
			continue
		}
		for _, t := range prj.Todos {
			tt = append(tt, t)
		}
		if len(paprojects[i].Todos) == 0 {
			continue
		}
		for _, t := range paprojects[i].Todos {
			tt = append(tt, t)
		}
		sortedProjects[prj.Projectno()] = tt.SortByClienttype()
	}
	c.Logger().Debugf("got all todos sorted by client")

	for _, cts := range sortedProjects {
		if len(clienttypes) == 0 {
			clienttypes = cts.Clienttypes()
		}
		for ct, tt := range *cts {
			if ct == "basecamp" {
				continue
			}
			tt.SortedByCounterpart((*cts)[ct])
		}
	}

	c.Logger().Debugf("sorted projects: %+v\n", sortedProjects)

	// // fill struct to better access data in template
	// for _, bcp := range bcprojects {
	// 	for _, pap := range paprojects {
	// 		if pap.Projectno == bcp.Projectno() {
	// 			projectsmap[pap.Projectno] = Projectpair{
	// 				Basecamp: bcp,
	// 				Proad:    pap,
	// 			}
	// 		}
	// 	}
	// }

	// return c.Render(200, render.JSON(paprojects))
	c.Set("clienttypes", clienttypes)
	c.Set("projects", sortedProjects)
	return c.Render(200, r.HTML("transfer/show.html"))
}
