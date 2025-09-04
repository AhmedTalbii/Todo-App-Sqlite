package main

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"strconv"
	"strings"

	_ "github.com/mattn/go-sqlite3"
)

type Todo struct {
	ID          int    `json:"id"`
	Title       string `json:"Title"`
	Description string `json:"Description"`
}

var DB *sql.DB

func initial() {
	DN, err := sql.Open("sqlite3", "./database.db")
	if err != nil {
		log.Fatal("Failed to connect to DB")
	}
	sqlmat := `CREATE TABLE IF NOT EXISTS TODO (
		id INTEGER PRIMARY KEY AUTOINCREMENT,
		title TEXT,
		description TEXT
	);`
	DN.Exec(sqlmat)
	DB = DN
}

func main() {
	initial()

	http.HandleFunc("/Posts", func(w http.ResponseWriter, r *http.Request) {
		// CORS headers
		w.Header().Set("Access-Control-Allow-Origin", "*")
		w.Header().Set("Access-Control-Allow-Methods", "POST, GET, DELETE, OPTIONS")
		w.Header().Set("Access-Control-Allow-Headers", "Content-Type")

		if r.Method == http.MethodOptions {
			w.WriteHeader(http.StatusOK)
			return
		}

		if r.Method == http.MethodPost {
			body, _ := io.ReadAll(r.Body)
			todo := &Todo{}
			json.Unmarshal(body, todo)

			_, err := DB.Exec(`INSERT INTO TODO (title, description) VALUES (?, ?)`, todo.Title, todo.Description)
			if err != nil {
				http.Error(w, "Insert error", http.StatusInternalServerError)
				return
			}
			w.Write([]byte(`{"status":"ok"}`))
			return
		}

		if r.Method == http.MethodGet {
			rows, err := DB.Query("SELECT id, title, description FROM TODO")
			if err != nil {
				http.Error(w, "Query error", http.StatusInternalServerError)
				return
			}
			defer rows.Close()

			var data []Todo
			for rows.Next() {
				var t Todo
				if err := rows.Scan(&t.ID, &t.Title, &t.Description); err != nil {
					http.Error(w, "Scan error", http.StatusInternalServerError)
					return
				}
				data = append(data, t)
			}

			f, _ := json.Marshal(data)
			w.Write(f)
			return
		}

		// else => if not matched
		http.Error(w, "Method Not Allowed", http.StatusMethodNotAllowed)
	})

	// DELETE endpoint like /Posts/3
	http.HandleFunc("/Posts/", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Access-Control-Allow-Origin", "*")
		w.Header().Set("Access-Control-Allow-Methods", "DELETE, OPTIONS")
		w.Header().Set("Access-Control-Allow-Headers", "Content-Type")

		if r.Method == http.MethodOptions {
			w.WriteHeader(http.StatusOK)
			return
		}

		if r.Method != http.MethodDelete {
			http.Error(w, "Method Not Allowed", http.StatusMethodNotAllowed)
			return
		}

		parts := strings.Split(r.URL.Path, "/")
		fmt.Println(parts)
		if len(parts) < 3 {
			http.Error(w, "ID is required", http.StatusBadRequest)
			return
		}

		id, err := strconv.Atoi(parts[2])
		if err != nil {
			http.Error(w, "Invalid ID", http.StatusBadRequest)
			return
		}

		_, err = DB.Exec("DELETE FROM TODO WHERE id = ?", id)
		if err != nil {
			http.Error(w, "Failed to delete", http.StatusInternalServerError)
			return
		}

		w.WriteHeader(http.StatusNoContent)
	})

	log.Fatal(http.ListenAndServe(":8080", nil))
}
