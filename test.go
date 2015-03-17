package main

import (
	sh "github.com/codeskyblue/go-sh"
	"fmt"
)

func main() {
	session := sh.NewSession()
	output, err := session.SetDir("/opt/wings").Command("pwd").Output()
	fmt.Println(string(output), err)
}
