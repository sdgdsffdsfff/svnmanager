package server

import (
	"king/helper"
	"king/utils/JSON"
	"king/utils/shell"
	"king/service"
	"king/model"
	"time"
)

func SvnUpCtrl() (JSON.Type, error){
	now := time.Now()

	num, list, err := shell.SvnUp()
	if err != nil {
		return nil, err
	}

	if service.SvnService.IsChanged(num) == false {
		return nil, helper.NewError("no change")
	}

	version := model.Version{
		Version: num,
		Time: now,
		List: JSON.Stringify(list),
	}

	if err := service.SvnService.UpdateVersion(&version); err != nil {
		return nil, err
	}

	if err := service.SvnService.SaveUpFile(list); err != nil {
		return nil, err
	}

	return JSON.Type{
		"Version": version,
		"List": list,
	}, nil
}
