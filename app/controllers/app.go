package controllers

import (
	"github.com/revel/revel"
)

// App controller
type App struct {
	*revel.Controller
}

// Index action: GET
func (c App) Index() revel.Result {
	// return c.RenderError(errors.New("nsa"))

	return c.Render()
}
