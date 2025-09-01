package actions

import "github.com/sharon-xa/gomux/internal/models"

func (a *Actions) CreateScript() error {
	err := a.DB.AddScript(&models.Script{})
	if err != nil {
		return err
	}
	return nil
}
