package host

import (
	"king/service/proc"
	"king/rpc"
	"time"
	"king/service/webSocket"
	"king/utils/JSON"
	"king/service/task"
	sh "github.com/codeskyblue/go-sh"
	"fmt"
)

var deploying bool
var shDir = "shells/deploy1/"

const (
	Start int = iota
	Message
	Error
	End
)

func broadcastAll(what int, message string){
	CallRpc("DeployStatue", &rpc.SimpleArgs{
		Id: Detail.Id,
		Message: message,
		What: what,
	})
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

		var output []byte
		var err error

		deploying = true
		this.Enable = false
		defer func(){
			deploying = false
		}()

		broadcastAll(Start, "mvn compiling..")
		output, err = sh.Command("sh", shDir+"compile.sh").SetTimeout(time.Second * 10).Output()
		if err == nil {
			fmt.Println(string(output))
			broadcastAll(Message, "killing java..")
			output, err = sh.Command("sh", shDir+"kill.sh").SetTimeout(time.Second * 30).Output()
			if err == nil {
				fmt.Println(string(output))
				broadcastAll(Message, "backup whole project..")
				output, err = sh.Command("sh", shDir+"backup.sh").SetTimeout(time.Second * 30).Output()
				if err == nil {
					fmt.Println(string(output))
					broadcastAll(Message, "execute mvn war:exploded")
					output, err = sh.Command("sh", shDir+"exploaded.sh").SetTimeout(time.Second * 30).Output()
					if err == nil {
						fmt.Println(string(output))
						broadcastAll(Message, "starting server..")
						output, err = sh.Command("sh", shDir+"startup.sh").SetTimeout(time.Second * 30).Output()
						if err == nil {
							fmt.Println(string(output))
							broadcastAll(Message, "everying ok!")
						}
					}
				}
			}
		}

		if err != nil {
			broadcastAll(Error, err.Error())
		}

		endDeploy()
	})
}
