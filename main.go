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
	products           = "Products"
	bestProducts       = "bestProducts"
	defaultSessionName = "session"
)

var (
	clientID = os.Getenv("ACUZIO_CLIENTID")
	secretID = os.Getenv("ACUZIO_SECRETID")
)

type (
	IndexTemplate struct {
		Pr         *ProductData
		NotVisited interface{}
	}
	CartProduct struct {
		Pr       Product
		Count    int
		Subtotal typePrice
	}

	CartProducts struct {
		Products []CartProduct
		Total    typePrice //`default:0`
	}

	typePrice int

	Product struct {
		Id          string    `bson:"_id" json:"id"`
		Name        string    `json:"name"`
		Currency    string    `json:"currency"`
		Price       typePrice `json:"price"`
		Discount    int       `json:"discount"`
		Title       string    `json:"title"`
		Description string    `json:"description"`
		Meta        string    `json:"meta"`
		Images      int       `json:"images"`
		// DiscountedPrice int
	}

	ProductData struct {
		Products []Product `json:"products"`
	}
)

func (t typePrice) String() string {
	v := int(t)
	return fmt.Sprintf("%d.%d", v/1000, (v/10)%100)
}

// Prices are in in thousands of a $,
// so we need to format the value of the price for humans
func (e Product) PriceWithDiscount() typePrice {
	discounted := (int(e.Price) * (100 - e.Discount) / 100)
	return typePrice(discounted)
}

func Index(c echo.Context) error {
	sess, err := session.Get(defaultSessionName, c)
	if err != nil {
		c.Logger().Error(err)
	}
	sess.Options = &sessions.Options{
		Path:     "/",
		MaxAge:   1800,
		HttpOnly: true,
	}

	bestProducts, err := loadData(bestProducts, "")
	var tmpl = IndexTemplate{
		Pr:         bestProducts,
		NotVisited: sess.Values["visited"],
	}
	sess.Values["visited"] = true
	sess.Save(c.Request(), c.Response().Writer)
	if err != nil {
		c.Logger().Error("error loading data: %v", err)
		return err
	}
	return c.Render(http.StatusOK, "indexTmpl", tmpl)
}

func (p *ProductData) CalculateDiscount() {
	// for _, e := range p.Products {
	// 	e.Title = "f2osd"
	// } jhkbkhbkhbk
}

func Shop(c echo.Context) error {
	order := c.QueryParam("order")
	var (
		orderedData *ProductData
		err         error
	)
	orderedData, err = loadData(products, order)
	if err != nil {
		c.Logger().Error(err)
		return err
	}
	return c.Render(http.StatusOK, "shopTmpl", orderedData)
}
func ProductDetails(c echo.Context) error {
	id := c.Param("id")
	shopProducts, err := loadData(products, "")

	if err != nil {
		c.Logger().Error("error Loading products: %v", err)
		return err
	}

	for _, e := range shopProducts.Products {
		if e.Id == id {
			if err != nil {
				c.Logger().Error("error injecting oldprice: ", err)
				return err
			}
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
	http.Redirect(c.Response(), c.Request(), "/cart", http.StatusMovedPermanently)
	return nil

}
func Contact(c echo.Context) error {
	return c.Render(http.StatusOK, "Contact", "")
}
func About(c echo.Context) error {
	return c.Render(http.StatusOK, "About", "")
}
func SendMessage(c echo.Context) error {
	email, name := c.FormValue("email"), c.FormValue("name")
	message, subject := c.FormValue("message"), c.FormValue("subject")

	title := subject + "    by: " + name + " email: " + email + "\n"

	coll := mongoClient.Database("Acuzio").Collection("Contacts")
	doc := bson.D{
		{Key: "title", Value: title},
		{Key: "body", Value: message},
		{Key: "name", Value: name},
		{Key: "email", Value: email},
	}

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

func PayPalCreateOrder(c echo.Context) error {
	// Initialize client
	p, err := paypal.NewClient(clientID, secretID, paypal.APIBaseLive)
	if err != nil {
		c.Logger().Error("Error creating paypal client at create order", err)
		return err
	}
	ctx, cancel := context.WithTimeout(c.Request().Context(), 10*time.Second)
	defer cancel()
	accessToken, err := p.GetAccessToken(ctx)
	_ = accessToken

	if err != nil {
		c.Logger().Error("Paypal client error:", err)
		return err
	}
	sess, err := session.Get(defaultSessionName, c)
	if err != nil {
		log.Error("error in creating a session")
		return err
	}
	cart, err := loadCart(sess)
	if err != nil {
		c.Logger().Error("error loading cart ar create order", cart)
		return err
	}
	var payer paypal.CreateOrderPayer
	var application paypal.ApplicationContext
	var Amount paypal.PurchaseUnitAmount

	Amount.Value = cart.Total.String()
	Amount.Currency = "USD"
	order, err := p.CreateOrder(ctx, paypal.OrderIntentCapture, []paypal.PurchaseUnitRequest{paypal.PurchaseUnitRequest{ReferenceID: "ref-id", Amount: &Amount}}, &payer, &application)
	if err != nil {
		c.Logger().Error(err)
		return err
	}
	return c.JSON(http.StatusOK, order)
}
func PayPalCaptureOrder(c echo.Context) error {
	// TODO: store payment information such as the transaction ID
	sess, err := session.Get(defaultSessionName, c)
	if err != nil {
		log.Error("error in creating a session")
		return err
	}
	cart, err := loadCart(sess)
	if err != nil {
		c.Logger().Error("error loading cart ar create order", cart)
		return err
	}
	orderID := c.Param("orderId")
	p, err := paypal.NewClient(clientID, secretID, paypal.APIBaseLive)
	if err != nil {
		c.Logger().Error("Client Paypal Error:", err)
		return err
	}
	ctx, cancel := context.WithTimeout(c.Request().Context(), 10*time.Second)
	defer cancel()
	accessToken, err := p.GetAccessToken(ctx)
	_ = accessToken
	if err != nil {
		c.Logger().Error("Token Paypal Error:", err)
		return err
	}
	order, err := p.GetOrder(ctx, orderID)
	if err != nil {
		c.Logger().Error(err)
		return err
	}
	items := make(map[string]int)
	for _, e := range cart.Products {
		items[e.Pr.Name] = e.Count
	}
	coll := mongoClient.Database("Acuzio").Collection("Orders")
	doc := bson.D{
		{Key: "time", Value: order.CreateTime},
		{Key: "transactionID", Value: order.ID},
		{Key: "status", Value: order.Status},
		{Key: "payerName", Value: order.Payer.Name},
		{Key: "payerEmail", Value: order.Payer.EmailAddress},
		{Key: "payerCountry", Value: order.Payer.Address.CountryCode},
		{Key: "shipping", Value: order.PurchaseUnits[0].Shipping},
		{Key: "items", Value: items},
	}
	coll.InsertOne(ctx, doc)
	capture, err := p.CaptureOrder(ctx, orderID, paypal.CaptureOrderRequest{})
	if err != nil {
		c.Logger().Error("Capture order Error:", err)
		return err
	}

	sess.Options.MaxAge = -1
	err = sess.Save(c.Request(), c.Response())
	if err != nil {
		c.Logger().Fatal("failed to delete session", err)
	}

	return c.JSON(http.StatusOK, capture)
}
func ThankYou(c echo.Context) error {
	return c.Render(http.StatusOK, "thankYou", "")
}

var (
	mongoUri    string = os.Getenv("ACUZIO_MONGOURI")
	mongoClient *mongo.Client
)

func connectMongo(ctx context.Context) (*mongo.Client, error) {
	ctx, cancel := context.WithTimeout(ctx, 10*time.Second)
	defer cancel()
	serverAPIOptions := options.ServerAPI(options.ServerAPIVersion1)
	clientOptions := options.Client().
		ApplyURI(mongoUri).
		SetServerAPIOptions(serverAPIOptions)
	client, err := mongo.Connect(ctx, clientOptions)
	if err != nil {
		return nil, err
	}
	return client, nil
}
func Subscribe(c echo.Context) error {
	email := c.FormValue("email")
	coll := mongoClient.Database("Acuzio").Collection("Emails")
	doc := bson.D{
		{Key: "email", Value: email},
	}

	ctx, cancel := context.WithTimeout(c.Request().Context(), 10*time.Second)
	defer cancel()
	result, err := coll.InsertOne(ctx, doc)
	if err != nil {
		c.Logger().Error("Cannot insert data in DB. Error: %v", err)
		return err
	}
	_ = result
	return c.Redirect(http.StatusMovedPermanently, "/")
}

func main() {
	globalCtx := context.Background()
	var err error

	if mongoClient, err = connectMongo(globalCtx); err != nil {
		panic(err)
	}
	type linkInfo struct {
		In   int
		Name string
	}
	e := echo.New()
	e.Renderer = &Template{
		templates: template.Must(template.New("t").Funcs(template.FuncMap{
			"Iterate": func(count int, name string) []linkInfo {
				var i int
				var Items []linkInfo
				for i = 1; i <= (count); i++ {
					Items = append(Items, linkInfo{
						In:   i,
						Name: name,
					})
				}
				return Items
			},
		}).ParseGlob("./static/*.html")),
	}

	e.Use(session.Middleware(sessions.NewCookieStore([]byte("secret"))))

	e.Static("/", "./static")
	e.GET("/shop", Shop)
	e.GET("/", Index)
	e.GET("/product/:id", ProductDetails)
	e.POST("/addToCart/:id", AddToCart)
	e.POST("/subscribe", Subscribe)
	e.POST("/sendMessage", SendMessage)

	e.GET("/about", About)
	e.GET("/contact", Contact)
	e.File("/submit-success", "static/submit-success.html")
	e.GET("/thank-you", ThankYou)
	//Chekout and payment
	e.GET("/checkout", Checkout)
	// //Paypal POST
	e.POST("/api/orders", PayPalCreateOrder)
	e.POST("/api/orders/capture/:orderId", PayPalCaptureOrder)

	//Cart Requests
	e.POST("/remove", Remove)
	e.GET("/cart", Cart)

	e.Logger.Fatal(e.Start(":" + getEnvDef("PORT", "8080")))

}

func getEnvDef(env string, def string) (res string) {
	res = os.Getenv(env)
	if res == "" {
		res = def
	}
	return
}
