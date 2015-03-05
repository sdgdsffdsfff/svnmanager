package server

import (
	"github.com/go-martini/martini"
	"github.com/martini-contrib/render"
	"king/config"
	"king/helper"
	"king/service"
	"king/utils/JSON"
	"king/model"
	"net/http"
	"strconv"
)

type HostCtrl struct{}
type clientList []*service.HostClient

func init() {
	config.AppendValue(config.Controller, &HostCtrl{})
}

func (ctn *HostCtrl) SetRouter(m *martini.ClassicMartini) {

	m.Group("/aj/client", func(r martini.Router) {

		r.Get("/list", func(rend render.Render) {

			result := JSON.Type{}

			groupDict := map[int64]clientList{}
			groups, err := service.Group.List()
			for _, group := range groups {
				groupDict[group.Id] = clientList{}
			}

			if err == nil {
				clients := service.Client.List()
				for _, client := range clients {
					if _, found := groupDict[client.Group]; found {
						groupDict[client.Group] = append(groupDict[client.Group], client)
					}
				}
			}

			if err != nil {
				rend.JSON(200, helper.Error(err))
				return
			}

			for _, group := range groups {
				g := JSON.Parse(group)
				g["Clients"] = groupDict[group.Id]
				result[strconv.FormatInt(group.Id, 10)] = g;
			}

			rend.JSON(200, helper.Success(result))
		})

		r.Get("/heartbeat", func(rend render.Render){
			list := service.Client.List()
			result := []JSON.Type{}

			for _, client := range list {
				result = append(result, JSON.Type{
					"Id": client.Id,
					"Status": client.Status,
				})
			}

			rend.JSON(200, helper.Success(result))
		})

		r.Post("/add", func(rend render.Render, req *http.Request){
			client := &model.WebServer{}
			body := JSON.FormRequest(req.Body)

			if err := JSON.ParseToStruct(JSON.Stringify(body), client); err != nil {
				rend.JSON(200, helper.Error(helper.ParamsError, err))
				return
			}

			if errType, err := service.Client.Add(client); err != nil {
				rend.JSON(200, helper.Error(errType, err))
				return
			}

			rend.JSON(200, helper.Success(client))
		})

		r.Post("/:id/update", func(rend render.Render, req *http.Request){
			body := JSON.FormRequest(req.Body)
			client := &model.WebServer{}
			if err := JSON.ParseToStruct(JSON.Stringify(body), client); err != nil {
				rend.JSON(200, helper.Error(helper.ParamsError))
				return
			}
			keys := JSON.GetKeys(body)

			if err := service.Client.Update(client, keys...); err != nil {
				rend.JSON(200, helper.Error(err))
				return
			}
			rend.JSON(200, helper.Success(service.Client.FindFromCache(client.Id)))
		})

		r.Post("/:id/change/group/:gid", func(rend render.Render, params martini.Params) {
			id := helper.Int64(params["id"])
			gid := helper.Int64(params["gid"])

			client := model.WebServer{Id: id, Group: gid}
			if err := service.Client.Update(&client, "Group"); err != nil {
				rend.JSON(200, helper.Error(err))
				return
			}

			rend.JSON(200, helper.Success())
		})
	})
}
