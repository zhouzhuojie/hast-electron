package main

import (
	"github.com/gopherjs/gopherjs/js"
	electron "github.com/oskca/gopherjs-electron"
	nodejs "github.com/oskca/gopherjs-nodejs"
)

func main() {
	nodejs.Require("events").Get("EventEmitter").Set("defaultMaxListeners", 0)
	app := electron.GetApp()
	app.On(electron.EvtAppReady, func(args ...*js.Object) {
		opt := electron.NewBrowserWindowOption()
		bw := electron.NewBrowserWindow(opt)
		bw.LoadURL("file://"+nodejs.DirName()+"/index.html", nil)
	})
}
