package utils

import(
	"net/http"
	"path/filepath"
	"os"
	"io/ioutil"
)

//fileUrl, where, [name]
func Download(fileUrl, where string, args ...string) error {

	name := ""
	if len(args) == 1 {
		name = args[0]
	}else if _, n := filepath.Split(fileUrl); n != "" {
		name = n
	}

	if name == "" {
		return nil
	}

	err := os.MkdirAll(where, 0775)
	if err != nil {
		return err
	}

	resp, err := http.Get(fileUrl)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	file, err := os.Create(where+name)
	if err != nil {
		return err
	}
	defer file.Close()

	content, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return err
	}

	_, err = file.Write(content)
	if err != nil {
		return err
	}

	return nil
}

func RemovePath(path string, where string) error {
	if _, err := os.Stat(where+path); os.IsNotExist(err) {
		return err
	}
	err := os.Remove(where+path)
	if err != nil {
		return err
	}
	return nil
}
