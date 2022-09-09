package main

import (
	"context"
	"fmt"
	"net/http"
	"sort"
	"strconv"
	"text/template"
	"time"

	"github.com/gorilla/sessions"
	"github.com/labstack/echo-contrib/session"
	"github.com/labstack/echo/v4"
	"github.com/labstack/gommon/log"
	"github.com/plutov/paypal/v4"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

const (
	productsJSON       = "products"
	bestProductsJSON   = "bestProducts"
	defaultSessionName = "session"
	clientID           = "AcQJf_BZi9tLGfptFCXm0hv_TkItaYow6I0VR2UthIppfaWDWKjAm2kmSJhOIEIQsklEEuEjUXJgCs0q"
	secretID           = "EFw4g1cEJIPzUYIxXWhvURDBBKPcqnIpqwSAyW-h4NWe3NBg28SYWgBfCi7UH-yh03SCgaK_mb2l9n7S"
)

type (
	CartProduct struct {
		Pr       Product
		Count    string
		Subtotal float64
	}

	CartProducts struct {
		Products []CartProduct
		Total    float64 `default:0`
	}

	Product struct {
		Name              string  `json:"name"`
		Price             string  `json:"price"`
		Currency          string  `json:"currency"`
		Id                string  `json:"id"`
		OldPrice          string  `json:"oldPrice"`
		Title             string  `json:"title"`
		Description       string  `json:"description"`
		Meta              string  `json:"meta"`
		PercentOfOldPrice float64 `json:"percent"`
	}

	ProductData struct {
		Products []Product `json:"products"`
	}
)

//create a Cookie

func Index(c echo.Context) error {

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
	sess.Values[id] = count
	sess.Save(c.Request(), c.Response())
	return c.Redirect(http.StatusMovedPermanently, "/cart")

}
func SendMessage(c echo.Context) error {
	email := c.FormValue("email")
	name := c.FormValue("name")
	message := c.FormValue("message")
	subject := c.FormValue("subject") + "    by: " + name + " email: " + email + "\n"

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
	coll := client.Database("Contacts").Collection("Contacts")
	doc := bson.D{{"subject", subject}, {"body", message}, {"name", name}, {"email", email}}
	result, _ := coll.InsertOne(context.TODO(), doc)
	fmt.Printf("Inserted document with _id: %v\n", result.InsertedID)
	return c.Redirect(http.StatusMovedPermanently, "/submit-success")
}

func Checkout(c echo.Context) error {
	sess, err := session.Get(defaultSessionName, c)
	if err != nil {
		log.Error("error in creating a session")
	}

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
				cart.Total = cart.Total + x.Subtotal
			}
		}
	}
	return c.Render(http.StatusOK, "checkoutTmpl", cart)
}

func Cart(c echo.Context) error {

	sess, err := session.Get(defaultSessionName, c)
	if err != nil {
		log.Error("error in creating a session")
	}

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
func Remove(c echo.Context) error {
	id := c.FormValue("id")
	sess, err := session.Get(defaultSessionName, c)
	if err != nil {
		log.Error("error creating session: %v", err)
	}
	delete(sess.Values, id)
	sess.Save(c.Request(), c.Response())

	return c.Redirect(http.StatusMovedPermanently, "/cart")
}

func PayPalOrder(c echo.Context) error {
	// Initialize client
	p, err := paypal.NewClient(clientID, secretID, paypal.APIBaseSandBox)
	if err != nil {
		panic(err)
	}

	// Retrieve access token
	_, err = p.GetAccessToken(context.Background())
	if err != nil {
		panic(err)
	}

	sess, err := session.Get(defaultSessionName, c)
	if err != nil {
		log.Error("error in creating a session")
	}

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
				cart.Total = cart.Total + x.Subtotal
			}
		}
	}
	totalString := fmt.Sprintf("%.3f", cart.Total)
	paypalAmount := paypal.PurchaseUnitAmount{}
	paypalAmount.Currency = "USD"
	paypalAmount.Value = totalString
	var orderPayer *paypal.CreateOrderPayer
	var application *paypal.ApplicationContext
	order, err := p.CreateOrder(context.Background(), paypal.OrderIntentCapture, []paypal.PurchaseUnitRequest{paypal.PurchaseUnitRequest{ReferenceID: "ref-id", Amount: &paypalAmount}}, orderPayer, application)
	fmt.Println(order)
	// if err != nil {
	// 	log.Error(err)
	// }
	return nil
}
func PayPalCaptureOrder(c echo.Context) error {
	//orderID := c.Param("id")
	p, err := paypal.NewClient(clientID, secretID, paypal.APIBaseSandBox)
	if err != nil {
		panic(err)
	}

	// Retrieve access token
	_, err = p.GetAccessToken(context.Background())
	if err != nil {
		panic(err)
	}
	//capture, err := p.CaptureOrder(orderID, paypal.CaptureOrderRequest{})

	return nil
}

var productData *ProductData
var bestProducts *ProductData

var uri string = "mongodb+srv://acuzio:uOBzJFvD4voHaWdb@cluster0.ynz5x4i.mongodb.net/?retryWrites=true&w=majority"

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
	e.Use(session.Middleware(sessions.NewCookieStore([]byte("secret"))))

	e.Static("/", "./static")
	e.GET("/shop", Shop)
	e.GET("/", Index)
	e.GET("/product/:id", ProductDetails)
	e.POST("/addToCart/:id", AddToCart)
	e.GET("/cart", Cart)
	e.POST("/sendMessage", SendMessage)
	e.POST("/remove", Remove)
	e.File("/about", "static/about.html")
	e.File("/contact", "static/contact.html")
	e.File("/submit-success", "static/submit-success.html")
	e.GET("/checkout", Checkout)

	//Paypal POST
	e.POST("/api/orders", PayPalOrder)
	e.POST("/api/orders/capture/:id", PayPalCaptureOrder)
	//e.POST("/confirm", conf)

	e.Logger.Fatal(e.Start("0.0.0.0:8080"))

}
