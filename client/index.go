package main

import "github.com/gopherjs/gopherjs/js"

func main() {
	a := &App{
		Document: js.Global.Get("document"),
		Remark:   js.Global.Get("remark"),
	}
	a.Bootstrap()
}

// App is the app struct
type App struct {
	*js.Object
	Document *js.Object
	Remark   *js.Object
}

// Bootstrap starts the app
func (a *App) Bootstrap() {
	a.renderMarkdown()
}

func (a *App) renderMath() {
	body := a.Document.Get("body")
	js.Global.Call("renderMathInElement", body, js.M{
		"delimiters": js.S{
			js.M{"left": "$$", "right": "$$", "display": true},
			js.M{"left": "$", "right": "$", "display": false},
			js.M{"left": "\\[", "right": "\\]", "display": true},
			js.M{"left": "\\(", "right": "\\)", "display": false},
		},
	})
}

func (a *App) renderMarkdown() {
	option := js.M{
		"container": a.Document.Call("getElementById", "h-slides-wrapper"),
	}
	a.Remark.Call("create", option, a.renderMath)
}
