package views

import (
	"io"
)

type reviewLoginView struct {
	htmlctx BaseHTMLContext
}

func NewReviewLoginView(htmlctx BaseHTMLContext) View {
	return &reviewLoginView{htmlctx: htmlctx}
}

func (view reviewLoginView) ContentType() string {
	return "text/html"
}

func (view reviewLoginView) Render(w io.Writer) error {
	return view.htmlctx.RenderUsing(w, "ui/reviews/login.gohtml", nil)
}
