package actions

import (
	"github.com/sharon-xa/gomux/internal/database"
)

type Actions struct {
	DB *database.DB
}

type IActions interface {
	CreateScript() error
}

func NewActions() IActions {
	return &Actions{
		DB: database.NewDB(),
	}
}
