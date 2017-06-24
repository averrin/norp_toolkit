package plugin

import "github.com/therecipe/qt/qml"

type PluginInterface interface {
	GetName() string
	GetVersion() string
	GetDescription() string
	Init() error
	StartUI(*qml.QQmlApplicationEngine)
}
