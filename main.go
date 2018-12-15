package main

import (
	"database/sql"
	"log"
	"net/http"
	"os"

	_ "github.com/go-sql-driver/mysql"
)

var db *sql.DB
var stmt *sql.Stmt

func hakaruHandler(w http.ResponseWriter, r *http.Request) {
	stmt, e := db.Prepare("INSERT INTO eventlog(at, name, value) values(NOW(), ?, ?)")
	if e != nil {
		panic(e.Error())
	}

	defer stmt.Close()

	name := r.URL.Query().Get("name")
	value := r.URL.Query().Get("value")

	_, _ = stmt.Exec(name, value)

	w.WriteHeader(200)
}

func logHandler(handler http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		log.Printf("%s %s %s\n", r.RemoteAddr, r.Method, r.URL)
		handler.ServeHTTP(w, r)
	})
}

func main() {
	dataSourceName := os.Getenv("HAKARU_DATASOURCENAME")
	if dataSourceName == "" {
		dataSourceName = "root:password@tcp(127.0.0.1:13306)/hakaru-db"
	}
	db, err := sql.Open("mysql", dataSourceName)
	if err != nil {
		panic(err.Error())
	}
	defer db.Close()

	http.HandleFunc("/hakaru", hakaruHandler)
	http.HandleFunc("/ok", func(w http.ResponseWriter, r *http.Request) { w.WriteHeader(200) })

	http.ListenAndServe(":8080", logHandler(http.DefaultServeMux))
}
