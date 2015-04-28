package client

import (
	"king/bootstrap"
	"king/enum/status"
	"king/helper"
	"king/model"
	"king/rpc"
	"king/service/task"
	//"king/service/webSocket"
	"king/utils"
	"king/utils/JSON"
	"king/utils/db"
	"sync"
	"time"
)

type ProcStat struct {
	CPUPercent float64		`json:"cpu_percent"`
	MEMPercent float64		`json:"mem_percent"`
}

type HostClient struct {
	*model.WebServer		`json:"web_server"`
	Status  status.Status	`json:"status"`
	Proc    *ProcStat		`json:"proc"`
	Message string			`json:"message"`
	Error   string			`json:"error"`
}

func (r *HostClient) SetBusy(yes ...bool) {
	isBusy := true
	if len(yes) > 0 {
		isBusy = yes[0]
	}
	if !isBusy {
		r.Status = status.Alive
	} else {
		r.Status = status.Busy
	}
}

func (r *HostClient) SetMessage(message ...string) {
	if len(message) > 0 {
		r.Message = message[0]
	} else {
		r.Message = ""
	}
}

func (r *HostClient) SetError(err ...string) {
	if len(err) > 0 {
		r.Error = err[0]
	} else {
		r.Error = ""
	}
}

func (r *HostClient) CallRpc(method string, params ...rpc.RpcInterface) (interface{}, error) {
	var param rpc.RpcInterface = &rpc.SimpleArgs{Id: r.Id}

	if len(params) > 0 && params[0] != nil {
		param = params[0]
		param.SetId(r.Id)
	}

	return rpc.Send(r.RpcIp(), "RpcClient."+method, param)
}

func (r *HostClient) RpcIp() string {
	return "http://" + r.InternalIp + ":" + r.Port + "/rpc"
}

func (r *HostClient) ReportMeUsage() interface{} {
	result, _ := r.CallRpc("Procstat")
	return result
}

func (r *HostClient) GetStatus() status.Status {
	isConnect, _ := utils.Ping(r.GetAvailableIp() + ":" + r.Port)
	if isConnect {
		return status.Alive
	}
	return status.Die
}

func (r *HostClient) GetAvailableIp() string {
	ip := ""
	if len(r.InternalIp) > 0 {
		ip = r.InternalIp
	} else if len(r.Ip) > 0 {
		ip = r.Ip
	}
	return ip
}

func (r *HostClient) GetUnDeployFiles() (JSON.Type, error) {
	c := model.WebServer{Id: r.Id}
	err := db.Orm().Read(&c)
	if err != nil {
		return nil, err
	}
	return JSON.Parse(c.UnDeployList), nil
}

func (r *HostClient) Deploy() (interface{}, error) {
	r.Message = "ready to deploy.."
	result, err := r.CallRpc("Deploy", rpc.SimpleArgs{Id: r.Id})
	if err != nil {
		return nil, err
	}
	r.Message = "deploying.."
	return result, nil
}

func (r *HostClient) Revert(path string) (interface{}, error) {
	r.SetBusy()
	r.SetMessage("revert to " + path)
	defer r.SetBusy(false)

	result, err := r.CallRpc("Revert", rpc.SimpleArgs{
		Id:      r.Id,
		Message: path,
	})

	if err != nil {
		r.SetError("revert to " + path + " failure")
		return nil, err
	}
	r.SetMessage("revert complete")
	time.Sleep(time.Second * 2)
	r.SetMessage()
	return result, nil
}

func (r *HostClient) RemoveBackup(path string) (interface{}, error) {
	r.SetBusy()
	r.SetMessage("removing " + path)
	defer r.SetBusy(false)

	result, err := r.CallRpc("RemoveBackup", rpc.SimpleArgs{
		Id:      r.Id,
		Message: path,
	})

	if err != nil {
		r.SetError("remove " + path + " failure")
		return nil, err
	}
	r.SetMessage("remove complete")
	time.Sleep(time.Second * 2)
	r.SetMessage()
	return result, nil
}

func (r *HostClient) UpdateUnDeployList(list *JSON.Type) error {
	db.Orm().Read(r.WebServer)
	newList := JSON.Extend(r.WebServer.UnDeployList, list)
	r.WebServer.UnDeployList = JSON.Stringify(newList)
	_, err := db.Orm().Update(r.WebServer, "UnDeployList")
	if err != nil {
		return err
	}
	return nil
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
		result, err := c.CallRpc(method, params...)
		if err != nil {
			results[helper.Itoa64(c.Id)] = helper.Error(err)
		} else {
			results[helper.Itoa64(c.Id)] = result
		}
		return false
	})
	return results
}

func Fetch() (HostMap, error) {
	var list []*model.WebServer
	hostMap = HostMap{}
	_, err := db.Orm().QueryTable("web_server").All(&list)
	if err == nil {
		for _, webServer := range list {
			hostMap[webServer.Id] = &HostClient{webServer, status.Die, &ProcStat{}, "", ""}
		}
	}
	return hostMap, err
}

//参数为空或者是[0]代表获取所有主机
func List(ids ...[]int64) HostMap {
	if len(ids) == 1 && len(ids[0]) > 0 && ids[0][0] != 0 {
		list := HostMap{}
		idList := ids[0]
		helper.AsyncMap(idList, func(key, value interface{}) bool {
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
			c.Status = status.Alive
			found = true
			return c.Id
		}
	}
	if !found {
		hostMap[client.Id] = &HostClient{client, status.Die, &ProcStat{}, "", ""}
	}
	return 0
}

func GetAliveList() HostMap {
	aliveHostMap := HostMap{}
	for id, c := range hostMap {
		if c.Status == status.Alive {
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
		s := c.GetStatus()
		if s == status.Die {
			c.Proc.CPUPercent = 0
			c.Proc.MEMPercent = 0
		} else if c.Status == status.Die && s == status.Alive {
			c.ReportMeUsage()
		}
		if c.Status == status.Busy && s == status.Alive {

		} else {
			c.Status = s
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

func Del(client *model.WebServer) error {
	id := client.Id
	if _, err := db.Orm().Delete(client); err != nil {
		return err
	}
	delete(hostMap, id)
	return nil
}

func Active(client *model.WebServer) (int64, error) {
	created, id, err := db.Orm().ReadOrCreate(client, "Ip", "InternalIp", "Port")
	if err != nil {
		return id, err
	}
	if created || id > 0 {
		id = FindOrAppend(client)
		return id, nil
	}

	return id, helper.NewError(helper.AppendString("already exisits client, id: ", id))
}

func StartTask() {
	if taskStarted {
		return
	}
	taskStarted = true
	task.Trigger("Heartbeat")
}

func StopTask() {
	if !taskStarted {
		return
	}
	taskStarted = false
}

func init() {
	bootstrap.Register(func() {
		if db.IsConnected() {
			Fetch()
		}
	})
}
