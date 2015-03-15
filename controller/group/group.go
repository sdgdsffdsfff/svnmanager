package group

import (
	"github.com/martini-contrib/render"
	"net/http"
	"king/helper"
	"github.com/antonholmquist/jason"
	"king/service/group"
)

func Add(rend render.Render, req *http.Request) {
	params, _ := jason.NewObjectFromReader(req.Body)
	name, _ := params.GetString("name")
	result, err := group.Add(name)
	if err != nil {
		rend.JSON(200, helper.Error(err))
		return
	}
	rend.JSON(200, helper.Success(result))
}
