package client

import (
	"github.com/go-martini/martini"
	"github.com/martini-contrib/render"
	"king/config"
	"net/http"
)

type ClientIndexCtrl struct{}

func init(){
	config.AppendValue(config.Controller, &ClientIndexCtrl{})
}

func (ctn *ClientIndexCtrl) SetRouter(m *martini.ClassicMartini) {
	m.Get("/client/index", func(rend render.Render, req *http.Request){
			rend.HTML(200, "client/index", "")
		})
}
