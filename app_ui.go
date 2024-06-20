package main

import (
	"github.com/ying32/govcl/vcl"
	"github.com/ying32/govcl/vcl/types"
)

type MainForm struct {
	*vcl.TForm

	devsCombo   *vcl.TComboBox
	renewButton *vcl.TButton

	domainMemo *vcl.TMemo
	resultMemo *vcl.TMemo

	searchButton   *vcl.TButton
	generateButton *vcl.TButton
}

func (f *MainForm) setupView() {
	label := vcl.NewLabel(f)
	label.SetParent(f)
	label.SetCaption("选择网卡")
	label.SetAutoSize(false)
	label.SetAlignment(types.TaCenter)

	f.devsCombo = vcl.NewComboBox(f)
	f.devsCombo.SetParent(f)
	f.devsCombo.SetStyle(types.CsDropDownList)

	f.renewButton = vcl.NewButton(f)
	f.renewButton.SetParent(f)
	f.renewButton.SetCaption("刷新网卡")

	label.SetBounds(20, 28, 80, 25)
	f.devsCombo.SetBounds(100, 24, 480, 25)
	f.renewButton.SetBounds(600, 24, 100, 25)

	f.domainMemo = vcl.NewMemo(f)
	f.domainMemo.SetParent(f)
	f.domainMemo.SetReadOnly(true)
	f.domainMemo.SetScrollBars(types.SsAutoBoth)

	f.resultMemo = vcl.NewMemo(f)
	f.resultMemo.SetParent(f)
	f.resultMemo.SetReadOnly(true)
	f.resultMemo.SetScrollBars(types.SsAutoBoth)

	f.searchButton = vcl.NewButton(f)
	f.searchButton.SetParent(f)
	f.searchButton.SetCaption("开始查询")

	f.generateButton = vcl.NewButton(f)
	f.generateButton.SetParent(f)
	f.generateButton.SetCaption("生成文件")

	f.domainMemo.SetBounds(20, 70, 330, 258)
	f.resultMemo.SetBounds(370, 70, 330, 258)
	f.searchButton.SetBounds(20, 348, 330, 32)
	f.generateButton.SetBounds(370, 348, 330, 32)
}

func (f *MainForm) enableView() {
	f.devsCombo.SetEnabled(true)
	f.renewButton.SetEnabled(true)

	f.searchButton.SetEnabled(true)
	f.generateButton.SetEnabled(true)
}

func (f *MainForm) disableView() {
	f.devsCombo.SetEnabled(false)
	f.renewButton.SetEnabled(false)

	f.searchButton.SetEnabled(false)
	f.generateButton.SetEnabled(false)
}

func (f *MainForm) OnFormCreate(_ vcl.IObject) {
	f.SetCaption("DNS查询工具")
	f.SetWidth(720)
	f.SetHeight(400)
	f.SetPosition(types.PoScreenCenter)
	f.SetBorderStyle(types.BsSingle)

	f.setupView()
}
