package websocket

import(
	"github.com/go-martini/martini"
	sockets "github.com/beatrichartz/martini-sockets"
	"king/service"
	"king/config"
)

type WebSocketController struct{}

func init() {
	config.AppendValue(config.Controller, &WebSocketController{})
}

func (ctn *WebSocketController) SetRouter(m *martini.ClassicMartini) {
	m.Get("/socket", sockets.JSON(service.Message{}), MainSocket)
}

func MainSocket( params martini.Params, receiver <-chan *service.Message, sender chan<- *service.Message, done <-chan bool, disconnect chan<- int, err <-chan error ) (int, string) {
	return service.WebSocket.Listen(params, receiver, sender, done, disconnect, err)
}
