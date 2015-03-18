package client

import (
	"king/service/client"
)

func deploy(host *client.HostClient) (interface{},error) {
	host.Message = "Ready to deploy.."
	result, err := client.CallRpc(host, "RpcClient.Deploy" , nil)
	if err != nil {
		return nil, err
	}
	host.Message = "Deploying.."
	return result, nil
}
