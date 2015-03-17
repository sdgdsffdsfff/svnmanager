package host

import (
	"king/service/proc"
	"king/rpc"
	"time"
	"king/service/webSocket"
	"king/utils/JSON"
	sh "github.com/codeskyblue/go-sh"
	"fmt"
	"bytes"
)

var deploying bool

func broadcastAll(message, err string){
	CallRpc("BroadCastAll", webSocket.Message{"deploy", JSON.Type{
		"clientId": Detail.Id,
		"message": message,
		"error": err,
	}})
}

func init(){
	Task("ProcStat", func(this *TaskCallback){
		cpu := proc.CPUPercent()
		mem := proc.MEMPercent()
		CallRpc("ReportUsage", rpc.UsageArgs{Detail.Id, cpu, mem})
	}, time.Second * 1)

	Task("Deploy", func(this *TaskCallback){

		if deploying {
			this.Enable = false
			return
		}

		var err error
		var output []byte

		deploying = true
		this.Enable = false
		defer func(){
			deploying = false
		}()

		broadcastAll("mvn client start", "")
		output, err = sh.Command("sh", "shells/mvn.sh").Output()
		if err != nil {
			broadcastAll("mvn clean:clean compile", err.Error())
			return
		}
		broadcastAll(string(output), "")

		broadcastAll("kill tomcat", "")
		output, err = sh.Command("ps", "aux").Command("grep", "java").Command("wc","-l").Output()
		if err != nil {
			broadcastAll("ps aux | grep java | wc -l", err.Error())
			return
		}
		broadcastAll(string(output), "")
	})
}
