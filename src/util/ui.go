package util

import (
	"html/template"
	"log"
	"net/http"
)

func RenderTemplate(w http.ResponseWriter, filename string, data any) {
	layout := template.Must(template.ParseGlob("../templates/layout.gohtml"))
	tmpl := template.Must(layout.ParseFiles("../templates/" + filename + ".gohtml"))
	err := tmpl.ExecuteTemplate(w, "layout.gohtml", data)
	if err != nil {
		http.Error(w, "Error rendering template", http.StatusInternalServerError)
		log.Println("Template execution error:", err)
	}
}
func RenderMessage(w http.ResponseWriter, message string) {
	layout := template.Must(template.ParseGlob("../templates/layout.gohtml"))
	tmpl := template.Must(layout.ParseFiles("../templates/message.gohtml"))
	data := map[string]any{
		"Message": message,
	}
	err := tmpl.ExecuteTemplate(w, "layout.gohtml", data)
	if err != nil {
		http.Error(w, "Error rendering template", http.StatusInternalServerError)
		log.Println("Template execution error:", err)
	}
}
