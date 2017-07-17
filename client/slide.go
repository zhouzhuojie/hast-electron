package main

import "github.com/gopherjs/gopherjs/js"

// Slide is the slide object
type Slide struct {
	mdContent string
	container *js.Object
}

// NewSlide creates a new Slide
func NewSlide() *Slide {
	s := &Slide{
		container: Document.Call("getElementById", "h-slides-wrapper"),
	}
	s.Render()
	return s
}

// SetContent sets the mdContent
func (s *Slide) SetContent(c string) {
	s.mdContent = c
}

// Render renders Slide
func (s *Slide) Render() {
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
