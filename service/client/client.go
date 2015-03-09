package client

import (
	"king/model"
	"king/utils/db"
	"king/utils"
	"king/rpc"
	"king/helper"
	"time"
	"king/utils/JSON"
	"king/bootstrap"
)

type Status int

const (
	Die Status = iota
	Connecting
	Alive
	Busy
)

type ProcStat struct {
	CPUPercent float64
	MEMPercent float64
}

type HostClient struct{
	*model.WebServer
	Status Status
	Proc *ProcStat
}

var heartbeatEnable bool
var procMonitorEnable bool
var hostList []*HostClient
var refreshDuration = time.Second * 3

func DisableHeartbeat() {
	heartbeatEnable = false
}

//多终端Rpc调用
//TODO
//队列调用，最大同时调用数
func BatchCallRpc(clients []*HostClient, method string, params interface{}) JSON.Type {
	results := JSON.Type{}
	helper.AsyncMap(clients, func(index int) bool {
		c := clients[index]
		result, err := CallRpc(c, method, params)
		if err != nil {
			results[ helper.Itoa64(c.Id) ] = helper.Error(err)
		} else {
			results[ helper.Itoa64(c.Id) ] = result
		}
		return false
	})
	return results
}

func CallRpc(client *HostClient, method string, params interface{})(interface{}, error) {
	return rpc.Send(RpcIp(client), method, params)
}

func Fetch() ([]*HostClient, error) {
	var list []*model.WebServer
	hostList = []*HostClient{}
	_, err := db.Orm().QueryTable("web_server").All(&list)
	if err == nil {
		for _, webServer := range list {
			hostList = append(hostList, &HostClient{webServer, Connecting, &ProcStat{}})
		}
	}
	return hostList, err
}

//参数为空或者是[0]代表获取所有主机
func List( ids ...[]int64 ) ([]*HostClient) {
	if len(ids) == 1 && len(ids[0]) > 0 && ids[0][0] != 0 {
		list := []*HostClient{}
		idList := ids[0]
		helper.AsyncMap(idList, func(i int) bool{
			if client := FindFromCache(idList[i]); client != nil {
				list = append(list, client)
			}
			return false
		})
		return list
	}

	return hostList
}

func Find(client model.WebServer) (*model.WebServer, error) {
	err := db.Orm().Read(&client)
	if err != nil {
		return &client, err
	}
	return &client, nil
}

//仅从缓存中查找
func FindFromCache(id int64) *HostClient {
	for _, client := range hostList {
		if client.Id == id {
			return client
		}
	}
	return nil
}

//仅向缓存列表里添加，ip与port不能重复
func FindOrAppend(client *model.WebServer) {
	found := false
	for _, c := range hostList {
		if c.Ip == client.Ip && c.Port == client.Port {
			c.Status = Alive
			found = true
		}
	}
	if !found {
		hostList = append(hostList , &HostClient{client, Connecting, &ProcStat{}})
	}
}

func GetAliveList() []*HostClient {
	aliveList := []*HostClient{}
	helper.Map(hostList, func(i int) bool{
		host := hostList[i]
		if host.Status == Alive {
			aliveList = append( aliveList, host )
		}
		return false
	})
	return aliveList
}

func Refresh() {
	if list := List(); len(list) > 0 {
		helper.AsyncMap(list, func(index int) bool {
			client := list[index]
			client.Status = GetClientStatus(client)
			return false
		})
	}
}

func Heartbeat() {
	heartbeatEnable = true
	for {
		if heartbeatEnable {
			Refresh()
			time.Sleep( refreshDuration )
		} else {
			break
		}
	}
}

//获取cpu与内存使用状况
func GetProcStat() {
	procMonitorEnable = true
	for {
		if procMonitorEnable {
			aliveList := GetAliveList()
			results := BatchCallRpc(aliveList, "RpcProcstat.Stat", nil)
			for key, value := range results {
				if c := FindFromCache(helper.Int64(key)); c != nil {
					proc := &ProcStat{}
					JSON.ParseToStruct(value, proc)
					c.Proc = proc
				}
			}
			time.Sleep( refreshDuration )
		}else{
			break
		}
	}
}

func Add(client *model.WebServer) (helper.ErrorType, error) {
	return Active(client)
}

func Update(client *model.WebServer, fields ...string) error {
	if _, err := db.Orm().Update(client, fields...); err != nil {
		return err
	}

	if c := FindFromCache(client.Id); c != nil {
		helper.ExtendStruct(c, client, fields...)
	}

	return nil
}

func Active(client *model.WebServer) (helper.ErrorType, error) {

	created, id, err := db.Orm().ReadOrCreate(client, "Ip", "Port");
	if  err != nil {
		return helper.DefaultError, err
	}
	if created || id > 0 {
		FindOrAppend(client)
	} else {
		return helper.ExistsError, helper.NewError( helper.AppendString("already exisits client, id: ", id))
	}

	return helper.DefaultError, nil
}

func Del(client *model.WebServer) (error) {
	_, err := db.Orm().Delete(client);
	return err
}

func GetClientStatus(client *HostClient) Status {
	isConnect, _ := utils.Ping( GetAvailableIp(client) + ":" + client.Port )
	if isConnect {
		return Alive
	}
	return Die
}

func GetAvailableIp(client *HostClient) string {
	ip := ""
	if len(client.InternalIp) > 0 {
		ip = client.InternalIp
	} else if len(client.Ip) > 0 {
		ip = client.Ip
	}
	return ip
}

func RpcIp(client *HostClient) string {
	return "http://" + client.InternalIp + ":" + client.Port + "/rpc"
}

func init(){
	bootstrap.Register(func(){
		if db.IsConnected() {
			Fetch()
			go GetProcStat()
			Heartbeat()
		}
	})
}
