package views

import (
	"chicken-review/webserver/pkg/models"
	"io"
)

type reviewDeleteView struct {
	htmlctx BaseHTMLContext
	review  *models.Review
}

func NewReviewDeleteView(htmlctx BaseHTMLContext, review *models.Review) View {
	return &reviewDeleteView{htmlctx: htmlctx, review: review}
}

func (view reviewDeleteView) ContentType() string {
	return "text/html"
}

func (view reviewDeleteView) Render(w io.Writer) error {
	return view.htmlctx.RenderUsing(w, "ui/reviews/delete.gohtml", view.review)
}
