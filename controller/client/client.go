package client

import (
	"github.com/martini-contrib/render"
	"github.com/go-martini/martini"
	"github.com/antonholmquist/jason"
	"king/service/client"
	"king/service/group"
	"king/utils/JSON"
	"king/helper"
	"net/http"
	"king/rpc"
	"king/model"
)

func List(rend render.Render){
	result := map[string]JSON.Type{}
	clientsDict := map[string]JSON.Type{}

	groups := group.List()
	clients := client.List()

	for id, g := range groups {
		result[helper.Itoa64(id)] = JSON.Parse(g)
		clientsDict[helper.Itoa64(id)] = JSON.Type{}
	}

	for id, c := range clients {
		if list, found := clientsDict[helper.Itoa64(c.Group)]; found {
			list[helper.Itoa64(id)] = c
		}
	}

	for id, g := range result {
		g["Clients"] = clientsDict[id]
	}

	rend.JSON(200, helper.Success(result))
}

func Check(rend render.Render, req *http.Request) {
	body, _ := jason.NewObjectFromReader(req.Body)
	clientsId, _ := body.GetInt64Array("clientsId")
	clientList := client.List(clientsId)
	results := JSON.Type{}
	for _, c := range clientList {
		result, err := c.CallRpc("CheckDeployPath", rpc.CheckDeployPathArgs{c.Id, c.DeployPath})
		if err != nil {
			results[ helper.Itoa64(c.Id) ] = helper.Error(err)
		} else {
			results[ helper.Itoa64(c.Id) ] = result
		}
	}
	rend.JSON(200, helper.Success(results))
}

func Add(rend render.Render, req *http.Request){
	c := &model.WebServer{}
	body := JSON.FormRequest(req.Body)

	if err := JSON.ParseToStruct(JSON.Stringify(body), c); err != nil {
		rend.JSON(200, helper.Error(helper.ParamsError, err))
		return
	}

	if _, err := client.Add(c); err != nil {
		rend.JSON(200, helper.Error(err))
		return
	}

	rend.JSON(200, helper.Success(c))
}

func Edit(rend render.Render, req *http.Request, params martini.Params){

	id := helper.Int64(params["id"])
	host, errResponse := getClientWithNoBusyOrJSONError(id)
	if host == nil {
		rend.JSON(200, errResponse)
		return
	}

	body := JSON.FormRequest(req.Body)
	c := &model.WebServer{}
	if err := JSON.ParseToStruct(JSON.Stringify(body), c); err != nil {
		rend.JSON(200, helper.Error(helper.ParamsError))
		return
	}
	keys := JSON.GetKeys(body)

	if err := client.Edit(c, keys...); err != nil {
		rend.JSON(200, helper.Error(err))
		return
	}
	rend.JSON(200, helper.Success(client.FindFromCache(c.Id)))
}

func Del(rend render.Render, params martini.Params){

	id := helper.Int64(params["id"])
	host, errResponse := getClientWithNoBusyOrJSONError(id)
	if host == nil {
		rend.JSON(200, errResponse)
		return
	}

	c := model.WebServer{Id: id}
	if err := client.Del(&c); err != nil {
		rend.JSON(200, helper.Error(err))
		return
	}
	rend.JSON(200, helper.Success())
}

func Move(rend render.Render, params martini.Params) {
	id := helper.Int64(params["id"])
	gid := helper.Int64(params["gid"])

	c := model.WebServer{Id: id, Group: gid}
	if err := client.Edit(&c, "Group"); err != nil {
		rend.JSON(200, helper.Error(err))
		return
	}

	rend.JSON(200, helper.Success())
}

func Update(rend render.Render, req *http.Request, params martini.Params){

	id := helper.Int64(params["id"])
	host, errResponse := getClientWithAliveOrJSONError(id)
	if host == nil {
		rend.JSON(200, errResponse)
		return
	}

	body, err := jason.NewObjectFromReader(req.Body)
	fileIds, err := body.GetInt64Array("fileIds")

	result, err := host.Update(fileIds)
	if err != nil {
		rend.JSON(200, helper.Error(err))
		return
	}

	rend.JSON(200, helper.Success(result))
}

func Deploy(rend render.Render, req *http.Request, params martini.Params){

	id := helper.Int64(params["id"])
	host, errResponse := getClientWithAliveOrJSONError(id)
	if host == nil {
		rend.JSON(200, errResponse)
		return
	}

	host.SetBusy()
	result, err := host.Deploy()
	if err != nil {
		rend.JSON(200, helper.Error( err ))
		return
	}

	rend.JSON(200, helper.Success(result))
}

func ShowLog(rend render.Render, params martini.Params) {

	id := helper.Int64(params["id"])
	host, errResponse := getClientOrJSONError(id)
	if host == nil {
		rend.JSON(200, errResponse)
		return
	}

	result, err := host.CallRpc("ShowLog", &rpc.SimpleArgs{Id: host.Id})
	if err != nil {
		rend.JSON(200, helper.Error(err))
		return
	}

	rend.JSON(200, helper.Success(result))
}


func getClientOrJSONError(id int64) (*client.HostClient, JSON.Type) {
	host := client.FindFromCache(id)
	if host == nil {
		return nil, helper.Error(helper.EmptyError, "Client is not found")
	}
	return host, nil
}

func getClientWithAliveOrJSONError(id int64) (*client.HostClient, JSON.Type) {
	host, err := getClientOrJSONError(id)
	if host != nil {
		switch host.Status {
		case client.Alive:
			return host, nil
		case client.Die:
			err = helper.Error(helper.OfflineError)
			break
		case client.Busy:
			err = helper.Error(helper.BusyError)
			break
		}
	}
	return nil, err
}

func getClientWithNoBusyOrJSONError(id int64) (*client.HostClient, JSON.Type) {
	host, err := getClientOrJSONError(id)
	if host != nil {
		switch host.Status {
		case client.Busy:
			err = helper.Error(helper.BusyError)
			break
		default:
			return host, nil
		}
	}
	return nil, err
}
