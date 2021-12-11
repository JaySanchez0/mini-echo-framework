package main

import (
	"app/server"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"strconv"
	"strings"
	"testing"
)

type Pet struct {
	Name string
}

var url string

func GetTestEndpoint(context server.Context) error {
	p := Pet{Name: "MePet"}
	return context.Json(http.StatusAccepted, &p)
}

func PostTestEndPoint(context server.Context) error {
	p := Pet{}
	context.Bind(&p)
	return context.Json(http.StatusAccepted, &p)
}

func GetQueryTestEndPoint(context server.Context) error {
	p := Pet{Name: context.Query["name"]}
	return context.Json(http.StatusAccepted, &p)
}

func GetParamTestEndPoint(context server.Context) error {
	p := Pet{Name: context.GetParam("name")}
	return context.Json(http.StatusAccepted, &p)
}

func TestMain(m *testing.M) {
	s := server.New()
	defer s.Stop()
	url = "http://localhost:" + strconv.Itoa(GetPort())
	s.Get("/pet", GetTestEndpoint)
	s.Get("/pet/:name", GetParamTestEndPoint)
	s.Post("/pet", PostTestEndPoint)
	s.Get("/mpet", GetQueryTestEndPoint)
	s.Start(GetPort())
	fmt.Println("Inicio correr pruebas")
	os.Exit(m.Run())
}

func TestShouldBeGetPet(t *testing.T) {
	res, _ := http.Get(url + "/pet")
	b, _ := io.ReadAll(res.Body)
	p := Pet{}
	if json.Unmarshal(b, &p) != nil || p.Name != "MePet" {
		t.Error("Invalid")
	}
}

func TestShouldBePost(t *testing.T) {
	cli := http.Client{}
	pt := Pet{Name: "Petti"}
	b, _ := json.Marshal(&pt)
	rq, _ := http.NewRequest("POST", url+"/pet", strings.NewReader(string(b)))
	res, _ := cli.Do(rq)
	bt, _ := io.ReadAll(res.Body)
	pt2 := Pet{}
	json.Unmarshal(bt, &pt)
	if json.Unmarshal(bt, &pt) != nil || pt2.Name != pt.Name {
		t.Error()
	}
}

func TestShouldBeMatchQuery(t *testing.T) {
	res, _ := http.Get(url + "/mpet?name=pablo")
	b, _ := io.ReadAll(res.Body)
	p := Pet{}
	if json.Unmarshal(b, &p) != nil || p.Name != "pablo" {
		t.Error()
	}
}

func TestShouldBeMapParams(t *testing.T) {
	res, _ := http.Get(url + "/pet/kev")
	b, _ := io.ReadAll(res.Body)
	p := Pet{}
	if json.Unmarshal(b, &p) != nil || p.Name != "kev" {
		t.Error()
	}

}

func GetPort() int {
	return 80
}
