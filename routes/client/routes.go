package client

import (
	"github.com/go-martini/martini"
	"github.com/martini-contrib/render"
	"king/controller"
	"king/helper"
)

type Routes struct{}

func init() {
	controller.AddController(&Routes{})
}

func (ctn *Routes) SetRouter(m *martini.ClassicMartini) {
	m.Get("/", func(rend render.Render) {
		rend.JSON(200, helper.Success("Is Running"))
	})
}
