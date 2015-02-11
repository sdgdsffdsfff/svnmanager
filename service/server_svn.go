package service

import (
	"king/model"
	"king/helper"
	"king/utils/JSON"
	"king/utils/db"
)

const (
	None int = iota
	Add
	Update
	Del
	Wait
)

type svnService struct {
	Version int
	isBusy bool
}

var SvnService = &svnService{0, false}

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

//保存或更新到未部署列表
func (r *svnService) SaveUpFile(list []JSON.Type) error {
	oldList, err := r.GetUnDeployFileList();
	if err != nil {
		return err
	}
	//未部署列表为空，直接把记录整体插入
	if oldList == nil {
		if _, err := db.Orm().InsertMulti(100, list); err != nil {
			return helper.NewError("SvnService.SaveUpFile InsertBatch error", err)
		}
	}else{
		helper.AsyncMap(list, func(i int) bool {
				newPath := list[i]
				found := 0

				helper.AsyncMap(oldList, func(index int) bool {
						//对相同路径下的不同动作进行更新
						oldPath := oldList[index]
				if oldPath.Path == newPath["Path"] && oldPath.Action != newPath["Acton"] {
					oldPath.Action = newPath["Action"].(int)
					if created, _, err := db.Orm().ReadOrCreate(&oldPath, "Path"); err == nil {
						if created {
							db.Orm().Update(&oldPath)
						}else {
							found++
							return true
						}
					}
				}
						return false
					})

				//TODO 更改调用方法
				if found == 0 {
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
