package views

import (
	"io"
	"webserver/pkg/models"
)

type reviewListView struct {
	htmlctx BaseHTMLContext
	reviews []*models.Review
}

func NewReviewListView(htmlctx BaseHTMLContext, reviews []*models.Review) View {
	return &reviewListView{htmlctx: htmlctx, reviews: reviews}
}

func (view reviewListView) ContentType() string {
	return "text/html"
}

func (view reviewListView) Render(w io.Writer) error {
	return view.htmlctx.RenderUsing(w, "ui/reviews/list.gohtml", view.reviews)
}
