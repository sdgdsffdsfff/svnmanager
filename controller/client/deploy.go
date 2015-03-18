package client

import (
	"king/service/client"
	"king/rpc"
)

func deploy(host *client.HostClient) (interface{},error) {
	host.Message = "Ready to deploy.."
	result, err := client.CallRpc(host, "Deploy" , rpc.SimpleArgs{Id: host.Id})
	if err != nil {
		return nil, err
	}
	host.Message = "Deploying.."
	return result, nil
}
