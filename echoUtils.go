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
	"go.mongodb.org/mongo-driver/mongo/options"
)

func loadData(collection string, order string) (*ProductData, error) {
	coll := mongoClient.Database("Acuzio").Collection(collection)

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	opt := options.Find()
	if order == "LtoH" {
		opt.SetSort(bson.D{{Key: "price", Value: 1}})
	} else if order == "HtoL" {
		opt.SetSort(bson.D{{Key: "price", Value: -1}})
	} else if order == "newest" || order == "" {
		// we do nothing to Options
	}
	docs, err := coll.Find(ctx, bson.D{}, opt)
	if err != nil {
		return nil, err
	}
	var p ProductData
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

	err := t.templates.ExecuteTemplate(w, name, data)
	if err != nil {
		fmt.Println("Template error:", err)
		return err
	}
	return nil
}

func PriceStringToInt(a string) (int, error) {
	a = strings.Replace(a, ".", "", -1)
	aInt, err := strconv.Atoi(a)
	if err != nil {
		return 0, err
	}
	return aInt, nil
}

func loadCart(sess *sessions.Session) (*CartProducts, error) {
	products, err := loadData(products, "")
	if err != nil {
		return nil, err
	}
	var cart CartProducts
	cart.Total = 0
	for key, value := range sess.Values {
		for _, p := range products.Products {
			if key == p.Id {
				count := value.(string) // TODO Handle reflection error
				countInt, err := strconv.Atoi(count)
				if err != nil {
					return &CartProducts{}, err
				}
				subtotal := typePrice(countInt) * (p.PriceWithDiscount())
				if err != nil {
					return nil, err
				}
				x := CartProduct{
					Pr:       p,
					Count:    countInt,
					Subtotal: subtotal,
				}
				cart.Products = append(cart.Products, x)
				cart.Total += subtotal
			}
		}
	}
	return &cart, nil
}
