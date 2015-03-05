package service

import (
	"king/model"
	"king/utils/db"
	"king/utils"
	"king/rpc"
	"king/helper"
	"time"
	"king/bootstrap"
)

type Status int

const (
	Die Status = iota
	Connecting
	Alive
	Busy
)

type HostClient struct{
	*model.WebServer
	Status Status
}

type clientService struct{
	list []*HostClient
	heartbeatEnable bool
}

var Client = &clientService{}

func (r *clientService) DisableHeartbeat() {
	r.heartbeatEnable = false
}

//多终端Rpc调用
//TODO
//队列调用，最大同时调用数
func (r *clientService) BatchCallRpc(clients []*HostClient, method string, params interface{}) []interface{} {
	var results []interface{}
	helper.AsyncMap(clients, func(index int) bool {
		client := clients[index]
		result, err := r.CallRpc(client, method, params)
		if err != nil {
			results = append(results, result)
		}
		return false
	})
	return results
}

func (r *clientService) CallRpc(client *HostClient, method string, params interface{})(interface{}, error) {
	return rpc.Send(r.RpcIp(client), method, params)
}

func (r *clientService) Fetch() ([]*HostClient, error) {
	var list []*model.WebServer
	r.list = []*HostClient{}
	_, err := db.Orm().QueryTable("web_server").All(&list)
	if err == nil {
		for _, webServer := range list {
			r.list = append(r.list, &HostClient{webServer, Connecting})
		}
	}
	return r.list, err
}

//参数为空或者是[0]代表获取所有主机
func (r *clientService) List( ids ...[]int64 ) ([]*HostClient) {
	if len(ids) == 1 && len(ids[0]) > 0 && ids[0][0] != 0 {
		list := []*HostClient{}
		idList := ids[0]
		helper.AsyncMap(idList, func(i int) bool{
			if client := r.FindFromCache(idList[i]); client != nil {
				list = append(list, client)
			}
			return false
		})
		return list
	}

	return r.list
}

func (r *clientService) Find(client model.WebServer) (*model.WebServer, error) {
	err := db.Orm().Read(&client)
	if err != nil {
		return &client, err
	}
	return &client, nil
}

//仅从缓存中查找
func (r *clientService) FindFromCache(id int64) *HostClient {
	for _, client := range r.list {
		if client.Id == id {
			return client
		}
	}
	return nil
}

//仅向缓存列表里添加，ip与port不能重复
func (r *clientService) FindOrAppend(client *model.WebServer) {
	found := false
	for _, c := range r.list {
		if c.Ip == client.Ip && c.Port == client.Port {
			c.Status = Alive
			found = true
		}
	}
	if !found {
		r.list = append(r.list , &HostClient{client, Connecting})
	}
}

func (r *clientService) Refresh() {
	if list := r.List(); len(list) > 0 {
		helper.AsyncMap(list, func(index int) bool {
			client := list[index]
			client.Status = r.GetClientStatus(client)
			return false
		})
	}
}

var refreshDuration = time.Second * 3
func (r *clientService) Heartbeat() {
	r.heartbeatEnable = true
	r.Fetch()
	for {
		if r.heartbeatEnable {
			r.Refresh()
			time.Sleep( refreshDuration )
		} else {
			break
		}
	}
}

func (r *clientService) Add(client *model.WebServer) (helper.ErrorType, error) {
	return r.Active(client)
}

func (r *clientService) Update(client *model.WebServer, fields ...string) error {
	if _, err := db.Orm().Update(client, fields...); err != nil {
		return err
	}

	if c := r.FindFromCache(client.Id); c != nil {
		helper.ExtendStruct(c, client, fields...)
	}

	return nil
}

func (r *clientService) Active(client *model.WebServer) (helper.ErrorType, error) {

	created, id, err := db.Orm().ReadOrCreate(client, "Ip", "Port");
	if  err != nil {
		return helper.DefaultError, err
	}
	if created || id > 0 {
		r.FindOrAppend(client)
	} else {
		return helper.ExistsError, helper.NewError( helper.AppendString("already exisits client, id: ", id))
	}

	return helper.DefaultError, nil
}

func (r *clientService) Del(client *model.WebServer) (error) {
	_, err := db.Orm().Delete(client);
	return err
}

func (r *clientService) GetClientStatus(client *HostClient) Status {
	isConnect, _ := utils.Ping( r.GetAvailableIp(client) + ":" + client.Port )
	if isConnect {
		return Alive
	}
	return Die
}

func (r *clientService) GetAvailableIp(client *HostClient) string {
	ip := ""
	if len(client.InternalIp) > 0 {
		ip = client.InternalIp
	} else if len(client.Ip) > 0 {
		ip = client.Ip
	}
	return ip
}

func (r *clientService) RpcIp(client *HostClient) string {
	return "http://" + client.InternalIp + ":" + client.Port + "/rpc"
}

func init(){
	bootstrap.Register(func(){
		if db.IsConnected() {
			Client.Heartbeat()
		}
	})
}
