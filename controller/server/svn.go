package server

import (
	"github.com/go-martini/martini"
	"github.com/martini-contrib/render"
	"king/config"
	"king/helper"
	"king/utils/JSON"
	"net/http"
	"king/service"
)

type SvnCtrl struct{}

func init(){
	config.AppendValue(config.Controller, &SvnCtrl{})
}

type Book struct {
	Name string
}

type A struct {
	Name  string
	Age   int
	Books []Book
}

func (ctn *SvnCtrl) SetRouter(m *martini.ClassicMartini) {

	m.Group("/server", func(r martini.Router){
		r.Get("/svn", func (rend render.Render, req *http.Request){
			result, err := service.SvnService.GetLastVersion()
			if err != nil {
				rend.HTML(500, "500", err)
				return
			}
			rend.HTML(200, "server/svn", result)
		})
	})

	m.Group("/aj/svn", func(r martini.Router){
		r.Post("/up", func(rend render.Render, req *http.Request) {
			result, err := SvnUpCtrl()
			if err != nil {
				rend.JSON(200, helper.Error(err))
				return
			}
			rend.JSON(200, helper.Success(result))
		})

		r.Get("/revert/:version", func (rend render.Render, params martini.Params){
			rend.JSON(200, JSON.Type{
				"code": params["version"],
			})
		})

		r.Get("/deploy", func(rend render.Render, req *http.Request){
			result, err := DeployCtrl()
			//报告错误原因
			if err != nil {
				rend.JSON(200, helper.Error(err, result))
				return
			}
			rend.JSON(200, helper.Success(result))
		})

		r.Get("/lastVersion", func(rend render.Render) {
			version, err := service.SvnService.GetLastVersion()
			if err != nil {
				rend.JSON(200, helper.Error(err))
				return
			}
			result := JSON.Parse(version)
			result["List"] = JSON.Parse(version.List)
			rend.JSON(200, helper.Success(result))
		})
	})
}
