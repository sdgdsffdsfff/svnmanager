package service

import (
	"king/model"
	"king/helper"
	"king/utils/JSON"
	"king/utils/db"
	"sync"
)

const (
	None int = iota
	Add
	Update
	Del
	Wait
)

type svnService struct {
	sync.Mutex
	Version int
	isBusy bool
	isLock bool
}

var Svn = &svnService{sync.Mutex{}, 0, false, false}

//获取最新版本
func (r *svnService) GetLastVersion() (model.Version, error) {
	var data model.Version
	if err := db.Orm().QueryTable("version").OrderBy("-id").Limit(1).One(&data); err != nil {
		return data, helper.NewError("SvnService.GetLastVersion", err)
	}
	r.Version = data.Version
	return data, nil
}

//是否更新
func (r *svnService) IsChanged(v int) bool {
	if version, err := r.GetLastVersion(); version.Id == 0 || err == nil && v > version.Version {
		return true
	}
	return false
}

//添加更新记录
func (r *svnService) UpdateVersion(data *model.Version) error{
	if _, err := db.Orm().Insert(data); err != nil {
		return helper.NewError("save version error", err)
	}
	return nil
}

//获取未部署列表
func (r *svnService) GetUnDeployFileList() ([]*model.UpFile, error) {
	var list []*model.UpFile
	_, err := db.Orm().QueryTable("up_file").All(&list)
	if err != nil {
		return list, err
	}
	return list, nil
}

func (r *svnService) ClearDeployFile() error {
	return db.Truncate("up_file");
}

func (r *svnService) GetLock() bool {
	if r.IsLock() {
		return false
	}
	r.Lock()
	r.isLock = true;
	return true
}

func (r *svnService) Release(){
	r.Unlock()
	r.isLock = false
}

func (r* svnService) IsLock() bool{
	return r.isLock
}

//保存或更新到未部署列表
func (r *svnService) SaveUpFile(list []JSON.Type) error {

	oldList, err := r.GetUnDeployFileList();
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
		helper.AsyncMap(list, func(i int) bool {
			newPath := list[i]
			action := newPath["Action"].(int)

			found := false
			helper.AsyncMap(oldList, func(index int) bool {
				//对相同路径下的不同动作进行更新
				oldPath := oldList[index]
				if oldPath.Path == newPath["Path"] && oldPath.Action != action {
					oldPath.Action = action
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
				JSON.ParseToStruct(newPath, up_file)
				db.Orm().Insert(up_file)
			}
			return false
		})
	}
	return nil
}

func (r *svnService) ParseAction(t string) int {
	switch t {
	case "A":
		return Add
	case "U":
		return Update
	case "D":
		return Del
	case "W":
		return Wait
	default:
		return None
	}
}
