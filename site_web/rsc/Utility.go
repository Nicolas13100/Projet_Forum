package API

import (
	"fmt"
	"html/template"
	"net/http"
	"strings"
)

// ////
func renderTemplate(w http.ResponseWriter, tmplName string, data interface{}) {
	// Taken from hangman
	tmpl, err := template.New(tmplName).Funcs(template.FuncMap{"join": join, "contains": containsString}).ParseFiles("site_web/Template/" + tmplName + ".html")
	if err != nil {
		fmt.Println("Error parsing template:", err)
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}

	err = tmpl.ExecuteTemplate(w, tmplName, data)
	if err != nil {
		fmt.Println("Error executing template:", err)
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}
}
func join(s []string, sep string) string {
	// same
	return strings.Join(s, sep)
}
func containsString(slice []string, str string) bool {
	for _, s := range slice {
		if s == str {
			return true
		}
	}
	return false
}

///////

func Init() {

}
