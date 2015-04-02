package config

import (
	"github.com/martini-contrib/render"
	"king/helper"
	"king/model"
	"king/utils/JSON"
	"king/utils/db"
	"net/http"
)

func Add(rend render.Render, req *http.Request) {
	cfg := model.Config{
		Name: "hahaha",
		Content: JSON.Stringify(JSON.Type{
			"Name":    "languid",
			"XX":      "jeremy",
			"isTest":  true,
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
}
