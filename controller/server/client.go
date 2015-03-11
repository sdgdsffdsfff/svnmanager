package server

import (
	"github.com/go-martini/martini"
	"github.com/martini-contrib/render"
	"king/config"
	"king/helper"
	"king/service/client"
	"king/service/group"
	"king/utils/JSON"
	"king/model"
	"net/http"
	"github.com/antonholmquist/jason"
	"king/rpc"
)

type HostCtrl struct{}

type listmap map[string]JSON.Type

func init() {
	config.AppendValue(config.Controller, &HostCtrl{})
}

func (ctn *HostCtrl) SetRouter(m *martini.ClassicMartini) {

	m.Group("/aj/client", func(r martini.Router) {

		r.Get("/list", func(rend render.Render) {
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
		})

		r.Get("/heartbeat", func(rend render.Render){
			list := client.List()
			result := []JSON.Type{}

			for _, c := range list {
				result = append(result, JSON.Type{
					"Id": c.Id,
					"Status": c.Status,
				})
			}

			rend.JSON(200, helper.Success(result))
		})

		r.Post("/check", func(rend render.Render, req *http.Request){
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
		})

		r.Post("/add", func(rend render.Render, req *http.Request){
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
		})

		r.Post("/:id/update", func(rend render.Render, req *http.Request){
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
		})

		r.Post("/:id/del", func(rend render.Render, params martini.Params){
			id :=helper.Int64(params["id"])
			c := model.WebServer{Id: id}
			if err := client.Del(&c); err != nil {
				rend.JSON(200, helper.Error(err))
				return
			}
			rend.JSON(200, helper.Success())
		})

		r.Post("/:id/change/group/:gid", func(rend render.Render, params martini.Params) {
			id := helper.Int64(params["id"])
			gid := helper.Int64(params["gid"])

			c := model.WebServer{Id: id, Group: gid}
			if err := client.Update(&c, "Group"); err != nil {
				rend.JSON(200, helper.Error(err))
				return
			}

			rend.JSON(200, helper.Success())
		})
	})
}
