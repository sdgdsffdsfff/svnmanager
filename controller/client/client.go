package client

import (
	"github.com/antonholmquist/jason"
	"github.com/go-martini/martini"
	"github.com/martini-contrib/render"
	"king/enum/status"
	"king/helper"
	"king/model"
	"king/rpc"
	"king/service/client"
	"king/service/group"
	"king/utils/JSON"
	"net/http"
	"king/config"
	"reflect"
	"fmt"
)

func fillClientList() map[string]JSON.Type {
	result := map[string]JSON.Type{}
	//组map
	clientsDict := map[string]JSON.Type{}

	//获取组
	groups := group.List()
	//获取主机
	clients := client.List()

	//创建组map
	for id, g := range groups {
		result[helper.Itoa64(id)] = JSON.Parse(g)
		clientsDict[helper.Itoa64(id)] = JSON.Type{}
	}

	for id, c := range clients {
		//添加到对应组map里
		if list, found := clientsDict[helper.Itoa64(c.Group)]; found {
			list[helper.Itoa64(id)] = c
		}
	}

	for id, g := range result {
		g["clients"] = clientsDict[id]
	}

	return result
}

func List(rend render.Render) {
	rend.JSON(200, helper.Success(fillClientList()))
}

func Refresh(rend render.Render) {
	_, err := client.Fetch()
	if err != nil {
		rend.JSON(200, helper.Error(err))
		return
	}
	rend.JSON(200, helper.Success(fillClientList()))
}

func Check(rend render.Render, req *http.Request) {
	body, _ := jason.NewObjectFromReader(req.Body)
	clientsId, _ := body.GetInt64Array("clientsId")
	clientList := client.List(clientsId)
	results := JSON.Type{}
	for _, c := range clientList {
		result, err := c.CallRpc("CheckDeployPath", rpc.CheckDeployPathArgs{c.Id, c.DeployPath})
		if err != nil {
			results[helper.Itoa64(c.Id)] = helper.Error(err)
		} else {
			results[helper.Itoa64(c.Id)] = result
		}
	}
	rend.JSON(200, helper.Success(results))
}

func Add(rend render.Render, req *http.Request) {
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

func Edit(rend render.Render, req *http.Request, params martini.Params) {

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

	c.Id = id;

	if err := client.Edit(c, keys...); err != nil {
		rend.JSON(200, helper.Error(err))
		return
	}
	rend.JSON(200, helper.Success(c))
}

func Del(rend render.Render, params martini.Params) {

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

func Update(rend render.Render, req *http.Request, params martini.Params) {
	id := helper.Int64(params["id"])
	host, errResponse := getClientWithNoBusyOrJSONError(id)
	if host == nil {
		rend.JSON(200, errResponse)
		return
	}

	//未部署列表
	unDeployFileList, err := host.GetUnDeployFiles()
	if err != nil || unDeployFileList == nil {
		rend.JSON(200, helper.Error(err))
		return
	}

	//选择要上传的文件
	body := JSON.FormRequest(req.Body)
	files := body["files"]

	upFile := rpc.UploadFileList{}
	JSON.ParseToStruct(JSON.Stringify(files), &upFile)

	uploadFiles := &rpc.UpdateArgs{
		Id: host.Id,
		ResPath: config.ResServer(),
		FileList: upFile,
		DeployPath: host.DeployPath,
	}

	//上传
	completeUploadList, err := host.CallRpc("Update", uploadFiles)
	if err != nil {
		rend.JSON(200, helper.Error(err))
		return
	}

	//上传ok就删除对应路径记录
	helper.Map(completeUploadList, func(key, value interface{}) bool {
		p := key.(string)
		v := reflect.ValueOf(value)
		if v.Kind() == reflect.Bool {
			delete(unDeployFileList, p)
		}
		return false
	})

	//更新记录
	host.WebServer.UnDeployList = JSON.Stringify(unDeployFileList)

	err = client.Edit(host.WebServer, "UnDeployList")
	if err != nil {
		fmt.Println(err)
	}
	rend.JSON(200, helper.Success(completeUploadList))
}

func Deploy(rend render.Render, req *http.Request, params martini.Params) {

	id := helper.Int64(params["id"])
	host, errResponse := getClientWithAliveOrJSONError(id)
	if host == nil {
		rend.JSON(200, errResponse)
		return
	}

	host.SetBusy()
	result, err := host.Deploy()
	if err != nil {
		rend.JSON(200, helper.Error(err))
		return
	}

	rend.JSON(200, helper.Success(result))
}

func Revert(rend render.Render, req *http.Request, params martini.Params) {

	id := helper.Int64(params["id"])
	host, errResponse := getClientWithAliveOrJSONError(id)
	if host == nil {
		rend.JSON(200, errResponse)
		return
	}

	body, _ := jason.NewObjectFromReader(req.Body)
	path, _ := body.GetString("path")

	_, err := host.Revert(path)
	if err != nil {
		rend.JSON(200, helper.Error(err))
		return
	}
	rend.JSON(200, helper.Success())
}

func RemoveBackup(rend render.Render, req *http.Request, params martini.Params) {
	id := helper.Int64(params["id"])
	host, errResponse := getClientWithAliveOrJSONError(id)
	if host == nil {
		rend.JSON(200, errResponse)
		return
	}

	body, _ := jason.NewObjectFromReader(req.Body)
	path, _ := body.GetString("path")

	_, err := host.RemoveBackup(path)
	if err != nil {
		rend.JSON(200, helper.Error(err))
		return
	}
	rend.JSON(200, helper.Success())
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

func GetBackupList(rend render.Render, params martini.Params) {

	id := helper.Int64(params["id"])
	host, errResponse := getClientWithNoBusyOrJSONError(id)
	if host == nil {
		rend.JSON(200, errResponse)
		return
	}

	result, err := host.CallRpc("GetBackupList")
	if err != nil {
		rend.JSON(200, helper.Error(err))
		return
	}

	rend.JSON(200, helper.Success(result))
}

func GetUnDeployFiles(rend render.Render, params martini.Params) {

	id := helper.Int64(params["id"])
	host, errResponse := getClientWithNoBusyOrJSONError(id)
	if host == nil {
		rend.JSON(200, errResponse)
		return
	}
	result, err := host.GetUnDeployFiles()
	if err != nil || result == nil || helper.Cap(result) == 0 {
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
		case status.Alive:
			return host, nil
		case status.Die:
			err = helper.Error(helper.OfflineError)
			break
		case status.Busy:
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
		case status.Busy:
			err = helper.Error(helper.BusyError)
			break
		default:
			return host, nil
		}
	}
	return nil, err
}
