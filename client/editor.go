package main

import "github.com/gopherjs/gopherjs/js"

// Editor is an Editor
type Editor struct {
	ace *js.Object

	contentCh chan string
	pageNumCh chan int
}

// NewEditor creates a new Editor
func NewEditor() *Editor {
	ace := Ace.Call("edit", "h-editor-wrapper")
	session := ace.Call("getSession")
	session.Call("setMode", "ace/mode/markdown")
	e := &Editor{
		ace:       ace,
		contentCh: make(chan string, 1),
	}
	e.StartSendingChanges()
	return e
}

// GetContentCh gets the ChangeCh
func (e *Editor) GetContentCh() chan string {
	return e.contentCh
}

// GetValue gets the current editor value
func (e *Editor) GetValue() string {
	return e.ace.Call("getValue").String()
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
			Console.Call("log", e.ace.Call("getCursorPosition"))
		}()
	})
}
