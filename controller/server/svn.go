package server

import (
	"github.com/go-martini/martini"
	"github.com/martini-contrib/render"
	"king/config"
	"king/helper"
	"king/utils/JSON"
	"net/http"
	"king/service/svn"
	_ "github.com/antonholmquist/jason"
	"github.com/antonholmquist/jason"
)

type SvnCtrl struct{
}

func init(){
	config.AppendValue(config.Controller, &SvnCtrl{})
}

func (ctn *SvnCtrl) SetRouter(m *martini.ClassicMartini) {

	m.Group("/server", func(r martini.Router){
		r.Get("/svn", func (rend render.Render, req *http.Request){
			result, err := svn.GetLastVersion()
			if err != nil {
				rend.HTML(500, "500", err)
				return
			}
			rend.HTML(200, "server/svn", result)
		})
	})

	m.Group("/aj/svn", func(r martini.Router){
		r.Post("/up", func(rend render.Render, req *http.Request) {

			//body, _ := jason.NewObjectFromReader(req.Body)
			//paths, _ := body.GetStringArray("paths")

			result, err := SvnUpCtrl()
			if err != nil {
				rend.JSON(200, helper.Error(err, result))
				return
			}

			rend.JSON(200, helper.Success(result))
		})

		r.Get("/revert/:version", func (rend render.Render, params martini.Params){
			rend.JSON(200, JSON.Type{
				"code": params["version"],
			})
		})

		r.Get("/lastVersion", func(rend render.Render) {
			version, err := svn.GetLastVersion()
			if err != nil {
				rend.JSON(200, helper.Error(err))
				return
			}
			result := JSON.Parse(version)
			result["List"] = JSON.Parse(version.List)
			rend.JSON(200, helper.Success(result))
		})

		r.Get("/undeploy/files", func(rend render.Render){
			list, err := svn.GetUnDeployFileList()
			if err != nil {
				rend.JSON(200, helper.Error(err))
			} else if len(list) == 0 {
				rend.JSON(200, helper.Error(helper.EmptyError) )
			}else{
				rend.JSON(200, helper.Success(list))
			}
		})
	})

	m.Post("/aj/deploy", func(rend render.Render, req *http.Request){
		body, err := jason.NewObjectFromReader(req.Body)
		if err != nil {
			rend.JSON(200, helper.Error(helper.ParamsError))
			return
		}
		//如果数据为[0]则表示up_file表中的全部文件
		filesId, err := body.GetInt64Array("filesId")
		if err != nil {
			rend.JSON(200, helper.Error(helper.ParamsError))
			return
		}
		clientsId, err := body.GetInt64Array("clientsId")
		if err != nil {
			rend.JSON(200, helper.Error(helper.ParamsError))
			return
		}

		result, err := DeployCtrl( filesId, clientsId )
		//报告错误原因
		if err != nil {
			rend.JSON(200, helper.Error(err, result))
			return
		}
		rend.JSON(200, helper.Success(result))
	})
}
