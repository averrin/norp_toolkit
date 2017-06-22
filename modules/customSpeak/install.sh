#!/bin/bash

pluginName="customSpeak"

echo "* Building plugin: $pluginName"
go build -buildmode=plugin -ldflags "-pluginpath=$pluginName -s" -o deploy/modules/$pluginName/$pluginName.so modules/$pluginName/*.go

cp modules/$pluginName/README.md deploy/modules/$pluginName/
cp modules/$pluginName/*.png deploy/modules/$pluginName/
cp modules/$pluginName/*.qml deploy/modules/$pluginName/
