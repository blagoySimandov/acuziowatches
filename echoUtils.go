package main

import (
	"encoding/json"
	"io"
	"log"
	"os"
	"text/template"

	"github.com/labstack/echo/v4"
)

func loadData() (*ProductData, error) {
	f, err := os.Open("./static/products.json")
	if err != nil {
		log.Println("Cannot open file")
		log.Println(err)
		return nil, err
	}
	defer f.Close()
	enc := json.NewDecoder(f)
	var p ProductData
	if err := enc.Decode(&p); err != nil {
		log.Println(err)
		return nil, err
	}
	return &p, nil
}

type Template struct {
	templates *template.Template
}

func (t *Template) Render(w io.Writer, name string, data interface{}, c echo.Context) error {
	return t.templates.ExecuteTemplate(w, name, data)
}
