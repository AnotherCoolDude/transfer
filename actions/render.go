package actions

import (
	"errors"
	"github.com/gobuffalo/buffalo/render"
	"github.com/gobuffalo/packr/v2"
	"io"
	"io/ioutil"
	"net/http"
)

var r *render.Engine
var assetsBox = packr.New("app:assets", "../public")

func init() {
	r = render.New(render.Options{
		// HTML layout to be used for all HTML requests:
		HTMLLayout: "application.plush.html",

		// Box containing all of the templates:
		TemplatesBox: packr.New("app:templates", "../templates"),
		AssetsBox:    assetsBox,

		// Add template helpers here:
		Helpers: render.Helpers{
			// for non-bootstrap form helpers uncomment the lines
			// below and import "github.com/gobuffalo/helpers/forms"
			// forms.FormKey:     forms.Form,
			// forms.FormForKey:  forms.FormFor,
		},
	})
}

// custom renderer

// responseRenderer

type responseRenderer struct {
	response          *http.Response
	unmarshalledBytes []byte
}

func (rr responseRenderer) ContentType() string {
	return "application/json; charset=utf-8"
}

func (rr responseRenderer) Render(writer io.Writer, data render.Data) error {
	var bb []byte
	if rr.response == nil && len(rr.unmarshalledBytes) == 0 {
		return errors.New("must provide at least one field in responseRenderer")
	}
	bb = rr.unmarshalledBytes
	if rr.response != nil {
		readbytes, err := ioutil.ReadAll(rr.response.Body)
		if err != nil {
			return err
		}
		defer rr.response.Body.Close()
		bb = readbytes
	}
	_, err := writer.Write(bb)
	if err != nil {
		return err
	}
	return nil
}
