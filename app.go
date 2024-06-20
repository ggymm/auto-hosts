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

	devices []*Device

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
	vcl.Application.Initialize()
	vcl.Application.SetMainFormOnTaskBar(true)

	// 创建主窗口
	vcl.Application.CreateForm(&a.ui)

	// 设置窗口显示事件
	a.ui.SetOnShow(func(sender vcl.IObject) {
		go func() {
			a.renderDevs()
			log.Info().Msg("init devices")
		}()

		go func() {
			a.domains = GetDomains()
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
			a.nameservers = GetNameservers()
			log.Info().Msg("init nameservers")
		}()

		// 配置按钮事件
		a.ui.renewButton.SetOnClick(func(sender vcl.IObject) {
			go func() {
				a.renderDevs()
			}()
		})

		a.ui.searchButton.SetOnClick(func(sender vcl.IObject) {
			// 获取当前选中的网卡
			i := a.ui.devsCombo.ItemIndex()
			if i < 0 {
				vcl.ShowMessage("请选择一个网卡")
				return
			}
			a.ui.disableView()

			go func() {
				defer a.ui.enableView()

				// 初始化扫描器
				dev := a.devices[i]
				a.scanner = NewScanner()
				err := a.scanner.Init(dev)
				if err != nil {
					log.Error().
						Str("1.dev", dev.String()).
						Err(errors.WithStack(err)).Msg("scanner init error")
					return
				}
				ret := a.scanner.Start(a.nameservers, a.domains)
				for d, ips := range ret {
					fmt.Println(d)
					for _, ip := range ips {
						fmt.Println(ip)
					}
				}

				// 使用 ping 分别测速
			}()
		})
	})

	// 启动应用
	vcl.Application.Run()
}

func (a *App) renderDevs() {
	a.lock.Lock()
	defer a.lock.Unlock()

	dev := GetDevices()
	vcl.ThreadSync(func() {
		a.ui.devsCombo.Clear()
		for _, d := range dev {
			a.ui.devsCombo.Items().Add(d.String())
		}
		a.ui.devsCombo.SetItemHeight(30)
	})
}

func (a *App) Run() {
	a.init()
	a.showUI()
}
