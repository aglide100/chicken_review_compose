package views

import (
	"chicken-review/webserver/ui"
	"fmt"
	"io"
	"io/ioutil"
	"text/template"
)

var DefaultBaseHTMLContext = BaseHTMLContext{
	GlobPattern: "ui/*.gohtml",
	HTML: func(bodyContent interface{}) ui.HTML {
		return ui.HTML{
			Head: ui.Head{
				FavIcoURL:   "",
				Title:       "Chicken Review",
				Author:      "",
				Description: "We review chicken restaurants",
			},
			Body: ui.Body{Content: bodyContent},
		}
	},
}

type BaseHTMLContext struct {
	GlobPattern string
	HTML        func(bodyContent interface{}) ui.HTML
}

func (htmlctx *BaseHTMLContext) RenderImage(w io.Writer, path string) error {

	content, err := ioutil.ReadFile(path[9:])
	if err != nil {
		return nil
	}

	w.Write(content)

	return nil
}

func (htmlctx *BaseHTMLContext) RenderUsing(w io.Writer, contentPattern string, bodyContent interface{}) error {
	baseT, err := template.ParseGlob(htmlctx.GlobPattern)
	if err != nil {
		return fmt.Errorf("parsing base html: %v", err)
	}
	contentT, err := template.Must(baseT.Clone()).ParseGlob(contentPattern)
	if err != nil {
		return fmt.Errorf("parsing reviews html: %v", err)
	}

	html := htmlctx.HTML(bodyContent)
	if err := contentT.ExecuteTemplate(w, "html", html); err != nil {
		return fmt.Errorf("executing template: %v", err)
	}
	return nil
}
