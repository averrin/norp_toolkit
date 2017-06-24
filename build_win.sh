#!/bin/bash
export PATH=$PATH:/c/Go/bin:/c/msys64/mingw64/bin
export GOPATH=~/Documents/go
export QT_MSYS2=true

windres icon.rc -o icon_windows.syso
./build.sh $1
cp -r deploy/modules deploy/windows
cd deploy/windows/
rm -rf Qt5Widgets.dll
find . -maxdepth 1 -type d ! -name modules -exec rm -rf {} +
cd -
