package database

import (
	"fmt"
	"os"
)

func checkDirectoryExistence(dirPath string) {
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
