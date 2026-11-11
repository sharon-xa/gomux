// Package store persists workspace specs and their generated scripts as files
// under ~/.local/share/gomux. It does not render anything itself.
package store

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

type Store struct {
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

// Init creates the gomux data directories if they don't exist yet.
func Init() {
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

	ensureDir(dataStateDir)
	ensureDir(dataDir)
}

func New() *Store {
	homeDir, _ := os.UserHomeDir()

	gomuxDirPath := path.Join(homeDir + GOMUX_DATA_PATH)
	dataStateDir := path.Join(gomuxDirPath, "dataState")
	dataDir := path.Join(gomuxDirPath, "data")

	return &Store{
		gomuxDirPath: gomuxDirPath,
		dataStateDir: dataStateDir,
		dataDir:      dataDir,
	}
}

func ensureDir(dirPath string) {
	if _, err := os.Stat(dirPath); os.IsNotExist(err) {
		err = os.Mkdir(dirPath, 0o755)
		if err != nil {
			fmt.Println(
				"Couldn't make direcotory inside gomux directory.\nFull Path:",
				dirPath,
			)
			panic(err)
		}
	}
}

func (s *Store) saveFile(filepath string, fileContent []byte, executable bool) error {
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
