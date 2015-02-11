package shell

import (
	"regexp"
	"king/utils/JSON"
	"strings"
	"king/helper"
	"strconv"
	"king/service"
)

var svnDir string = "/home/languid/svn/project/king"

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
func SvnUp(paths ...string) (int, []JSON.Type, error){

	var path string

	if len(paths) > 0 {

	}else{
		path = svnDir
	}

	out, err := Cmd("svn up " + path)
	if err != nil {
		return -1, nil, helper.NewError("svn up command error", err)
	}

	lines := strings.Split(helper.Trim(out), "\n")

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
						"Action": service.SvnService.ParseAction(action),
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
