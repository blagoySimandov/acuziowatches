package main

import (
	"net/http"
	"strconv"
	"text/template"
	"time"

	"github.com/labstack/echo/v4"
	"github.com/labstack/gommon/log"
)

const productsJSON string = "./static/products.json"
const bestProductsJSON string = "./static/bestProducts.json"

type Product struct {
	ProductName  string
	ProductPrice int
	Count        int
}
type ProductData struct {
	Products []struct {
		Name     string `json:"name"`
		Price    string `json:"price"`
		Currency string `json:"currency"`
		Img      string `json:"img"`
		Id       string `json:"id"`
		OldPrice string `json:"oldPrice"`
	} `json:"products"`
}

func (p Product) ProductTotal() int {
	return p.Count * p.ProductPrice
}

//create a Cookie

func writeCookie(c echo.Context, name string, value string) {
	cookie := new(http.Cookie)
	cookie.Name = name
	cookie.Value = value
	//Make the Cookie expire after 2 hours
	cookie.Expires = time.Now().Add(2 * time.Hour)
	c.SetCookie(cookie)
}

func Index(c echo.Context) error {
	/*count := c.FormValue("count")

	id := "1"
	fmt.Printf("Quantity of wathces ordered: %v\n", count)
	writeCookie(c, "count", count)
	writeCookie(c, "id", id)*/
	return c.Render(http.StatusOK, "indexTmpl", bestProducts)
}

func Shop(c echo.Context) error {
	return c.Render(http.StatusOK, "shopTmpl", productData)

}
func ProductDetails(c echo.Context) error {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		log.Error(err)
	}

	return c.Render(http.StatusOK, "productTmpl", productData.Products[id])
}

var productData *ProductData
var bestProducts *ProductData

func main() {
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
	e.Static("/", "./static")
	e.GET("/shop", Shop)
	e.GET("/", Index)
	e.GET("/product/:id", ProductDetails)
	//e.POST("/confirm", conf)

	e.Logger.Fatal(e.Start(":8080"))

}
