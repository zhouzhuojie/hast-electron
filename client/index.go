package main

import (
	"encoding/json"
	"time"

	"github.com/gopherjs/gopherjs/js"
	vue "github.com/oskca/gopherjs-vue"

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

	// IpcRenderer is the ipc that can talk to electron main
	IpcRenderer = js.Global.Call("require", "electron").Get("ipcRenderer")

	// FS is the fs module in electron main
	FS = js.Global.Call("require", "electron").Get("remote").Call("require", "fs")

	// HomePath is the path to the home is the os
	HomePath = js.Global.Call("require", "electron").Get("remote").Call("require", "os").Call("homedir").String()
)

// Doc is the data storage for a single doc
type Doc struct {
	UUID      string
	Title     string
	Content   string
	UpdatedAt time.Time
}

// Corpus is the data storage for docs
type Corpus struct {
	Docs     map[string]*Doc
	Filename string
}

// UpsertDoc upserts a doc to the corpus
func (c *Corpus) UpsertDoc(
	uuid string,
	title string,
	content string,
) {
	d := &Doc{}
	d.UUID = uuid
	d.Title = title
	d.Content = content
	d.UpdatedAt = time.Now()

	c.Docs[d.UUID] = d
}

// Save saves the whole corpus
func (c *Corpus) Save() {
	b, _ := json.Marshal(c)
	FS.Call("writeFileSync", c.Filename, string(b))
	Console.Call("log", "saved to "+c.Filename)
}

// Load loads the whole corpus
func (c *Corpus) Load() {
	tmpC := &Corpus{}
	s := FS.Call("readFileSync", c.Filename).String()
	json.Unmarshal([]byte(s), tmpC)
	if tmpC.Docs == nil {
		tmpC.Docs = make(map[string]*Doc)
	}
	c.Docs = tmpC.Docs
}

// VueModel is a vue model
type VueModel struct {
	*js.Object
	FullScreenMode bool `js:"fullScreenMode"`
}

// ToggleFullScreenMode toggles the full screen mode
func (vm *VueModel) ToggleFullScreenMode() {
	vm.FullScreenMode = !vm.FullScreenMode
	if vm.FullScreenMode {
		IpcRenderer.Call("send", "ipc_full_screen")
	} else {
		IpcRenderer.Call("send", "ipc_un_full_screen")
	}
}

// App is the app struct
type App struct {
	S  *Slide
	E  *Editor
	VM *VueModel
	C  *Corpus
}

// Bootstrap starts the app
func (a *App) Bootstrap() {
	a.VM = &VueModel{Object: js.Global.Get("Object").New()}
	a.VM.FullScreenMode = false
	vue.New("#app", a.VM)

	a.C = &Corpus{}
	a.C.Filename = HomePath + "/.hast_data"
	a.C.Docs = make(map[string]*Doc)

	a.S = NewSlide()
	a.E = NewEditor()
	a.startSyncEditorToSlides()
	a.startAutoSaving()
}

func (a *App) startSyncEditorToSlides() {
	contentCh := a.E.GetContentCh()
	pageNumCh := a.E.GetPageNumCh()
	go func() {
		for {
			content := <-contentCh
			a.S.SetContent(content)
			a.S.Render()
			a.C.UpsertDoc(
				"uuid1",
				"title_uuid1",
				content,
			)
		}
	}()
	go func() {
		for {
			a.S.GotoPage(<-pageNumCh)
		}
	}()
}

func (a *App) startAutoSaving() {
	go func() {
		for range time.Tick(3 * time.Second) {
			a.C.Save()
		}
	}()
}
