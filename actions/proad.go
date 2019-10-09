package actions

import (
	"github.com/gobuffalo/buffalo"
	"net/http"
)

var (
	paClient = defaultProadclient()
)

// ProadShow default implementation.
func ProadShow(c buffalo.Context) error {
	resp, err := paClient.do("GET", "projects", http.NoBody, map[string]string{ /*"projectno": "ACCO-0001-0002"*/ })
	if err != nil {
		return c.Error(404, err)
	}
	return c.Render(200, responseRenderer{response: resp})
	// return c.Render(200, r.HTML("proad/show.html"))
}
