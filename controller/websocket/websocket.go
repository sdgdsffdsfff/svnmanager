package websocket

import(
	"github.com/go-martini/martini"
	sockets "github.com/beatrichartz/martini-sockets"
	"time"
	"sync"
	"king/config"
	"king/helper"
	"king/utils/JSON"
)


/***********************
Websocket register
************************/

var callbacks = map[string]func()JSON.Type{}

func Register(name string, callback func()JSON.Type){
	if _, found := callbacks[name]; found {
		return
	}
	callbacks[name] = callback
}

/***********************
Websocket processer
************************/

type Message struct {
	Method string `json:"method"`
	Params interface {} `json:"data"`
}

type Client struct {
	in         <-chan *Message
	out        chan<- *Message
	done       <-chan bool
	err        <-chan error
	disconnect chan<- int
	message *Message // 为无限通知预留的调用状态，为最后一次请求内容
}

type Chat struct {
	sync.Mutex
	clients []*Client
	loopEnable bool
}

func (c *Chat) appendClient(client *Client) {
	c.Lock()
	c.clients = append(c.clients, client)
//  告诉大家我来了
//	for _, xc := range c.clients {
//		if xc != client {
//			xc.ount <- &Message{};
//		}
//	}
	c.Unlock()
}

// the chat
var chat *Chat

func (c *Chat) removeClient(client *Client) {
	c.Lock()
	defer c.Unlock()

	for index, xc := range c.clients {
		if xc == client {
			c.clients = append(c.clients[:index], c.clients[(index+1):]...)
		}
	}
}

func (c *Chat) emit(client *Client, msg *Message) {
	c.Lock()
	method := msg.Method
	if _, found := callbacks[method]; method != "" && found {
		client.out <- &Message{method, callbacks[method]()}
	} else {
		client.out <- &Message{method, helper.Error("method undefined")}
	}
	c.Unlock()
}

func (c *Chat) sendMessageToAllClient(msg *Message){
	c.Lock()
	for _, xc := range c.clients {
		xc.out <- msg
	}
	defer c.Unlock()
}

/*
主动推送方法，将会根据method不停向客户端推送内容
 */
func (c *Chat) loopPushFrame(){
	for {
		c.Lock()
		helper.AsyncMap(c.clients, func(i int) bool{
			client := c.clients[i]
			if client != nil && client.message != nil {
				if method := client.message.Method; method != "" {
					client.out <- &Message{method, callbacks[method]()}
				}
			}
			return false
		})
		c.Unlock()
		time.Sleep(3*time.Second)
	}
}

/***********************
Websocket Controller
************************/

type WebSocketController struct{}

func init() {
	config.AppendValue(config.Controller, &WebSocketController{})
	chat = &Chat{sync.Mutex{}, make([]*Client, 0), true}
}

func (ctn *WebSocketController) SetRouter(m *martini.ClassicMartini) {
	m.Get("/socket", sockets.JSON(Message{}), MainSocket)
}

func MainSocket( params martini.Params, receiver <-chan *Message, sender chan<- *Message, done <-chan bool, disconnect chan<- int, err <-chan error ) (int, string) {
	client := &Client{receiver, sender, done, err, disconnect, nil}
	chat.appendClient(client)
// 暂停使用无限通知，让客户端自己通过循环调用
//	if chat.loopEnable {
//		go chat.loopPushFrame()
//		chat.loopEnable = false
//	}
	for {
		select {
		case <-client.err:
			// Don't try to do this:
			// client.out <- &Message{"system", "system", "There has been an error with your connection"}
			// The socket connection is already long gone.
			// Use the error for statistics etc
		case msg := <-client.in:
			chat.emit(client, msg)
		case <-client.done:
			chat.removeClient(client)
			return 200, "OK"
		}
	}
}
