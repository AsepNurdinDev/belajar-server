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

func main() {
	var err error
	db, err = sql.Open("mysql", "asep:281205@tcp(127.0.0.1:3306)/app_db")
	if err != nil {
		log.Fatal(err)
	}

	http.HandleFunc("/webhook", func(w http.ResponseWriter, r *http.Request) {
		if r.Method == "POST" {
			go func() {
				exec.Command("bash", "deploy.sh").Run()
			}()
			w.Write([]byte("Deploy triggered"))
		}
	})

	http.Handle("/static/", http.StripPrefix("/static/", http.FileServer(http.Dir("static"))))

	http.HandleFunc("/", homeHandler)
	http.HandleFunc("/login", loginHandler)
	http.HandleFunc("/register", registerHandler)
	http.HandleFunc("/dashboard", dashboardHandler)

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
	tmpl, _ := template.ParseFiles("templates/dashboard.html")
	tmpl.Execute(w, nil)
}
