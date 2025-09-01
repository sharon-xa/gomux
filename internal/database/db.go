package database

import (
	"fmt"
	"io/fs"
	"os"
	"path"
)

const (
	GOMUX_DATA_PATH   string = "/.local/share/gomux/"
	GOMUX_CONFIG_PATH string = "/.config/gomux/"
	GOMUX_CACHE_PATH  string = "/.cache/gomux/"
)

type DB struct {
	gomuxDirPath string
	dataStateDir string
	dataDir      string
}

/*
dataState/
	|___ example/
	|				|___ example.tmux.json
	| 			|___ example.tmux.json
	|___ example/
					|___ example.tmux.json

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

	gomuxDirFullPath := homeDir + GOMUX_DATA_PATH
	dataStateDir := gomuxDirFullPath + "dataState/"
	dataDir := gomuxDirFullPath + "data/"

	if _, err := os.Stat(gomuxDirFullPath); os.IsNotExist(err) {
		err = os.MkdirAll(gomuxDirFullPath, 0o755)
		if err != nil {
			fmt.Println("Couldn't make all direcotories in the following path:", gomuxDirFullPath)
			panic(err)
		}
	}

	checkDirectoryExistence(dataStateDir)
	checkDirectoryExistence(dataDir)
}

func NewDB() *DB {
	homeDir, _ := os.UserHomeDir()

	gomuxDirPath := path.Join(homeDir + GOMUX_DATA_PATH)
	dataStateDir := path.Join(gomuxDirPath, "dataState")
	dataDir := path.Join(gomuxDirPath, "data")

	return &DB{
		gomuxDirPath: gomuxDirPath,
		dataStateDir: dataStateDir,
		dataDir:      dataDir,
	}
}

func (db *DB) saveFile(filepath string, fileContent []byte, executable bool) error {
	dir, _ := path.Split(filepath)

	if dir != "" {
		err := os.MkdirAll(dir, 0755)
		if err != nil {
			return fmt.Errorf("failed to create directory %s: %w", dir, err)
		}
	}

	perm := fs.FileMode(0644)
	if executable {
		perm = fs.FileMode(0755)
	}

	err := os.WriteFile(filepath, fileContent, perm)
	if err != nil {
		return fmt.Errorf("failed to write file %s: %w", filepath, err)
	}

	return nil
}
