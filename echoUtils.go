package main

import (
	"context"
	"io"
	"text/template"

	"github.com/labstack/echo/v4"
	"go.mongodb.org/mongo-driver/bson"
)

func loadData(collection string) (*ProductData, error) {
	var p ProductData

	coll := mongoClient.Database("Products").Collection(collection)
	docs, err := coll.Find(context.TODO(), bson.D{})
	if err != nil {
		return nil, err
	}
	if err := docs.All(context.Background(), &p.Products); err != nil {
		return nil, err
	}
	return &p, nil
}

type Template struct {
	templates *template.Template
}

// func PercentFiller(arr []*Product) {
// 	for _, e := range arr {
// 		price, err := strconv.ParseFloat(e.Price, 64)
// 		if err != nil {
// 			fmt.Printf("errorparsing: %v", err)
// 		}
// 		oldPrice, err := strconv.ParseFloat(e.OldPrice, 64)
// 		if err != nil {
// 			fmt.Printf("parse error: %v", err)
// 		}
// 		e.PercentOfOldPrice = math.Round(((price / oldPrice) * 100))
// 	}
// }

func (t *Template) Render(w io.Writer, name string, data interface{}, c echo.Context) error {
	return t.templates.ExecuteTemplate(w, name, data)
}
