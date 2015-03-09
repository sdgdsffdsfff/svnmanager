package server

import (
	"king/helper"
	"king/utils/JSON"
	"king/service/svn"
	"king/service/client"
	"king/service/webSocket"
	"king/rpc"
)

func DeployCtrl(filesId []int64, clientsId []int64) (JSON.Type, error) {
	results := JSON.Type{}
	errorCount := 0

	fileList, err := svn.GetUnDeployFileList(filesId)
	if err != nil {
		return results, err
	}

	svn.GetLock()
	clientList := client.List(clientsId)
	webSocket.Notify("Locking control!")

	helper.AsyncMap(clientList, func(i int) bool {
		c := clientList[i]
		result, err := client.CallRpc(c, "RpcDeploy.Deploy", rpc.DeployArgs{fileList, c.DeployPath})
		if err != nil {
			errorCount++
			return false
		}
		c.Version = svn.Version
		err = client.Update(c.WebServer, "Version")
		results[helper.Itoa64(c.Id)] = JSON.Type{
			"Version": c.Version,
			"result": result,
			"error": err,
		}
		return false
	})

	svn.Release()

	//TODO
	//版本同步确认后再清空
	//部署没有错误
	if errorCount == 0 {
		//清空未部署列表
		if err := svn.ClearDeployFile(); err != nil {
			return nil, err
		}
	}else {
		return results, helper.NewError("deplpy error")
	}

	return results, nil
}
