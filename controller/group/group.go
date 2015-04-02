package group

import (
	"github.com/antonholmquist/jason"
	"github.com/martini-contrib/render"
	"king/helper"
	"king/service/group"
	"net/http"
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
