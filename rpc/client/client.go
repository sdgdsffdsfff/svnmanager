package client

import (
	"king/rpc"
	"king/utils/JSON"
	"net/http"
	"king/service/proc"
	"king/helper"
	"king/service/svn"
	"king/config"
	"net/url"
	"path/filepath"
	"king/utils"
	"king/model"
)


//rpc method
type RpcClient struct{}

func (h *RpcClient) CheckDeployPath(r *http.Request, args *JSON.Type, reply *rpc.RpcReply) error {
	reply.Response = true
	return nil
}

func (h *RpcClient) ProcStat(r *http.Request, args *JSON.Type, reply *rpc.RpcReply) error {
	cpu := proc.CPUPercent()
	mem := proc.MEMPercent()
	reply.Response = JSON.Type{
		"CPUPercent" : cpu,
		"MEMPercent" : mem,
	}
	return nil
}

func (h *RpcClient) Deploy(r *http.Request, args *rpc.DeployArgs, reply *rpc.RpcReply) error {
	results := []JSON.Type{}
	fileList := args.FileUrl
	deployPath := args.DeployPath

	if len(args.FileUrl) > 0 {
		helper.AsyncMap(fileList, func(key, value interface{}) bool{
			var err error
			file := value.(*model.UpFile)

			//添加和更新直接下载覆盖
			if file.Action == svn.Add || file.Action == svn.Update{
				fileUrl := config.ResServer() + file.Path

				//解析URL错误
				Url, err := url.Parse(fileUrl)
				if err == nil {
					dir, name := filepath.Split(Url.Path)
					path := deployPath + dir

					//下载错误
					err = utils.Download(fileUrl, path, name)
				}
			}else if file.Action == svn.Del {

				//删除文件错误
				err = utils.RemovePath(file.Path, deployPath)
			}

			results = append(results, JSON.Type{
				"UpFile": file,
				"error": err,
			})
			return false
		})
	}

	reply.Response = results
	return nil
}

func init(){
	rpc.AddCtrl(new(RpcClient))
}
