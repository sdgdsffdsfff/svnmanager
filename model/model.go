package model

import (
	"time"
)

type Version struct {
	Id      int64 		`json:"id"`
	Version int 		`json:"version"`      //版本号
	Time    time.Time 	`json:"time"`//更新时间
	Comment string 		`json:"comment"`  //更新内容备注
	List    string 		`json:"list"` //更新内容
}

type Config struct {
	Id      int64		`json:"id"`
	Name    string 		`json:"name"` //配置名
	Content string 		`json:"content"` //配置内容
}

type WebServer struct {
	Id           int64	`json:"id"`
	Ip           string	`json:"ip"`
	Port         string `json:"port"`//web端口
	Name         string `json:"name"`//备注
	Version      int    `json:"version"`//版本号
	Status       int    `json:"status"`//状态
	Group        int64  `json:"group"`//分组
	InternalIp   string `json:"internalIp"`//内网IP
	DeployPath   string `json:"deployPath"`//部署地址
	UnDeployList string `json:"unDeployList"`//未部署列表
}

type Group struct {
	Id   int64			`json:"id"`
	Name string			`json:"name"`
	Desc string			`json:"desc"`
}
