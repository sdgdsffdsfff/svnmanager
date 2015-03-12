package websocket

import (
	"king/service/webSocket"
	"king/service/client"
)

func init(){
	webSocket.OnOut(func(length int) {
		if length == 0 {
			client.SetHeartEnable(false)
			client.SetProcMonitorEnable(false)
		}
	})
}
