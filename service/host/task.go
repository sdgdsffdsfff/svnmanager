package host

import (
	"king/service/proc"
	"king/rpc"
	"time"
	"king/service/webSocket"
	"king/utils/JSON"
	sh "github.com/codeskyblue/go-sh"
	"os/exec"
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
		var stdout io.Reader
		var session *sh.Session

		deploying = true
		this.Enable = false
		defer func(){
			deploying = false
		}()

		ch := make(chan error)
		go func(){
			broadcastAll("mvn client start", "")
			cmd := exec.Command("sh", "shells/mvn.sh")

			if stdout, err = cmd.StdoutPipe(); err == nil {
				if err = cmd.Run(); err == nil {
					buf := new(bytes.Buffer)
					buf.ReadFrom(stdout)
					fmt.Println(buf.String())
				}
			}
			ch <- err
		}()

		if	err = <-ch; err != nil {
			broadcastAll("mvn clean:clean compile", err.Error())
			return
		}

		broadcastAll("kill tomcat", "")
		session = sh.Command("ps", "aux").Command("grep", "java").Command("wc","-l")
		if output, err = session.Output(); err == nil {
			broadcastAll(string(output), "")
		}
		if err != nil {
			broadcastAll("ps aux | grep java | wc -l", err.Error())
			return
		}

	})
}
