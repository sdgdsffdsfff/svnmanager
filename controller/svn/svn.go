package svn


import (
	"github.com/go-martini/martini"
	"github.com/martini-contrib/render"
	"king/helper"
	"king/utils/JSON"
	"net/http"
	"king/service/svn"
	_ "github.com/antonholmquist/jason"
)

func Update(rend render.Render, req *http.Request){
	result, err := update()
	if err != nil {
		rend.JSON(200, helper.Error(err, result))
		return
	}

	rend.JSON(200, helper.Success(result))
}

func Revert(rend render.Render, params martini.Params){
	rend.JSON(200, JSON.Type{
		"code": params["version"],
	})
}

func GetLastVersion(rend render.Render){
	version, err := svn.GetLastVersion()
	if err != nil {
		rend.JSON(200, helper.Error(err))
		return
	}
	result := JSON.Parse(version)
	result["List"] = JSON.Parse(version.List)
	rend.JSON(200, helper.Success(result))
}

func GetUndeployFiles(rend render.Render){
	list, err := svn.GetUnDeployFileList()
	if err != nil {
		rend.JSON(200, helper.Error(err))
	} else if len(list) == 0 {
		rend.JSON(200, helper.Error(helper.EmptyError) )
	}else{
		rend.JSON(200, helper.Success(list))
	}
}
