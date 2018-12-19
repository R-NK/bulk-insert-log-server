package main

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"time"

	_ "github.com/go-sql-driver/mysql"
	"github.com/gocraft/dbr"
)

type Record struct {
	At    time.Time `db:"at"`
	Name  string    `db:"name"`
	Value string    `db:"value"`
}

var (
	db        *dbr.Connection
	queueChan = make(chan Record, 1000)
)

func ticker() {
	t := time.NewTicker(1 * time.Second) //1秒周期の ticker
	queue := make([]Record, 0, 1000)

	for {
		select {
		case <-t.C:
			if len(queue) == 0 {
				continue
			}
			sess := db.NewSession(nil)
			stmt := sess.InsertInto("eventlog").Columns("at", "name", "value")
			for _, value := range queue {
				stmt.Record(value)
			}
			result, err := stmt.Exec()
			if err != nil {
				fmt.Println(err)
			} else {
				count, _ := result.RowsAffected()
				fmt.Println("records count: ", count)
				queue = make([]Record, 0, 1000)
			}
		case record := <-queueChan:
			queue = append(queue, record)
		}
	}
}

func hakaruHandler(w http.ResponseWriter, r *http.Request) {
	name := r.URL.Query().Get("name")
	value := r.URL.Query().Get("value")

	record := Record{At: time.Now(), Name: name, Value: value}
	queueChan <- record

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
	db, _ = dbr.Open("mysql", dataSourceName, nil)

	go ticker()

	http.HandleFunc("/hakaru", hakaruHandler)
	http.HandleFunc("/ok", func(w http.ResponseWriter, r *http.Request) { w.WriteHeader(200) })

	http.ListenAndServe(":8081", logHandler(http.DefaultServeMux))
}
