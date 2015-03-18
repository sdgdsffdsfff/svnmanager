package host

import (
	"king/service/proc"
	"king/rpc"
	"time"
	"king/service/webSocket"
	"king/utils/JSON"
	"king/service/task"
	sh "github.com/codeskyblue/go-sh"
)

var deploying bool

const (
	Start int = iota
	Log
	Error
	End
)

func broadcastAll(types int, message string){
	CallRpc("BroadCastAll", webSocket.Message{"deploy", JSON.Type{
		"clientId": Detail.Id,
		"type": types,
		"message": message,
	}})
}

func endDeploy(){
	CallRpc("EndDeploy", webSocket.Message{"deploy", JSON.Type{
		"clientId": Detail.Id,
		"type": End,
	}})
}

func init(){
	task.New("ProcStat", func(this *task.Task){

		if IsConnected == false {
			return
		}
		cpu := proc.CPUPercent()
		mem := proc.MEMPercent()

		CallRpc("ReportUsage", rpc.UsageArgs{Detail.Id, cpu, mem})

	}, time.Second * 1)

	task.New("Deploy", func(this *task.Task){

		if IsConnected == false {
			return
		}

		if deploying {
			this.Enable = false
			return
		}

		deploying = true
		this.Enable = false
		defer func(){
			deploying = false
		}()

		broadcastAll(Start, "starting deploy")
		_, err := sh.Command("sh", "shells/auto_deploy.sh").SetTimeout(time.Second * 10).Output()
		if err != nil {
			broadcastAll(Error, err.Error())
		}

		endDeploy()
	})
}
