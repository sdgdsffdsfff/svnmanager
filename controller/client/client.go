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
		result, err := client.CallRpc(c, "RpcClient.CheckDeployPath", rpc.CheckDeployPathArgs{c.DeployPath})
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

func Edit(rend render.Render, req *http.Request){
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
	id :=helper.Int64(params["id"])
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
	body, err := jason.NewObjectFromReader(req.Body)

	fileIds, err := body.GetInt64Array("fileIds")

	host := client.FindFromCache(id)
	if host == nil {
		rend.JSON(200, helper.Error(helper.EmptyError, "Client is not found"))
		return
	}

	result, err := update(host, fileIds)
	if err != nil {
		rend.JSON(200, helper.Error(err))
		return
	}

	rend.JSON(200, helper.Success(result))
}

func Deploy(rend render.Render, req *http.Request, params martini.Params){

	id := helper.Int64(params["id"])

	host := client.FindFromCache(id)
	if host == nil {
		rend.JSON(200, helper.Error(helper.EmptyError, "Client is not found"))
		return
	}

	result, err := deploy(host)
	if err != nil {
		rend.JSON(200, helper.Error( err ))
		return
	}

	rend.JSON(200, helper.Success(result))
}
