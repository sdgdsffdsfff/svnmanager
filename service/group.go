package service

import (
	"king/model"
	"king/utils/db"
	"king/helper"
)

type groupService struct{}

var Group = &groupService{}

func (r *groupService) Add(name string) (*model.Group, error) {
	var group model.Group;
	group.Name = name
	_, err := db.Orm().Insert(&group)
	return &group, err
}

func (r *groupService) List() ([]*model.Group, error) {
	var groups []*model.Group
	if _, err := db.Orm().QueryTable("group").All(&groups); err != nil {
		return groups, helper.NewError("GroupService.List", err)
	}
	groups = append(groups, &model.Group{ 0, "Ungrouped", "" })
	return groups, nil
}

