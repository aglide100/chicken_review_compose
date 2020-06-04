package views

import (
	"io"
	"path/filepath"
)

type reviewGetImageView struct {
	htmlctx BaseHTMLContext
	path    string
	format  string
}

func NewReviewGetImageView(htmlctx BaseHTMLContext, path string, format string) View {
	return &reviewGetImageView{htmlctx: htmlctx, path: path, format: format}
}

func (view reviewGetImageView) ContentType() string {
	var contentType string
	var ext string
	switch view.format {
	case "ReviewImage":
		view.path = view.path[9:]
		ext = filepath.Ext(view.path)
	case "":

	}

	switch ext {
	case ".png":
		contentType = "image/png"
	case ".jpg":
		contentType = "image/jpeg"
	case ".gif":
		contentType = "image/gif"
	}

	return contentType
}

func (view reviewGetImageView) Render(w io.Writer) error {
	return view.htmlctx.RenderImage(w, view.path)
}
