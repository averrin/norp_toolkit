package dicespy

import (
	"encoding/json"
	"fmt"
	"html/template"
	"io"
	"io/ioutil"
	"path"

	"github.com/jinzhu/configor"
	"github.com/labstack/echo"
	"github.com/labstack/echo/middleware"
	"github.com/therecipe/qt/qml"
	// "log"
	"net/http"
	"strings"

	"golang.org/x/net/websocket"
)

const avatarRoot string = "https://app.roll20.net"
const root string = "modules/dicespy"

var config = ConfigStruct{}
var rolls []*Roll
var players map[string]string

func StartUI(view *qml.QQmlApplicationEngine) {
}

func Init() error {
	return configor.Load(&config, path.Join(root, "config.yml"))
}

func (t *Template) Render(w io.Writer, name string, data interface{}, c echo.Context) error {
	return t.templates.ExecuteTemplate(w, name, data)
}

func result(c echo.Context) error {
	return c.Render(http.StatusOK, c.Param("name"), struct {
		Rolls  []*Roll
		Config ConfigStruct
	}{rolls, config})
}

var socket *websocket.Conn

func wsHandler(c echo.Context) error {
	websocket.Handler(func(ws *websocket.Conn) {
		socket = ws
		defer ws.Close()
		for {
			websocket.Message.Receive(ws, nil)
		}
	}).ServeHTTP(c.Response(), c.Request())
	return nil
}

func Serve() {

	e := echo.New()
	e.Use(middleware.CORSWithConfig(middleware.CORSConfig{
		AllowOrigins: []string{"*"},
		AllowMethods: []string{echo.GET, echo.PUT, echo.POST, echo.DELETE},
	}))
	go e.Start(":1323")
	t := &Template{
		templates: template.Must(template.ParseGlob(path.Join(root, "templates/*.html"))),
	}
	e.Renderer = t
	e.File("/script", "payload.js")
	e.GET("/display/:name", result)
	// e.GET("/ws", wsHandler)
	e.Static("/templates", "templates")

	e.POST("/players", func(c echo.Context) error {
		readPlayers(c.Request())
		fmt.Println(players)
		return c.String(http.StatusOK, "OK")
	})

	e.POST("/roll", func(c echo.Context) error {
		configor.Load(&config, "config.yml")
		roll := readRoll(c.Request())
		fmt.Println(roll)
		for len(rolls) >= config.HistoryCount {
			rolls = rolls[1:]
		}
		rolls = append(rolls, roll)
		message := ""
		for _, r := range rolls {
			r.Message = renderRoll(r)
			message += r.Message + "\n\n"
		}

		ioutil.WriteFile("roll.txt",
			[]byte(message), 0644)

		if socket != nil {
			websocket.Message.Send(socket, "Hello, Client!")
		}
		return c.String(http.StatusOK, "OK")
	})
	fmt.Println("")
	fmt.Println("-------")
	fmt.Println("")
	fmt.Println("Exec `$.getScript('http://127.0.0.1:1323/script');` in roll20.net WebInspector console")
	fmt.Println("Use `http://127.0.0.1:1323/display/basic` as OBS BrowserSource")
	fmt.Println("")
	fmt.Println("-------")
	fmt.Println("")
	// http.ListenAndServe(":1323", handler)
}

func renderRoll(roll *Roll) string {
	results := roll.Rolls[0].Results
	roll.Results = results
	roll.Skill = strings.TrimSpace(roll.Rolls[len(roll.Rolls)-1].Text)
	message := fmt.Sprintf("%v:", roll.Player)
	if roll.Skill != "" {
		message += fmt.Sprintf("\n%v", roll.Skill)
	}
	message += "\n("
	for i, d := range results {
		if i < len(results)-1 {
			message += fmt.Sprintf("%v, ", d.V)
		} else {
			message += fmt.Sprintf("%v", d.V)
		}
	}
	message += ")"

	if len(roll.Rolls) >= 3 {
		roll.Mod = strings.TrimSpace(roll.Rolls[len(roll.Rolls)-2].Expr)
		if roll.Mod != "" {
			message += fmt.Sprintf(" %v", roll.Mod)
		}
	}
	message += fmt.Sprintf(" = %v", roll.Total)
	return message
}

func readRoll(req *http.Request) *Roll {
	decoder := json.NewDecoder(req.Body)
	var rw RollWrapper
	err := decoder.Decode(&rw)
	if err != nil {
		panic(err)
	}
	defer req.Body.Close()
	var r Roll
	err = json.Unmarshal([]byte(rw.D.Content), &r)
	r.Player = players[rw.D.Playerid]
	r.OrigRoll = rw.D.OrigRoll
	r.Avatar = fmt.Sprintf("%v/users/avatar/%v/200", avatarRoot, strings.Split(rw.D.Avatar, "/")[3])

	return &r
}

func readPlayers(req *http.Request) {
	decoder := json.NewDecoder(req.Body)
	decoder.Decode(&players)
}
