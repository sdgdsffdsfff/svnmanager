package host

import (
	"king/bootstrap"
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
	"fmt"
	"github.com/golang/glog"
	"log"
)

type TaskCallback struct {
	IsRunning bool
	Enable bool
	Watch func()
	Duration time.Duration
}

func (r *TaskCallback) Stop(){
	r.Enable = false
}

func (r *TaskCallback) Start() {
	r.Enable = true
}

func (r *TaskCallback) Quit(){
	r.IsRunning = false
}

var Detail = model.WebServer{}

var defaultDuration = time.Second * 5
var taskList = map[string]*TaskCallback{}
var lock sync.Mutex = sync.Mutex{}
var IsConnected = false
var taskLoopEnabled = false

func Task(name string, callback func(*TaskCallback), duration ...time.Duration){
	d := defaultDuration

	if len(duration) > 0 {
		d = duration[0]
	}

	task := &TaskCallback{
		IsRunning: false,
		Enable: false,
		Duration: d,
	}

	task.Watch = func(){
		if task.IsRunning {
			return
		}
		task.IsRunning = true
		for {
			if task.IsRunning == false {
				break
			}
			if taskLoopEnabled && task.Enable && IsConnected {
				callback(task)
			}
			time.Sleep(task.Duration)
		}
	}

	taskList[name] = task
}

func Trigger(name string){

	StartTask()

	if method, found:= taskList[name]; found {
		method.Enable = true
		method.IsRunning = true
	}
}

var reActiveTimes time.Duration = 5
var retryTimes time.Duration = 0
func Active(){
	result, err := rpc.Send(config.MasterRpc(), "RpcServer.Active", Detail)
	if err != nil {
		retryTimes++
		if retryTimes > reActiveTimes {
			time.Sleep(time.Minute * 1)
		} else {
			time.Sleep(time.Second * 2 * retryTimes)
		}
		glog.Errorln(err)
		Active()
	} else {
		Detail.Id = int64(result.(float64))
		IsConnected = true
		log.Println("already connect to server")
	}
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

func CallRpc(method string, params interface{})(interface{}, error) {
	result, err := rpc.Send(config.MasterRpc(), "RpcServer."+method, params)
	if err != nil {
		IsConnected = false
		return nil, err
	}
	return result, nil
}

func StartTask(){
	taskLoopEnabled = true
}
func StopTask(){
	taskLoopEnabled = false
}

func init(){
	bootstrap.Register(func(){
		StartTask()
		helper.AsyncMap(taskList, func(key, value interface{}) bool {
			value.(*TaskCallback).Watch()
			return false
		})
		fmt.Println("can not reach")
	})
}
