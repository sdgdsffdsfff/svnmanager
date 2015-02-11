package main

import (
	. "gopkg.in/godo.v1"
)

func Tasks(p *Project){

	p.Task("default", func(){
			Bash("echo Hello $USER!")
		})

	p.Task("server", W{"king/server.go"}, func(){

		})
}

func main(){
	Godo(Tasks)
}
