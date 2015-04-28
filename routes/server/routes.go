package server

import (
	sockets "github.com/beatrichartz/martini-sockets"
	"github.com/go-martini/martini"
	"github.com/martini-contrib/render"
	"king/controller"
	"king/controller/client"
	"king/controller/config"
	"king/controller/group"
	"king/controller/master"
	ctrlWebSocket "king/controller/webSocket"
	"king/service/webSocket"
)

type Routes struct{}

func init() {
	controller.AddController(&Routes{})
}

func (ctn *Routes) SetRouter(m *martini.ClassicMartini) {

	m.Get("/", func(rend render.Render) {
		rend.HTML(200, "index", nil)
	})

	m.Get("/socket", sockets.JSON(webSocket.Message{}), ctrlWebSocket.Socket)

	m.Group("/aj/client", func(r martini.Router) {
		r.Get("/list", client.List)
		r.Get("/:id/backuplist", client.GetBackupList)
		r.Get("/:id/log", client.ShowLog)
		r.Get("/:id/unDeploy", client.GetUnDeployFiles)
		r.Post("/check", client.Check)
		r.Post("/add", client.Add)
		r.Post("/refresh", client.Refresh)
		r.Post("/:id/edit", client.Edit)
		r.Post("/:id/del", client.Del)
		r.Post("/:id/change/group/:gid", client.Move)
		r.Post("/:id/update", client.Update)
		r.Post("/:id/deploy", client.Deploy)
		r.Post("/:id/revert", client.Revert)
		r.Post("/:id/removebackup", client.RemoveBackup)
	})

	m.Group("/aj/group", func(r martini.Router) {
		r.Post("/add", group.Add)
	})

	m.Group("/aj", func(r martini.Router) {
		r.Get("/error", master.ShowError)
		r.Get("/lastVersion", master.GetLastVersion)
		r.Post("/update", master.Update)
		r.Post("/compile", master.Compile)
	})

	m.Group("/aj/config", func(r martini.Router) {
		r.Get("/aj/config/add", config.Add)
	})

}
