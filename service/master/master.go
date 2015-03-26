package master

import (
	"king/model"
	"king/helper"
	"king/utils/JSON"
	"king/utils/db"
	"github.com/astaxie/beego/orm"
	"king/enum/status"
	"fmt"
	"king/service/webSocket"
)

var Version int
var LastVersion model.Version
var isLock bool
var Message string
var Error bool
var ErrorLog string
var Status status.Status

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
func UpdateVersion(data *model.Version) error{
	if _, err := db.Orm().Insert(data); err != nil {
		return helper.NewError("save version error", err)
	}
	LastVersion = *data
	Version = LastVersion.Version
	return nil
}

func DeployMessage( message string ) error {
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

//获取未部署列表
//参数为空或者是[0]代表获取所有文件
func GetUnDeployFileList( ids ...[]int64 ) ([]*model.UpFile, error) {
	var list []*model.UpFile
	var qs orm.QuerySeter

	//如果filter不连在QueryTable后面写会无效
	if len(ids) == 1 && len(ids[0]) > 0 && ids[0][0] != 0 {
		qs = db.Orm().QueryTable("up_file").Filter("id__in", ids[0])
	} else {
		qs = db.Orm().QueryTable("up_file")
	}
	if _, err := qs.All(&list); err != nil {
		return list, err
	}
	return list, nil
}

func ClearDeployFile() error {
	return db.Truncate("up_file");
}

func Lock() bool {
	if IsLock() {
		return false
	}
	isLock = true;
	webSocket.LockMaster()
	return true
}

func Unlock(){
	isLock = false
	webSocket.UnlockMaster()
}

func IsLock() bool{
	return isLock
}

//保存或更新到未部署列表
func SaveUpFile(list []JSON.Type) error {

	oldList, err := GetUnDeployFileList();
	if err != nil {
		return err
	}
	//未部署列表为空，直接把记录整体插入
	if oldList == nil {
		// http://beego.me/docs/mvc/model/object.md#insertmulti
		// 并列插入100条，顺序插入改为1
		bulk := 100
		if _, err := db.Orm().InsertMulti(bulk, list); err != nil {
			return helper.NewError("SvnService.SaveUpFile InsertBatch error", err)
		}
	}else{
		//并发更新文件列表
		helper.AsyncMap(list, func(key, value interface{}) bool {
			newPath := value.(JSON.Type)
			action := newPath["Action"].(int)

			found := false
			helper.AsyncMap(oldList, func(ikey, ivalue interface{}) bool {
				//对相同路径下的不同动作进行更新
				oldPath := ivalue.(*model.UpFile)
				if oldPath.Path == newPath["Path"] {
					oldPath.Version = Version
					if oldPath.Action != action {
						oldPath.Action = action
					}
					db.Orm().Update(oldPath)
					found = true
					//跳出循环
					return true
				}
				return false
			})

			//未找到相同路径的记录则插入
			if found == false {
				up_file := &model.UpFile{}
				newPath["Version"] = Version
				JSON.ParseToStruct(newPath, up_file)
				db.Orm().Insert(up_file)
			}
			return false
		})
	}
	return nil
}
