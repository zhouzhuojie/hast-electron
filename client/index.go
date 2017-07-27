package main

import (
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

	// NeDB is a db for doc storage
	NeDB = js.Global.Call("require", "electron").Get("remote").Call("require", "nedb")

	// Moment js
	Moment = js.Global.Get("moment")

	// Elasticlunr js
	Elasticlunr = js.Global.Get("elasticlunr")
)

// NewObject creates a new js.Object
func NewObject() *js.Object {
	return js.Global.Get("Object").New()
}

// Doc is the data storage for a single doc
type Doc struct {
	*js.Object
	ID              int64  `js:"_id"`
	Title           string `js:"title"`
	Content         string `js:"content"`
	UpdatedAt       int64  `js:"updated_at"`
	HighlighCurrent bool   `js:"highlight_current"`
	TimeAgo         string `js:"time_ago"`
}

// Eq is the equal function of Docs
func (d *Doc) Eq(that *Doc) bool {
	return d.ID == that.ID && d.Title == that.Title && d.Content == that.Content && d.UpdatedAt == that.UpdatedAt
}

// NewDoc news a doc
func NewDoc(id int64) *Doc {
	doc := &Doc{Object: NewObject()}
	doc.Content = "# Hello World!"
	doc.UpdatedAt = time.Now().UnixNano()
	if id == 0 {
		doc.ID = time.Now().UnixNano()
	} else {
		doc.ID = id
	}
	return doc
}

// Corpus is the data storage for docs
type Corpus struct {
	DB         *js.Object
	CurrentDoc *Doc
	Index      *js.Object
}

// NewCorpus creates a new Corpus
func NewCorpus(filename string) *Corpus {
	c := &Corpus{
		DB: NeDB.New(js.M{
			"filename": filename,
			"autoload": true,
		}),
	}
	allDocs := c.GetAll()
	if len(allDocs) == 0 {
		c.CurrentDoc = NewDoc(0)
	} else {
		c.CurrentDoc = allDocs[0]
	}
	return c
}

// UpsertDoc upserts a doc to the corpus
func (c *Corpus) UpsertDoc(
	id int64,
	title string,
	content string,
) {
	d := NewDoc(id)
	d.Title = title
	d.Content = content

	c.CurrentDoc = d
	c.DB.Call("update", js.M{"_id": id}, d, js.M{"upsert": true})
}

// Reindex reindex all the docs
func (c *Corpus) Reindex(docs []*Doc) {
	c.Index = Elasticlunr.Invoke()
	c.Index.Call("addField", "title")
	c.Index.Call("addField", "content")
	c.Index.Call("setRef", "id")
	for _, doc := range docs {
		c.Index.Call("addDoc", js.M{
			"id":      doc.ID,
			"title":   doc.Title,
			"content": doc.Content,
		})
	}
}

// Search searches the corpus
func (c *Corpus) Search(q string) []int64 {
	Console.Call("log", q)
	ids := make([]int64, 0)
	results := c.Index.Call("search", q)
	for i := 0; i < results.Length(); i++ {
		ids = append(ids, results.Index(i).Get("ref").Int64())
	}
	Console.Call("log", ids)
	return ids
}

// GetAll gets all the docs
func (c *Corpus) GetAll() []*Doc {
	ch := make(chan []*Doc)
	exec := c.DB.Call("find", js.M{}).Call("sort", js.M{"_id": 1})
	exec.Call("exec", func(err *js.Object, data *js.Object) {
		n := data.Length()
		docs := make([]*Doc, 0)
		for i := 0; i < n; i++ {
			d := data.Index(i)
			doc := &Doc{Object: d}
			doc.TimeAgo = Moment.Invoke(doc.UpdatedAt / 1e6).Call("fromNow").String()
			docs = append(docs, doc)
		}
		ch <- docs
	})
	docs := <-ch
	c.Reindex(docs)
	return docs
}

// App is the app struct
type App struct {
	*js.Object

	S  *Slide
	E  *Editor
	C  *Corpus
	VM *vue.ViewModel

	RefreshDocsFunc *js.Object

	// Vue data binding
	FullScreenMode bool   `js:"fullScreenMode"`
	Docs           []*Doc `js:"allDocs"`
	SearchStr      string `js:"searchStr"`
}

// ToggleFullScreenMode toggles the full screen mode
func (a *App) ToggleFullScreenMode() {
	a.FullScreenMode = !a.FullScreenMode
	if a.FullScreenMode {
		IpcRenderer.Call("send", "ipc_full_screen")
	} else {
		IpcRenderer.Call("send", "ipc_un_full_screen")
	}
}

// Bootstrap starts the app
func (a *App) Bootstrap() {
	a.Object = NewObject()
	a.FullScreenMode = false
	a.Docs = make([]*Doc, 0)
	a.SearchStr = ""
	a.VM = vue.New("#app", a)

	a.C = NewCorpus(HomePath + "/.hast_data")
	a.S = NewSlide()
	a.E = NewEditor()
	a.E.SetValue(a.C.CurrentDoc.Content)

	// init some debounce functions
	a.RefreshDocsFunc = Lodash.Call("throttle", func() {
		go func() {
			a.Docs = a.C.GetAll()
			a.SetCurrentDoc(a.C.CurrentDoc.ID)
		}()
	}, 1)
	a.RefreshDocsFunc.Invoke()

	a.startSyncEditorToSlides()
}

// CreateDoc creates a new page
func (a *App) CreateDoc() {
	a.C.CurrentDoc = NewDoc(0)
	a.E.SetValue(a.C.CurrentDoc.Content)
}

// SetCurrentDoc sets the current doc
func (a *App) SetCurrentDoc(id int64) {
	for _, doc := range a.Docs {
		doc.HighlighCurrent = false
		if doc.ID == id {
			a.C.CurrentDoc = doc
			a.C.CurrentDoc.HighlighCurrent = true
		}
	}
}

// SetCurrentDocAndReload sets the current doc
func (a *App) SetCurrentDocAndReload(id int64) {
	a.SetCurrentDoc(id)
	a.E.SetValue(a.C.CurrentDoc.Content)
}

// FilterSearchResult filters the docs based on query
func (a *App) FilterSearchResult() {
	go func() {
		query := a.SearchStr
		a.Docs = a.C.GetAll()
		if len(query) <= 1 {
			return
		}
		ids := a.C.Search(query)
		docs := make([]*Doc, len(ids))
		for i, id := range ids {
			for _, doc := range a.Docs {
				if doc.ID == id {
					docs[i] = doc
				}
			}
		}
		Console.Call("log", docs)
		a.Docs = docs
	}()
}

func (a *App) startSyncEditorToSlides() {
	contentCh := a.E.GetContentCh()
	go func() {
		for {
			content := <-contentCh
			a.S.SetContent(content)
			a.S.Render()
			a.C.UpsertDoc(a.C.CurrentDoc.ID, a.S.GetTitle(), content)
			a.RefreshDocsFunc.Invoke()
		}
	}()

	pageNumCh := a.E.GetPageNumCh()
	go func() {
		for {
			a.S.GotoPage(<-pageNumCh)
		}
	}()
}
