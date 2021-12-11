package main

import (
	"app/server"
	"fmt"
	"net/http"
)

type People struct {
	Name string
	Test string
}

func main() {
	ser := server.New()
	ser.Get("/notes", func(c server.Context) error {
		p := People{Name: c.Query["name"], Test: c.Query["age"]}
		fmt.Println(c.Query)
		return c.Json(http.StatusAccepted, &p)
	})

	ser.Get("/notes/:id", func(c server.Context) error {
		p := People{Name: c.GetParam("id")}
		fmt.Println(c.Query)
		return c.Json(http.StatusAccepted, &p)
	})

	ser.Post("/notes", func(c server.Context) error {
		p := People{}
		c.Bind(&p)
		return c.Json(http.StatusAccepted, &p)
	})
	ser.Start(80)
}
