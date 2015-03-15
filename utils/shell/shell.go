package shell

import (
	"os/exec"
	"strings"
)

type cmd struct {
	cmd *exec.Cmd
	next *cmd
}

func (r *cmd) Pipe( command string ) *cmd {
	r.next = Cmd(command)
	return r.next
}

func (r *cmd) Output() (string, error) {
	if r.next != nil {
		for r.next != nil {
			out, err := r.cmd.StdoutPipe()
			if err != nil {
				return "", err
			}
			r.next.cmd.Stdin = out
			r.cmd.Start()
			r.cmd.Wait()
		}

		out, err := r.cmd.Output()
		if err != nil {
			return "", err
		}
		return string(out), nil
	} else {
		out, err := r.cmd.Output()
		if err != nil {
			return "", err
		}
		return string(out), nil
	}
}

func Cmd(command string) *cmd {
	parts := strings.Fields( command )
	head := parts[0]
	parts = parts[1:len(parts)]

	c := &cmd{exec.Command(head, parts...), nil}

	return c
}
