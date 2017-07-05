import QtQuick 2.6
import QtQuick.Controls 2.0
import QtQuick.Layouts 1.3

Page {
    id: page
    Column{
        spacing: 2
        width: parent.width
        anchors.margins: 20

        Label {
            text: "<h1>CustomSpeak plugin</h1><br>" +
                "This tool listen voice channel your participating and draw per-user images when they speaks or not."

            anchors.left: parent.left
            anchors.right: parent.right
            horizontalAlignment: Label.AlignHLeft
            wrapMode: Label.Wrap
        }


        TextField {
            id: email
            placeholderText: "Discord email"
            width: 250
            text: initEmail
            anchors.bottomMargin: 4
        }

        TextField {
            id: password
            placeholderText: "Discord password"
            echoMode: TextInput.Password
            width: 250
            text: initPassword
            anchors.bottomMargin: 4
        }

        Label {
            id: connectError
            color: '#881111'
            visible: text != ""
            Connections {
                target: customSpeak
                onError: {
                    connectError.text = err
                }
            }
        }

        Switch {
            id: connectStatus
            checked: false
            enabled: email.text != "" && password.text != ""
            text: "Connected"
            Connections {
                target: customSpeak
                onOffline: {
                    connectStatus.checked = false
                }
            }
            onClicked: {
                if (checked) {
                    customSpeak.serve(email.text, password.text)
                } else {
                    customSpeak.disconnect()
                }
            }
        }

        Label {
            id: guildName
            color: "#111188"
            Connections {
                target: customSpeak
                onSetguild: {
                    guildName.text = "Guild: " + guild
                    guildName.visible = guild != ""
                }
            }
        }
        Label {
            id: channelName
            color: "#111188"
            Connections {
                target: customSpeak
                onSetchannel: {
                    channelName.text = "Channel: " + channel
                    channelName.visible = channel != ""
                }
            }
        }
    }
}
