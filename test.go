package main

import (
	sh "github.com/codeskyblue/go-sh"
	"fmt"
)

func main() {
	session := sh.NewSession()
	output, err := session.Command("sh", "shells/mvn.sh").Output()
	fmt.Println(string(output), err)
}
