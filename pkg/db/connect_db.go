package db

import (
	"chicken-review/webserver/pkg/models"
	"database/sql"
	"fmt"
	"log"

	_ "github.com/lib/pq"
)

type Database struct {
	conn *sql.DB
}

func ConnectDB(host string, port int, user, password, dbname string) (*Database, error) {
	psqInfo := fmt.Sprintf("host=%s port=%d user=%s password=%s dbname=%s sslmode=disable",
		host, port, user, password, dbname)

	db, err := sql.Open("postgres", psqInfo)
	if err != nil {
		return nil, fmt.Errorf("connecting to db: %v", err)
	}
	return &Database{conn: db}, nil
}

func (db *Database) SearchReviews(name string, subject string) ([]*models.Review, error) {

	const q = `
SELECT 
	ID,
	Title,
	Date,
	Author
FROM 
	review
WHERE
	$1 IN ($2)
`
	var reviews []*models.Review
	reviews = nil

	switch subject {
	case "Title", "Date", "Author":
		// ok
	default:
		log.Printf("It is not subject type!")
		return reviews, fmt.Errorf("What?")
	}

	var (
		ID     int64
		Title  string
		Date   string
		Author string
	)

	rows, err := db.conn.Query(q, subject, name)
	if err != nil {
		return reviews, fmt.Errorf("querying: %v", err)
	}
	defer rows.Close()

	for rows.Next() {
		err := rows.Scan(&ID, &Title, &Date, &Author)
		if err != nil {
			return nil, fmt.Errorf("scanning rows: %v", err)
		}

		review := &models.Review{
			ID:     ID,
			Title:  Title,
			Date:   Date,
			Author: Author,
		}
		reviews = append(reviews, review)
	}
	err = rows.Err()
	if err != nil {
		return nil, fmt.Errorf("iterating over rows: %v", err)
	}
	return reviews, nil
}

func (db *Database) GetLastInsertReviewID() (int64, error) {
	const q = `
SELECT
	ID
FROM 
	review
WHERE
	ID = (SELECT MAX (ID) FROM review)
`

	var id int

	err := db.conn.QueryRow(q).Scan(
		&id)
	if err == sql.ErrNoRows {
		log.Printf("There are no rows", err)
		return 1, nil
	}
	if err != nil {
		return 1, fmt.Errorf("Scanning :%v", err)
	}
	//log.Printf("id :%v", id)
	ID64 := int64(id + 1)

	return ID64, nil
}

func (db *Database) CreateReview(newReview *models.Review) (*models.Review, error) {

	q := `
INSERT INTO review (
	Title,
	Author,
	DefaultPictureURL,
	StoreName,
	Date,
	PhoneNumber,
	Comment,
	Score
) VALUES ($1, $2, $3, $4, $5, $6, $7, $8)`

	res, err := db.conn.Exec(q,
		newReview.Title,
		newReview.Author,
		newReview.DefaultPictureURL,
		newReview.StoreName,
		newReview.Date,
		newReview.PhoneNumber,
		newReview.Comment,
		newReview.Score)
	if err != nil {
		return nil, fmt.Errorf("inserting: %v", err)
	}
	id, err := res.LastInsertId()
	if err != nil {
		return nil, fmt.Errorf("reading ID: %v", err)
	}
	newReview.ID = id
	return newReview, nil
}

func (db *Database) UpdateReview(updateReview *models.Review, id int) (*models.Review, error) {
	const q = `
UPDATE review
SET 
	Title = $2,
	Author = $3,
	DefaultPictureURL = $4,
	StoreName = $5,
	Date = $6,
	PhoneNumber = $7,
	Comment = $8,
	Score = $9 
WHERE ID = $1
`

	_, err := db.conn.Exec(q, id, updateReview.Title, updateReview.Author, updateReview.DefaultPictureURL, updateReview.StoreName, updateReview.Date, updateReview.PhoneNumber, updateReview.Comment, updateReview.Score)
	if err != nil {
		return updateReview, fmt.Errorf("updating: %v", err)
	}

	return updateReview, nil
}

func (db *Database) DeleteReview(id int) error {
	const q = `
DELETE 
FROM review	
WHERE
	ID=$1
	`

	_, err := db.conn.Exec(q, id)
	if err != nil {
		return fmt.Errorf("deleting: %v", err)
	}

	return nil

}

func (db *Database) GetReview(id int) (*models.Review, bool, error) {
	const q = `
SELECT
	ID,
	Title,
	Author,
	DefaultPictureURL, 
	StoreName, 
	Date, 
	PhoneNumber, 
	Comment, 
	Score
FROM review
WHERE
	ID=$1
	`

	review := new(models.Review)

	err := db.conn.QueryRow(q, id).Scan(
		&review.ID,
		&review.Title,
		&review.Author,
		&review.DefaultPictureURL,
		&review.StoreName,
		&review.Date,
		&review.PhoneNumber,
		&review.Comment,
		&review.Score)
	if err == sql.ErrNoRows {
		return nil, false, nil
	}
	if err != nil {
		return nil, false, fmt.Errorf("querying: %v", err)
	}

	return review, true, nil
}

func (db *Database) ListReviews(orderType string, pagenumber int) ([]*models.Review, error) {

	const q = `
SELECT 
	ID,
	Title, 
	Date, 
	Author 
FROM review
ORDER BY ID ASC
	`
	var (
		ID     int64
		Title  string
		Date   string
		Author string
	)

	var allReviews []*models.Review
	allReviews = nil
	rows, err := db.conn.Query(q)

	if err != nil {
		return allReviews, fmt.Errorf("There are no reviews : %v", err)
	}

	for rows.Next() {
		err := rows.Scan(&ID, &Title, &Date, &Author)
		if err != nil {
			return nil, fmt.Errorf("rows err : %v", err)
		}

		Review := &models.Review{
			ID:     ID,
			Title:  Title,
			Date:   Date,
			Author: Author,
		}

		allReviews = append(allReviews, Review)
	}
	err = rows.Err()
	if err != nil {
		log.Printf("nothing!")
	}

	return allReviews, nil
}
