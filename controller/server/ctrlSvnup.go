package server

import (
	"king/helper"
	"king/utils/JSON"
	"king/utils/shell"
	"king/service"
	"king/model"
	"time"
)

func SvnUpCtrl() (model.Version, error){
	now := time.Now()
	version := model.Version{}


	num, list, err := shell.SvnUp()
	if err != nil {
		return version, err
	}

	if service.Svn.IsChanged(num) == false {
		return version, helper.NewError("no change")
	}

	version = model.Version{
		Version: num,
		Time: now,
		List: JSON.Stringify(list),
	}

	if err := service.Svn.UpdateVersion(&version); err != nil {
		return version, err
	}

	if err := service.Svn.SaveUpFile(list); err != nil {
		return version, err
	}

	msg := &service.Message{"svnup", helper.Success(version)}
	service.WebSocket.NotifyAll(msg)

	return version, nil
}
