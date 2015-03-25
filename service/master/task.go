package master

import (
	"king/service/task"
	sh "github.com/codeskyblue/go-sh"
	"time"
	"king/service/webSocket"
)

var compiling = false

func init() {


	task.New("compile", func(this *task.Task) interface{} {

		if compiling {
			return nil
		}

		Lock()
		SetBusy()
		SetMessage("Compiling")

		this.Enable = false
		compiling = true

		var err error
		var output []byte

		defer func(){
			compiling = false
			Unlock()
			SetBusy(false)
			if err != nil {
				SetError(true, string(output))
			}
		}()

		this.Enable = false

		output, err = sh.Command("sh", "shells/compile.sh").SetTimeout( time.Minute * 1 ).Output()
		if err != nil {
			webSocket.Notify("Compile failur!", webSocket.Warning)
		}

		return nil
	})
}
