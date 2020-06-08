package views

import (
	"io"
	"chicken_review_webserver/pkg/models"
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
