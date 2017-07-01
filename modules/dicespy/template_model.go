package dicespy

import "github.com/therecipe/qt/core"

const (
	Title = int(core.Qt__UserRole) + 1<<iota
	Link
)

type TemplateModel struct {
	core.QAbstractListModel

	_ func() `constructor:"init"`

	_ map[int]*core.QByteArray `property:"roles"`
	_ []*Template              `property:"templates"`

	_ func(*Template)                    `slot:"addTemplate"`
	_ func(row int, title, linkn string) `slot:"editTemplate"`
	_ func(row int)                      `slot:"removeTemplate"`
}

type Template struct {
	core.QObject

	_ string `property:"title"`
	_ string `property:"link"`
}

func init() {
	Template_QRegisterMetaType()
}

func (m *TemplateModel) init() {
	m.SetRoles(map[int]*core.QByteArray{
		Title: core.NewQByteArray2("title", len("title")),
		Link:  core.NewQByteArray2("link", len("link")),
	})

	m.ConnectData(m.data)
	m.ConnectRowCount(m.rowCount)
	m.ConnectColumnCount(m.columnCount)
	m.ConnectRoleNames(m.roleNames)

	m.ConnectAddTemplate(m.addTemplate)
	m.ConnectEditTemplate(m.editTemplate)
	m.ConnectRemoveTemplate(m.removeTemplate)
}

func (m *TemplateModel) data(index *core.QModelIndex, role int) *core.QVariant {
	if !index.IsValid() {
		return core.NewQVariant()
	}

	if index.Row() >= len(m.Templates()) {
		return core.NewQVariant()
	}

	var p = m.Templates()[index.Row()]

	switch role {
	case Title:
		{
			return core.NewQVariant14(p.Title())
		}

	case Link:
		{
			return core.NewQVariant14(p.Link())
		}

	default:
		{
			return core.NewQVariant()
		}
	}
}

func (m *TemplateModel) rowCount(parent *core.QModelIndex) int {
	return len(m.Templates())
}

func (m *TemplateModel) columnCount(parent *core.QModelIndex) int {
	return 1
}

func (m *TemplateModel) roleNames() map[int]*core.QByteArray {
	return m.Roles()
}

func (m *TemplateModel) addTemplate(p *Template) {
	m.BeginInsertRows(core.NewQModelIndex(), len(m.Templates()), len(m.Templates()))
	m.SetTemplates(append(m.Templates(), p))
	m.EndInsertRows()
}

func (m *TemplateModel) editTemplate(row int, title string, link string) {
	var p = m.Templates()[row]

	if title != "" {
		p.SetTitle(title)
	}

	if link != "" {
		p.SetLink(link)
	}

	var pIndex = m.Index(row, 0, core.NewQModelIndex())
	m.DataChanged(pIndex, pIndex, []int{Title, Link})
}

func (m *TemplateModel) removeTemplate(row int) {
	m.BeginRemoveRows(core.NewQModelIndex(), row, row)
	m.SetTemplates(append(m.Templates()[:row], m.Templates()[row+1:]...))
	m.EndRemoveRows()
}
