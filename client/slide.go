package main

import "github.com/gopherjs/gopherjs/js"

const (
	containerID = "h-slides-wrapper"
)

// Slide is the slide object
type Slide struct {
	mdContent string
	container *js.Object
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
		"container": s.container,
		"source":    s.mdContent,
	}
	Remark.Call("create", option, s.renderMath)
}
