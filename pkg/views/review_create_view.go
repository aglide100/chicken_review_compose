package views

import (
	"io"
)

type reviewCreateView struct {
	htmlctx BaseHTMLContext
}

func NewReviewCreateView(htmlctx BaseHTMLContext) View {
	return &reviewCreateView{htmlctx: htmlctx}
}

func (view reviewCreateView) ContentType() string {
	return "text/html"
}

func (view reviewCreateView) Render(w io.Writer) error {
	return view.htmlctx.RenderUsing(w, "ui/reviews/create.gohtml", nil)

}
