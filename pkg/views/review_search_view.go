package views

import (
	"chicken_review_webserver/pkg/models"
	"io"
)

type reviewSearchView struct {
	htmlctx BaseHTMLContext
	reviews []*models.Review
}

func NewReviewSearchView(htmlctx BaseHTMLContext, reviews []*models.Review) View {
	return &reviewSearchView{htmlctx: htmlctx, reviews: reviews}
}

func (view reviewSearchView) ContentType() string {
	return "text/html"
}

func (view reviewSearchView) Render(w io.Writer) error {
	return view.htmlctx.RenderUsing(w, "ui/reviews/search.gohtml", view.reviews)
}
