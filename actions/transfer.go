package actions

import (
	"net/http"
	"strconv"

	"github.com/mitchellh/mapstructure"

	"github.com/AnotherCoolDude/transfer/models"
	"github.com/gobuffalo/buffalo"
)

// TransferShow default implementation.
func TransferShow(c buffalo.Context) error {

	var bcprojects []models.BCProject
	err := BCClient.unmarshal("/projects.json", query{}, &bcprojects)
	if err != nil {
		c.Error(404, err)
	}
	unmarshalChan := make(chan asyncUnmarshal)
	counter := 0
	for _, p := range bcprojects {
		if p.Projectno() == "" {
			continue
		}
		go PAClient.async("GET", "projects", http.NoBody, map[string]string{"projectno": p.Projectno()}, unmarshalChan)
		counter++
	}
	paprojects := []models.PAProject{}
	result := make([]asyncUnmarshal, counter)

	for i := range result {
		result[i] = <-unmarshalChan
		if result[i].err != nil {
			return c.Error(404, result[i].err)
		}
		var p []models.PAProject
		err = mapstructure.Decode(result[i].model, &p)
		if err != nil {
			return c.Error(404, err)
		}
		paprojects = append(paprojects, p...)
	}

	counter = 0
	for _, p := range paprojects {
		go PAClient.async("GET", "tasks", http.NoBody, query{"project": strconv.Itoa(p.Urno)}, unmarshalChan)
		counter++
	}
	result = make([]asyncUnmarshal, counter)
	tt := []models.PATodo{}
	for i := range result {
		result[i] = <-unmarshalChan
		if result[i].err != nil {
			return c.Error(404, err)
		}
		var t []models.PATodo
		err = mapstructure.Decode(result[i].model, &t)
		if err != nil {
			return c.Error(404, err)
		}
		tt = append(tt, t...)
	}

	// return c.Render(200, render.JSON(paprojects))
	c.Set("basecamp", bcprojects)
	c.Set("proad", paprojects)
	return c.Render(200, r.HTML("transfer/show.html"))
}
