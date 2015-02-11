package server

import (
	"github.com/go-martini/martini"
	"github.com/martini-contrib/render"
	"king/config"
	"king/rpc"
	"fmt"
	"king/utils/JSON"
)

type TestCtrl struct{}

func init() {
	config.AppendValue(config.Controller, &TestCtrl{})
}

func (ctn *TestCtrl) SetRouter(m *martini.ClassicMartini) {
	m.Get("/test", func(rend render.Render) {
			result, err := rpc.Send(
				"http://127.0.0.1:4000/rpc",
				"RpcState.Alive",
				JSON.Type{
				"Name": "languid",
			},
			)

			fmt.Println(result, err)

			rend.JSON(200, result)
		})
}
