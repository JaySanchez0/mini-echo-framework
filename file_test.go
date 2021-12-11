package main

import (
	"app/server"
	"encoding/json"
	"io"
	"net/http"
	"testing"
)

type Pet struct {
	Name string
}

func GetTestEndpoint(context server.Context) error {
	p := Pet{Name: "MePet"}
	return context.Json(http.StatusAccepted, &p)
}

func TestMain(t *testing.M) {
	s := server.New()
	s.Get("/pet", GetTestEndpoint)
	s.Start(80)
	t.Run()
	s.Stop()
}

func TestSHouldBeGetPet(t *testing.T) {
	res, _ := http.Get("http://localhost/pet")
	b, _ := io.ReadAll(res.Body)
	p := Pet{}
	if json.Unmarshal(b, &p) != nil || p.Name != "MePet" {
		t.Error("Invalid")
	}
}
