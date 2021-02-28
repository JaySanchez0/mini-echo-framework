package main

import (
	"fmt"
	"net/http"
	"server/server"
)

type People struct {
	Name string `json:"name"`
	Age  int    `json:"age"`
}

func main() {
	x := server.New()
	x.Get("/", func(c server.Context) error {
		return c.String(http.StatusAccepted, "Hola")
	})
	x.Get("/html", func(c server.Context) error {
		return c.Html(http.StatusAccepted, "<h1>Hola</h1>")
	})
	x.Get("/json", func(c server.Context) error {
		x := []int{1, 2, 3}
		return c.Json(http.StatusAccepted, x)
	})
	x.Post("/post", func(c server.Context) error {
		p := People{}
		e := c.Bind(&p)
		if e == nil {
			fmt.Println("Convirtio")
			fmt.Println("Name")
			fmt.Println(p.Name)
			fmt.Println("Age")
			fmt.Println(p.Age)
		} else {
			fmt.Println(e.Error())
		}

		return c.Json(http.StatusAccepted, p)
	})
	x.Start(":80")
}
