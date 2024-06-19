package main

import (
	_ "github.com/ying32/govcl/pkgs/winappres"

	"auto-hosts/log"
)

func init() {
	log.Init()
}

func main() {
	app := NewApp()
	app.init()

	app.showUI()
}
