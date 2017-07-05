import QtQuick 2.7
import QtQuick.Controls 2.1
import QtQuick.Controls.Universal 2.0
import QtQuick.Layouts 1.3

Page {
    id: page
    Column {
        width: parent.width
        anchors.margins: 20

        Label {
            text: "<h1>DiceSpy plugin</h1><br>" +
                "This tool listen dice rolls in your roll20 game and render pretty html view for streaming in OBS."

            anchors.left: parent.left
            anchors.right: parent.right
            horizontalAlignment: Label.AlignHLeft
            wrapMode: Label.Wrap
        }

        GridLayout {
            columns: 2

            Switch {
                id: connectStatus
                checked: false
                text: "Connected"
                Connections {
                    target: diceSpy
                    onOffline: {
                        connectStatus.checked = false
                    }
                }
                onClicked: {
                    if (checked) {
                        diceSpy.serve()
                    } else {
                        diceSpy.disconnect()
                    }
                }
            }
            Item {
                Layout.fillHeight: true
                Layout.fillWidth: true
                Rectangle { anchors.fill: parent; color: "#ffffff" }
            }

            TextField {
                Layout.preferredWidth: 350
                readOnly: true
                text: injectScript
                selectByMouse: true
            }
            Row {
                Button {
                    width: 70
                    text: "Copy"
                    anchors.rightMargin: 4
                    onClicked: {
                        diceSpy.copyscript()
                    }
                }

                Button {
                    width: 70
                    text: "Roll"
                    enabled: connectStatus.checked
                    onClicked: {
                        diceSpy.roll()
                    }
                }
            }


            Label {
                text: "Available templates:"
                horizontalAlignment: Text.AlignRight
                Layout.minimumWidth: 350
            }
            ComboBox {
                id: tpl
                Layout.minimumWidth: 140
                textRole: "title"
                model: templateModel
            }

            Label {
                text: "Rolls in history:"
                horizontalAlignment: Text.AlignRight
                Layout.minimumWidth: 350
            }

            SpinBox {
                to: 20
                value: initHistorySize
                /* editable: true */
                onValueChanged: {
                    diceSpy.sethistory(value);
                }
            }


            TextField {
                Layout.preferredWidth: 350
                readOnly: true
                text: templateModel.data(templateModel.index(tpl.currentIndex, 0), Qt.UserRole + 2)
                selectByMouse: true
            }
            Row {
                Layout.preferredWidth: 140
                Button {
                    width: 70
                    text: "Copy"
                    anchors.rightMargin: 4
                    onClicked: {
                        diceSpy.copylink(templateModel.data(templateModel.index(tpl.currentIndex, 0), Qt.UserRole + 2))
                    }
                }
                Button {
                    width: 70
                    text: "View"
                    enabled: connectStatus.checked
                    onClicked: {
                        diceSpy.viewlink(templateModel.data(templateModel.index(tpl.currentIndex, 0), Qt.UserRole + 2))
                    }
                }
            }

        }
    }
}
