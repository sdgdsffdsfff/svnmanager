package master

import (
	"github.com/go-martini/martini"
	"github.com/martini-contrib/render"
	"king/helper"
	"king/utils/JSON"
	"net/http"
	"king/service/master"
	_ "github.com/antonholmquist/jason"
	"king/service/task"
)

func Update(rend render.Render, req *http.Request){

	if master.IsLock() {
		rend.JSON(200, helper.Error(helper.BusyError))
		return
	}

	result, err := update()
	if err != nil {
		rend.JSON(200, helper.Error(err, result))
		return
	}

	rend.JSON(200, helper.Success(result))
}

func Compile(rend render.Render, req *http.Request) {

	if master.IsLock() {
		rend.JSON(200, helper.Error(helper.BusyError))
		return
	}

	master.Lock()
	master.SetBusy()
	master.SetMessage("ready to compile")

	task.Trigger("compile")
	rend.JSON(200, helper.Success())
}

func Revert(rend render.Render, params martini.Params){
	rend.JSON(200, JSON.Type{
		"code": params["version"],
	})
}

func GetLastVersion(rend render.Render){
	version, err := master.GetLastVersion()
	if err != nil {
		rend.JSON(200, helper.Error(err))
		return
	}
	result := JSON.Parse(version)
	result["List"] = JSON.Parse(version.List)
	rend.JSON(200, helper.Success(result))
}

func GetUndeployFiles(rend render.Render){
	list, err := master.GetUnDeployFileList()
	if err != nil {
		rend.JSON(200, helper.Error(err))
	} else if len(list) == 0 {
		rend.JSON(200, helper.Error(helper.EmptyError) )
	}else{
		rend.JSON(200, helper.Success(list))
	}
}

func ShowError(rend render.Render) {
	if master.Error {
		rend.JSON(200, helper.Success(master.ErrorLog))
	} else {
		rend.JSON(200, helper.Error())
	}
}
