package main

import (
	"context"
	"io"
	"log"
	"text/template"
	"time"

	"github.com/labstack/echo/v4"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

func loadData(path string) (*ProductData, error) {
	var p ProductData
	serverAPIOptions := options.ServerAPI(options.ServerAPIVersion1)
	clientOptions := options.Client().
		ApplyURI(uri).
		SetServerAPIOptions(serverAPIOptions)
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	client, err := mongo.Connect(ctx, clientOptions)
	if err != nil {
		log.Fatal(err)
	}
	coll := client.Database("Products").Collection(path)
	if err = coll.FindOne(context.TODO(), bson.D{}).Decode(&p); err != nil {
		panic(err)
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
