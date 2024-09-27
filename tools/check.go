package tools

import (
	"os"
)

func CheckFileOrFolderExits(filepath string) bool {
	_, err := os.Stat(filepath)
	return !os.IsNotExist(err)
}
