package main

import (
	"context"
	"io"
	"strconv"
	"strings"
	"text/template"

	"github.com/gorilla/sessions"
	"github.com/labstack/echo/v4"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

func loadData(collection string, order string) (*ProductData, error) {
	coll := mongoClient.Database("Products").Collection(collection)
	if order == "LtoH" {
		options.Find().SetSort(bson.D{{"price", +1}})
		return findFiles(coll)
	} else if order == "HtoL" {
		// Sort by `price` field descending
		options.Find().SetSort(bson.D{{"price", -1}})
		return findFiles(coll)
	}

	options.Find().SetSort(bson.D{{}})
	return findFiles(coll)
}

func findFiles(coll *mongo.Collection) (*ProductData, error) {
	var p ProductData
	docs, err := coll.Find(context.Background(), bson.D{})
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

func MultiplyStringPrices(a, b string) (string, error) {
	a, b = strings.Replace(a, ".", "", -1), strings.Replace(a, ".", "", -1)
	aInt, err := strconv.Atoi(a)
	if err != nil {
		return "", err
	}
	bInt, err := strconv.Atoi(b)
	if err != nil {
		return "", err
	}
	product := aInt * bInt
	productString := strconv.Itoa(product)
	index := len(productString) - 3 // 000000.00 two places decimal
	productString = productString[:index] + "." + productString[index:]
	return productString, nil
}

func loadCart(sess *sessions.Session) (*CartProducts, error) {
	products, err := loadData("products", "")
	if err != nil {
		return nil, err
	}
	var cart CartProducts
	for key, count := range sess.Values {
		for _, p := range products.Products {
			if key == p.Id {
				subtotal, err := MultiplyStringPrices(count.(string), p.Price)
				if err != nil {
					return nil, err
				}
				x := CartProduct{
					Pr:       p,
					Count:    count.(string),
					Subtotal: subtotal,
				}
				cart.Products = append(cart.Products, x)
			}
		}
	}
	return &cart, nil
}
