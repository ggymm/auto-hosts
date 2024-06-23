package main

import (
	"github.com/ying32/govcl/vcl"
	"github.com/ying32/govcl/vcl/types"
)

type MainForm struct {
	*vcl.TForm

	urlEdit     *vcl.TEdit
	renewButton *vcl.TButton

	domainMemo *vcl.TMemo
	resultMemo *vcl.TMemo

	searchButton   *vcl.TButton
	generateButton *vcl.TButton
}

func (f *MainForm) setupView() {
	enFont := vcl.NewFont()
	//enFont.SetName("JetBrains Mono Medium")
	enFont.SetName("Microsoft YaHei UI")

	zhFont := vcl.NewFont()
	zhFont.SetName("Microsoft YaHei UI")

	f.urlEdit = vcl.NewEdit(f)
	f.urlEdit.SetParent(f)
	f.urlEdit.SetFont(enFont)
	f.urlEdit.SetText("https://public-dns.info/nameservers.txt")
	f.urlEdit.SetAutoSelect(false)

	f.renewButton = vcl.NewButton(f)
	f.renewButton.SetParent(f)
	f.renewButton.SetFont(zhFont)
	f.renewButton.SetCaption("同步 nameserver 列表")

	f.urlEdit.SetBounds(20, 24, 500, 25)
	f.renewButton.SetBounds(540, 24, 160, 25)

	f.domainMemo = vcl.NewMemo(f)
	f.domainMemo.SetParent(f)
	f.domainMemo.SetFont(enFont)
	f.domainMemo.SetReadOnly(true)
	f.domainMemo.SetWordWrap(false)
	f.domainMemo.SetScrollBars(types.SsAutoVertical)

	f.resultMemo = vcl.NewMemo(f)
	f.resultMemo.SetParent(f)
	f.resultMemo.SetFont(enFont)
	f.resultMemo.SetReadOnly(true)
	f.resultMemo.SetWordWrap(false)
	f.resultMemo.SetScrollBars(types.SsAutoVertical)

	f.searchButton = vcl.NewButton(f)
	f.searchButton.SetParent(f)
	f.searchButton.SetFont(zhFont)
	f.searchButton.SetCaption("开始查询")

	f.generateButton = vcl.NewButton(f)
	f.generateButton.SetParent(f)
	f.generateButton.SetFont(zhFont)
	f.generateButton.SetCaption("生成文件")

	f.domainMemo.SetBounds(20, 70, 330, 258)
	f.resultMemo.SetBounds(370, 70, 330, 258)
	f.searchButton.SetBounds(20, 348, 330, 32)
	f.generateButton.SetBounds(370, 348, 330, 32)
}

func (f *MainForm) enableView() {
	f.urlEdit.SetEnabled(true)
	f.renewButton.SetEnabled(true)

	f.searchButton.SetEnabled(true)
	f.generateButton.SetEnabled(true)
}

func (f *MainForm) disableView() {
	f.urlEdit.SetEnabled(false)
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
