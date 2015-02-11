package server

import (
	"github.com/go-martini/martini"
	"github.com/martini-contrib/render"
	"king/config"
)

type IndexCtrl struct{}

func init() {
	config.AppendValue(config.Controller, &IndexCtrl{})
}

func (ctn *IndexCtrl) SetRouter(m *martini.ClassicMartini) {
	m.Get("/", func(rend render.Render) {
			rend.HTML(200, "index", nil)
		})
}
