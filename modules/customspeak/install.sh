#!/bin/bash

pluginName="customspeak"

echo "* Building plugin: $pluginName"
# go build -buildmode=plugin -ldflags "-pluginpath=$pluginName -s" -o deploy/modules/$pluginName/$pluginName.so modules/$pluginName/*.go

mkdir -p deploy/modules/$pluginName
cp modules/$pluginName/README.md deploy/modules/$pluginName/
cp modules/$pluginName/*.png deploy/modules/$pluginName/
cp -r modules/$pluginName/qml deploy/modules/$pluginName/
