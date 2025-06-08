package database

import (
	"fmt"
	"os"

	"github.com/sharon-xa/gomux/internal/models"
)

const (
	gomuxDataPath   string = "/.local/share/gomux/"
	gomuxConfigPath string = "/.config/gomux/"
	gomuxCachePath  string = "/.cache/gomux/"
)

type DB struct {
	File             *models.File
	gomuxDirFullPath string
	dataDirFullPath  string
	dataStateFile    string
}

/*
dataState.json
data/
	|___ example/
	|				|___ example.tmux.sh
	| 			|___ example.tmux.sh
	|___ example/
					|___ example.tmux.sh

the folder itself is called "gomux/"
*/

func InitDB() {
	homeDir, err := os.UserHomeDir()
	if err != nil {
		fmt.Println("Couldn't get the home directory")
		panic(err)
	}

	gomuxDirFullPath := homeDir + gomuxDataPath
	dataStateFile := gomuxDirFullPath + "dataState.json"
	dataDirFullPath := gomuxDirFullPath + "data/"

	if _, err := os.Stat(gomuxDirFullPath); os.IsNotExist(err) {
		err = os.MkdirAll(gomuxDirFullPath, 0o755)
		if err != nil {
			fmt.Println("Couldn't make all direcotories in the following path:", gomuxDirFullPath)
			panic(err)
		}
	}

	if _, err := os.Stat(dataStateFile); os.IsNotExist(err) {
		err := os.WriteFile(dataStateFile, []byte(""), 0o644)
		if err != nil {
			fmt.Println("Couldn't create the dataState file: \"dataState.json\"")
			panic(err)
		}
	}

	if _, err := os.Stat(dataDirFullPath); os.IsNotExist(err) {
		err = os.Mkdir(dataDirFullPath, 0o755)
		if err != nil {
			fmt.Println(
				"Couldn't make direcotory inside gomux directory.\nFull Path:",
				gomuxDirFullPath,
			)
			panic(err)
		}
	}
}

func NewDB() *DB {
	homeDir, _ := os.UserHomeDir()
	gomuxDirFullPath := homeDir + gomuxDataPath
	dataStateFile := gomuxDirFullPath + "dataState.json"
	dataDirFullPath := gomuxDirFullPath + "data/"

	return &DB{
		gomuxDirFullPath: gomuxDirFullPath,
		dataStateFile:    dataStateFile,
		dataDirFullPath:  dataDirFullPath,
	}
}

// This function will create a tmux script file
func (db *DB) Create(f *models.File) error {
	return nil
}
