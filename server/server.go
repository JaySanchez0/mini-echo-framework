package server

import (
	"bufio"
	"encoding/json"
	"fmt"
	"net"
	"net/http"
	"strconv"
	"strings"
)

type Echo struct {
	handlers []Cond
	con      net.Conn
}

type Cond struct {
	Path   string
	Method string
	f      func(Context) error
}

type Context struct {
	Method      string
	Path        string
	Query       map[string]string
	Headers     map[string]string
	MatchingUrl string
	body        []byte
	con         net.Conn
}

func (context *Context) MatchPath(path string) bool {
	// User url
	p1 := strings.Split(context.Path, "/")
	// Url to match response
	p2 := strings.Split(path, "/")
	if len(p1) != len(p2) {
		return false
	}
	for i := 0; i < len(p1); i++ {
		//fmt.Println(": ok - " + p1[i])
		//fmt.Println("p1[i][0] - " + string(p2[i]))
		if p1[i] != p2[i] && p2[i][0] != ':' {
			fmt.Println("Resp 2")
			fmt.Println(p1[i])
			fmt.Println(p2[i])
			return false
		}
	}
	return true
}

func (context *Context) GetParam(name string) string {
	p1 := strings.Split(context.MatchingUrl, "/")
	p2 := strings.Split(context.Path, "/")
	for i := 0; i < len(p1); i++ {
		if p1[i] != p2[i] && p1[i][1:] == name {
			return p2[i]
		}
	}
	return ""
}

func (context *Context) Json(status int, data interface{}) error {
	b, e := json.Marshal(data)
	w := ""
	if e == nil {
		w = "HTTP/1.1 " + strconv.Itoa(status) + " Ok\nContent-Type:application/json\n\n" + string(b)
	} else {
		w = "HTTP/1.1 " + strconv.Itoa(http.StatusInternalServerError) + " InternalServerError\nContent-Type:application/json\n\n" + string(b)
	}
	//fmt.Println(w)
	fmt.Fprint(context.con, w)
	context.con.Close()
	return e
}

func (context *Context) Bind(obj interface{}) error {
	fmt.Println(string(context.body))
	return json.Unmarshal(context.body, obj)
}

func (echo *Echo) buildRequest(con net.Conn, headers string, body string) Context {
	splitHeaders := strings.Split(headers, "\n")
	first := strings.Split(splitHeaders[0], " ")
	headersMap := map[string]string{}
	for i := 1; i < len(splitHeaders); i++ {
		currentLine := splitHeaders[i]
		li := strings.Split(currentLine, ":")
		//fmt.Println(currentLine)
		//fmt.Print(li)
		if len(li) == 2 {
			headersMap[li[0]] = li[1]
		}
	}
	if len(first) == 3 {
		//fmt.Println("--- Start first ----------")
		//fmt.Println(first)
		//fmt.Println("--- ENd first ----------")
		uripath := strings.Split(first[1], "?")
		query := map[string]string{}
		if len(uripath) == 2 {
			q := strings.Split(uripath[1], "&")
			for _, w := range q {
				it := strings.Split(w, "=")
				query[it[0]] = it[1]
			}
		}
		return Context{
			Method:  first[0],
			Path:    uripath[0],
			Query:   query,
			Headers: headersMap,
			body:    []byte(body),
			con:     con,
		}
	}
	return Context{}
}

func (echo *Echo) processRequest(cli net.Conn) {
	reader := bufio.NewReader(cli)
	byteli := make([]byte, reader.Size()-1)
	cli.Read(byteli)
	res := string(byteli)
	// Separa el cuerpo y el body, y si encuentra espacios los limpia
	li := strings.Split(res, "\n\r\n\r")
	headerBody := make([]string, 2)
	// litmp contiene 0 - headers y 1 - body
	headerBody[0] = li[0]
	if len(li) >= 2 {
		headerBody[1] = "" //Body
		for i := 1; i < len(li); i++ {
			headerBody[1] = headerBody[1] + li[i] // Construye body
		}
		headerBody[1] = strings.Trim(strings.Trim(headerBody[1], " "), "\n\r")
	} else {
		headerBody[1] = ""
	}
	headerBody[1] = strings.Replace(headerBody[1], string('\x00'), "", -1)
	fmt.Println(headerBody[1])
	c := echo.buildRequest(cli, headerBody[0], headerBody[1])
	isInvoke := false
	for _, p := range echo.handlers {
		if p.Method == c.Method && c.MatchPath(p.Path) {
			c.MatchingUrl = p.Path
			isInvoke = true
			p.f(c)
		}
	}
	if !isInvoke && c.Method != "" {
		// Request no vacio
		w := "HTTP/1.1 " + strconv.Itoa(http.StatusNotFound) + " InternalServerError\nContent-Type:text/plain\n\n" + "404 not found"
		fmt.Fprint(c.con, w)
		c.con.Close()
	}
}

func (echo *Echo) Start(port int) {
	server, _ := net.Listen("tcp", ":80")
	echo.listenApp(server)
}

func (echo *Echo) listenApp(server net.Listener) {
	for {
		cli, _ := server.Accept()
		go echo.processRequest(cli)

	}
}

func (echo *Echo) Get(path string, f func(Context) error) {
	echo.handlers = append(echo.handlers, Cond{Method: "GET", Path: path, f: f})
}

func (echo *Echo) Post(path string, f func(Context) error) {
	echo.handlers = append(echo.handlers, Cond{Method: "POST", Path: path, f: f})
}

func (echo *Echo) Put(path string, f func(Context) error) {
	echo.handlers = append(echo.handlers, Cond{Method: "PUT", Path: path, f: f})
}

func (echo *Echo) Delete(path string, f func(Context) error) {
	echo.handlers = append(echo.handlers, Cond{Method: "DELETE", Path: path, f: f})
}

func (echo *Echo) Stop() {
	echo.con.Close()
}

func New() *Echo {
	return &Echo{}
}
