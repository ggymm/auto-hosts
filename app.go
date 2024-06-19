package main

import (
	"os"
	"path/filepath"
	"runtime"
	"strings"

	"github.com/pkg/errors"
	"github.com/ying32/govcl/vcl"

	"auto-hosts/log"
)

type App struct {
	wd string
	ui *MainForm

	scanner *Scanner

	devices []*Device

	domains     []string
	nameservers []string
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
		panic(errors.WithStack(err))
	}
}

func (a *App) showUI() {
	vcl.DEBUG = false
	vcl.Application.Initialize()
	vcl.Application.SetMainFormOnTaskBar(true)

	// 创建主窗口
	vcl.Application.CreateForm(&a.ui)

	// 设置窗口显示事件
	a.ui.SetOnShow(func(sender vcl.IObject) {
		go func() {
			a.devices = GetDevices()
			vcl.ThreadSync(func() {
				for _, d := range a.devices {
					a.ui.devices.Items().Add(d.String())
				}
			})
			log.Info().Msg("devices loaded")
		}()

		go func() {
			a.domains = GetDomains()
			vcl.ThreadSync(func() {
				for _, domain := range a.domains {
					a.ui.domain.Lines().Add(domain)
				}
			})
			log.Info().Msg("domains loaded")
		}()
		go func() {
			a.nameservers = GetNameservers()
			log.Info().Msg("nameservers loaded")
		}()
	})

	// 启动应用
	vcl.Application.Run()
}
