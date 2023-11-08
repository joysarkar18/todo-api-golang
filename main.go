package main

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"net/http"

	_ "github.com/go-sql-driver/mysql"
	"github.com/gorilla/mux"
)

type Todo struct {
	ID        int    `json:"id"`
	Title     string `json:"title"`
	Completed bool   `json:"completed"`
}

var db *sql.DB

func InitDB() {
	dataSourceName := "root:JoySarkar@456@tcp(localhost:3306)/todo_db"
	var err error
	db, err = sql.Open("mysql", dataSourceName)
	if err != nil {
		fmt.Println(err.Error())
		panic(err)
	}

	err = db.Ping()
	if err != nil {
		fmt.Println(err.Error())
		panic(err)
	}
	fmt.Println("Connected to the database")
}

func CreateTodo(w http.ResponseWriter, r *http.Request) {
	var newTodo Todo
	json.NewDecoder(r.Body).Decode(&newTodo)

	_, err := db.Exec("INSERT INTO todos (title, completed) VALUES (?, ?)", newTodo.Title, newTodo.Completed)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusCreated)
}
func GetTodo(w http.ResponseWriter, r *http.Request) {
	rows, err := db.Query("SELECT * FROM todos")
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	defer rows.Close()

	var todos []Todo
	for rows.Next() {
		var todo Todo
		err := rows.Scan(&todo.ID, &todo.Title, &todo.Completed)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		todos = append(todos, todo)
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(todos)
}
func UpdateTodo(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)
	todoID := params["id"]

	var updatedTodo Todo
	json.NewDecoder(r.Body).Decode(&updatedTodo)

	_, err := db.Exec("UPDATE todos SET title=?, completed=? WHERE id=?", updatedTodo.Title, updatedTodo.Completed, todoID)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
}
func DeleteTodo(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)
	todoID := params["id"]

	_, err := db.Exec("DELETE FROM todos WHERE id=?", todoID)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
}
func main() {

	r := mux.NewRouter()

	r.HandleFunc("/todos", CreateTodo).Methods("POST")
	r.HandleFunc("/todos", GetTodo).Methods("GET")
	r.HandleFunc("/todos/{id}", UpdateTodo).Methods("PUT")
	r.HandleFunc("/todos/{id}", DeleteTodo).Methods("DELETE")

	http.Handle("/", r)

	InitDB()

	// Start the server
	http.ListenAndServe(":8080", nil)

}
