package commands

import (
	"encoding/hex"
	"encoding/json"
	"io/ioutil"
	"os"
)

type UploadCommand struct {
	FilePath string
	Data     string
	Append   bool
}

func (u *UploadCommand) createFile(data []byte) error {
	return ioutil.WriteFile(u.FilePath, data, 0644)
}
func (u *UploadCommand) appendToFile(data []byte) error {
	file, err := os.OpenFile(u.FilePath, os.O_APPEND|os.O_WRONLY, 0644)
	if err != nil {
		return err
	}
	defer file.Close()
	if _, err := file.Write(data); err != nil {
		return err
	}
	return nil
}

func (u *UploadCommand) Save() error {
	binData, _ := hex.DecodeString(u.Data)
	if u.Append {
		return u.appendToFile(binData)
	} else {
		return u.createFile(binData)
	}
}

func SaveFile(payload []byte) error {
	var u UploadCommand
	err := json.Unmarshal(payload, &u)
	if err != nil {
		return err
	} else {
		return u.Save()
	}
}
