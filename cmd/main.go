package main

import (
	"auto-hosts"
	_ "github.com/ying32/govcl/pkgs/winappres"
)

func main() {
	app := autohosts.NewApp()
	app.Run()
}
