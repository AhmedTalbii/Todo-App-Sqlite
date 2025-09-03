package main

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"

	_ "github.com/mattn/go-sqlite3"
)

type Todo struct {
	Title       string `json:"Title"`
	Description string `json:"Description"`
}

var DB *sql.DB

func initial() {
	DN, err := sql.Open("sqlite3", "./database.db")
	if err != nil {
		log.Fatal("zaba")
	}
	sqlmat := `CREATE TABLE TODO (
  id INTEGER PRIMARY KEY AUTOINCREMENT,
  title TEXT,
  description TEXT
);

	`
	DN.Exec(sqlmat)
	DB = DN
}

func main() {
	initial()
	// mux := http.NewServeMux()
	http.HandleFunc("/Posts", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Access-Control-Allow-Origin", "*")
		w.Header().Set("Access-Control-Allow-Methods", "POST, GET, OPTIONS")
		w.Header().Set("Access-Control-Allow-Headers", "Content-Type")

		if r.Method == http.MethodOptions {
			w.WriteHeader(http.StatusOK) // respond OK to preflight
		}
		if r.Method == http.MethodPost {
			body, _ := io.ReadAll(r.Body)
			todo := &Todo{}
			json.Unmarshal(body, todo)
			// fmt.Println(todo)

			_, err := DB.Exec(
				`INSERT INTO TODO (title, description) VALUES (?, ?)`,
				todo.Title, todo.Description,
			)
			if err != nil {
				fmt.Println("insert error:", err)
				return
			}

		}
	if r.Method == http.MethodGet {
    rows, err := DB.Query("SELECT id, title, description FROM TODO")
    if err != nil {
        fmt.Println("query error:", err)
        return
    }
    defer rows.Close()

    for rows.Next() {
        var id int
        var title, description string
        if err := rows.Scan(&id, &title, &description); err != nil {
            fmt.Println("scan error:", err)
            return
        }
        fmt.Println(id, title, description)
    }
}

		w.WriteHeader(http.StatusOK)
		w.Write([]byte(`{"status":"ok"}`))
	})

	log.Fatal(http.ListenAndServe(":8080", nil))
}
