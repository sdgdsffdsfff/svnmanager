package server

import (
	"king/helper"
	"king/utils/JSON"
	"king/service"
	"king/rpc"
)

func DeployCtrl(filesId []int64, clientsId []int64) (JSON.Type, error) {
	results := JSON.Type{}
	errorCount := 0

//	if len(filesId) == 0 || len(clientsId) == 0 {
//		return nil, helper.NewError("no file need update")
//	}

	fileList, err := service.Svn.GetUnDeployFileList(filesId)
	if err != nil {
		return results, err
	}

	service.Svn.Lock()
	clientList := service.Client.List(clientsId)

	helper.AsyncMap(clientList, func(i int) bool {
		client := clientList[i]
		result, err := service.Client.CallRpc(client, "RpcDeploy.Deploy", rpc.DeployArgs{fileList, client.DeployPath})
		if err != nil {
			errorCount++
			return false
		}
		client.Version = service.Svn.Version
		err = service.Client.Update(client.WebServer, "Version")
		results[helper.Itoa64(client.Id)] = JSON.Type{
			"Version": client.Version,
			"result": result,
			"error": err,
		}
		return false
	})

	service.Svn.Release()

	//TODO
	//版本同步确认后再清空
	//部署没有错误
	if errorCount == 0 {
		//清空未部署列表
		if err := service.Svn.ClearDeployFile(); err != nil {
			return nil, err
		}
	}else {
		return results, helper.NewError("deplpy error")
	}

	return results, nil
}
