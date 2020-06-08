package controllers

import (
	"chicken_review_webserver/pkg/db"
	"chicken_review_webserver/pkg/models"
	"chicken_review_webserver/pkg/views"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"regexp"
	"strconv"
	"strings"
)

type ReviewController struct {
	db *db.Database
}

func NewReviewController(db *db.Database) *ReviewController {
	return &ReviewController{db: db}
}

func findString(resp http.ResponseWriter, req *http.Request, str string) (id int, orderType string, pagenumber int) {
	var matches []string
	Pnumber := 0

	var showReviewPattern = regexp.MustCompile("^/reviews/([0-9]+$)")
	var deleteReviewPattern = regexp.MustCompile("^/delete/([0-9]+$)")
	var updateReviewPattern = regexp.MustCompile("^/update/([0-9]+$)")
	var uploadUpdateReviewPattern = regexp.MustCompile("^/reviews/update/upload/([0+9+])$")
	var listReviewPattern = regexp.MustCompile("^reviews/([A-Z]+)-pagenumber=([0-9]+)$")

	switch str {
	case "Show":
		matches = showReviewPattern.FindStringSubmatch(req.URL.Path)
	case "Delete":
		matches = deleteReviewPattern.FindStringSubmatch(req.URL.Path)
	case "Update":
		matches = updateReviewPattern.FindStringSubmatch(req.URL.Path)
	case "UploadUpdate":
		matches = uploadUpdateReviewPattern.FindStringSubmatch(req.URL.Path)
	case "List":
		matches = listReviewPattern.FindStringSubmatch(req.URL.Path)
	}

	if len(matches) != 2 {
		//http.Error(resp, "no ID provided", http.StatusBadRequest)

		if str == "List" {
			//var err error
			/*
				orderType = matches[3]
				pagenumber, err = strconv.Atoi(matches[20]) // 20
				if err != nil {
					log.Printf("PageNumber is not numeric: %v", err)
				}
			*/

			return 0, orderType, pagenumber
		}

		return
	}

	idStr := matches[1]
	id, err := strconv.Atoi(idStr)
	if err != nil {
		//http.Error(resp, fmt.Sprintf("ID is not numeric: %v", err.Error()), http.StatusBadRequest)
		return
	}

	return id, "", Pnumber
}

func (hdl *ReviewController) GetImage(resp http.ResponseWriter, req *http.Request) {
	log.Printf("receive request to get image")

	view := views.NewReviewGetImageView(views.DefaultBaseHTMLContext, req.URL.Path, "ReviewImage") // subdivided to ReviewImage and Favico
	resp.Header().Set("Content-Type", view.ContentType())
	err := view.Render(resp)
	if err != nil {
		log.Printf("failed to render: %v", err)
	}

}

func (hdl *ReviewController) Create(resp http.ResponseWriter, req *http.Request) {
	log.Printf("receive request to create a review")
	view := views.NewReviewCreateView(views.DefaultBaseHTMLContext)

	resp.Header().Set("Content-Type", view.ContentType())
	err := view.Render(resp)
	if err != nil {
		log.Printf("failed to render: %v", err)
	}
}

func (hdl *ReviewController) Revise(resp http.ResponseWriter, req *http.Request) {
	log.Printf("receive request to update a reivew")

	id, _, _ := findString(resp, req, "Update")

	review, ok, err := hdl.db.GetReview(id)
	if err != nil {
		log.Printf("finding review: %v", err)
		http.Error(resp, err.Error(), http.StatusInternalServerError)
		return
	}
	var view views.View
	if !ok {
		view = views.NewNotFoundView(views.DefaultBaseHTMLContext)
	} else {
		view = views.NewReviewUpdateView(views.DefaultBaseHTMLContext, review)
	}

	resp.Header().Set("Content-Type", view.ContentType())
	err = view.Render(resp)
	if err != nil {
		log.Printf("failed to render: %v", err)
	}
}

func (hdl *ReviewController) Delete(resp http.ResponseWriter, req *http.Request) {
	log.Printf("receive request to delete a review")

	id, _, _ := findString(resp, req, "Delete")
	log.Printf("Delete ID :%v", id)
	hdl.db.DeleteReview(id)

	http.Redirect(resp, req, "/reviews", 301)
}

const (
	KiB = 1 << 10
	MiB = 1024 * KiB

	maxImageSize = 20 * MiB
)

func SaveImage(resp http.ResponseWriter, req *http.Request, hdl *ReviewController) (string, bool, error) {
	// It is just save at local, if you can change the other way to save file. you should be changed this code

	if err := req.ParseMultipartForm(10 << 20); err != nil {
		return "", true, fmt.Errorf("parsing multipart form: %v", err)
	}

	file, fileHeader, err := req.FormFile("image")
	if err != nil {
		return "", false, fmt.Errorf("looking up image from form file: %v", err)
	}
	defer file.Close()

	imageBytes, err := ioutil.ReadAll(io.LimitReader(file, maxImageSize))
	if err != nil {
		return "", true, fmt.Errorf("can't read image data: %v", err)
	}

	fileType := http.DetectContentType(imageBytes)
	switch fileType {
	case "image/jpeg", "image/png", "image/gif", "image/webp":
		fileType = strings.Replace(fileType, "image/", ".", 1)
		// ok
	default:
		return "", true, fmt.Errorf("invalid image format: %q", fileType)
	}
	basePath := "ui/img/"

	//reviewid := id

	reviewid, err := hdl.db.GetLastInsertReviewID()
	if err != nil {
		log.Printf("Can't get LastInsertId!", err)
	}

	IDstr := strconv.FormatInt(reviewid, 10)
	log.Printf("reivewID: %v, IDstr: %v", reviewid, IDstr) //
	//currentReviewBasePath := filepath.Join(basePath, IDstr)
	currentReviewBasePath := basePath

	err = os.MkdirAll(currentReviewBasePath, os.ModePerm)
	if err != nil {
		return "", true, fmt.Errorf("creatig directory for image: %v", err)
	}

	imageFilename := fileHeader.Filename
	log.Printf("file name :%v", imageFilename)
	IDstr += fileType
	//currentReviewImagePath := filepath.Join(currentReviewBasePath, IDstr, imageFilename) // Can't create directory
	currentReviewImagePath := filepath.Join(currentReviewBasePath, IDstr)

	if err := ioutil.WriteFile(currentReviewImagePath, imageBytes, 0644); err != nil {
		return "", true, fmt.Errorf("creating image file on disk: %v", err)
	}

	log.Printf("img path: %v", currentReviewImagePath)
	return currentReviewImagePath, true, nil
}

func SaveReview(resp http.ResponseWriter, req *http.Request, hdl *ReviewController, ReviewType string) (*models.Review, error, bool, string) {

	path, ok, err := SaveImage(resp, req, hdl)
	if !ok {
		log.Printf("There are no image!")
	} else if err != nil {
		log.Printf("saving image: %v", err)
		http.Error(resp, "can't save image", http.StatusInternalServerError)
		return nil, err, false, ""
	}

	checklist := [...]string{
		"store_name",
		"date",
		"phone_number",
		"author",
		"title",
		"comment",
	}

	blacklist := [...]string{
		"<",
		">",
		"$",
		"<script>",
		"<style>",
		"$func",
	}

	review := &models.Review{
		StoreName:         req.PostFormValue("store_name"),
		Date:              req.PostFormValue("date"),
		PhoneNumber:       req.PostFormValue("phone_number"),
		Author:            req.PostFormValue("author"),
		Title:             req.PostFormValue("title"),
		DefaultPictureURL: path,
		Comment:           req.PostFormValue("comment"),
		//UpdateDate:        req.PostFormValue("write_date"),
	}

	checklistnum := 6
	blacklistnum := 6

	for i := 0; i < checklistnum; i++ {
		for k := 0; k < blacklistnum; k++ {
			result := strings.Replace(req.PostFormValue(checklist[i]), blacklist[k], "Alert", -1)
			if result != req.PostFormValue(checklist[i]) {
				return nil, err, true, result
			}
		}
	}
	switch ReviewType {
	case "Save":

	case "Update":

	}

	//log.Printf("Title : %v", review.Title)
	return review, nil, false, ""
}

func (hdl *ReviewController) Save(resp http.ResponseWriter, req *http.Request) {
	log.Printf("receive request to save a review")

	review, err, xss, str := SaveReview(resp, req, hdl, "Save")
	if err != nil {
		if xss {
			log.Printf("%v", str)
			http.Redirect(resp, req, "/Don't-Use-Script-or-Css-at-review", 301)
		} else {
			log.Fatal("Can't save review %v:", err)
		}
	}

	// return to default page
	if !xss {
		hdl.db.CreateReview(review)
		http.Redirect(resp, req, "/reviews", 301)
	}
}

func (hdl *ReviewController) Update(resp http.ResponseWriter, req *http.Request) {
	log.Printf("receive request to update a review")

	id, _ := strconv.Atoi(req.PostFormValue("id"))

	review, err, xss, str := SaveReview(resp, req, hdl, "Save")
	if err != nil {
		if xss {

			log.Printf("%v", str)
			http.Redirect(resp, req, "/Don't-Use-Script-or-Css-at-review", 301)
		} else {
			log.Fatal("Can't save review")
		}
	}

	if !xss {
		hdl.db.UpdateReview(review, id)
		http.Redirect(resp, req, "/reviews", 301)
	}
}

func (hdl *ReviewController) Search(resp http.ResponseWriter, req *http.Request) {
	log.Printf("receive request to search review")

	req.ParseForm()

	var (
		name    string
		subject string
	)
	name = req.FormValue("name")
	subject = req.FormValue("subject")

	reviews, err := hdl.db.SearchReviews(name, subject)
	if err != nil {
		log.Printf("listing reviews: %v", err)
		http.Error(resp, err.Error(), http.StatusInternalServerError)
		return
	}

	view := views.NewReviewSearchView(views.DefaultBaseHTMLContext, reviews)
	resp.Header().Set("Content-Type", view.ContentType())
	err = view.Render(resp)
	if err != nil {
		log.Printf("failed to render : %v", err)
	}
}

func (hdl *ReviewController) Get(resp http.ResponseWriter, req *http.Request) {
	log.Printf("receive request to get a review")

	id, _, _ := findString(resp, req, "Show")

	review, ok, err := hdl.db.GetReview(id)
	if err != nil {
		log.Printf("finding review: %v", err)
		http.Error(resp, err.Error(), http.StatusInternalServerError)
		return
	}
	var view views.View
	if !ok {
		view = views.NewNotFoundView(views.DefaultBaseHTMLContext)
	} else {
		view = views.NewReviewShowView(views.DefaultBaseHTMLContext, review)
	}
	resp.Header().Set("Content-Type", view.ContentType())
	err = view.Render(resp)
	if err != nil {
		log.Printf("failed to render: %v", err)
	}
}

func (hdl *ReviewController) List(resp http.ResponseWriter, req *http.Request) {
	log.Printf("receive request to list reviews")

	_, orederType, pagenumber := findString(resp, req, "List")

	reviews, err := hdl.db.ListReviews(orederType, pagenumber)
	if err != nil {
		log.Printf("listing reviews: %v", err)
		http.Error(resp, err.Error(), http.StatusInternalServerError)
		return
	}

	view := views.NewReviewListView(views.DefaultBaseHTMLContext, reviews)
	resp.Header().Set("Content-Type", view.ContentType())
	err = view.Render(resp)
	if err != nil {
		log.Printf("failed to render : %v", err)
	}
}

func (hdl *ReviewController) Login(resp http.ResponseWriter, req *http.Request) {
	log.Printf("recevie request to login view")

	view := views.NewReviewLoginView(views.DefaultBaseHTMLContext)
	resp.Header().Set("Content-Type", view.ContentType())
	err := view.Render(resp)
	if err != nil {
		log.Printf("faild to render : %v", err)
	}
}
