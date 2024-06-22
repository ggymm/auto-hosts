package main

import (
	"fmt"
	"os"
	"path/filepath"
	"runtime"
	"strings"
	"sync"

	"github.com/pkg/errors"
	"github.com/ying32/govcl/vcl"

	"auto-hosts/log"
)

type App struct {
	lock *sync.Mutex

	wd string
	ui *MainForm

	scanner *Scanner

	domains     []string
	nameservers []string
}

func NewApp() *App {
	return &App{
		lock: &sync.Mutex{},
	}
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
	vcl.Application.SetScaled(true)
	vcl.Application.Initialize()
	vcl.Application.SetMainFormOnTaskBar(true)

	// 创建主窗口
	vcl.Application.CreateForm(&a.ui)

	// 设置窗口显示事件
	a.ui.SetOnShow(func(sender vcl.IObject) {
		go func() {
			a.domains = LoadDomains()
			vcl.ThreadSync(func() {
				for _, domain := range a.domains {
					if len(domain) > 0 {
						a.ui.domainMemo.Lines().Add(domain)
					}
				}
			})
			log.Info().Msg("init domains")
		}()
		go func() {
			a.nameservers = LoadNameservers()
			log.Info().Msg("init nameservers")
		}()

		a.ui.renewButton.SetOnClick(func(sender vcl.IObject) {
			a.ui.disableView()
			go func() {
				defer a.ui.enableView()

				GetNameservers()
				a.nameservers = LoadNameservers()
			}()
		})

		a.ui.searchButton.SetOnClick(func(sender vcl.IObject) {
			a.ui.disableView()
			go func() {
				defer a.ui.enableView()

				a.scanner = NewScanner()
				ret := a.scanner.Run(a.domains, a.nameservers)
				for d, ips := range ret {
					fmt.Println(d)
					for _, ip := range ips {
						fmt.Println(ip)
					}
				}
			}()
		})
	})

	// 启动应用
	vcl.Application.Run()
}

func (a *App) Run() {
	a.init()
	a.showUI()
}
