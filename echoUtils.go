package main

import (
	"context"
	"fmt"
	"io"
	"strconv"
	"strings"
	"text/template"
	"time"

	"github.com/gorilla/sessions"
	"github.com/labstack/echo/v4"
	"github.com/labstack/gommon/log"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

func loadData(collection string, order string) (*ProductData, error) {
	coll := mongoClient.Database("Products").Collection(collection)
	var (
		opt *options.FindOptions
	)
	collation := options.Collation{
		Locale:          "en_US",
		NumericOrdering: true,
	}
	if order == "LtoH" {
		opt = options.Find().SetSort(bson.D{{Key: "price", Value: 1}}).SetCollation(&collation)
		return findFiles(coll, opt)
	} else if order == "HtoL" {
		// Sort by `price` field descending
		opt = options.Find().SetSort(bson.D{{Key: "price", Value: -1}}).SetCollation(&collation)
		return findFiles(coll, opt)
	}

	return findFiles(coll, nil)
}

func findFiles(coll *mongo.Collection, opt *options.FindOptions) (*ProductData, error) {
	var (
		p    ProductData
		docs *mongo.Cursor
		err  error
	)
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)

	defer cancel()
	if opt == nil {
		docs, err = coll.Find(ctx, bson.D{})
	} else {
		docs, err = coll.Find(ctx, bson.D{}, opt)
	}
	if err != nil {
		log.Error("error finding the docs:", err)
		return nil, err
	}
	if err := docs.All(ctx, &p.Products); err != nil {
		log.Error("error converting the cursor ro array:", err)
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

func PriceStringToInt(a string) (int, error) {
	a = strings.Replace(a, ".", "", -1)
	aInt, err := strconv.Atoi(a)
	if err != nil {
		return 0, err
	}
	return aInt, nil
}

func MultiplyStringPrices(a, b string) (string, error) {
	var (
		err  error
		aInt int
		bInt int
	)
	aInt, err = PriceStringToInt(a)
	if err != nil {
		return "", err
	}
	bInt, err = PriceStringToInt(b)
	if err != nil {
		return "", err
	}
	product := aInt * bInt
	productString := strconv.Itoa(product)
	productString = AddDecimalPointToString(productString)
	return productString, nil
}

func AddStringPrices(a, b string) (string, error) {
	var (
		err  error
		aInt int
		bInt int
	)
	aInt, err = PriceStringToInt(a)
	if err != nil {
		return "", err
	}
	bInt, err = PriceStringToInt(b)
	if err != nil {
		return "", err
	}
	result := aInt + bInt
	resultString := strconv.Itoa(result)
	resultString = AddDecimalPointToString(resultString)
	return resultString, nil
}
func AddDecimalPointToString(a string) string {
	index := len(a) - 2 // 000000.00 two places decimal
	a = a[:index] + "." + a[index:]
	return a
}
func loadCart(sess *sessions.Session) (*CartProducts, error) {
	products, err := loadData("products", "")
	if err != nil {
		return nil, err
	}
	var cart CartProducts
	cart.Total = "0"
	for key, count := range sess.Values {
		for _, p := range products.Products {
			if key == p.Id {
				subtotal, err := MultiplyStringPrices(count.(string), p.Price)
				fmt.Println(subtotal)
				if err != nil {
					return nil, err
				}
				x := CartProduct{
					Pr:       p,
					Count:    count.(string),
					Subtotal: subtotal,
				}
				cart.Products = append(cart.Products, x)
				cart.Total, err = AddStringPrices(subtotal, cart.Total)
				if err != nil {
					return nil, err
				}
			}
		}
	}
	return &cart, nil
}
