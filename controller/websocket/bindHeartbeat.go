package websocket

import (
	"king/service/webSocket"
	"king/service/client"
	"king/utils/JSON"
	"king/helper"
)

func init(){
	webSocket.BindWebSocketMethod("heartbeat", func() JSON.Type {
		list := client.List()
		result := []JSON.Type{}

		for _, c := range list {
			result = append(result, JSON.Type{
				"Id": c.Id,
				"Status": c.Status,
			})
		}

		return helper.Success(result)
	})
}
