package utils

import(
	"os"
)

func PathEnable(path string) error {
	if _, err := os.Stat(path); err != nil {
		return err
	}
	return nil
}
