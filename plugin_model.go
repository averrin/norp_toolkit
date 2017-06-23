package main

import "github.com/therecipe/qt/core"

const (
	Title = int(core.Qt__UserRole) + 1<<iota
	UISource
	Icon
)

type PluginModel struct {
	core.QAbstractListModel

	_ func() `constructor:"init"`

	_ map[int]*core.QByteArray `property:"roles"`
	_ []*Plugin                `property:"plugins"`

	_ func(*Plugin)                               `slot:"addPlugin"`
	_ func(row int, title, uisource, icon string) `slot:"editPlugin"`
	_ func(row int)                               `slot:"removePlugin"`
}

type Plugin struct {
	core.QObject

	_ string `property:"title"`
	_ string `property:"uisource"`
	_ string `property:"icon"`
}

func init() {
	Plugin_QRegisterMetaType()
}

func (m *PluginModel) init() {
	m.SetRoles(map[int]*core.QByteArray{
		Title:    core.NewQByteArray2("title", len("title")),
		UISource: core.NewQByteArray2("uisource", len("uisource")),
		Icon:     core.NewQByteArray2("icon", len("icon")),
	})

	m.ConnectData(m.data)
	m.ConnectRowCount(m.rowCount)
	m.ConnectColumnCount(m.columnCount)
	m.ConnectRoleNames(m.roleNames)

	m.ConnectAddPlugin(m.addPlugin)
	m.ConnectEditPlugin(m.editPlugin)
	m.ConnectRemovePlugin(m.removePlugin)
}

func (m *PluginModel) data(index *core.QModelIndex, role int) *core.QVariant {
	if !index.IsValid() {
		return core.NewQVariant()
	}

	if index.Row() >= len(m.Plugins()) {
		return core.NewQVariant()
	}

	var p = m.Plugins()[index.Row()]

	switch role {
	case Title:
		{
			return core.NewQVariant14(p.Title())
		}

	case UISource:
		{
			return core.NewQVariant14(p.Uisource())
		}

	case Icon:
		{
			return core.NewQVariant14(p.Icon())
		}

	default:
		{
			return core.NewQVariant()
		}
	}
}

func (m *PluginModel) rowCount(parent *core.QModelIndex) int {
	return len(m.Plugins())
}

func (m *PluginModel) columnCount(parent *core.QModelIndex) int {
	return 1
}

func (m *PluginModel) roleNames() map[int]*core.QByteArray {
	return m.Roles()
}

func (m *PluginModel) addPlugin(p *Plugin) {
	m.BeginInsertRows(core.NewQModelIndex(), len(m.Plugins()), len(m.Plugins()))
	m.SetPlugins(append(m.Plugins(), p))
	m.EndInsertRows()
}

func (m *PluginModel) editPlugin(row int, title string, uisource string, icon string) {
	var p = m.Plugins()[row]

	if title != "" {
		p.SetTitle(title)
	}

	if uisource != "" {
		p.SetUisource(uisource)
	}

	if icon != "" {
		p.SetIcon(icon)
	}

	var pIndex = m.Index(row, 0, core.NewQModelIndex())
	m.DataChanged(pIndex, pIndex, []int{Title, UISource, Icon})
}

func (m *PluginModel) removePlugin(row int) {
	m.BeginRemoveRows(core.NewQModelIndex(), row, row)
	m.SetPlugins(append(m.Plugins()[:row], m.Plugins()[row+1:]...))
	m.EndRemoveRows()
}
