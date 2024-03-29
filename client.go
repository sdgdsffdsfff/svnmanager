package main

import (
	"king/bootstrap"
	"king/config"
	_ "king/routes/client"
	_ "king/rpc/client"
	"king/service/host"
	"log"
	"os"
)

func main() {
	port := "4000"
	argLen := len(os.Args)

	if argLen > 1 {
		port = os.Args[1]
	}

	log.Println("Running client on port", port)

	//设置Rpc地址
	config.Set("rpc", config.GetString("master")+"/rpc")

	bootstrap.Start(port, func() {
		//在服务器上添加自己，必须确定唯一属性
		host.Detail.Ip = "192.168.1.111"
		host.Detail.InternalIp = "192.168.1.111"
		host.Detail.Port = port
		host.Connect()
	})
}
