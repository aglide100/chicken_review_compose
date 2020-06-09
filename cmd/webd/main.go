package main

import (
	"fmt"
	"log"
	"net"
	"net/http"
	"os"
	"strconv"

	"github.com/aglide100/chicken_review_webserver/pkg/controllers"
	"github.com/aglide100/chicken_review_webserver/pkg/router"

	"github.com/aglide100/chicken_review_webserver/pkg/db"
)

/*

func main() {

	if err := realMain(); err != nil {
		fmt.Errorf("%v", err)
	}
}
*/
func main() {
	log.Printf("start realMain")

	listenAddr := os.Getenv("LISTEN_ADDR")
	listenPort := os.Getenv("LISTEN_PORT")
	dbAddr := os.Getenv("DB_ADDR")
	dbPort := os.Getenv("DB_PORT")
	dbUser := os.Getenv("DB_USER")
	dbPassword := os.Getenv("DB_PASSWORD")
	dbName := os.Getenv("DB_NAME")

	addr := net.JoinHostPort(listenAddr, listenPort)

	dbport, _ := strconv.Atoi(dbPort)
	myDB, err := db.ConnectDB(dbAddr, dbport, dbUser, dbPassword, dbName)
	if err != nil {
		fmt.Errorf("connecting to DB: %v", err)
	}

	defaultCtrl := &controllers.DefaultController{}
	notFoundCtrl := &controllers.NotFoundController{}
	reviewsCtrl := controllers.NewReviewController(myDB)

	rtr := router.NewRouter(notFoundCtrl)

	rtr.AddRule("default", "GET", "^/$", defaultCtrl.ServeHTTP)
	rtr.AddRule("reviews", "GET", "/login$", reviewsCtrl.Login)

	rtr.AddRule("reviews", "GET", "^/reviews/?$", reviewsCtrl.List)
	rtr.AddRule("reviews", "GET", "^reviews/([A-Z]{1,3	})-pagenumber=([0-9]+)$", reviewsCtrl.List)
	rtr.AddRule("reviews", "GET", "^/reviews/([0-9]+)$", reviewsCtrl.Get)

	rtr.AddRule("reviews", "GET", "^/reviews/create$", reviewsCtrl.Create)
	rtr.AddRule("reviews", "POST", "^/reviews/create/upload", reviewsCtrl.Save)

	rtr.AddRule("reviews", "GET", "^/update/([0-9]+)$", reviewsCtrl.Revise)
	rtr.AddRule("reviews", "POST", "^/reviews/update/upload/", reviewsCtrl.Update)

	rtr.AddRule("reviews", "GET", "^/delete/([0-9]+)$", reviewsCtrl.Delete)

	rtr.AddRule("reviews", "GET", "^/reviews/search/", reviewsCtrl.Search)

	rtr.AddRule("reviews", "GET", "^/img", reviewsCtrl.GetImage)

	//rtr.AddRule("reviews", "GET", "^/reviews/ui/img/([0-9]+)/[a-z0-9A-Z_+.-.\\s.-]+.(?i)(img|jpg|jpeg|png|gif)$", reviewsCtrl.GetImage)
	rtr.AddRule("reviews", "GET", "^/reviews/ui/img/[a-z0-9A-Z_+.-.\\s.-]+.(?i)(img|jpg|jpeg|png|gif)$", reviewsCtrl.GetImage)

	log.Println("tcp listen start addr: %v", addr)
	ln, err := net.Listen("tcp", addr)
	log.Println("declare listener")
	if err != nil {
		fmt.Errorf("creating network listener: %v", err)
	}
	defer ln.Close()

	srv := http.Server{Handler: rtr}
	log.Printf("listening on address %q", ln.Addr().String())

	err = srv.Serve(ln)
	log.Printf("starting server at address %q", ln.Addr().String())
	if err != nil {
		fmt.Errorf("serving: %v", err)
	}

}
