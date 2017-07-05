package dicespy

import (
	"encoding/json"
	"fmt"
	template "html/template"
	"io"
	"io/ioutil"
	"log"
	// "os"
	"path"
	"strconv"

	yaml "gopkg.in/yaml.v1"

	"net/http"
	"strings"

	"github.com/atotto/clipboard"
	"github.com/jinzhu/configor"
	"github.com/labstack/echo"
	"github.com/labstack/echo/middleware"
	"github.com/therecipe/qt/core"
	"github.com/therecipe/qt/qml"
	"github.com/skratchdot/open-golang/open"
	// "github.com/therecipe/qt/webengine"
	// "github.com/therecipe/qt/widgets"
	"github.com/wirepair/gcd"
	// "github.com/wirepair/gcd/gcdapi"
	"golang.org/x/sys/windows/registry"
	"golang.org/x/net/websocket"
)

const avatarRoot string = "https://app.roll20.net"
const root string = "modules/dicespy"
const port string = "1323"
// const injectScript = "$.getScript('http://127.0.0.1:" + port + "/script')"
var injectScript string

var config = ConfigStruct{}
var rolls []*Roll
var players map[string]string

var e *echo.Echo
var bridge *DsMocBridge

//go:generate qtmoc
type DsMocBridge struct {
	core.QObject

	_ func()            `signal:"serve"`
	_ func()            `signal:"disconnect"`
	_ func()            `signal:"offline"`
	_ func(err string)  `signal:"error"`
	_ func()            `signal:"copyscript"`
	_ func(link string) `signal:"copylink"`
	_ func(link string) `signal:"viewlink"`
	_ func()            `signal:"roll"`
	_ func(size int)    `signal:"sethistory"`
}

func saveConfig() {
	d, _ := yaml.Marshal(&config)
	err := ioutil.WriteFile(path.Join(root, "config.yml"), d, 0644)
	if err != nil {
		fmt.Println(err)
	}
}

func StartUI(view *qml.QQmlApplicationEngine) {
	bridge = NewDsMocBridge(nil)
	bridge.ConnectServe(func() {
		err := Serve()
		if err != nil {
			bridge.Error(fmt.Sprintf("Connection error: %v", err))
			bridge.Offline()
		} else {
			bridge.Error("")
		}
	})
	bridge.ConnectDisconnect(func() {
		e.Close()
	})

	bridge.ConnectCopyscript(func() {
		clipboard.WriteAll(injectScript)
	})
	bridge.ConnectCopylink(func(link string) {
		clipboard.WriteAll(link)
	})
	bridge.ConnectSethistory(func(size int) {
		config.HistoryCount = size
		saveConfig()
	})
	bridge.ConnectViewlink(func(link string) {
		// widgets.NewQApplication(len(os.Args), os.Args)
		// var window = widgets.NewQMainWindow(nil, 0)

		// var centralWidget = widgets.NewQWidget(nil, 0)
		// centralWidget.SetLayout(widgets.NewQVBoxLayout())

		// var wview = webengine.NewQWebEngineView(nil)
		// wview.Load(core.NewQUrl3(link, 0))
		// centralWidget.Layout().AddWidget(wview)

		// var rbutton = widgets.NewQPushButton2("Reload", nil)
		// rbutton.ConnectClicked(func(checked bool) {
		// 	wview.Reload()
		// })
		// centralWidget.Layout().AddWidget(rbutton)
		// var rollButton = widgets.NewQPushButton2("Test roll", nil)
		// rollButton.ConnectClicked(func(checked bool) {
		// 	processRoll(getTestRoll())
		// })
		// centralWidget.Layout().AddWidget(rollButton)

		// window.SetCentralWidget(centralWidget)
		// window.Show()
		open.Run(link)
	})
	bridge.ConnectRoll(func() {
		processRoll(getTestRoll())
	})

	view.RootContext().SetContextProperty("diceSpy", bridge)
	// view.RootContext().SetContextProperty2("injectScript", core.NewQVariant14(injectScript))
	view.RootContext().SetContextProperty2("initHistorySize", core.NewQVariant14(strconv.Itoa(config.HistoryCount)))
	files, err := ioutil.ReadDir(path.Join(root, "templates"))
	if err != nil {
		log.Fatal(err)
	}

	model := NewTemplateModel(nil)
	for _, file := range files {
		if strings.HasSuffix(file.Name(), ".html") && !strings.Contains(file.Name(), "_content") {
			t := NewTemplate(nil)
			t.SetTitle(file.Name())
			t.SetLink(fmt.Sprintf("http://127.0.0.1:%v/display/%v", port, strings.Replace(file.Name(), ".html", "", 1)))
			model.AddTemplate(t)

		}
	}
	view.RootContext().SetContextProperty("templateModel", model)

}

func Init() error {
	return configor.Load(&config, path.Join(root, "config.yml"))
}

func (t *MyTemplate) Render(w io.Writer, name string, data interface{}, c echo.Context) error {
	templates := template.Must(template.ParseGlob(path.Join(root, "templates/*.html")))
	return templates.ExecuteTemplate(w, name, data)
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

func processRoll(roll *Roll) {
	configor.Load(&config, path.Join(root, "config.yml"))
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

}


func Serve() error {

	e = echo.New()
	e.Use(middleware.CORSWithConfig(middleware.CORSConfig{
		AllowOrigins: []string{"*"},
		AllowMethods: []string{echo.GET, echo.PUT, echo.POST, echo.DELETE},
	}))
	t := &MyTemplate{}
	e.Renderer = t
	e.File("/script", "payload.js")
	e.GET("/display/:name", result)
	e.GET("/ws", wsHandler)
	e.Static("/templates", path.Join(root, "templates"))

	e.POST("/players", func(c echo.Context) error {
		readPlayers(c.Request())
		fmt.Println(players)
		return c.String(http.StatusOK, "OK")
	})

	e.POST("/roll", func(c echo.Context) error {
		roll := readRoll(c.Request())
		processRoll(roll)
		return c.String(http.StatusOK, "OK")
	})
	go e.Start(":" + port)
	go openRoll20()
	return nil
}

func openRoll20() {
	k, err := registry.OpenKey(registry.LOCAL_MACHINE, `SOFTWARE\Microsoft\Windows\CurrentVersion\App Paths\chrome.exe`, registry.QUERY_VALUE)
	if err != nil {
		log.Fatal(err)
	}
	defer k.Close()

	s, _, err := k.GetStringValue("Path")
	if err != nil {
		log.Fatal(err)
	}

	debugger := gcd.NewChromeDebugger()
	debugger.AddFlags([]string{"--allow-running-insecure-content"})
	// start process, specify a tmp profile path so we get a fresh profiled browser
	// set port 9222 as the debug port
	debugger.StartProcess(s+"\\chrome.exe", "", "9222")
	// defer debugger.ExitProcess()          // exit when done
	targets, err := debugger.GetTargets() // get the 'targets' or tabs/background processes
	if err != nil {
		log.Fatalf("error getting targets: %s\n", err)
	}
	target := targets[0] // take the first one

	// get the Page API and enable it
	if _, err := target.Page.Enable(); err != nil {
		log.Fatalf("error getting page: %s\n", err)
	}
	ret, err := target.Page.Navigate("http://roll20.net", "", "") // navigate
	if err != nil {
		log.Fatalf("Error navigating: %s\n", err)
	}
	log.Printf("ret: %#v\n", ret)
}

func getResult(roll *Roll, t string) RollResult {
	for _, result := range roll.Rolls {
		if result.Type == t {
			return result
		}
	}
	return RollResult{}
}

func renderRoll(roll *Roll) string {
	results := getResult(roll, "R").Results
	roll.Results = results
	roll.Skill = strings.TrimSpace(getResult(roll, "C").Text)
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

	modResult := getResult(roll, "M")
	if modResult.Type != "" {
		roll.Mod = strings.TrimSpace(modResult.Expr)
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
