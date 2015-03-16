package client

import(
	sysConfig "king/config"
	"github.com/go-martini/martini"
	"github.com/martini-contrib/render"
	"king/helper"
)

type Routes struct {}

func init() {
	sysConfig.AppendValue(sysConfig.Controller, &Routes{})
}

func (ctn *Routes) SetRouter(m *martini.ClassicMartini){
	m.Get("/", func(rend render.Render) {
		rend.JSON(200, helper.Success("Is Running"))
	})
}
