package master

import (
	"fmt"
	"king/enum/status"
	"king/helper"
	"king/model"
	"king/service/webSocket"
	"king/utils/JSON"
	"king/utils/db"
)

var Version int
var LastVersion model.Version
var isLock bool
var Message string
var Error bool
var ErrorLog string
var Status status.Status
var UnDeployList JSON.Type

//获取最新版本
func GetLastVersion() (model.Version, error) {
	LastVersion = model.Version{}
	if err := db.Orm().QueryTable("version").OrderBy("-id").Limit(1).One(&LastVersion); err != nil {
		return LastVersion, helper.NewError("SvnService.GetLastVersion", err)
	}
	Version = LastVersion.Version
	return LastVersion, nil
}

//是否更新
func IsChanged(v int) bool {
	if version, err := GetLastVersion(); version.Id == 0 || err == nil && v > version.Version {
		return true
	}
	return false
}

//添加更新记录
func UpdateVersion(data *model.Version) error {
	if _, err := db.Orm().Insert(data); err != nil {
		return helper.NewError("save version error", err)
	}
	LastVersion = *data
	Version = LastVersion.Version
	return nil
}

func DeployMessage(message string) error {
	LastVersion.Comment = message
	if _, err := db.Orm().Update(&LastVersion, "Comment"); err != nil {
		fmt.Println("error", err)
		return err
	}
	return nil
}

func SetBusy(yes ...bool) {
	isBusy := true
	if len(yes) > 0 {
		isBusy = yes[0]
	}
	if !isBusy {
		Status = status.Alive
	} else {
		Status = status.Busy
	}
}

func SetMessage(message ...string) {
	if len(message) > 0 {
		Message = message[0]
	} else {
		Message = ""
	}
}

func SetError(params ...interface{}) {
	err := false
	msg := ""

	if len(params) > 0 {
		err = params[0].(bool)
		if len(params) > 1 {
			msg = params[1].(string)
		}
	}

	Error = err
	ErrorLog = msg
}

func SetUnDeployFile(params ...JSON.Type) {
	if len(params) > 0 {
		UnDeployList = params[0]
	} else {
		UnDeployList = nil
	}
}

func Lock() bool {
	if IsLock() {
		return false
	}
	isLock = true
	webSocket.LockMaster()
	return true
}

func Unlock() {
	isLock = false
	webSocket.UnlockMaster()
}

func IsLock() bool {
	return isLock
}
