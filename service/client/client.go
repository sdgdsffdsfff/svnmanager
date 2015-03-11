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

type HostMap map[int64]*HostClient

var heartbeatEnable bool
var procMonitorEnable bool
var hostMap HostMap
var refreshDuration = time.Second * 3

func DisableHeartbeat() {
	heartbeatEnable = false
}

//多终端Rpc调用
//TODO
//队列调用，最大同时调用数
func BatchCallRpc(clients HostMap, method string, params interface{}) JSON.Type {
	results := JSON.Type{}
	helper.AsyncMap(clients, func(key, value interface{}) bool {
		c := value.(*HostClient)
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

func Fetch() (HostMap, error) {
	var list []*model.WebServer
	hostMap = HostMap{}
	_, err := db.Orm().QueryTable("web_server").All(&list)
	if err == nil {
		for _, webServer := range list {
			hostMap[webServer.Id] = &HostClient{webServer, Connecting, &ProcStat{}}
		}
	}
	return hostMap, err
}

//参数为空或者是[0]代表获取所有主机
func List( ids ...[]int64 ) (HostMap) {
	if len(ids) == 1 && len(ids[0]) > 0 && ids[0][0] != 0 {
		list := HostMap{}
		idList := ids[0]
		helper.AsyncMap(idList, func(key, value interface{}) bool{
			if c := FindFromCache(value.(int64)); c != nil {
				list[c.Id] = c
			}
			return false
		})
		return list
	}

	return hostMap
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
	if c, found := hostMap[id]; found {
		return c
	}
	return nil
}

//仅向缓存列表里添加，ip与port不能重复
func FindOrAppend(client *model.WebServer) {
	found := false
	for _, c := range hostMap {
		if c.Ip == client.Ip && c.Port == client.Port {
			c.Status = Alive
			found = true
			break
		}
	}
	if !found {
		hostMap[client.Id] = &HostClient{client, Connecting, &ProcStat{}}
	}
}

func GetAliveList() HostMap {
	aliveHostMap := HostMap{}
	for id, c := range hostMap {
		if c.Status == Alive {
			aliveHostMap[id] = c
		}
	}
	return aliveHostMap
}

func Count() int {
	return helper.Cap(hostMap)
}

func Refresh() {
	if list := List(); len(list) > 0 {
		for _, c := range hostMap {
			c.Status = GetClientStatus(c)
		}
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
			results := BatchCallRpc(aliveList, "RpcClient.ProcStat", nil)
			for key, value := range results {
				if c := FindFromCache(helper.Int64(key)); c != nil && value != nil {
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
	id := client.Id
	if _, err := db.Orm().Delete(client); err != nil {
		return err
	}
	delete(hostMap, id)
	return nil
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
