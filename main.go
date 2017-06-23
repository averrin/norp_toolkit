package main

import (
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"path"
	"path/filepath"
	"plugin"

	"github.com/labstack/echo"
	"github.com/labstack/echo/middleware"
	"github.com/therecipe/qt/core"
	"github.com/therecipe/qt/gui"
	"github.com/therecipe/qt/qml"
	"github.com/therecipe/qt/quickcontrols2"
)

var plugins map[string]*plugin.Plugin
var root string

const pluginDir = "modules"

func main() {
	fmt.Println("Starting NoRP Toolkit")
	fmt.Println("")

	root, _ = filepath.Abs(filepath.Dir(os.Args[0]))
	plugins = map[string]*plugin.Plugin{}

	files, err := ioutil.ReadDir(path.Join(root, pluginDir))
	if err != nil {
		log.Println("No plugins found!")
	}
	model := NewPluginModel(nil)

	var h = NewPlugin(nil)
	h.SetTitle("Home")
	h.SetUisource("qrc:icon.png")
	h.SetIcon("qrc:qml/home.qml")
	model.AddPlugin(h)

	for _, file := range files {
		if file.IsDir() {
			pp := path.Join(root, pluginDir, file.Name(), file.Name()+".so")
			fmt.Printf("Loading: %v…\n", pp)
			p, err := plugin.Open(pp)
			if err == nil {
				name, err := p.Lookup("GetName")
				if err == nil {
					n := name.(func() string)()
					plugins[n] = p
					fmt.Printf("Plugin %v loaded.\n", n)
					d, _ := p.Lookup("GetDescription")
					fmt.Println(d.(func() string)())
					v, _ := p.Lookup("GetVersion")
					fmt.Printf("Version: %v\n\n", v.(func() string)())
					init, err := p.Lookup("Init")
					if err == nil {
						init.(func())()
					}

					var pl = NewPlugin(nil)
					pl.SetTitle(n)
					pl.SetUisource(fmt.Sprintf("file:modules/%v/%v.qml", n, n))
					pl.SetIcon(fmt.Sprintf("file:modules/%v/icon.png", n))
					model.AddPlugin(pl)
				}
			}
		}
	}

	e := echo.New()
	for _, p := range plugins {
		serve, err := p.Lookup("PrepeareServer")
		if err == nil {
			serve.(func(*echo.Echo))(e)
		}
	}

	e.Use(middleware.CORSWithConfig(middleware.CORSConfig{
		AllowOrigins: []string{"*"},
		AllowMethods: []string{echo.GET, echo.PUT, echo.POST, echo.DELETE},
	}))
	go e.Start(":1323")
	startUI(model)
}

func startUI(model *PluginModel) {
	core.QCoreApplication_SetAttribute(core.Qt__AA_EnableHighDpiScaling, true)
	gui.NewQGuiApplication(len(os.Args), os.Args)
	quickcontrols2.QQuickStyle_SetStyle("universal")
	view := qml.NewQQmlApplicationEngine(nil)

	view.RootContext().SetContextProperty("PluginModel", model)

	view.Load(core.NewQUrl3("qrc:/qml/main.qml", 0))

	gui.QGuiApplication_Exec()
}
