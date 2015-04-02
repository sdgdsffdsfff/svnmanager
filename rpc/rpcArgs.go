package rpc

import "king/model"

type RpcInterface interface {
	SetId(int64)
	GetId() int64
}

type SimpleArgs struct {
	Id      int64
	Message string
	IsError bool
	Data    interface{}
	What    int
}

/***** Server Rpc Arguments ******/
type UpdateArgs struct {
	Id         int64
	FileUrl    []*model.UpFile
	DeployPath string
}

type UsageArgs struct {
	Id         int64
	CPUPercent float64
	MEMPercent float64
}

type MessageArgs struct {
	Id     int64
	Method string
	Params interface{}
}

/***** Client Rpc Arguments ******/

type CheckDeployPathArgs struct {
	Id   int64
	Path string
}

type ActiveArgs struct {
	Ip         string
	InternalIp string
	Port       string
	Id         int64
}

/***** Interface method ******/

func (r SimpleArgs) GetId() int64 {
	return r.Id
}
func (r SimpleArgs) SetId(id int64) {
	r.Id = id
}

func (r UpdateArgs) GetId() int64 {
	return r.Id
}
func (r UpdateArgs) SetId(id int64) {
	r.Id = id
}

func (r UsageArgs) GetId() int64 {
	return r.Id
}
func (r UsageArgs) SetId(id int64) {
	r.Id = id
}

func (r MessageArgs) GetId() int64 {
	return r.Id
}
func (r MessageArgs) SetId(id int64) {
	r.Id = id
}

func (r CheckDeployPathArgs) GetId() int64 {
	return r.Id
}
func (r CheckDeployPathArgs) SetId(id int64) {
	r.Id = id
}
