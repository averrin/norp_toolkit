import QtQuick 2.7
import QtQuick.Controls 2.1
import QtQuick.Controls.Universal 2.0

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
        Row {
            width: parent.width
            TextField {
                width: 350
                readOnly: true
                text: injectScript
                selectByMouse: true
            }
            Button {
                text: "Copy"
                onClicked: {
                    diceSpy.copyscript()
                }
            }
        }

        Row {
            width: parent.width
            Label {
                text: "Available templates:"
            }
            ComboBox {
                id: tpl
                textRole: "text"
                model: ListModel{
                    id: tplModel
                    ListElement {text: "basic.html"; link: "http://127.0.0.1:1323/display/basic"}
                    ListElement {text: "complex.html"; link: "http://127.0.0.1:1323/display/complex"}
                }
            }
        }
        Row {
            width: parent.width
            TextField {
                width: 350
                readOnly: true
                text: tplModel.get(tpl.currentIndex).link
                selectByMouse: true
            }
            Button {
                text: "Copy"
                onClicked: {
                    diceSpy.copylink(tplModel.get(tpl.currentIndex).link)
                }
            }
        }
    }
}
