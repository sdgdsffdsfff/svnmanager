package host

import (
	sh "github.com/codeskyblue/go-sh"
	. "king/enum/deploy"
	"king/helper"
	"king/rpc"
	"king/service/proc"
	"king/service/task"
	"log"
	"time"
)

var deploying bool
var shDir = "shells/deploy/"

func broadcastAll(what int, message string) {
	CallRpc("DeployStatue", &rpc.SimpleArgs{
		Id:      Detail.Id,
		Message: message,
		What:    what,
	})
}

func init() {
	task.New("host.ProcStat", func(this *task.Task) interface{} {

		if IsConnected == false {
			return nil
		}
		cpu := proc.CPUPercent()
		mem := proc.MEMPercent()
		CallRpc("ReportUsage", rpc.UsageArgs{Detail.Id, cpu, mem})

		return nil
	}, time.Second*1)

	task.New("host.Deploy", func(this *task.Task) interface{} {

		if IsConnected == false {
			return nil
		}

		if deploying {
			this.Enable = false
			return nil
		}

		var err error

		deploying = true
		this.Enable = false
		defer func() {
			deploying = false
			if err != nil {
				log.Println(err)
				broadcastAll(Error, err.Error())
			} else {
				broadcastAll(Message, "deploy complete")
				time.Sleep(time.Second * 2)
				broadcastAll(Finish, "")
			}
		}()

		broadcastAll(Start, "mvn compiling..")
		err = helper.GetCMDOutputWithComplete(sh.Command("sh", shDir+"compile.sh").SetTimeout(time.Second * 60).Output())
		if err != nil {
			err = helper.NewError("mvn compile error!", err)
			return nil
		}

		broadcastAll(Message, "killing java..")
		if err = helper.GetCMDOutputWithComplete(sh.Command("sh", shDir+"kill.sh").SetTimeout(time.Second * 30).Output()); err != nil {
			err = helper.NewError("kill java error!", err)
			return nil
		}

		broadcastAll(Message, "backup whole project..")
		if err = helper.GetCMDOutputWithComplete(sh.Command("sh", shDir+"backup.sh").SetTimeout(time.Second * 30).Output()); err != nil {
			err = helper.NewError("backup error!", err)
			return nil
		}

		broadcastAll(Message, "execute mvn war:exploded")
		if err = helper.GetCMDOutputWithComplete(sh.Command("sh", shDir+"exploded.sh").SetTimeout(time.Second * 30).Output()); err != nil {
			err = helper.NewError("mvn war:exploded error!", err)
			return nil
		}

		broadcastAll(Message, "starting server..")
		if err = helper.GetCMDOutputWithComplete(sh.Command("sh", shDir+"startup.sh").SetTimeout(time.Second * 30).Output()); err != nil {
			err = helper.NewError("server start error!", err)
			return nil
		}

		return nil
	})
}
