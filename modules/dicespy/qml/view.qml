import QtQuick 2.2
import QtQuick.Controls 1.1
import QtWebView 1.1
import QtQuick.Layouts 1.1
import QtQuick.Controls.Styles 1.2

ApplicationWindow {
    visible: true
    width: 520
    height: 480
    title: webView.title
    toolBar: ToolBar {
        id: navigationBar
        RowLayout {
            anchors.fill: parent
            spacing: 0

            ToolButton {
                id: reloadButton
                tooltip: webView.loading ? qsTr("Stop"): qsTr("Refresh")
                iconSource: webView.loading ? "images/stop-32.png" : "images/refresh-32.png"
                onClicked: webView.loading ? webView.stop() : webView.reload()
                Layout.preferredWidth: navigationBar.height
                style: ButtonStyle {
                    background: Rectangle { color: "transparent" }
                }
            }
            Item { Layout.preferredWidth: 5 }
            ToolButton {
                id: rollButton
                tooltip: 'Test roll'
                iconSource: "images/roll-32.png"
                onClicked: {
                    diceSpy.roll();
                }
                Layout.preferredWidth: navigationBar.height
                style: ButtonStyle {
                    background: Rectangle { color: "transparent" }
                }
            }
            Item { Layout.preferredWidth: 10 }
        }
    }
    WebView {
        id: webView
        anchors.fill: parent
        url: templateLink
        onLoadingChanged: {
            if (loadRequest.errorString)
                console.error(loadRequest.errorString);
        }
    }
}
