package server

import (
	"encoding/json"
	"fmt"
	"net"
	"strconv"
	"strings"
)

type Context struct {
	con  net.Conn
	json string
}

func (context *Context) String(status int, data string) error {
	fmt.Fprintln(context.con, "HTTP/1.1 200 Ok\nContent-Type:text/plain\n\n"+data)
	fmt.Println(strconv.Itoa(status))
	context.con.Close()
	return nil
}

func (context *Context) Html(status int, data string) error {
	fmt.Fprintln(context.con, "HTTP/1.1 200 Ok\nContent-Type:text/html\n\n"+data)
	fmt.Println(strconv.Itoa(status))
	context.con.Close()
	return nil
}

func (context *Context) Json(status int, data interface{}) error {
	d, er := json.Marshal(data)
	if er == nil {
		fmt.Fprintln(context.con, "HTTP/1.1 200 Ok\nContent-Type:application/json\n\n"+string(d))
	} else {
		fmt.Fprintf(context.con, "HTTP/1.1 404 Not Found\n\n")
	}
	context.con.Close()
	return nil
}

func (context *Context) Bind(data interface{}) error {
	return json.Unmarshal([]byte(context.json), data)
}

type Oper struct {
	Method string
	Path   string
}

type Server struct {
	methods map[Oper]func(Context) error
}

func (server *Server) Get(path string, f func(Context) error) {
	server.methods[Oper{Method: "GET", Path: path}] = f
}

func (server *Server) Post(path string, f func(Context) error) {
	server.methods[Oper{Method: "POST", Path: path}] = f
}

func (server *Server) onListen(s net.Conn) {
	x := make([]byte, 1024)
	s.Read(x)
	fmt.Println(string(x))
	w := strings.Split(string(x), "\n")
	wi := strings.Split(w[0], " ")
	w[len(w)-1] = strings.Replace(w[len(w)-1], "\x00", "", -1)
	fmt.Println()
	var c Context
	if wi[0] != "GET" {
		c = Context{con: s, json: w[len(w)-1]}
	} else {
		c = Context{con: s}
	}
	f := server.methods[Oper{Method: wi[0], Path: wi[1]}]
	if f != nil {
		f(c)
	} else {
		fmt.Fprintf(s, "HTTP/1.1 404 Not Found\n\n")
		s.Close()
	}
}

func (server *Server) Start(port string) {
	s, e := net.Listen("tcp", port)
	if e == nil {
		for {
			fmt.Println("Esperando")
			client, e := s.Accept()
			if e == nil {
				go server.onListen(client)
			}
		}
	} else {
		fmt.Println("No se pudo iniciar")
	}
}

func New() *Server {
	x := &Server{methods: make(map[Oper]func(Context) error)}
	return x
}
