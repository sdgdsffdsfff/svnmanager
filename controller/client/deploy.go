package client

import (
	"king/service/client"
)

func deploy(host *client.HostClient) (interface{},error) {
	result, err := client.CallRpc(host, "RpcClient.Deploy" , nil)
	if err != nil {
		return nil, err
	}
	return result, nil
}
