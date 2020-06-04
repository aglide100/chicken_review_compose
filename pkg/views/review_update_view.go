package views

import (
	"chicken-review/webserver/pkg/models"
	"io"
)

type reviewUpdateView struct {
	htmlctx BaseHTMLContext
	review  *models.Review
}

func NewReviewUpdateView(htmlctx BaseHTMLContext, review *models.Review) View {
	return &reviewUpdateView{htmlctx: htmlctx, review: review}
}

func (view reviewUpdateView) ContentType() string {
	return "text/html"
}

func (view reviewUpdateView) Render(w io.Writer) error {

	return view.htmlctx.RenderUsing(w, "ui/reviews/update.gohtml", view.review)
}
