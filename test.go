package main

import (
	sh "github.com/codeskyblue/go-sh"
	"fmt"
)

func main() {
	session := sh.NewSession()
	output, err := session.SetDir("/opt/wings").Command("mvn clean:clean compile").Output()
	fmt.Println(string(output), err)
}
