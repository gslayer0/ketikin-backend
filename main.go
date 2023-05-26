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
	Id      string
}

func main() {

	db, err := sql.Open("mysql", "root:secret@(127.0.0.1:3306)/ketikin")
	if err != nil {
		panic(err.Error())
	}

	defer db.Close()

	err = db.Ping()
	if err != nil {
		panic(err.Error())
	}

	r := mux.NewRouter()
	tmpl := template.Must(template.ParseFiles("public/index.html"))
	r.HandleFunc("/{title}", func(w http.ResponseWriter, r *http.Request) {

		switch {
		case r.Method == "GET":
			req := mux.Vars(r)
			noteId := req["title"]
			var (
				id   string
				teks string
			)
			query := `SELECT id, teks FROM note WHERE id = ?`
			err := db.QueryRow(query, noteId).Scan(&id, &teks)
			switch {
			case err == sql.ErrNoRows:
				_, err := db.Exec(`INSERT INTO note (id, teks) VALUES(?,?)`, noteId, "")
				if err != nil {
					panic(err.Error())
				}
				data := PageData{
					Writing: "",
					Id:      noteId,
				}
				tmpl.Execute(w, data)

			case err != nil:
				panic(err.Error())
			default:
				data := PageData{
					Writing: teks,
					Id:      id,
				}
				tmpl.Execute(w, data)
			}
		case r.Method == "POST":
			writing := r.FormValue("teks")
			id := r.FormValue("id")
			_, err := db.Exec(`UPDATE note set teks = ? where id = ?`, writing, id)
			if err != nil {
				panic(err.Error())
			}
			data := PageData{
				Writing: writing,
				Id:      id,
			}
			tmpl.Execute(w, data)
		}
	})
	err = http.ListenAndServe(":8080", r)
	if err != nil {
		fmt.Println(err.Error())
	}

}
