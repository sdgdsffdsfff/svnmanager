package client

import (
	"net/http"
	"king/utils/JSON"
	"king/service"
	"king/rpc"
	"king/utils"
	"king/model"
	"path/filepath"
	"net/url"
	"king/config"
	"king/helper"
	"fmt"
)

//rpc method
type RpcDeploy struct {}

func (h *RpcDeploy) Deploy(r *http.Request, args *rpc.DeployArgs, reply *rpc.RpcReply) error {
	errors := deployFileList(args.FileUrl)
	if len(errors) > 0 {
		reply.Response = errors
	}else {
		reply.Response = true
	}
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

func deployFileList(list []*model.UpFile) []JSON.Type {
	deployDir := config.Get("deployDir").(string)
	length := len(list)
	errors := []JSON.Type{}

	if length == 0 {
		return errors
	}

	helper.AsyncMap(list, func(index int) bool{
		var err error
		file := list[index]
		if file.Action == service.Add || file.Action == service.Update{
			fileUrl := config.ResServer() + file.Path
			//解析URL错误
			Url, err := url.Parse(fileUrl)
			if err == nil {
				dir, name := filepath.Split(Url.Path)
				path := deployDir + dir
				//下载错误
				fmt.Println(path)
				err = utils.Download(fileUrl, path, name)
			}
		}else if file.Action == service.Del {
			//删除文件错误
			err = utils.RemovePath(file.Path)
		}

		if err != nil {
			errors = append(errors, JSON.Type{
				"UpFile": file,
				"error": err,
			})
		}
		return false
	})


	return errors
}
