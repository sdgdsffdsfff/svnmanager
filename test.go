package main

import (
sh "github.com/codeskyblue/go-sh"
)

func main() {
	session := sh.NewSession()
	session.SetDir("/opt/wings").Command("mvn clean:clean compile").Run()
}
