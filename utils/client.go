package utils

import(
	"net/http"
	"net"
	"time"
	"fmt"
)



var Client *http.Client

func init(){

	Client = &http.Client{
		Transport: &http.Transport{
			Dial: func(netw, addr string) (net.Conn, error) {
				c, err := net.DialTimeout(netw, addr, time.Second*20)
				if err != nil {
					fmt.Println("dail timeout", err)
					return nil, err
				}
				return c, nil
			},
			MaxIdleConnsPerHost:   10,
			ResponseHeaderTimeout: time.Second * 20,
		},
	}
}

