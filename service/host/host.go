package host

import (
	"time"
	"king/rpc"
	"king/config"
	"sync"
	"king/utils/JSON"
	"king/helper"
	"king/service/svn"
	"king/model"
	"net/url"
	"path/filepath"
	"king/utils"
	"github.com/golang/glog"
	sh "github.com/codeskyblue/go-sh"
	"log"
)

var Detail = rpc.ActiveArgs{}

var lock sync.Mutex = sync.Mutex{}
var IsConnected = false
var reActiveTimes time.Duration = 5
var retryTimes time.Duration = 0

func Connect(){
	result, err := CallRpc("Active", Detail)
	if err != nil {
		retryTimes++
		if retryTimes > reActiveTimes {
			time.Sleep(time.Minute * 1)
		} else {
			time.Sleep(time.Second * 2 * retryTimes)
		}
		glog.Errorln(err)
		Connect()
	} else {
		Detail.Id = int64(result.(float64))
		log.Println("already connect to server")
	}
}

func Active(id int64){
	if IsConnected {
		return
	}
	Detail.Id = id
	IsConnected = true
}

func Update(fileList []*model.UpFile, deployPath string) []JSON.Type {
	lock.Lock()
	defer lock.Unlock()

	results := []JSON.Type{}

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

	return results
}

func ShowLog() (string, error){
	output, err := sh.Command("sh", "shells/log.sh").Output()
	if err != nil {
		return "", err
	}
	return string(output), nil
}

func CallRpc(method string, params interface{})(interface{}, error) {
	result, err := rpc.Send(config.MasterRpc(), "RpcServer."+method, params)
	if err != nil {
		IsConnected = false
		log.Println("lose connect")
		return nil, err
	}
	return result, nil
}
