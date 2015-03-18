package client

import (
	"king/model"
	"king/utils/db"
	"king/utils"
	"king/rpc"
	"king/helper"
	"king/utils/JSON"
	"king/bootstrap"
	"sync"
	"king/service/task"
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
	Message string
}

type HostMap map[int64]*HostClient

var hostMap HostMap
var lock sync.Mutex = sync.Mutex{}
var taskStarted = false

//多终端Rpc调用
//TODO
//队列调用，最大同时调用数
func BatchCallRpc(clients HostMap, method string, params ...rpc.RpcInterface) JSON.Type {
	results := JSON.Type{}
	helper.AsyncMap(clients, func(key, value interface{}) bool {
		c := value.(*HostClient)
		result, err := CallRpc(c, method, params...)
		if err != nil {
			results[ helper.Itoa64(c.Id) ] = helper.Error(err)
		} else {
			results[ helper.Itoa64(c.Id) ] = result
		}
		return false
	})
	return results
}

func CallRpc(client *HostClient, method string, params ...rpc.RpcInterface)(interface{}, error) {
	var param rpc.RpcInterface = &rpc.SimpleArgs{Id: client.Id}

	if len(params) > 0 && params[0] != nil {
		param = params[0]
		param.SetId(client.Id)
	}
	return rpc.Send(RpcIp(client), "RpcClient."+method, param)
}

func Fetch() (HostMap, error) {
	var list []*model.WebServer
	hostMap = HostMap{}
	_, err := db.Orm().QueryTable("web_server").All(&list)
	if err == nil {
		for _, webServer := range list {
			hostMap[webServer.Id] = &HostClient{webServer, Die, &ProcStat{}, ""}
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

//TODO 合并
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
func FindOrAppend(client *model.WebServer) int64 {
	found := false
	for _, c := range hostMap {
		if c.Ip == client.Ip && c.Port == client.Port {
			c.Status = Alive
			found = true
			return c.Id
		}
	}
	if !found {
		hostMap[client.Id] = &HostClient{client, Connecting, &ProcStat{}, ""}
	}
	return 0
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

func UpdateUsage(id int64, cpu, mem float64) {
	c := FindFromCache(id)
	if c != nil {
		c.Proc.CPUPercent = cpu
		c.Proc.MEMPercent = mem
	}
}

func Count() int {
	return helper.Cap(hostMap)
}

func Refresh() {
	helper.AsyncMap(hostMap, func(key, value interface{}) bool {
		c := value.(*HostClient)
		status := GetClientStatus(c)
		if status == Die {
			c.Proc.CPUPercent = 0
			c.Proc.MEMPercent = 0
		} else if c.Status == Die && status == Alive {
			ReportMeUsage(c)
		}
		if c.Status == Busy && status == Alive {

		} else {
			c.Status = status
		}
		return false
	})
}

func Add(client *model.WebServer) (int64, error) {
	return Active(client)
}

func Edit(client *model.WebServer, fields ...string) error {
	if _, err := db.Orm().Update(client, fields...); err != nil {
		return err
	}

	if c := FindFromCache(client.Id); c != nil {
		helper.ExtendStruct(c, client, fields...)
	}

	return nil
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

func Active(client *model.WebServer) (int64, error) {
	created, id, err := db.Orm().ReadOrCreate(client, "Ip", "InternalIp", "Port");
	if  err != nil {
		return id, err
	}
	if created || id > 0 {
		id = FindOrAppend(client)
		return id, nil
	}

	return id, helper.NewError( helper.AppendString("already exisits client, id: ", id))
}

func SetBusy(id int64, yes ...bool) {
	if client := FindFromCache(id); client != nil {
		isBusy := true
		if len(yes) > 0 {
			isBusy = yes[0]
		}
		if !isBusy {
			client.Status = Alive
		} else {
			client.Status = Busy
		}
	}
}

func SetMessage(id int64, message ...string) {
	if client := FindFromCache(id); client != nil {
		if len(message) > 0 {
			client.Message = message[0]
		} else {
			client.Message = ""
		}
	}
}

func ReportMeUsage(client ...*HostClient) interface{} {
	if len(client) > 0 {
		result, _ := CallRpc(client[0], "Procstat")
		return result
	} else {
		return BatchCallRpc(hostMap, "Procstat")
	}
}

func StartTask(){
	if taskStarted {
		return
	}
	taskStarted = true
	task.Trigger("Heartbeat")
}

func StopTask(){
	if !taskStarted {
		return
	}
	taskStarted = false
}

func RpcIp(client *HostClient) string {
	return "http://" + client.InternalIp + ":" + client.Port + "/rpc"
}

func init(){
	bootstrap.Register(func(){
		if db.IsConnected() {
			Fetch()
		}
	})
}
