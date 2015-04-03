package controller

import (
	"github.com/go-martini/martini"
	"github.com/golang/glog"
)

var ctrlMaps []interface{}

// Append a value to the key's slice
func AddController(value interface{}) {
	ctrlMaps = append(ctrlMaps, value)
}

func MappingController(m *martini.ClassicMartini) {
	glog.Infoln("mapping controller")
	for _, ctn := range ctrlMaps {
		var a = ctn.(IController)
		a.SetRouter(m)
	}
}

type IController interface {
	SetRouter(m *martini.ClassicMartini)
}
