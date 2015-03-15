package svn

import (
	"king/helper"
	"king/utils/JSON"
	"king/utils/shell"
	"king/service/svn"
	"king/service/webSocket"
	"king/model"
	"time"
)

func update() (model.Version, error){
	now := time.Now()
	version := model.Version{}


	num, list, err := shell.SvnUp()
	if err != nil {
		return version, err
	}

	if svn.IsChanged(num) == false {
		return version, helper.NewError("no change")
	}

	version = model.Version{
		Version: num,
		Time: now,
		List: JSON.Stringify(list),
	}

	if err := svn.UpdateVersion(&version); err != nil {
		return version, err
	}

	if err := svn.SaveUpFile(list); err != nil {
		return version, err
	}

	webSocket.BroadCastAll(&webSocket.Message{
	"svnup",
	helper.Success(version),
})

	return version, nil
}
