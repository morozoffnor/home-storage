package frontend

import (
	"fmt"
	"html/template"
	"net/http"
)

func (c *Container) GetAllInHome(w http.ResponseWriter, r *http.Request) {
	homeID, ok := r.Context().Value("home_id").(int)
	if !ok {
		http.Error(w, http.StatusText(400), 400)
	}

	type container struct {
		ID          int
		Name        string
		Description string
		Category    string
		Location    string
		CreatedAt   string
		ItemCount   int
		HomeID      int
	}

	var containers []*container
	// home ctx for frontend
	cont, err := c.db.Container.GetAllInHome(homeID)
	if err != nil {
		http.Error(w, http.StatusText(500), 500)
		return
	}

	for _, v := range cont {

		itemsCount, err := c.db.Container.ItemsCount(v.ID)
		if err != nil {
			http.Error(w, http.StatusText(500), 500)
			return
		}
		containers = append(containers, &container{
			ID:          v.ID,
			Name:        v.Name,
			Description: v.Description,
			Category:    v.Category,
			Location:    v.Location,
			CreatedAt:   v.CreatedAt.Format("2 January, 2006"),
			ItemCount:   itemsCount,
			HomeID:      v.HomeID,
		})
	}

	tmpl, err := template.ParseFiles("templates/container_block.html")
	if err != nil {
		fmt.Printf("error while parsing html file: %v\n", err)
		http.Error(w, "error while parsing html file", http.StatusInternalServerError)
	}
	err = tmpl.Execute(w, containers)
	if err != nil {
		fmt.Printf("error while exucuting html template: %v\n", err)
		http.Error(w, "error filing the template", http.StatusInternalServerError)
	}

}
