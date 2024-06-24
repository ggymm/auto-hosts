package main

import (
	_ "github.com/ying32/govcl/pkgs/winappres"
)

func init() {
	log.Init()
}

func main() {
	app := NewApp()
	app.Run()
}
