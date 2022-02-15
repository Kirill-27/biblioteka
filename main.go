package main

import (
	"database/sql"
	"fmt"
	"github.com/go-chi/chi"
	_ "github.com/lib/pq"
	"net/http"
)

const (
	host     = "localhost"
	port     = 5435
	user     = "dashboard"
	password = "dashboard"
	dbname   = "dashboard"
)

func main() {

	psqlInfo := fmt.Sprintf("host=%s port=%d user=%s "+
		"password=%s dbname=%s sslmode=disable",
		host, port, user, password, dbname)
	fmt.Println(psqlInfo)
	db, err := sql.Open("postgres", psqlInfo)
	if err != nil {
		panic(err)
	}
	defer db.Close()
	r := chi.NewRouter()
	r.Route("/tags", func(r chi.Router) {
		r.Get("/", handlers.GetOfferList)
		r.Get("/{id}", exchangeRatesKeeper.CalculatePrise)
	})

	http.ListenAndServe(":8384", r)
	fmt.Println("Successfully connected!")

}
