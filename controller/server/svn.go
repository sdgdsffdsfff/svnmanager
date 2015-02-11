package server

import (
	"github.com/go-martini/martini"
	"github.com/martini-contrib/render"
	"king/config"
	"king/helper"
	"king/utils/JSON"
	"king/utils/shell"
	"net/http"
	"king/service"
	"king/model"
	"time"
	"king/rpc"
	"reflect"
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
			result, err := svnUpCtrl()
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
			result, err := deploy()
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

func svnUpCtrl() (JSON.Type, error){
	now := time.Now()

	num, list, err := shell.SvnUp()
	if err != nil {
		return nil, err
	}

	if service.SvnService.IsChanged(num) == false {
		return nil, helper.NewError("not change")
	}

	version := model.Version{
		Version: num,
		Time: now,
		List: JSON.Stringify(list),
	}

	if err := service.SvnService.UpdateVersion(&version); err != nil {
		return nil, err
	}

	if err := service.SvnService.SaveUpFile(list); err != nil {
		return nil, err
	}

	return JSON.Type{
		"Version": version,
		"List": list,
	}, nil
}

func deploy() ([]JSON.Type, error) {
	results := []JSON.Type{}
	errorList := []JSON.Type{}

	fileList, err := service.SvnService.GetUnDeployFileList()
	if err != nil {
		return nil, err
	}
	if len(fileList) == 0 {
		return nil, helper.NewError("no file need update")
	}

	list := service.ClientService.List()
	results = service.ClientService.BatchCall(
		list,
		"RpcDeploy.Deploy",
		rpc.DeployArgs{fileList},
	)
	helper.AsyncMap(results, func(index int) bool {
			var client model.WebServer
			var err interface{}
			//RpcDeploy.Deploy的返回结果
			//正确Response返回为true,错误返回[]JSON.Type的报错列表
			res := results[index]
			JSON.ParseToStruct(res["client"], &client)

			if r := res["result"]; r != nil {
				//部署报错
				switch reflect.TypeOf(r).Kind() {
				case reflect.Slice:
					errorList = append(errorList, res)
			}
		}
			if err = res["error"]; err == nil {
				client.Version = service.SvnService.Version
				//更新失败
				if err = service.ClientService.Update(&client); err != nil {
					res["error"] = err.(error).Error()
			}
			}
			//常规报错
			if err != nil {
				errorList = append(errorList, res)
			}

			return false
		})

	//部署没有错误
	if len(errorList) == 0 {
		//清空未部署列表
		if err := service.SvnService.ClearDeployFile(); err != nil {
			return nil, err
		}
	}else {
		return errorList, helper.NewError("deplpy error")
	}

	return nil, nil
}
