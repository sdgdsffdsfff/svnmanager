package websocket

import (
	"king/service"
	"king/utils/JSON"
	"king/helper"
)

func init(){
	Register("heartbeat", func() JSON.Type {
		list := service.ClientService.List()
		result := []JSON.Type{}

		for _, client := range list {
			result = append(result, JSON.Type{
				"Id": client.Id,
				"Status": client.Status,
			})
		}

		return helper.Success(result)
	})
}
