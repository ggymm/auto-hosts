package main

import (
	"github.com/ying32/govcl/vcl"
	"github.com/ying32/govcl/vcl/types"
)

type MainForm struct {
	*vcl.TForm

	fix     *vcl.TComboBox
	devices *vcl.TComboBox

	domain *vcl.TMemo
	result *vcl.TMemo
}

func (f *MainForm) fixComboBox() {
	f.fix.SetFocus()
}

func (f *MainForm) setupMenu() {
	menu := vcl.NewMainMenu(f)

	item := vcl.NewMenuItem(f)
	item.SetCaption("DNS服务器(&N)")
	menu.Items().Add(item)

	item = vcl.NewMenuItem(f)
	item.SetCaption("域名列表(&D)")
	menu.Items().Add(item)

	item = vcl.NewMenuItem(f)
	item.SetCaption("帮助(&H)")
	menu.Items().Add(item)

	item = vcl.NewMenuItem(f)
	item.SetCaption("退出(&Q)")
	item.SetOnClick(func(sender vcl.IObject) {
		f.Close()
	})
	menu.Items().Add(item)
}

func (f *MainForm) setupContent() {
	label := vcl.NewLabel(f)
	label.SetParent(f)
	label.SetCaption("选择设备网卡")
	label.SetAutoSize(false)
	label.SetAlignment(types.TaLeftJustify)

	combo := vcl.NewComboBox(f)
	combo.SetParent(f)
	combo.SetStyle(types.CsDropDownList)
	combo.SetOnSelect(func(sender vcl.IObject) {
	})
	combo.SetOnCloseUp(func(_ vcl.IObject) {
		f.fixComboBox()
	})
	f.devices = combo

	button := vcl.NewButton(f)
	button.SetParent(f)
	button.SetCaption("刷新网卡列表")
	button.SetOnClick(func(sender vcl.IObject) {
		go func() {
			dev := GetDevices()
			vcl.ThreadSync(func() {
				for _, d := range dev {
					combo.Items().Add(d.String())
				}
			})
		}()
	})

	label.SetBounds(20, 24, 100, 25)
	combo.SetBounds(120, 20, 420, 25)
	button.SetBounds(560, 20, 140, 25)

	domain := vcl.NewMemo(f)
	domain.SetParent(f)
	domain.SetReadOnly(true)
	domain.SetScrollBars(types.SsAutoBoth)
	f.domain = domain

	result := vcl.NewMemo(f)
	result.SetParent(f)
	result.SetReadOnly(true)
	result.SetScrollBars(types.SsAutoBoth)
	f.result = result

	button = vcl.NewButton(f)
	button.SetParent(f)
	button.SetCaption("开始查询")
	button.SetOnClick(func(sender vcl.IObject) {

	})

	domain.SetBounds(20, 65, 330, 248)
	result.SetBounds(370, 65, 330, 248)
	button.SetBounds(20, 330, 680, 32)
}

func (f *MainForm) OnFormCreate(_ vcl.IObject) {
	f.SetCaption("DNS查询工具")
	f.SetWidth(720)
	f.SetHeight(400)
	f.SetPosition(types.PoScreenCenter)
	f.SetBorderStyle(types.BsSingle)
	f.SetOnShow(func(_ vcl.IObject) {
		f.fixComboBox()
	})
	f.SetDoubleBuffered(true)

	f.setupMenu()
	f.setupContent()

	f.fix = vcl.NewComboBox(f)
	f.fix.SetParent(f)
	f.fix.SetBounds(0, 0, 0, 0)
	f.fix.SetStyle(types.CsDropDownList)
}
