#!/bin/bash
export PATH=$PATH:/c/Go/bin:/c/msys64/mingw64/bin
export GOPATH=~/Documents/go
export QT_MSYS2=true
export QT_WEBKIT=true
export QT_VERSION=5.8.0

windres icon.rc -o icon_windows.syso
./build.sh $1
cp -r deploy/modules deploy/windows
cd deploy/windows/
cd -
