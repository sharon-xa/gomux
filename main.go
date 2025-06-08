package main

import (
	"github.com/sharon-xa/gomux/cmd"
	"github.com/sharon-xa/gomux/internal/database"
)

func main() {
	database.InitDB()
	cmd.Execute()
}
