package main

import (
	"context"
	"fmt"
	"net/http"
	"os"
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
	products           = "products"
	bestProducts       = "bestProducts"
	defaultSessionName = "session"
	clientID           = "AcQJf_BZi9tLGfptFCXm0hv_TkItaYow6I0VR2UthIppfaWDWKjAm2kmSJhOIEIQsklEEuEjUXJgCs0q"
	secretID           = "EFw4g1cEJIPzUYIxXWhvURDBBKPcqnIpqwSAyW-h4NWe3NBg28SYWgBfCi7UH-yh03SCgaK_mb2l9n7S"
)

type (
	CartProduct struct {
		Pr       Product
		Count    string
		Subtotal string
	}

	CartProducts struct {
		Products []CartProduct
		Total    float64 //`default:0`
	}

	Product struct {
		Id          string `bson:"_id" json:"id"`
		Name        string `json:"name"`
		Currency    string `json:"currency"`
		Price       string `json:"price"`
		OldPrice    string `json:"oldPrice"`
		Title       string `json:"title"`
		Description string `json:"description"`
		Meta        string `json:"meta"`
	}

	ProductData struct {
		Products []Product `json:"products"`
	}
)

//create a Cookie

func Index(c echo.Context) error {
	bestProducts, err := loadData(bestProducts, "")
	if err != nil {
		c.Logger().Error("error loading data: %v", err)
		return err
	}
	return c.Render(http.StatusOK, "indexTmpl", bestProducts)
}

func (p *ProductData) CalculateDiscount() {
	// for _, e := range p.Products {
	// 	e.Title = "f2osd"
	// }
}

func Shop(c echo.Context) error {
	order := c.QueryParam("order")

	var (
		orderedData *ProductData
		err         error
	)
	if order == "HtoL" {
		orderedData, err = loadData(products, "HtoL")
		if err != nil {
			c.Logger().Error(err)
			return err
		}

	} else if order == "LtoH" {
		orderedData, err = loadData(products, "LtoH")
		if err != nil {
			c.Logger().Error(err)
			return err
		}
	} else {
		orderedData, err = loadData(products, "")
		if err != nil {
			c.Logger().Error(err)
			return err
		}
	}
	return c.Render(http.StatusOK, "shopTmpl", orderedData)
}
func ProductDetails(c echo.Context) error {
	id := c.Param("id")
	products, err := loadData(products, "")
	if err != nil {
		c.Logger().Error("error Loading products: %v", err)
		return err
	}
	for _, e := range products.Products {
		if e.Id == id {
			fmt.Println("found")
			return c.Render(http.StatusOK, "productTmpl", e)
		}
	}
	return c.Render(http.StatusNotFound, "404", "")
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
	http.Redirect(c.Response(), c.Request(), "/cart", 301)
	return nil

}
func SendMessage(c echo.Context) error {
	email, name := c.FormValue("email"), c.FormValue("name")
	message, subject := c.FormValue("message"), c.FormValue("subject")

	title := subject + "    by: " + name + " email: " + email + "\n"

	coll := mongoClient.Database("Contacts").Collection("Contacts")
	doc := bson.D{{Key: "title", Value: title}, {Key: "body", Value: message}, {Key: "name", Value: name}, {Key: "email", Value: email}}

	ctx, cancel := context.WithTimeout(c.Request().Context(), 10*time.Second)
	defer cancel()
	result, err := coll.InsertOne(ctx, doc)
	if err != nil {
		c.Logger().Error("Cannot insert data in DB. Error: %v", err)
		return err
	}
	_ = result
	// c.Logger().Printf(("Inserted document with _id: %v\n", result.InsertedID)
	return c.Redirect(http.StatusMovedPermanently, "/submit-success")
}

func Checkout(c echo.Context) error {
	sess, err := session.Get(defaultSessionName, c)
	if err != nil {
		log.Error("error in creating a session")
	}

	cart, err := loadCart(sess)
	if err != nil {
		c.Logger().Error("error loading cart: %v", err)
	}
	return c.Render(http.StatusOK, "checkoutTmpl", cart)
}

func Cart(c echo.Context) error {

	sess, err := session.Get(defaultSessionName, c)
	if err != nil {
		log.Error("error in creating a session")
		return err
	}

	cart, err := loadCart(sess)
	if err != nil {
		c.Logger().Error("error loading cart: %v", err)
		return err
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
	_, err = p.GetAccessToken(c.Request().Context())
	if err != nil {
		panic(err)
	}

	sess, err := session.Get(defaultSessionName, c)
	if err != nil {
		log.Error("error in creating a session")
	}

	cart, err := loadCart(sess)
	if err != nil {
		c.Logger().Error("error loading cart: %v", err)
	}
	totalString := fmt.Sprintf("%.3f", cart.Total)
	_ = totalString
	paypalAmount := paypal.PurchaseUnitAmount{}
	paypalAmount.Currency = "USD"
	paypalAmount.Value = totalString
	var orderPayer paypal.CreateOrderPayer
	var application paypal.ApplicationContext
	order, err := p.CreateOrder(context.Background(), paypal.OrderIntentCapture,
		[]paypal.PurchaseUnitRequest{
			{ReferenceID: "ref-id", Amount: &paypalAmount},
		},
		&orderPayer, &application)

	fmt.Println(order)
	// if err != nil {
	// 	log.Error(err)
	// }
	return Cart(c)
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

var (
	uri string = "mongodb+srv://acuzio:uOBzJFvD4voHaWdb@cluster0.ynz5x4i.mongodb.net/?retryWrites=true&w=majority"

	mongoClient *mongo.Client
)

func connectMongo(ctx context.Context) (*mongo.Client, error) {
	ctx, cancel := context.WithTimeout(ctx, 10*time.Second)
	defer cancel()
	// serverAPIOptions := options.ServerAPI(options.ServerAPIVersion1)
	// clientOptions :=
	// 	SetServerAPIOptions(serverAPIOptions)
	client, err := mongo.Connect(ctx, options.Client().ApplyURI(uri))
	if err != nil {
		return nil, err
	}
	return client, nil
}

func main() {
	globalCtx := context.Background()
	var err error

	if mongoClient, err = connectMongo(globalCtx); err != nil {
		panic(err)
	}

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

	e.POST("/sendMessage", SendMessage)

	e.File("/about", "static/about.html")
	e.File("/contact", "static/contact.html")
	e.File("/submit-success", "static/submit-success.html")
	e.GET("/checkout", Checkout)

	//Cart Requests
	e.POST("/remove", Remove)
	e.GET("/cart", Cart)

	//Paypal POST
	e.POST("/api/orders", PayPalOrder)
	e.POST("/api/orders/capture/:id", PayPalCaptureOrder)
	//e.POST("/confirm", conf)

	e.Logger.Fatal(e.Start(":" + getDefEnv("PORT", "8080")))

}

func getDefEnv(env string, def string) (res string) {
	res = os.Getenv(env)
	if res == "" {
		res = def
	}
	return
}
