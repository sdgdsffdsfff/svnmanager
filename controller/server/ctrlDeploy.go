package server

import (
	"king/helper"
	"king/utils/JSON"
	"king/service"
	"king/model"
	"king/rpc"
	"reflect"
)

func DeployCtrl(filesId []int64, clientsId []int64) ([]JSON.Type, error) {
	results := []JSON.Type{}
	errorList := []JSON.Type{}

//	if len(filesId) == 0 || len(clientsId) == 0 {
//		return nil, helper.NewError("no file need update")
//	}

	fileList, err := service.Svn.GetUnDeployFileList(filesId)
	if err != nil {
		return results, err
	}

	service.Svn.Lock()
	clientList := service.Client.List(clientsId)

	results = service.Client.BatchCall(
		clientList,
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

	service.Svn.Release()

	//TODO
	//版本同步确认后再清空
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
