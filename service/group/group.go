package group

import (
	"king/model"
	"king/utils/db"
	"king/bootstrap"
)

type GroupMap map[int64]*model.Group

var groupMap GroupMap

func Add(name string) (*model.Group, error) {
	var group model.Group;
	group.Name = name
	_, err := db.Orm().Insert(&group)
	return &group, err
}

func Fetch() (GroupMap, error) {
	var list []*model.Group
	groupMap = GroupMap{}
	_, err := db.Orm().QueryTable("group").All(&list)
	if err == nil {
		for _, group := range list {
			groupMap[group.Id] = group
		}
	}
	groupMap[0] = &model.Group{0, "Ungrouped", "" }
	return groupMap, err
}

func List() GroupMap {
	return groupMap
}

func init(){
	bootstrap.Register(func(){
		if db.IsConnected() {
			Fetch()
		}
	})
}
