package main

import (
	"fmt"
	"os/exec"
)

func main() {
	output, err := exec.Command("pwd").Output()
	fmt.Println(string(output), err)
	//err  exec: "pwd": executable file not found in $PATH;

	output, err = exec.Command("ls", "-l").Output()
	fmt.Println(string(output), err)
	//exec: "ls": executable file not found in $PATH
}
