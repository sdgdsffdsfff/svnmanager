package common

import(
	"github.com/go-martini/martini"
	"github.com/martini-contrib/render"
	"net/http"
	"king/config"
)

type HomeController struct {}


func init() {
	config.AppendValue(config.Controller, &HomeController{})

}

func (ctn *HomeController) SetRouter(m *martini.ClassicMartini) {
	m.Get("/", func(rend render.Render, req *http.Request) {
			rend.HTML(200, "index", "")
	})
}

