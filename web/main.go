package main

import (
	"errors"
	"fmt"
	"html/template"
	"log"
	"net/http"
	"os"
	"regexp"
)

var templates = template.Must(template.ParseFiles("edit.html", "view.html"))
var validPath = regexp.MustCompile("^/(edit|save|view)/([a-zA-Z0-9]+)$")

type Page struct {
	Title string
	Body  []byte
}

func (p *Page) save() error {
	filename := p.Title + ".txt"
	return os.WriteFile(filename, p.Body, 0600)
}

func loadPage(title string) (*Page, error) {
	filename := title + ".txt"
	body, err := os.ReadFile(filename)
	if err != nil {
		return nil, err
	}
	return &Page{
		Title: title,
		Body:  body,
	}, nil
}
func getTitle(path string) (string, error) {
	m := validPath.FindStringSubmatch(path)
	if m == nil {
		return "", errors.New("invalid Page title")
	}
	return m[2], nil
}

func renderTemplete(w http.ResponseWriter, tmpl string, p *Page) error {
	file := fmt.Sprintf("%s.html", tmpl)
	return templates.ExecuteTemplate(w, file, p)
}

func viewHandler(w http.ResponseWriter, r *http.Request, title string) {
	p, err := loadPage(title)
	if err != nil {
		http.Redirect(w, r, fmt.Sprintf("/edit/%s", title), http.StatusFound)
		return
	}
	err = renderTemplete(w, "view", p)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}

func editHandler(w http.ResponseWriter, r *http.Request, title string) {
	p, err := loadPage(title)
	if err != nil {
		p = &Page{
			Title: title,
		}
	}
	err = renderTemplete(w, "edit", p)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}

func saveHandler(w http.ResponseWriter, r *http.Request, title string) {
	body := []byte(r.FormValue("body"))
	p, err := loadPage(title)
	if err != nil {
		p = &Page{
			Title: title,
			Body:  body,
		}
	} else {
		p.Body = []byte(r.FormValue("body"))
	}
	if err := p.save(); err != nil {
		http.Error(w, err.Error(), 500)
		return
	}
	http.Redirect(w, r, fmt.Sprintf("/view/%s", title), http.StatusFound)
}

func makeHandler(fn func(w http.ResponseWriter, r *http.Request, title string)) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		title, err := getTitle(r.URL.Path)
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}
		fn(w, r, title)
	}
}

func main() {
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("Home"))
	})
	http.HandleFunc("/view/", makeHandler(viewHandler))
	http.HandleFunc("/edit/", makeHandler(editHandler))
	http.HandleFunc("/save/", makeHandler(saveHandler))
	log.Fatal(http.ListenAndServe(":8000", nil))
}
