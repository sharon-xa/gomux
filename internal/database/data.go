package database

import (
	"log"
	"path/filepath"

	"github.com/sharon-xa/gomux/internal/models"
)

func (db *DB) CreateScript(script *models.Script) error {
	jsonScript, err := script.Jsonify()
	if err != nil {
		log.Println("Failed to jsonify script: ", err)
		return err
	}

	filename := script.Session.Name + ".tmux.sh"
	folderAndFile := filepath.Join(script.Session.Type, filename)

	filePath := filepath.Join(db.dataStateDir, folderAndFile)
	err = db.saveFile(filePath, jsonScript, false)
	if err != nil {
		return err
	}

	scriptAsString := script.Stringify()
	scriptFilePath := filepath.Join(db.dataDir, folderAndFile)

	err = db.saveFile(scriptFilePath, []byte(scriptAsString), true)
	if err != nil {
		return err
	}

	return nil
}
