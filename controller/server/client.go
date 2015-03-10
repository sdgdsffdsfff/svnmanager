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
	"strconv"
	"github.com/antonholmquist/jason"
)

type HostCtrl struct{}
type clientList []*client.HostClient

func init() {
	config.AppendValue(config.Controller, &HostCtrl{})
}

func (ctn *HostCtrl) SetRouter(m *martini.ClassicMartini) {

	m.Group("/aj/client", func(r martini.Router) {

		r.Get("/list", func(rend render.Render) {

			result := JSON.Type{}

			groupDict := map[int64]clientList{}
			groups, err := group.List()
			for _, g := range groups {
				groupDict[g.Id] = clientList{}
			}

			if err == nil {
				clients := client.List()
				for _, c := range clients {
					if _, found := groupDict[c.Group]; found {
						groupDict[c.Group] = append(groupDict[c.Group], c)
					}
				}
			}

			if err != nil {
				rend.JSON(200, helper.Error(err))
				return
			}

			for _, g := range groups {
				gdata := JSON.Parse(g)
				gdata["Clients"] = groupDict[g.Id]
				result[strconv.FormatInt(g.Id, 10)] = gdata;
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
				result, err := client.CallRpc(c, "RpcClient.CheckDeployPath", JSON.Type{
					"DeployPath": c.DeployPath,
				})
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
