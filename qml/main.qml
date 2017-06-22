import QtQuick 2.7
import QtQuick.Controls 2.1
import QtQuick.Layouts 1.3

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
            Layout.fillWidth: true
            Layout.minimumWidth: 150
            id: listView
            currentIndex: -1
            anchors.fill: parent

            delegate: ItemDelegate {
                width: parent.width
                text: model.title
                highlighted: ListView.isCurrentItem
                onClicked: {
                    if (listView.currentIndex != index) {
                        listView.currentIndex = index
                        stackView.push(model.source)
                    }
                }
            }

            model: ListModel {
                ListElement { title: "diceSpy"; source: "file:modules/diceSpy/diceSpy.qml" }
                ListElement { title: "customSpeak"; source: "file:modules/customSpeak/customSpeak.qml" }
            }

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
                    text: "Qt Quick Controls 2 provides a set of controls that can be used to build complete interfaces in Qt Quick."
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
