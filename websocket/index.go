package websocket

import (
	"log"
	"github.com/googollee/go-socket.io"
	"king/utils/JSON"
)

func GetServer() *socketio.Server{
	server, err := socketio.NewServer(nil)
	if err != nil {
		log.Fatal(err)
	}

	server.On("connection", func(so socketio.Socket) {
		log.Println("on connection")
		so.Join("chat")
		so.On("chat message", func(msg string){
			data := JSON.Stringify(JSON.Type{
					"A": 10,
					"X": []string{
					"l", "a", "n",
				},
			})
			so.Emit("chat talk", "hahahahahahah")
			so.Emit("chat message", data)
		})

		go func(){
			for {}
		}()

		so.On("disconnection", func(){
			log.Println("on disconnect")
		})

	})
	server.On("error", func(so socketio.Socket, err error){
		log.Println("error:", err)
	})

	return server
}
