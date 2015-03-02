package server

import (
	"github.com/go-martini/martini"
	"github.com/martini-contrib/render"
	"king/config"
	"net/http"
	"king/helper"
	"github.com/antonholmquist/jason"
	"king/service"
)

type GroupCtrl struct{}

func init() {
	config.AppendValue(config.Controller, &GroupCtrl{})
}

func (ctn *GroupCtrl) SetRouter(m *martini.ClassicMartini) {
	m.Post("/aj/group/add", func(rend render.Render, req *http.Request) {
			params, _ := jason.NewObjectFromReader(req.Body)
			name, _ := params.GetString("name")
			result, err := service.Group.Add(name)
			if err != nil {
				rend.JSON(200, helper.Error(err))
				return
			}
			rend.JSON(200, helper.Success(result))
		})

}
