package main

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/averrin/norp_toolkit/modules/customspeak"
	"github.com/averrin/norp_toolkit/modules/dicespy"
	"github.com/averrin/norp_toolkit/plugin"

	"github.com/therecipe/qt/core"
	"github.com/therecipe/qt/gui"
	"github.com/therecipe/qt/qml"
)

var plugins []plugin.PluginInterface
var root string

const pluginDir = "modules"

func main() {
	fmt.Println("Starting NoRP Toolkit")
	fmt.Println("")

	root, _ = filepath.Abs(filepath.Dir(os.Args[0]))
	dc := dicespy.NewPlugin()
	cs := customspeak.NewPlugin()
	plugins = []plugin.PluginInterface{
		dc,
		cs,
	}

	model := NewPluginModel(nil)

	for _, plugin := range plugins {
		n := plugin.GetName()
		err := plugin.Init()
		if err == nil {
			var pl = NewPlugin(nil)
			pl.SetTitle(n)
			pl.SetUisource(fmt.Sprintf("qrc:/qml/%v.qml", n))
			pl.SetIcon(fmt.Sprintf("file:%v/%v/%v/icon.png", root, pluginDir, n))
			model.AddPlugin(pl)
		}
	}

	startUI(model)
}

func startUI(model *PluginModel) {
	core.QCoreApplication_SetAttribute(core.Qt__AA_EnableHighDpiScaling, true)
	gui.NewQGuiApplication(len(os.Args), os.Args)
	// quickcontrols2.QQuickStyle_SetStyle("universal")
	view := qml.NewQQmlApplicationEngine(nil)

	view.RootContext().SetContextProperty("PluginModel", model)

	view.Load(core.NewQUrl3("qrc:/qml/main.qml", 0))

	for _, plugin := range plugins {
		plugin.StartUI(view)
	}

	gui.QGuiApplication_Exec()
}
