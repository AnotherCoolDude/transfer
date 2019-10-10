package actions

import (
	"github.com/gobuffalo/buffalo"
	"net/http"
)

// BasecampAuth checks wether the basecamp client is valid before performing any requests
func BasecampAuth(handler buffalo.Handler) buffalo.Handler {
	return func(c buffalo.Context) error {
		if BCClient.IsValid() {
			return handler(c)
		}
		return c.Redirect(http.StatusTemporaryRedirect, BCClient.AuthCodeURL())
	}
}
