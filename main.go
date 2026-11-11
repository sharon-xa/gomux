package main

import (
	"github.com/sharon-xa/gomux/cmd"
	"github.com/sharon-xa/gomux/internal/store"
)

func main() {
	store.Init()
	cmd.Execute()
}
