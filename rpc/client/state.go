package client

import (
	"king/rpc"
	"king/utils/JSON"
	"net/http"
)


//rpc method
type RpcState struct{}

func (h *RpcState) Alive(r *http.Request, args *JSON.Type, reply *rpc.RpcReply) error {
	reply.Response = true
	return nil
}

func init(){
	rpc.AddCtrl(new(RpcState))
}
