package main

import (
	"os"
	"path/filepath"
	"runtime"
	"strings"

	"github.com/ying32/govcl/vcl"
)

type App struct {
	// app 工作目录
	wd string

	// app ui 主窗口
	ui *MainForm
}

func NewApp() *App {
	return &App{}
}

func (a *App) init() {
	dir := ""
	exe, err := os.Executable()
	if err != nil {
		panic(err)
	}
	path := filepath.Base(exe)
	if !strings.HasPrefix(exe, os.TempDir()) && !strings.HasPrefix(path, "___") {
		dir = filepath.Dir(exe)
	} else {
		_, filename, _, ok := runtime.Caller(0)
		if ok {
			// 需要根据当前文件所处目录，修改相对位置
			dir = filepath.Dir(filename)
		}
	}
	a.wd = filepath.Join(dir, "data")

	// 设置 app 工作目录
	err = os.Chdir(a.wd)
	if err != nil {
		panic(err)
	}

	a.initUI()
}

func (a *App) initUI() {
	vcl.Application.Initialize()
	vcl.Application.SetMainFormOnTaskBar(true)

	// 创建主窗口
	vcl.Application.CreateForm(&a.ui)

	// 设置窗口显示事件
	a.ui.SetOnShow(func(sender vcl.IObject) {
		go func() {
			dev := GetDevices()
			vcl.ThreadSync(func() {
				for _, d := range dev {
					a.ui.devices.Items().Add(d.String())
				}
			})
		}()
		go func() {
			ls := GetDomains()
			vcl.ThreadSync(func() {
				for _, l := range ls {
					a.ui.domain.Lines().Add(l)
				}
			})
		}()
	})

	// 启动应用
	vcl.Application.Run()
}
