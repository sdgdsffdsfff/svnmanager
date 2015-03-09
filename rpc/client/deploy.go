package client

import (
	"net/http"
	"king/utils/JSON"
	"king/service/svn"
	"king/rpc"
	"king/utils"
	"king/model"
	"path/filepath"
	"net/url"
	"king/config"
	"king/helper"
)

//rpc method
type RpcDeploy struct {}

func (h *RpcDeploy) Deploy(r *http.Request, args *rpc.DeployArgs, reply *rpc.RpcReply) error {
	results := deployFileList(args.FileUrl, args.DeployPath)
	reply.Response = results
	return nil
}

//controller service
type deploy struct{
	List []*JSON.Type
}

func (r deploy) String() string {
	return JSON.Stringify(r.Data())
}

func (r deploy) Data() interface {} {
	return r.List;
}

var DeployService *deploy = &deploy{}

func init(){
	rpc.AddCtrl(new(RpcDeploy))
}

func deployFileList(list []*model.UpFile, deployPath string) []JSON.Type {
	results := []JSON.Type{}
	if len(list) == 0 {
		return results
	}

	helper.AsyncMap(list, func(index int) bool{
		var err error
		file := list[index]

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

	return results
}
