package rpc

import "king/model"


/***** Server Rpc Arguments ******/

type UpdateArgs struct {
	FileUrl []*model.UpFile
	DeployPath string
}

type UsageArgs struct {
	Id int64
	CPUPercent float64
	MEMPercent float64
}

type MessageArgs struct {
	Method string `json:"method"`
	Params interface {} `json:"data"`
}

/***** Client Rpc Arguments ******/

type CheckDeployPathArgs struct {
	Path string
}

type ActiveArgs struct {
	Id int64
}

type DeployArgs struct {
	Id int64
}
