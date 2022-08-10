package main

import (
	"fmt"
	"net/http"
	"time"

	"github.com/labstack/echo/v4"
)

type order struct {
	ID    int `json:"id"`
	Count int `json:"count"`
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

func confirm(w http.ResponseWriter, r *http.Request) {
	if err := r.ParseForm(); err != nil {
		http.Error(w, fmt.Sprintf("Error: %v", err), http.StatusInternalServerError)
		return
	}
	fmt.Fprintf(w, "Website post: %v\n", r.PostForm)
	fmt.Printf("Quantity of wathces ordered: %v\n", r.FormValue("count"))
	return
}
func conf(c echo.Context) error {
	count := c.FormValue("count")

	id := "1"
	fmt.Printf("Quantity of wathces ordered: %v\n", count)
	writeCookie(c, "count", count)
	writeCookie(c, "id", id)
	return c.
}

func main() {
	fileServer := http.FileServer(http.Dir("./static"))
	http.Handle("/", fileServer)

	e := echo.New()
	e.POST("/confirm", conf)

	e.Static("/", "static")

	e.Logger.Fatal(e.Start(":8080"))

}
