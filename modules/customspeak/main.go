package customspeak

import "github.com/therecipe/qt/qml"

// import "C"

// Name is plugin ID
const Name string = "customspeak"

// Version is version
const Version string = "0.3.0"

//Description is description
const Description string = "Custom Discord layer for OBS"

type Plugin struct{}

// GetName is name getter
func (p Plugin) GetName() string {
	return Name
}

// GetVersion is version getter
func (p Plugin) GetVersion() string {
	return Version
}

// GetDescription is desc getter
func (p Plugin) GetDescription() string {
	return Description
}

func NewPlugin() *Plugin {
	return &Plugin{}
}

func (p Plugin) Init() error {
	return Init()
}

func (p Plugin) StartUI(view *qml.QQmlApplicationEngine) {
	StartUI(view)
}

func (p Plugin) Close() {
	Close()
}
