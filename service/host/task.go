package host

import (
	"king/service/proc"
	"king/rpc"
	"time"
	"king/service/webSocket"
	"king/utils/JSON"
	sh "github.com/codeskyblue/go-sh"
	"king/helper"
)

var deploying bool

const (
	Message int = iota
	Log
	Error
)

func broadcastAll(types int, message string){
	CallRpc("BroadCastAll", webSocket.Message{"deploy", JSON.Type{
		"clientId": Detail.Id,
		"type": types,
		"message": message,
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

		broadcastAll(Message,"starting deploy")
		session := sh.Command("sh", "shells/auto_deploy.sh")
		session.WaitTimeout( time.Minute * 1 )
		output, err = session.Output()
		if err != nil {
			broadcastAll(Error, err.Error())
			return
		}
		broadcastAll(Message, string(output))
	})
}
