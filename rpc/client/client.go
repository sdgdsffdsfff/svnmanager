package client

import (
	"king/rpc"
	"net/http"
	"king/service/task"
	"king/utils"
	"king/service/host"
)

//rpc method
type RpcClient struct{}

func (h *RpcClient) CheckDeployPath(r *http.Request, args *rpc.CheckDeployPathArgs, reply *rpc.RpcReply) error {
	if err := utils.PathEnable(args.Path); err != nil {
		reply.Response = err.Error()
	} else {
		reply.Response = true
	}
	return nil
}

func (h *RpcClient) Procstat(r *http.Request, args *rpc.SimpleArgs, reply *rpc.RpcReply) error {
	host.Active(args.Id)

	task.Trigger("ProcStat")
	reply.Response = true
	return nil
}

func (h *RpcClient) Update(r *http.Request, args *rpc.UpdateArgs, reply *rpc.RpcReply) error {
	host.Active(args.Id)

	reply.Response = host.Update(args.FileUrl, args.DeployPath)
	return nil
}

func (h *RpcClient) Deploy(r *http.Request, args *rpc.SimpleArgs, reply *rpc.RpcReply) error {
	host.Active(args.Id)

	task.Trigger("Deploy")
	reply.Response = true
	return nil
}

func (h *RpcClient) Revert(r *http.Request, args *rpc.SimpleArgs, reply *rpc.RpcReply) error {
	host.Active(args.Id)

	err := host.Revert(args.Message)
	if err != nil {
		return err
	}
	reply.Response = true
	return nil
}

func (h *RpcClient) GetBackupList(r *http.Request, args *rpc.SimpleArgs, reply *rpc.RpcReply) error {
	host.Active(args.Id)
	result, err := host.GetBackupList()
	if err != nil {
		return err
	}
	reply.Response = result
	return nil
}

func (h *RpcClient) ShowLog(r *http.Request, args *rpc.SimpleArgs, reply *rpc.RpcReply) error {
	host.Active(args.Id)
	output, err := host.ShowLog()
	reply.Response = output
	return err
}

func init(){
	rpc.AddCtrl(new(RpcClient))
}
