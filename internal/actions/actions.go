package actions

import (
	"github.com/sharon-xa/gomux/internal/database"
	"github.com/sharon-xa/gomux/internal/models"
)

// Make sure to populate the Script struct fully before calling this function.
// This function will take the model and turn it into a tmux script string
// and store it in a file.
func CreateScript(s *models.Script) {
	f := &models.File{
		Name:    s.Session.SessionName,
		Type:    s.Session.SessionType,
		Content: s.Stringify(),
	}

	db := database.NewDB()
	err := db.Create(f)
	if err != nil {
	}
}
