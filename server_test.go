package main

import (
	"bufio"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"server/server"
	"strings"
	"testing"
)

func TestMain(m *testing.M) {
	fmt.Println("Setup")
	x := server.New()
	x.Get("/", func(c server.Context) error {
		return c.String(http.StatusAccepted, "Hola Mundo")
	})
	x.Get("/json", func(c server.Context) error {
		p := People{Name: "Juan", Age: 13}
		return c.Json(http.StatusAccepted, p)
	})
	x.Start(getPort())
	code := m.Run()
	fmt.Println("Ejecuto Todas Las Pruebas")
	x.Stop()
	os.Exit(code)
}

func TestShouldBeAccepted(t *testing.T) {
	fmt.Println("Ejecuto Prueba")
	res, _ := http.Get("http://localhost" + getPort())
	buf := new(strings.Builder)
	io.Copy(buf, res.Body)
	if strings.ReplaceAll(buf.String(), "\n", "") != "Hola Mundo" {
		t.Fail()
	}
}

func TestShouldBeGiveACorrectJSON(t *testing.T) {
	p := People{Name: "Juan", Age: 13}
	res, _ := http.Get("http://localhost" + getPort() + "/json")
	buf := bufio.NewReader(res.Body)
	bt := make([]byte, buf.Size())
	p1 := People{}
	buf.Read(bt)
	js := strings.Replace(string(bt), "\x00", "", -1)
	json.Unmarshal([]byte(js), &p1)
	if p != p1 {

		fmt.Println(string(bt))
		t.Fail()
	}

}

func getPort() string {
	return ":3032"
}
