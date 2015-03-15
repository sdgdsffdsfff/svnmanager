package server

import (
	sysConfig "king/config"
	"github.com/go-martini/martini"
	"king/controller/client"
	"king/controller/group"
	"king/controller/svn"
	"king/controller/config"
	ctrlWebSocket "king/controller/webSocket"
	"github.com/martini-contrib/render"
	sockets "github.com/beatrichartz/martini-sockets"
	"king/service/webSocket"
)

type ServerRoutes struct {}

func init() {
	sysConfig.AppendValue(sysConfig.Controller, &ServerRoutes{})
}

func (ctn *ServerRoutes) SetRouter(m *martini.ClassicMartini){

	m.Get("/", func(rend render.Render) {
		rend.HTML(200, "index", nil)
	})

	m.Get("/socket", sockets.JSON(webSocket.Message{}), ctrlWebSocket.Socket)

	m.Group("/aj/client", func(r martini.Router) {
		r.Get("/list",  client.List)
		r.Post("/check", client.Check)
		r.Post("/add", client.Add)
		r.Post("/:id/update", client.Update)
		r.Post("/:id/del", client.Del)
		r.Post("/:id/change/group/:gid", client.Move)
		r.Post("/:id/update")
	})

	m.Group("/aj/group", func(r martini.Router){
		r.Post("/add", group.Add)
	})

	m.Group("/aj/svn", func(r martini.Router){
		r.Get("/lastVersion", svn.GetLastVersion)
		r.Get("/undeploy/files", svn.GetUndeployFiles)
		r.Post("/update", svn.Update)
		r.Post("/revert/:version", svn.Revert)
	})

	m.Group("/aj/config", func(r martini.Router){
		r.Get("/aj/config/add", config.Add)
	})


}
