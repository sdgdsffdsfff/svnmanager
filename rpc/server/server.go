package server

import (
	"net/http"
	"king/utils/JSON"
	"king/rpc"
	"king/model"
	"king/service/client"
	"king/helper"
	"king/service/webSocket"
)

//rpc method
type RpcServer struct {}

func (h *RpcServer) Status(r *http.Request, args *JSON.Type, reply *rpc.RpcReply) error {
	reply.Response = JSON.Type{
		"Ip": "ok",
		"Status": "ok",
	}
	return nil
}

//客户端启动通知，保存客户端入库
func (h *RpcServer) Active(r *http.Request, args *model.WebServer, reply *rpc.RpcReply) error {
	id, err := client.Active(args)
	if err != nil {
		return helper.NewError("add client error", err)
	}
	reply.Response =
	return nil
}

func (h *RpcServer) Message(r *http.Request, args *JSON.Type, reply *rpc.RpcReply) error {
	reply.Response = args
	return nil
}

func (h *RpcServer) ReportUsage(r *http.Request, args *rpc.UsageArgs, reply *rpc.RpcReply) error {
	if args.Id > 0 {
		client.UpdateUsage(args.Id, args.CPUPercent, args.MEMPercent)
	}
	reply.Response = true
	return nil
}

func (h *RpcServer) BroadCastAll(r *http.Request, args *webSocket.Message, reply *rpc.RpcReply) error {
	webSocket.BroadCastAll(args)
	reply.Response = true
	return nil
}

//controller service
type report struct{
	List []*JSON.Type
}

func (r report) String() string {
	return JSON.Stringify(r.Data())
}

func (r report) Data() interface {} {
	return r.List;
}

var ReportService *report = &report{}

func init(){
	rpc.AddCtrl(new(RpcServer))
}
