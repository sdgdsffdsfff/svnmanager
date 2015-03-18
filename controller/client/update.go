package client

import (
	"king/service/client"
	"king/service/svn"
	"king/utils/JSON"
	"king/service/webSocket"
	"king/rpc"
)

func update(host *client.HostClient, fileIds []int64) (JSON.Type, error){
	result := JSON.Type{}

	fileList, err := svn.GetUnDeployFileList(fileIds)
	if err != nil {
		return result, err
	}

	webSocket.BroadCastAll(&webSocket.Message{"lock", nil})

	data, err := client.CallRpc(host, "RpcClient.Update", rpc.UpdateArgs{host.Id,fileList, host.DeployPath})
	if err != nil {
		return result, err
	}

	host.Version = svn.Version
	err = client.Edit(host.WebServer, "Version")

	result["Version"] = host.Version
	result["Rpc"] = data
	result["Error"] = err

	return result, nil
}
