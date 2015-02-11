package shell

import (
	"fmt"
	"os/exec"
	"log"
	"bytes"
	"strings"
)

func Cmd(command string) (string, error) {
	commands := strings.Split(strings.Trim(command, " "), " ")

	cmd := exec.Command(commands[0], commands[1:]...)

	var out bytes.Buffer
	var stderr bytes.Buffer

	cmd.Stdout = &out
	cmd.Stderr = &stderr

	err := cmd.Run()
	if err != nil {
		log.Println(command, "error:", fmt.Sprint(err) + ":" + stderr.String())
		return "", err
	}

	return out.String(), nil
}
