import QtQuick 2.15
import QtQuick.Controls 2.15
import SddmQtComponents 0.1 as SDDM

SDDM.Theme {
    id: theme

    property int padding: 20
    property color bgColor: "#0a1628"
    property color accentColor: "#5dade2"
    property color textColor: "#e8f4fc"

    background: Rectangle {
        color: theme.bgColor
        Image {
            source: "/usr/share/roxyos/assets/roxy_lock.png"
            fillMode: Image.PreserveAspectCrop
            anchors.fill: parent
            opacity: 0.5
        }
    }

    login: Column {
        spacing: 20
        anchors.centerIn: parent
        width: 300

        SDDM.TextField {
            id: username
            width: parent.width
            height: 50
            placeholderText: qsTr("Username")
            font.family: "JetBrainsMono Nerd Font"
            font.pointSize: 14
            color: theme.textColor
            backgroundColor: Qt.rgba(0, 0, 0, 0.5)
            borderColor: theme.accentColor
            selectedColor: theme.accentColor
            selectionColor: Qt.rgba(93, 173, 226, 0.3)
        }

        SDDM.PasswordField {
            id: password
            width: parent.width
            height: 50
            placeholderText: qsTr("Password")
            font.family: "JetBrainsMono Nerd Font"
            font.pointSize: 14
            color: theme.textColor
            backgroundColor: Qt.rgba(0, 0, 0, 0.5)
            borderColor: theme.accentColor
            submitOnEnter: true
        }

        SDDM.Button {
            id: loginButton
            width: parent.width
            height: 50
            text: qsTr("Login")
            font.family: "JetBrainsMono Nerd Font"
            font.pointSize: 14
            backgroundColor: theme.accentColor
            textColor: theme.bgColor
            activeFocusOnTab: true
        }
    }

    session: Column {
        spacing: 10
        anchors {
            right: parent.right
            bottom: parent.bottom
            margins: theme.padding
        }

        SDDM.ComboBox {
            id: session
            width: 200
            height: 40
            font.family: "JetBrainsMono Nerd Font"
            font.pointSize: 12
            color: theme.textColor
            backgroundColor: Qt.rgba(0, 0, 0, 0.5)
            borderColor: theme.accentColor
        }
    }

    power: Column {
        spacing: 10
        anchors {
            left: parent.left
            bottom: parent.bottom
            margins: theme.padding
        }

        SDDM.Button {
            text: qsTr("Shutdown")
            font.family: "JetBrainsMono Nerd Font"
            backgroundColor: Qt.rgba(139, 74, 43, 0.8)
            textColor: theme.textColor
        }

        SDDM.Button {
            text: qsTr("Reboot")
            font.family: "JetBrainsMono Nerd Font"
            backgroundColor: Qt.rgba(26, 82, 118, 0.8)
            textColor: theme.textColor
        }
    }

    userList: Column {
        spacing: 10
        anchors {
            left: parent.left
            top: parent.top
            margins: theme.padding
        }

        Repeater {
            model: userModel
            delegate: SDDM.UserListItem {
                icon: model.icon
                name: model.name
                isCurrentUser: model.currentUser
                font.family: "JetBrainsMono Nerd Font"
                color: theme.textColor
            }
        }
    }
}
