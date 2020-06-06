package db

import "webserver/pkg/models"

type DB interface {
	ConnectDB(host, port, user, password, dbname string) (*Database, error)
	SearchReviews(subject string, name string) ([]*models.Review, error)
	CreateReview(newReview *models.Review) (*models.Review, error)
	UpdateReview(Review *models.Review, ID int) (*models.Review, error)
	DeleteReview(id int)
	GetReview(id int) (*models.Review, bool, error)
	ListReviews(orderType string, pageNumber int) ([]*models.Review, error)
	GetLastInsertReviewID() (int64, error)
}
