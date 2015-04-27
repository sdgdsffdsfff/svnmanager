package host

import (
	sh "github.com/codeskyblue/go-sh"
	"github.com/golang/glog"
	"io/ioutil"
	"king/config"
	actionEnum "king/enum/action"
	"king/helper"
	"king/rpc"
	"king/utils"
	"king/utils/JSON"
	"log"
	"path/filepath"
	"strings"
	"sync"
	"time"
	"path"
)

var Detail = rpc.ActiveArgs{}

var lock sync.Mutex = sync.Mutex{}
var IsConnected = false
var reActiveTimes time.Duration = 5
var retryTimes time.Duration = 0
var backupPath = "/opt/bak"
var deployPath = "/home/languid/svn/download/test"
var webrootPath = "/usr/local/tomcat6/webapps"

func Connect() {
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

func Active(id int64) {
	if IsConnected {
		return
	}
	Detail.Id = id
	IsConnected = true
}

func Update(u *rpc.UpdateArgs) JSON.Type {
	lock.Lock()
	defer lock.Unlock()

	results := JSON.Type{}

	helper.AsyncMap(u.FileList, func(key, value interface{}) bool {
		var err error
		p := key.(string)
		filePath := "/" + p
		action := value.(int)

		//添加和更新直接下载覆盖
		if action == actionEnum.Add || action == actionEnum.Update {

			fileUrl := u.ResPath + filePath
			//下载错误
			err = utils.Download(fileUrl, path.Join(deployPath, filePath))
		} else if action == actionEnum.Del {
			//删除文件错误
			utils.RemovePath( path.Join(deployPath, filePath) )
		}

		if err != nil {
			results[p] = err.Error()
		} else {
			results[p] = true
		}

		return false
	})

	return results
}

func Revert(path string) error {
	root := filepath.Join(backupPath, path)
	err := utils.PathEnable(root)
	if err != nil {
		return err
	}

	_, err = sh.Command("rm", "-rf", deployPath).Output()
	if err != nil {
		return err
	}

	_, err = sh.Command("cp", "-r", root, deployPath).Output()
	if err != nil {
		return err
	}

	return nil
}

func RemoveBackup(path string) error {
	root := filepath.Join(backupPath, path)
	err := utils.PathEnable(root)
	if err != nil {
		return err
	}
	_, err = sh.Command("rm", "-rf", root).Output()
	if err != nil {
		return err
	}

	return nil
}

func ShowLog() (string, error) {
	output, err := sh.Command("sh", "shells/log.sh").Output()
	if err != nil {
		return "", err
	}
	return string(output), nil
}

func GetBackupList() ([]string, error) {

	var result []string
	var name string

	files, err := ioutil.ReadDir(backupPath)
	if err != nil {
		return result, err
	}

	for _, f := range files {
		if f.IsDir() {
			name = f.Name()
			if strings.Index(name, "ROOT") == 0 {
				result = append(result, name)
			}
		}
	}

	return result, nil
}

func CallRpc(method string, params interface{}) (interface{}, error) {
	result, err := rpc.Send(config.MasterRpc(), "RpcServer."+method, params)
	if err != nil {
		IsConnected = false
		log.Println("lose connect")
		return nil, err
	}
	return result, nil
}
