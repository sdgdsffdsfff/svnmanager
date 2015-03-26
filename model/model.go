package model

import (
	"time"
)

type Version struct {
	Id            	int64
	Version 		int 		//版本号
    Time 			time.Time 	//更新时间
	Comment 		string 		//更新内容备注
	List        	string 		//更新内容
}

type Config struct {
	Id          	int64
	Name 			string 		//配置名
	Content 		string 		//配置内容
}

type UpFile struct {
	Id         		int64
	Path       		string  //更新到的文件
	Action     		int    	//U A D 操作
 	Version	   		int 	//文件版本号（非整体用于区分文件）
}

type WebServer struct {
	Id        		int64
	Ip     			string
	Port   			string 	//web端口
	Name   			string 	//备注
	Version 		int 	//版本号
	Status 			int 	//状态
	Group     		int64   //分组
	InternalIp  	string 	//内网IP
  	DeployPath		string  //部署地址
    BackupPath 		string  //备份地址
  	UndeployList	string  //未部署列表
}

type Group struct {
	Id        int64
	Name      string
	Desc	  string
}
