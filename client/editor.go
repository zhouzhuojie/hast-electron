package main

import "github.com/gopherjs/gopherjs/js"

// Editor is an Editor
type Editor struct {
	aceSession *js.Object
	changeCh   chan string
}

// NewEditor creates a new Editor
func NewEditor() *Editor {
	ace := Ace.Call("edit", "h-editor-wrapper")
	session := ace.Call("getSession")
	session.Call("setMode", "ace/mode/markdown")
	e := &Editor{
		aceSession: session,
		changeCh:   make(chan string, 1),
	}
	e.StartSendingChanges()
	return e
}

// GetChangeCh gets the ChangeCh
func (e *Editor) GetChangeCh() chan string {
	return e.changeCh
}

// GetValue gets the current editor value
func (e *Editor) GetValue() string {
	return e.aceSession.Call("getValue").String()
}

// StartSendingChanges starts to send changes to the channel
func (e *Editor) StartSendingChanges() {
	e.aceSession.Call("on", "change", func() {
		go func() {
			e.changeCh <- e.GetValue()
		}()
	})
}
