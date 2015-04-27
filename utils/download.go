package utils

import (
	"net/http"
	"os"
	"io"
	"strings"
)

func Download(fileUrl, where string, args ...string) error {
	resp, err := http.Get(fileUrl)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	token := strings.Split(where, "/")
	name := token[len(token)-1]
	dir := ""
	if where != name {
		dir = where[:(len(where) - len(name))]
	}
	if len(args) > 0 {
		name = args[0]
	}

	err = os.MkdirAll(dir, 0775)
	if err != nil {
		return err
	}

	file, err := os.Create(where)
	if err != nil {
		return err
	}
	defer file.Close()

	_, err = io.Copy(file, resp.Body)
	if err != nil {
		return err
	}

	return nil
}

func RemovePath(path string) error {
	if _, err := os.Stat(path); os.IsNotExist(err) {
		return err
	}
	err := os.Remove(path)
	if err != nil {
		return err
	}
	return nil
}
