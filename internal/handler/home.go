package handler

import (
	"html/template"
	"net/http"
	"scrable3/internal/dto"
)

type homeHandler struct{}

func NewHomeHandler() http.Handler {
	return &homeHandler{}
}

func (h *homeHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	tmpl, err := template.ParseFiles("views/home/index.html")
	if err != nil {
		http.Error(w, err.Error(), 500)
		return
	}

	data := dto.HomePageData{Title: "Scrable3D home page"}

	err = tmpl.Execute(w, data)
	if err != nil {
		http.Error(w, err.Error(), 500)
	}
}
