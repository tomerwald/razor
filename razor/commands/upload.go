package commands

import (
	"encoding/hex"
	"encoding/json"
	"io/ioutil"
	"log"
	"os"
)

type UploadCommand struct {
	FilePath string
	Data     string
	Append   bool
}

func (u *UploadCommand) createFile(data []byte) {
	err := ioutil.WriteFile(u.FilePath, data, 0644)
	if err != nil {
		log.Fatal(err)
	}
}
func (u *UploadCommand) appendToFile(data []byte) {
	file, err := os.OpenFile(u.FilePath, os.O_APPEND|os.O_WRONLY, 0644)
	if err != nil {
		log.Fatal(err)
	}
	defer file.Close()
	if _, err := file.Write(data); err != nil {
		log.Fatal(err)
	}
}

func (u *UploadCommand) Save() []byte {
	binData, _ := hex.DecodeString(u.Data)
	if u.Append {
		u.appendToFile(binData)
	} else {
		u.createFile(binData)
	}
	return nil
}

func SaveFile(payload []byte) []byte {
	var u UploadCommand
	err := json.Unmarshal(payload, &u)
	if err != nil {
		log.Fatal(err)
		return nil
	} else {
		return u.Save()
	}
}
