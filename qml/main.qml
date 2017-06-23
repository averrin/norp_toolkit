import QtQuick 2.7
import QtQuick.Controls 2.1
import QtQuick.Layouts 1.3
import QtQuick.Controls.Universal 2.0

ApplicationWindow {
    id: window

    visible: true
    title: "Hello World Example"
    minimumWidth: 400
    minimumHeight: 600

    RowLayout {
        anchors.fill: parent
        ListView {

            Layout.fillHeight: true
            width: 64
            id: listView
            currentIndex: -1

            delegate: ItemDelegate {
                width: parent.width
                height: parent.width
                highlighted: ListView.isCurrentItem
                onClicked: {
                    if (listView.currentIndex != index) {
                        listView.currentIndex = index
                        stackView.push(model.uisource)
                    }
                }
                contentItem: Image {
                    /* fillMode: Image.Pad */
                    width: parent.width
                    height: parent.width
                    fillMode: Image.PreserveAspectFit
                    horizontalAlignment: Image.AlignHCenter
                    verticalAlignment: Image.AlignVCenter
                    source: model.icon
                }
            }

            model: PluginModel
            ScrollIndicator.vertical: ScrollIndicator { }
        }

        StackView {
            Layout.fillHeight: true
            Layout.fillWidth: true
            Layout.minimumWidth: 250
            id: stackView

            initialItem: Pane {
                id: pane

                Label {
                    text: "NoRoleplaying Toolkit provides a set of utilites for making your roll20 streams much better."
                    anchors.margins: 20
                    anchors.left: parent.left
                    anchors.right: parent.right
                    horizontalAlignment: Label.AlignHCenter
                    verticalAlignment: Label.AlignVCenter
                    wrapMode: Label.Wrap
                }

            }
        }
    }
}
