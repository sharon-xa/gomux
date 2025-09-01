package actions

import (
	"github.com/sharon-xa/gomux/internal/database"
	"github.com/sharon-xa/gomux/internal/models"
)

type Actions struct {
	DB *database.DB
}

type IActions interface {
	CreateScript(*models.Script) error
}

func NewActions() IActions {
	return &Actions{
		DB: database.NewDB(),
	}
}

func (a *Actions) CreateScript(s *models.Script) error {
	err := a.DB.CreateScript(s)
	if err != nil {
		return err
	}
	return nil
}
