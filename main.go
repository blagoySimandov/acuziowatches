package main

import (
	"fmt"
	"net/http"
	"sort"
	"strconv"
	"text/template"

	"github.com/gorilla/sessions"
	"github.com/labstack/echo-contrib/session"
	"github.com/labstack/echo/v4"
	"github.com/labstack/gommon/log"
)

const (
	productsJSON       = "./static/products.json"
	bestProductsJSON   = "./static/bestProducts.json"
	defaultSessionName = "session"
)

type (
	CartProduct struct {
		Pr       Product
		Count    string
		Subtotal float64
	}

	CartProducts struct {
		Products []CartProduct
	}

	Product struct {
		Name     string `json:"name"`
		Price    string `json:"price"`
		Currency string `json:"currency"`
		Id       string `json:"id"`
		OldPrice string `json:"oldPrice"`
	}

	ProductData struct {
		Products []Product `json:"products"`
	}
)

//create a Cookie

func Index(c echo.Context) error {
	/*count := c.FormValue("count")

	id := "1"
	fmt.Printf("Quantity of wathces ordered: %v\n", count)
	writeCookie(c, "count", count)
	writeCookie(c, "id", id)*/

	return c.Render(http.StatusOK, "indexTmpl", bestProducts)
}

func Shop(c echo.Context) error {
	order := c.QueryParam("order")
	var orderedData ProductData
	orderedData.Products = append([]Product{}, productData.Products...)
	if order == "HtoL" {
		sort.Slice(orderedData.Products, func(i, j int) bool {
			iPrice, err := strconv.ParseFloat(orderedData.Products[i].Price, 32)
			if err != nil {
				log.Error(err)
			}
			jPrice, err := strconv.ParseFloat(orderedData.Products[j].Price, 32)
			if err != nil {
				log.Error(err)
			}
			return iPrice > jPrice
		})
	} else if order == "LtoH" {
		sort.Slice(orderedData.Products, func(i, j int) bool {
			iPrice, err := strconv.ParseFloat(orderedData.Products[i].Price, 32)
			if err != nil {
				log.Error(err)
			}
			jPrice, err := strconv.ParseFloat(orderedData.Products[j].Price, 32)
			if err != nil {
				log.Error(err)
			}
			return iPrice < jPrice
		})
	} else {
		orderedData.Products = append([]Product{}, productData.Products...)
	}
	return c.Render(http.StatusOK, "shopTmpl", orderedData)
}
func ProductDetails(c echo.Context) error {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		log.Error(err)
	}

	return c.Render(http.StatusOK, "productTmpl", productData.Products[id])
}

func AddToCart(c echo.Context) error {
	id := c.Param("id")
	count := c.FormValue("count")
	//session
	sess, err := session.Get(defaultSessionName, c)
	if err != nil {
		log.Error("error in creating a session")
	}
	sess.Options = &sessions.Options{
		Path:     "/",
		MaxAge:   3600,
		HttpOnly: true,
	}
	fmt.Println(id)
	fmt.Println(sess.Values)
	sess.Values[id] = count
	sess.Save(c.Request(), c.Response())
	fmt.Println("just  before redirect")
	fmt.Println(sess.Values)
	return c.Redirect(http.StatusMovedPermanently, "/cart")

}

func Cart(c echo.Context) error {

	sess, err := session.Get(defaultSessionName, c)
	if err != nil {
		log.Error("error in creating ta session")
	}
	fmt.Println(sess.Values)
	//fmt.Println(productData.Products[0].Count)

	var cart CartProducts
	for key, count := range sess.Values {
		for _, p := range productData.Products {
			if key == p.Id {
				priceFloat, err := strconv.ParseFloat(p.Price, 64)
				if err != nil {
					log.Error(err)
				}
				countFloat, err := strconv.ParseFloat(count.(string), 64)
				if err != nil {
					log.Error(err)
				}
				x := CartProduct{
					Pr:       p,
					Count:    count.(string),
					Subtotal: countFloat * priceFloat,
				}
				cart.Products = append(cart.Products, x)
			}
		}
	}

	return c.Render(http.StatusOK, "cartTmpl", cart)
}

var productData *ProductData
var bestProducts *ProductData

//var cartData *CartData //Cart Data = Product Data + a field for quantity "Count"

func main() {

	// localCd, err := loadData(productsJSON)
	// if err != nil {
	// 	panic(err)
	// }

	localBpd, err := loadData(bestProductsJSON)
	if err != nil {
		panic(err)
	}
	bestProducts = localBpd

	localPd, err := loadData(productsJSON)
	if err != nil {
		panic(err)
	}
	productData = localPd

	e := echo.New()
	e.Renderer = &Template{
		templates: template.Must(template.ParseGlob("./static/*.html")),
	}
	e.Use(session.Middleware(sessions.NewCookieStore([]byte("secret"))))

	e.Static("/", "./static")
	e.GET("/shop", Shop)
	e.GET("/", Index)
	e.GET("/product/:id", ProductDetails)
	e.POST("/addToCart/:id", AddToCart)
	e.GET("/cart", Cart)
	e.File("/about", "static/about.html")
	e.File("/contact", "static/contact.html")
	//e.POST("/confirm", conf)

	e.Logger.Fatal(e.Start("0.0.0.0:8080"))

}
