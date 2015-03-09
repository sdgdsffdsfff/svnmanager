package client

import (
	"king/rpc"
	"king/utils/JSON"
	"net/http"
	"king/service/proc"
)


//rpc method
type RpcProcstat struct{}

func (h *RpcProcstat) Stat(r *http.Request, args *JSON.Type, reply *rpc.RpcReply) error {
	cpu := proc.CPUPercent()
	mem := proc.MEMPercent()
	reply.Response = JSON.Type{
		"CPUPercent" : cpu,
		"MEMPercent" : mem,
	}
	return nil
}

func init(){
	rpc.AddCtrl(new(RpcProcstat))
}
