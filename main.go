package main

import (
	"database/sql"
	"fmt"
	"html/template"
	"net/http"

	_ "github.com/go-sql-driver/mysql"

	"github.com/gorilla/mux"
)

type PageData struct {
	Writing string
}

func main() {

	db, err := sql.Open("mysql", "root:my-secret-pw@(127.0.0.1:3306)/ketikin")
	if err != nil {
		panic(err.Error())
	}

	err = db.Ping()
	if err != nil {
		panic(err.Error())
	}

	defer db.Close()

	r := mux.NewRouter()
	tmpl := template.Must(template.ParseFiles("public/index.html"))
	r.HandleFunc("/{title}", func(w http.ResponseWriter, r *http.Request) {

		req := mux.Vars(r)
		noteId := req["title"]

		var (
			id   string
			teks string
		)

		query := `SELECT id, teks FROM note WHERE id = ?`
		err := db.QueryRow(query, noteId).Scan(&id, &teks)
		if err != nil {
			data := PageData{
				Writing: err.Error(),
			}
			tmpl.Execute(w, data)
			return
		}

		data := PageData{
			Writing: teks,
		}
		tmpl.Execute(w, data)
	})
	err = http.ListenAndServe(":8080", r)
	if err != nil {
		fmt.Println(err.Error())
	}

}
