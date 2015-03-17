package main

import (
"fmt"
sh "github.com/codeskyblue/go-sh"
)

func main() {
	session := sh.NewSession()
	err := session.SetDir("/opt/wings").Command("mvn clean:clean compile").Run()
}
