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

	if errType, err := client.Add(c); err != nil {
		rend.JSON(200, helper.Error(errType, err))
		return
	}

	rend.JSON(200, helper.Success(c))
}

func Update(rend render.Render, req *http.Request){
	body := JSON.FormRequest(req.Body)
	c := &model.WebServer{}
	if err := JSON.ParseToStruct(JSON.Stringify(body), c); err != nil {
		rend.JSON(200, helper.Error(helper.ParamsError))
		return
	}
	keys := JSON.GetKeys(body)

	if err := client.Update(c, keys...); err != nil {
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
	if err := client.Update(&c, "Group"); err != nil {
		rend.JSON(200, helper.Error(err))
		return
	}

	rend.JSON(200, helper.Success())
}

func Deploy(rend render.Render, req *http.Request, params martini.Params){
	body, err := jason.NewObjectFromReader(req.Body)
	if err != nil {
		rend.JSON(200, helper.Error(helper.ParamsError))
		return
	}
	//如果数据为[0]则表示up_file表中的全部文件
	filesId, err := body.GetInt64Array("filesId")
	if err != nil {
		rend.JSON(200, helper.Error(helper.ParamsError))
		return
	}
	clientsId, err := body.GetInt64Array("clientsId")
	if err != nil {
		rend.JSON(200, helper.Error(helper.ParamsError))
		return
	}

	message, err := body.GetString("message")

	result, err := deploy( filesId, clientsId, message )
	//报告错误原因
	if err != nil {
		rend.JSON(200, helper.Error(err, result))
		return
	}
	rend.JSON(200, helper.Success(result))
}
