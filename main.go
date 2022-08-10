package main

import (
	"net/http"
	"text/template"
	"time"

	"github.com/labstack/echo/v4"
)

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
	return c.Render(http.StatusOK, "indexTmpl", productData)
}

func Shop(c echo.Context) error {
	return c.Render(http.StatusOK, "shopTmpl", productData)

}

var productData *ProductData

func main() {

	localPd, err := loadData()
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
	//e.POST("/confirm", conf)

	e.Logger.Fatal(e.Start(":8080"))

}
