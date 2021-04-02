package controllers

import (
	"github.com/revel/revel"
)

// Dashboard controller
type Dashboard struct {
	Admin
}

// Index action: GET
func (c Dashboard) Index() revel.Result {

	return c.Render()
}
