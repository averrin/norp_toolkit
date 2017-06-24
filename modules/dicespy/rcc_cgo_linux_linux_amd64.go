package dicespy

/*
#cgo CFLAGS: -pipe -O2 -Wall -W -D_REENTRANT -fPIC -DQT_NO_DEBUG -DQT_GUI_LIB -DQT_CORE_LIB
#cgo CXXFLAGS: -pipe -O2 -std=gnu++11 -Wall -W -D_REENTRANT -fPIC -DQT_NO_DEBUG -DQT_GUI_LIB -DQT_CORE_LIB
#cgo CXXFLAGS: -I../../modules -I. -I/opt/Qt5.8.0/5.8/gcc_64/include -I/opt/Qt5.8.0/5.8/gcc_64/include/QtGui -I/opt/Qt5.8.0/5.8/gcc_64/include/QtCore -I. -isystem /usr/include/libdrm -I/opt/Qt5.8.0/5.8/gcc_64/mkspecs/linux-g++
#cgo LDFLAGS: -Wl,-O1 -Wl,-rpath,/opt/Qt5.8.0/5.8/gcc_64/lib
#cgo LDFLAGS:  -L/opt/Qt5.8.0/5.8/gcc_64/lib -lQt5Gui -lQt5Core -lGL -lpthread
*/
import "C"
