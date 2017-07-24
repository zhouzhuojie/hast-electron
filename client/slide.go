package main

import "github.com/gopherjs/gopherjs/js"

const (
	containerID = "h-slides-wrapper"
)

// Slide is the slide object
type Slide struct {
	mdContent string
	container *js.Object
	slideshow *js.Object
}

// NewSlide creates a new Slide
func NewSlide() *Slide {
	s := &Slide{
		container: Document.Call("getElementById", containerID),
	}
	s.Render()
	return s
}

// SetContent sets the mdContent
func (s *Slide) SetContent(c string) {
	s.mdContent = c
}

func (s *Slide) refreshContainer() {
	parent := s.container.Get("parentNode")
	s.container.Call("remove")
	c := Document.Call("createElement", "div")
	c.Set("id", containerID)
	parent.Call("appendChild", c)
	s.container = c
}

// Render renders Slide
func (s *Slide) Render() {
	s.refreshContainer()
	s.renderMarkdown()
}

// GetTitle gets the current Title using h1 tag
func (s *Slide) GetTitle() string {
	h1s := Document.Call("querySelectorAll", ".remark-slide-content h1")
	if h1s.Length() == 0 {
		return ""
	}
	return h1s.Index(0).Get("textContent").String()
}

// GotoPage sets the slide to the specific page number
func (s *Slide) GotoPage(pageNum int) {
	s.slideshow.Call("gotoSlide", pageNum)
}

func (s *Slide) renderMath() {
	js.Global.Call("renderMathInElement", s.container, js.M{
		"delimiters": js.S{
			js.M{"left": "$$", "right": "$$", "display": true},
			js.M{"left": "$", "right": "$", "display": false},
			js.M{"left": "\\[", "right": "\\]", "display": true},
			js.M{"left": "\\(", "right": "\\)", "display": false},
		},
	})
}

func (s *Slide) renderMarkdown() {
	option := js.M{
		"container":      s.container,
		"source":         s.mdContent,
		"highlightStyle": "monokai",
	}
	s.slideshow = Remark.Call("create", option, s.renderMath)
	js.Global.Set("s", s.slideshow)
}
