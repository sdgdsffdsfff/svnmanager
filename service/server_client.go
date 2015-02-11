package service

import(
	"king/model"
	"king/utils/db"
	"king/utils"
	"king/rpc"
	"king/utils/JSON"
	"king/helper"
	"time"
	"king/bootstrap"
)

var RefreshDuration = time.Second * 3

type Status int

const (
	Die Status = iota
	Connecting
	Alive
	Busy
)

type Client struct{
	*model.WebServer
	Status Status
}

type clientService struct{
	list []*Client
	heartbeatEnable bool
}

func (r *clientService) DisableHeartbeat() {
	r.heartbeatEnable = false
}

//多终端Rpc调用
func (r *clientService) BatchCall(clients []*Client, method string, params interface{}) []JSON.Type {
	results := []JSON.Type{}
	helper.AsyncMap(clients, func(index int) bool {
		client := clients[index]
		res, err := rpc.Send(r.RpcIp(client), method, params)
		result := JSON.Type{"client": client, "result": res, "error": nil}
		if err != nil {
			result["error"] = err.Error()
		}
		results = append(results, result)
		return false
	})
	return results
}

func (r *clientService) Fetch() ([]*Client, error) {
	var list []*model.WebServer
	r.list = []*Client{}
	_, err := db.Orm().QueryTable("web_server").All(&list)
	if err == nil {
		for _, webServer := range list {
			r.list = append(r.list, &Client{webServer, Connecting})
		}
	}
	return r.list, err
}

func (r *clientService) List() ([]*Client) {
	return r.list
}

func (r *clientService) Find(client model.WebServer) (*model.WebServer, error) {
	err := db.Orm().Read(&client)
	if err != nil {
		return &client, err
	}
	return &client, nil
}

func (r *clientService) FindFromCache(id int) *Client {
	for _, client := range r.list {
		if client.Id == id {
			return client
		}
	}
	return nil
}

func (r *clientService) FindOrAppend(client *model.WebServer) {
	found := false
	for _, c := range r.list {
		if c.Ip == client.Ip && c.Port == client.Port {
			c.Status = Alive
			found = true
		}
	}
	if !found {
		r.list = append(r.list , &Client{client, Connecting})
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

func (r *clientService) Heartbeat() {
	r.heartbeatEnable = true
	r.Fetch()
	for {
		if r.heartbeatEnable {
			r.Refresh()
			time.Sleep( RefreshDuration )
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

	//用反射自动填充
	if c := r.FindFromCache(client.Id); c != nil {
		c.Name = client.Name
		c.Ip = client.Ip
		c.InternalIp = client.InternalIp
		c.DeployPath = client.DeployPath
		c.Port = client.Port
	}

	return nil
}

func (r *clientService) Active(client *model.WebServer) (helper.ErrorType, error) {
	created, id, err := db.Orm().ReadOrCreate(client, "Ip", "Port");
	if  err != nil {
		return helper.DefaultError, err
	}

	if created {
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

func (r *clientService) GetClientStatus(client *Client) Status {
	isConnect, _ := utils.Ping( r.GetAvailableIp(client) + ":" + client.Port )
	if isConnect {
		return Alive
	}
	return Die
}

func (r *clientService) GetAvailableIp(client *Client) string {
	ip := ""
	if len(client.InternalIp) > 0 {
		ip = client.InternalIp
	} else if len(client.Ip) > 0 {
		ip = client.Ip
	}
	return ip
}

func (r *clientService) RpcIp(client *Client) string {
	return "http://" + client.InternalIp + ":" + client.Port + "/rpc"
}

var ClientService = &clientService{}

func init(){
	bootstrap.Register(func(){
		if db.IsConnected() {
			ClientService.Heartbeat()
		}
	})
}
