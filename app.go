package main

import (
	"fmt"
	"os"
	"path/filepath"
	"runtime"
	"slices"
	"strings"
	"sync"
	"time"

	"github.com/ggymm/ping"
	"github.com/pkg/errors"
	"github.com/ying32/govcl/vcl"
)

type App struct {
	mu *sync.Mutex

	wd   string
	view *MainForm

	scanner *Scanner

	hosts       []string
	domains     []string
	nameservers []string
}

func NewApp() *App {
	return &App{
		mu:      &sync.Mutex{},
		scanner: NewScanner(),
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

func (a *App) showView() {
	vcl.DEBUG = false
	vcl.Application.SetScaled(true)
	vcl.Application.Initialize()
	vcl.Application.SetMainFormOnTaskBar(true)

	// 创建主窗口
	vcl.Application.CreateForm(&a.view)

	// 设置窗口显示事件
	a.view.SetOnShow(func(sender vcl.IObject) {
		go func() {
			a.hosts = LoadHosts()
			vcl.ThreadSync(func() {
				for _, host := range a.hosts {
					if len(host) > 0 {
						a.view.resultMemo.Lines().Add(host)
					}
				}
			})
			log.Info().Msg("init hosts")
		}()
		go func() {
			a.domains = LoadDomains()
			vcl.ThreadSync(func() {
				for _, domain := range a.domains {
					if len(domain) > 0 {
						a.view.domainMemo.Lines().Add(domain)
					}
				}
			})
			log.Info().Msg("init domains")
		}()
		go func() {
			a.nameservers = LoadNameservers()
			log.Info().Msg("init nameservers")
		}()

		a.view.renewButton.SetOnClick(func(sender vcl.IObject) {
			a.view.disableView()
			go func() {
				defer a.view.enableView()

				//RenewNameservers()
				//a.nameservers = LoadNameservers()
			}()
		})

		a.view.searchButton.SetOnClick(func(sender vcl.IObject) {
			a.view.disableView()
			a.view.resultMemo.Clear()

			go func() {
				// 重置 view 按钮状态
				defer func() {
					a.view.enableView()
					a.view.searchButton.SetCaption("开始查询")
				}()

				l := len(a.domains)
				for i, domain := range a.domains {
					log.Info().Str("domain", domain).Msg("scan ip")
					vcl.ThreadSync(func() {
						a.view.searchButton.SetCaption(fmt.Sprintf("正在查询 （%d/%d）", i+1, l))
					})

					ip := ""
					list := a.scanner.Scan(domain, a.nameservers)
					if len(list) > 0 {
						wg := &sync.WaitGroup{}
						for _, item := range list {
							wg.Add(1)
							log.Info().
								Str("1.ip", item.ip).
								Str("3.domain", domain).Msg("ping ip")

							go func(item *Info) {
								defer wg.Done()

								p, _ := ping.NewPinger(item.ip)
								p.Count = 4 // 尝试次数
								p.Timeout = 1 * time.Second
								p.SetPrivileged(true)
								err := p.Run()
								if err != nil {
									return
								}
								stats := p.Statistics()
								if stats.PacketsRecv != 0 {
									item.rtt = stats.AvgRtt
								} else {
									item.rtt = 99 * time.Second
								}
							}(item)
						}
						wg.Wait()

						// 排序
						slices.SortFunc(list, func(i, j *Info) int {
							if i.rtt < j.rtt {
								return -1
							} else {
								return 1
							}
						})
						ip = list[0].ip + " " + domain

						// 保存到文件
						ips := make([]string, 0)
						for _, item := range list {
							ips = append(ips, item.String())
						}
						err := writeLines(fmt.Sprintf("ips/%s.txt", domain), ips)
						if err != nil {
							log.Error().
								Str("domain", domain).
								Err(errors.WithStack(err)).Msg("write domain ips error")
							continue
						}
					} else {
						ip = "unknown"
					}

					vcl.ThreadSync(func() {
						a.view.resultMemo.Lines().Add(ip)
					})
				}
			}()
		})

		// 生成 hosts 文件
		a.view.generateButton.SetOnClick(func(sender vcl.IObject) {
			go func() {
				list := make([]string, 0)
				lines := a.view.resultMemo.Lines()
				count := lines.Count()
				for i := int32(0); i < count; i++ {
					list = append(list, lines.S(i))
				}
				err := writeLines(hostsFile, list)
				if err != nil {
					log.Error().Err(errors.WithStack(err)).Msg("write hosts error")
					return
				}
				vcl.ThreadSync(func() {
					vcl.ShowMessage("生成 hosts 文件成功")
				})
			}()
		})
	})

	// 启动应用
	vcl.Application.Run()
}

func (a *App) Run() {
	a.init()
	a.showView()
}
