package controllers

import (
	"github.com/revel/revel"
)

// App controller
type Admin struct {
	*revel.Controller
}

// Index action: GET
func (c Admin) Index() revel.Result {
	// return c.RenderError(errors.New("nsa"))

	return c.Render()
}
