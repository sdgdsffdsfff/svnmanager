package host

import (
	"king/service/proc"
	"king/rpc"
	"time"
	"king/service/webSocket"
	"king/utils/JSON"
)

func init(){
	Task("ProcStat", func(this *TaskCallback){
		cpu := proc.CPUPercent()
		mem := proc.MEMPercent()
		CallRpc("ReportUsage", rpc.UsageArgs{Detail.Id, cpu, mem})
	}, time.Second * 1)

	deploying := false
	Task("Deploy", func(this *TaskCallback){
		if( deploying ){
			return
		}

		deploying = true
		this.Enable = false

		CallRpc("BroadCastAll", webSocket.Message{"deploy", JSON.Type{
			"clientId": Detail.Id,
			"message": "Step 1",
		}})

		time.Sleep( time.Second * 3 )

		CallRpc("BroadCastAll", webSocket.Message{"deploy", JSON.Type{
			"clientId": Detail.Id,
			"message": "Step 2",
		}})

		time.Sleep( time.Second * 3 )

		CallRpc("BroadCastAll", webSocket.Message{"deploy", JSON.Type{
			"clientId": Detail.Id,
			"message": "Step 3",
		}})

		time.Sleep( time.Second * 3 )

		CallRpc("BroadCastAll", webSocket.Message{"deploy", JSON.Type{
			"clientId": Detail.Id,
			"message": "Step 4",
		}})

		deploying = false
	})
}
