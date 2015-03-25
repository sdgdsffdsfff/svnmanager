package master

import (
	"king/helper"
	"king/utils/JSON"
	"king/service/master"
	"king/service/webSocket"
	sh "github.com/codeskyblue/go-sh"
	"king/model"
	"time"
	"strings"
	"regexp"
	"strconv"
)

var svnDir string = "/home/languid/svn/wings"

func update() (model.Version, error){
	now := time.Now()
	version := model.Version{}

	num, list, err := svnup()
	if err != nil {
		return version, err
	}

	if master.IsChanged(num) == false {
		return version, helper.NewError("no change")
	}

	version = model.Version{
		Version: num,
		Time: now,
		List: JSON.Stringify(list),
	}

	if err := master.UpdateVersion(&version); err != nil {
		return version, err
	}

	if err := master.SaveUpFile(list); err != nil {
		return version, err
	}

	webSocket.BroadCastAll(&webSocket.Message{
		"svnup",
		helper.Success(version),
	})

	return version, nil
}

/**
	return JSON.Type{
		"reversion": "10344",
		"list": []JSON.Type{
			JSON.Type{
				"path": "/_res/js/"
				"type": "U"
			}
		}
	}
 */

func svnup(paths ...string) (int, []JSON.Type, error){

	var path string

	if len(paths) > 0 {

	}else{
		path = svnDir
	}

	output, err := sh.Command("svn", "up", path).Output()
	if err != nil {
		return -1, nil, helper.NewError("svn up command error", err)
	}

	lines := strings.Split(helper.Trim(string(output)), "\n")

	list := []JSON.Type{}
	regLine := regexp.MustCompile(`^([U|D|A])\s+(.*)`)
	version := getVersion( lines[len(lines)-1] )

	for _, line := range lines {
		if matches := regLine.FindAllStringSubmatch(line, -1); matches != nil {
			for _, match := range matches {
				action := match[1]
				path := match[2]
				path = path[len(svnDir):]
				list = append(list, JSON.Type{
						"Action": master.ParseAction(action),
						"Path": path,
					})
			}
		}
	}

	return version, list, nil
}

func getVersion(str string) int{
	vIndex := strings.LastIndex(str, " ")+1
	n, _ := strconv.Atoi(str[vIndex:len(str)-1])
	return n
}
