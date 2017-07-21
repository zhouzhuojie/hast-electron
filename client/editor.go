package main

import (
	"strings"

	"github.com/gopherjs/gopherjs/js"
)

const (
	pageDivider = "---"
)

// Editor is an Editor
type Editor struct {
	ace *js.Object

	contentCh chan string
	pageNumCh chan int
}

// NewEditor creates a new Editor
func NewEditor() *Editor {
	ace := Ace.Call("edit", "h-editor-wrapper")
	ace.Call("setTheme", "ace/theme/chrome")
	ace.Call("setFontSize", 15)
	session := ace.Call("getSession")
	session.Call("setMode", "ace/mode/markdown")
	session.Call("setUseWrapMode", true)
	e := &Editor{
		ace:       ace,
		contentCh: make(chan string, 1),
		pageNumCh: make(chan int, 1),
	}
	e.StartSendingChanges()
	return e
}

// GetContentCh gets the contentCh
func (e *Editor) GetContentCh() chan string {
	return e.contentCh
}

// GetPageNumCh gets the pageNumCh
func (e *Editor) GetPageNumCh() chan int {
	return e.pageNumCh
}

// GetValue gets the current editor value
func (e *Editor) GetValue() string {
	return e.ace.Call("getValue").String()
}

// SetValue sets the current editor value
func (e *Editor) SetValue(val string) {
	e.ace.Call("setValue", val, -1)
}

func (e *Editor) getCursorPageNum() int {
	cursorRowNum := e.ace.Call("getCursorPosition").Get("row").Int()
	lines := strings.Split(e.GetValue(), "\n")
	pageNum := 1
	for i, line := range lines {
		if i == cursorRowNum {
			break
		}
		if line == pageDivider {
			pageNum++
		}
	}
	return pageNum
}

// StartSendingChanges starts to send changes to the channel
func (e *Editor) StartSendingChanges() {
	e.ace.Call("on", "change", func() {
		go func() {
			e.contentCh <- e.GetValue()
		}()
	})

	e.ace.Call("getSelection").Call("on", "changeCursor", func() {
		go func() {
			e.pageNumCh <- e.getCursorPageNum()
		}()
	})
}
