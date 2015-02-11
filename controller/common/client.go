package common

import (
	"github.com/go-martini/martini"
	"github.com/martini-contrib/render"
	"king/config"
)

type ClientController struct {}

func init() {
	config.AppendValue(config.Controller, &ClientController{})
}

func (ctn *ClientController) SetRouter(m *martini.ClassicMartini) {
	m.Group("/client", func(r martini.Router) {
		r.Get("/push/notify", ctn.PushNotify)
		r.Get("/push/status", ctn.PushStatus)
	})
}

//client do deply return back status
func (ctn *ClientController) PushStatus(render render.Render) {
	//deal client nofity
}

//client
func (ctn *ClientController) PushNotify(render render.Render) {
	//deal client nofity
}
