package server

import (
	"github.com/go-martini/martini"
	"github.com/martini-contrib/render"
	"king/config"
	"net/http"
	"king/helper"
	"king/utils/db"
	"king/model"
	"king/utils/JSON"
)

type ConfigCtrl struct{}

func init(){
	config.AppendValue(config.Controller, &ConfigCtrl{})
}

func (ctn *ConfigCtrl) SetRouter(m *martini.ClassicMartini) {
	m.Get("/aj/config/add", func(rend render.Render, req *http.Request){
		cfg := model.Config{
			Name: "hahaha",
			Content: JSON.Stringify(JSON.Type{
				"Name": "languid",
				"XX": "jeremy",
				"isTest": true,
				"clients": []int{1, 2, 3, 4, 5},
			}),
		}

	if _, err := db.Orm().Insert(&cfg); err != nil {
		rend.JSON(200, helper.Error(err))
		return
	}

		result := JSON.Parse(cfg)
		result["Content"] = JSON.Parse(result["Content"])

		rend.JSON(200, helper.Success(result))
	})
}
