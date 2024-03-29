package master

import (
	sh "github.com/codeskyblue/go-sh"
	"king/service/task"
	"king/service/webSocket"
	"time"
)

var compiling = false

func init() {

	task.New("master.Compile", func(this *task.Task) interface{} {

		if compiling {
			return nil
		}

		SetMessage("Compiling")

		this.Enable = false
		compiling = true

		var err error
		var output []byte

		defer func() {
			compiling = false
			Unlock()
			SetBusy(false)
			SetMessage()
			if err != nil {
				SetError(true, string(output))
			}
		}()

		this.Enable = false

		output, err = sh.Command("sh", "shells/compile.sh").SetTimeout(time.Minute * 1).Output()
		if err != nil {
			webSocket.Notify("Compile failur!", webSocket.Warning)
		}

		return nil
	})
}
