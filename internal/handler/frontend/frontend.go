package frontend

import (
	"fmt"
	"html/template"
	"net/http"

	"github.com/morozoffnor/home-storage/internal/auth"
	"github.com/morozoffnor/home-storage/internal/database"
)

type FrontendHandler struct {
	db   *database.Database
	auth *auth.Auth
}

func New(db *database.Database) *FrontendHandler {
	return &FrontendHandler{
		db: db,
	}
}

func (f *FrontendHandler) HomePage(w http.ResponseWriter, r *http.Request) {
	// if !f.auth.Jwt.CheckToken(r) {

	// }
	tmpl, err := template.ParseFiles("templates/index.html")
	if err != nil {
		fmt.Printf("error while parsing html file: %v\n", err)
		http.Error(w, "error while parsing html file", http.StatusInternalServerError)
	}
	err = tmpl.Execute(w, struct{ Text string }{"text"})
	if err != nil {
		fmt.Printf("error while exucuting html template: %v\n", err)
		http.Error(w, "error filing the template", http.StatusInternalServerError)
	}
}

func (f *FrontendHandler) LoginPage(w http.ResponseWriter, r *http.Request) {
	tmpl, err := template.ParseFiles("templates/login.html")
	if err != nil {
		fmt.Printf("error while parsing html file: %v\n", err)
		http.Error(w, "error while parsing html file", http.StatusInternalServerError)
	}
	err = tmpl.Execute(w, "")
	if err != nil {
		fmt.Printf("error while exucuting html template: %v\n", err)
		http.Error(w, "error filing the template", http.StatusInternalServerError)
	}
}
