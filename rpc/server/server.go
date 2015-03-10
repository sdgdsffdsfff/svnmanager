package server

import (
	"net/http"
	"king/utils/JSON"
	"king/rpc"
	"king/model"
	"king/service/client"
	"king/helper"
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

	if _, err := client.Active(args); err != nil {
		return helper.NewError("add client error", err)
	}
	reply.Response = true
	return nil
}

func (h *RpcServer) Message(r *http.Request, args *JSON.Type, reply *rpc.RpcReply) error {
	reply.Response = args
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

func (r *report) updateClientStatus(c *JSON.Type){
	//	for k, v := range r.List {
	//		if v.Ip == c.Ip {
	//			r.List[k] = c
	//			return
	//		}
	//	}
	//	r.List = append(r.List, c)
}

var ReportService *report = &report{}

func init(){
	rpc.AddCtrl(new(RpcServer))
	//WebSocketService.Exports("ClientStatus", ReportService)
}
