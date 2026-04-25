package main

import (
	"database/sql"
	"fmt"
	"html/template"
	"log"
	"net/http"
	"os/exec"

	_ "github.com/go-sql-driver/mysql"
)

var db *sql.DB

type User struct {
	ID       int
	Username string
	Email    string
}

func main() {
	var err error
	db, err = sql.Open("mysql", "asep:281205@tcp(127.0.0.1:3306)/app_db")
	if err != nil {
		log.Fatal(err)
	}

	// webhook CI/CD
	http.HandleFunc("/webhook", func(w http.ResponseWriter, r *http.Request) {
		if r.Method == "POST" {
			go func() {
				exec.Command("bash", "deploy.sh").Run()
			}()
			w.Write([]byte("Deploy triggered"))
		}
	})

	// static file
	http.Handle("/public/", http.StripPrefix("/public/", http.FileServer(http.Dir("public"))))

	// routes
	http.HandleFunc("/", homeHandler)
	http.HandleFunc("/login", loginHandler)
	http.HandleFunc("/register", registerHandler)
	http.HandleFunc("/dashboard", dashboardHandler)
	http.HandleFunc("/admin", adminHandler)
	http.HandleFunc("/add", addUserHandler)
	http.HandleFunc("/delete", deleteUserHandler)

	fmt.Println("Server running on :8080")
	http.ListenAndServe(":8080", nil)
}

func homeHandler(w http.ResponseWriter, r *http.Request) {
	tmpl, _ := template.ParseFiles("templates/index.html")
	tmpl.Execute(w, nil)
}

func loginHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method == "POST" {
		email := r.FormValue("email")
		password := r.FormValue("password")

		var id int
		err := db.QueryRow("SELECT id FROM users WHERE email=? AND password=?", email, password).Scan(&id)

		if err == nil {
			http.Redirect(w, r, "/dashboard", http.StatusSeeOther)
			return
		}
	}

	tmpl, _ := template.ParseFiles("templates/login.html")
	tmpl.Execute(w, nil)
}

func registerHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method == "POST" {
		username := r.FormValue("username")
		email := r.FormValue("email")
		password := r.FormValue("password")

		_, err := db.Exec("INSERT INTO users(username, email, password) VALUES(?,?,?)",
			username, email, password)

		if err == nil {
			http.Redirect(w, r, "/login", http.StatusSeeOther)
			return
		}
	}

	tmpl, _ := template.ParseFiles("templates/register.html")
	tmpl.Execute(w, nil)
}

func dashboardHandler(w http.ResponseWriter, r *http.Request) {
	rows, err := db.Query("SELECT id, username, email FROM users")
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	defer rows.Close()

	var users []User

	for rows.Next() {
		var u User
		rows.Scan(&u.ID, &u.Username, &u.Email)
		users = append(users, u)
	}

	tmpl, _ := template.ParseFiles("templates/dashboard.html")
	tmpl.Execute(w, users)
}

func adminHandler(w http.ResponseWriter, r *http.Request) {
	rows, err := db.Query("SELECT id, username, email FROM users")
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	defer rows.Close()

	var users []User

	for rows.Next() {
		var u User
		rows.Scan(&u.ID, &u.Username, &u.Email)
		users = append(users, u)
	}

	tmpl, _ := template.ParseFiles("templates/admin.html")
	tmpl.Execute(w, users)
}

func addUserHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method == "POST" {
		username := r.FormValue("username")
		email := r.FormValue("email")
		password := r.FormValue("password")

		db.Exec("INSERT INTO users(username,email,password) VALUES(?,?,?)",
			username, email, password)
	}

	http.Redirect(w, r, "/dashboard", http.StatusSeeOther)
}

func deleteUserHandler(w http.ResponseWriter, r *http.Request) {
	id := r.URL.Query().Get("id")

	db.Exec("DELETE FROM users WHERE id=?", id)

	http.Redirect(w, r, "/dashboard", http.StatusSeeOther)
}