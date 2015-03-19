package host

import (
	"king/service/proc"
	"king/rpc"
	"time"
	"king/service/task"
	sh "github.com/codeskyblue/go-sh"
	"king/helper"
	. "king/enum/deploy"
	"log"
)

var deploying bool
var shDir = "shells/deploy1/"

func broadcastAll(what int, message string){
	CallRpc("DeployStatue", &rpc.SimpleArgs{
		Id: Detail.Id,
		Message: message,
		What: what,
	})
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
			log.Println(string(output))
			if err != nil {
				broadcastAll(Error, err.Error())
			} else {
				broadcastAll(Finish, "deploy complete")
			}
		}()

		broadcastAll(Start, "mvn compiling..")
		output, err = sh.Command("sh", shDir+"compile.sh").SetTimeout(time.Second * 10).Output()
		if err != nil {
			err = helper.NewError("mvn compile error!", err)
			return
		}
		broadcastAll(Message, "killing java..")
		output, err = sh.Command("sh", shDir+"kill.sh").SetTimeout(time.Second * 30).Output()
		if err != nil {
			err = helper.NewError("kill java error!", err)
			return
		}

		broadcastAll(Message, "backup whole project..")
		output, err = sh.Command("sh", shDir+"backup.sh").SetTimeout(time.Second * 30).Output()
		if err != nil {
			err = helper.NewError("backup error!", err)
			return
		}

		broadcastAll(Message, "execute mvn war:exploded")
		output, err = sh.Command("sh", shDir+"exploaded.sh").SetTimeout(time.Second * 30).Output()
		if err != nil {
			err = helper.NewError("mvn war:exploded error!", err)
			return
		}

		broadcastAll(Message, "starting server..")
		output, err = sh.Command("sh", shDir+"startup.sh").SetTimeout(time.Second * 30).Output()
		if err != nil {
			err = helper.NewError("server start error!", err)
			return
		}
	})
}
