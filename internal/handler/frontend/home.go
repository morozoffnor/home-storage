package frontend

import (
	"fmt"
	"html/template"
	"net/http"
)

func (h *Home) GetAll(w http.ResponseWriter, r *http.Request) {

	type FrontendHome struct {
		ID              int
		Name            string
		Description     string
		CreatedAt       string
		ContainersCount int
		ItemsCount      int
	}
	var frontendHomes []*FrontendHome
	homes, err := h.db.Home.GetAll()
	if err != nil {
		if err.Error() == "no rows in result set" {
			fmt.Printf("error while getting homes: %v\n", err)
			http.Error(w, "no homes", http.StatusNotFound)
			return
		} else {
			fmt.Printf("error while getting homes: %v\n", err)
			http.Error(w, "error while getting homes", http.StatusInternalServerError)
			return
		}
	}

	for _, v := range homes {
		contCount, err := h.db.Home.ContainersCount(v.ID)
		if err != nil {
			fmt.Printf("error while getting containers count: %v\n", err)
			http.Error(w, "error while getting containers count", http.StatusInternalServerError)
		}
		itemsCount, err := h.db.Home.ItemsCount(v.ID)
		if err != nil {
			fmt.Printf("error while getting items count: %v\n", err)
			http.Error(w, "error while getting items count", http.StatusInternalServerError)
		}
		frontendHomes = append(frontendHomes, &FrontendHome{
			ID:              v.ID,
			Name:            v.Name,
			Description:     v.Description,
			CreatedAt:       v.CreatedAt.Format("2 January, 2006"),
			ContainersCount: contCount,
			ItemsCount:      itemsCount,
		})
	}
	tmpl, err := template.ParseFiles("templates/home_block.html")
	if err != nil {
		fmt.Printf("error while parsing html file: %v\n", err)
		http.Error(w, "error while parsing html file", http.StatusInternalServerError)
	}
	err = tmpl.Execute(w, frontendHomes)
	if err != nil {
		fmt.Printf("error while exucuting html template: %v\n", err)
		http.Error(w, "error filing the template", http.StatusInternalServerError)
	}
}
