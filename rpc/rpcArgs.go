package rpc

import "king/model"

type DeployArgs struct {
	FileUrl []*model.UpFile
	DeployPath string
}

type CheckDeployPathArgs struct {
	Path string
}
