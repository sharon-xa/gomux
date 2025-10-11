package actions

import (
	"github.com/sharon-xa/gomux/internal/database"
	"github.com/sharon-xa/gomux/internal/models"
)

type Actions struct {
	DB            *database.DB
	CurrentScript *models.Script
}

type IActions interface {
	CreateScriptAllInOne(*models.Script) error
	CreateScript() error
	AddSession(*models.Session) error
	AddWindow(*models.Window) error
	AddPane(*models.Pane) error
}

func NewActions() IActions {
	return &Actions{
		DB:            database.NewDB(),
		CurrentScript: &models.Script{},
	}
}
