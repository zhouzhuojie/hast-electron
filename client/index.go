package main

import (
	"time"

	"github.com/gopherjs/gopherjs/js"
	"github.com/oskca/gopherjs-vue"

	_ "github.com/zhouzhuojie/hast-electron/client/static"
)

func main() {
	a := &App{}
	a.Bootstrap()
}

var (
	// Document represents the global document
	Document = js.Global.Get("document")

	// Remark represents the global remarkjs
	Remark = js.Global.Get("remark")

	// Ace represents the global Ace js object
	Ace = js.Global.Get("ace")

	// Console is the console
	Console = js.Global.Get("console")
)

// ViewModel is the Vue's ViewModel for the #app
type ViewModel struct {
	*js.Object
}

// App is the app struct
type App struct {
	VM *ViewModel
	S  *Slide
	E  *Editor
}

// Bootstrap starts the app
func (a *App) Bootstrap() {
	vue.New("#app", a.VM)
	a.VM = &ViewModel{Object: js.Global.Get("Object").New()}
	a.S = NewSlide()
	a.E = NewEditor()
	a.syncEditorToSlides()
}

func (a *App) syncEditorToSlides() {
	changeCh := a.E.GetChangeCh()
	go func() {
		value := ""
		for {
			select {
			case value = <-changeCh:
			case <-time.After(300 * time.Millisecond):
				a.S.SetContent(value)
				a.S.Render()
			}
		}
	}()
}
