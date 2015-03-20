package host

import (
	"king/service/proc"
	"king/rpc"
	"time"
	"king/service/task"
	sh "github.com/codeskyblue/go-sh"
	"king/helper"
	. "king/enum/deploy"
	"strings"
	"log"
)

var deploying bool
var shDir = "shells/deploy/"

func broadcastAll(what int, message string){
	CallRpc("DeployStatue", &rpc.SimpleArgs{
		Id: Detail.Id,
		Message: message,
		What: what,
	})
}

func getOutput(output []byte, err error) error {
	if err != nil {
		return err
	}
	lines := strings.Split(helper.Trim(string(output)), "\n")
	lastLine := lines[len(lines)-1]
	if lastLine != "complete" {
		return helper.NewError(lastLine)
	}
	return nil
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

		var err error

		deploying = true
		this.Enable = false
		defer func(){
			deploying = false
			if err != nil {
				log.Println(err)
				broadcastAll(Error, err.Error())
			} else {
				broadcastAll(Message, "deploy complete")
				time.Sleep( time.Second * 2 )
				broadcastAll(Finish, "")
			}
		}()

		broadcastAll(Start, "mvn compiling..")
		err = getOutput(sh.Command("sh", shDir+"compile.sh").SetTimeout(time.Second * 60).Output())
		if err != nil {
			err = helper.NewError("mvn compile error!", err)
			return
		}

		broadcastAll(Message, "killing java..")
		if err = getOutput(sh.Command("sh", shDir+"kill.sh").SetTimeout(time.Second * 30).Output()); err != nil {
			err = helper.NewError("kill java error!", err)
			return
		}

		broadcastAll(Message, "backup whole project..")
		if err = getOutput(sh.Command("sh", shDir+"backup.sh").SetTimeout(time.Second * 30).Output()); err != nil {
			err = helper.NewError("backup error!", err)
			return
		}

		broadcastAll(Message, "execute mvn war:exploded")
		if err = getOutput(sh.Command("sh", shDir+"exploded.sh").SetTimeout(time.Second * 30).Output()); err != nil {
			err = helper.NewError("mvn war:exploded error!", err)
			return
		}

		broadcastAll(Message, "starting server..")
		if err = getOutput(sh.Command("sh", shDir+"startup.sh").SetTimeout(time.Second * 30).Output()); err != nil {
			err = helper.NewError("server start error!", err)
			return
		}
	})
}
