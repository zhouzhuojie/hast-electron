package main

import (
	"github.com/gopherjs/gopherjs/js"

	_ "github.com/zhouzhuojie/hast-electron/client/static"
)

func main() {
	a := &App{}
	a.Bootstrap()
}

var (
	// Document represents the global document
	Document = js.Global.Get("document")

	// Console is the console
	Console = js.Global.Get("console")

	// Lodash is the lodash lib
	Lodash = js.Global.Get("_")

	// Remark represents the global remarkjs
	Remark = js.Global.Get("remark")

	// Ace represents the global Ace js object
	Ace = js.Global.Get("ace")
)

// App is the app struct
type App struct {
	S *Slide
	E *Editor
}

// Bootstrap starts the app
func (a *App) Bootstrap() {
	a.S = NewSlide()
	a.E = NewEditor()
	a.syncEditorToSlides()
}

func (a *App) syncEditorToSlides() {
	changeCh := a.E.GetChangeCh()
	go func() {
		for {
			value := <-changeCh
			a.S.SetContent(value)
			a.S.Render()
		}
	}()
}
