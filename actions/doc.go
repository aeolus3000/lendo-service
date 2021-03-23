package actions

import (
	"net/http"

	"github.com/gobuffalo/buffalo"
)

//DocHandler is a default handler to serve up
// a doc page.
func DocHandler(c buffalo.Context) error {
	return c.Render(http.StatusOK, r.HTML("swagger.html"))
}
