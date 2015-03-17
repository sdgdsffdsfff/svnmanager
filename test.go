package main

import (
"fmt"
sh "github.com/codeskyblue/go-sh"
)

func main() {
	session := sh.NewSession()
	err := session.SetDir("/").Command("pwd").Run()
	fmt.Print(err)
}
