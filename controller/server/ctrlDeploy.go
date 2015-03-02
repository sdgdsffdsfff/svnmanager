package server

import (
	"king/helper"
	"king/utils/JSON"
	"king/service"
	"king/model"
	"king/rpc"
	"reflect"
)

func DeployCtrl() ([]JSON.Type, error) {
	results := []JSON.Type{}
	errorList := []JSON.Type{}

	fileList, err := service.Svn.GetUnDeployFileList()
	if err != nil {
		return nil, err
	}
	if len(fileList) == 0 {
		return nil, helper.NewError("no file need update")
	}

	service.Svn.Lock()
	list := service.Client.List()
	results = service.Client.BatchCall(
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
			client.Version = service.Svn.Version
			//更新失败
			if err = service.Client.Update(&client); err != nil {
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
		if err := service.Svn.ClearDeployFile(); err != nil {
			return nil, err
		}
	}else {
		return errorList, helper.NewError("deplpy error")
	}

	return nil, nil
}
