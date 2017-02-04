package util

import (
	"encoding/json"
	"io/ioutil"
	"os"
)

func Marshal(filename string, inter interface{}) {
	configBytes, _ := json.Marshal(inter)
	_ = ioutil.WriteFile(filename, configBytes, 0755)
}

func Unmarshal(filename string, inter interface{}) {
	configBytes, _ := ioutil.ReadFile(filename)
	json.Unmarshal(configBytes, inter)
}
func CreateFileIfNotExist(file string) bool {
	if _, err := os.Stat(file); err != nil {
		if os.IsNotExist(err) {
			// file does not exist
			_, _ = os.Create(file)
		} else {
			// other error
		}
		return false
	}
	return true
}

func CreateDirectoryIfNotExist(dir string) bool {
	if _, err := os.Stat(dir); err != nil {
		if os.IsNotExist(err) {
			// file does not exist
			os.MkdirAll(dir, os.ModePerm)
		} else {
			// other error
		}
		return false
	}
	return true
}
