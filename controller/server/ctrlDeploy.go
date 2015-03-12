package server

import (
	"king/helper"
	"king/utils/JSON"
	"king/service/svn"
	"king/service/client"
	"king/service/webSocket"
	"king/rpc"
)

func DeployCtrl(filesId []int64, clientsId []int64, message string) (JSON.Type, error) {
	results := JSON.Type{}
	errorCount := 0

	fileList, err := svn.GetUnDeployFileList(filesId)
	if err != nil {
		return results, err
	}

	svn.GetLock()
	webSocket.BroadCastAll(&webSocket.Message{"lock", nil})

	svn.DeployMessage(message)

	clients := client.List(clientsId)
	helper.AsyncMap(clients, func(key, value interface{}) bool {
		c := value.(*client.HostClient)
		result, err := client.CallRpc(c, "RpcClient.Deploy", rpc.DeployArgs{fileList, c.DeployPath})
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
	webSocket.BroadCastAll(&webSocket.Message{"unlock", nil})

	//假设所有主机必须版本一致
	if errorCount == 0 {
		if len(clients) == client.Count() {
			if err := svn.ClearDeployFile(); err != nil {
				return nil, err
			}
		}
	}else {
		return results, helper.NewError("deplpy error")
	}

	return results, nil
}
