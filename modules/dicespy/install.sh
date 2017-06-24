#!/bin/bash

pluginName="dicespy"

echo "* Building plugin: $pluginName"
# go build -buildmode=plugin -ldflags "-pluginpath=$pluginName -s" -o deploy/modules/$pluginName/$pluginName.so modules/$pluginName/*.go

mkdir -p deploy/modules/$pluginName
cp modules/$pluginName/*.png deploy/modules/$pluginName/
cp modules/$pluginName/README.md deploy/modules/$pluginName/
cp -r modules/$pluginName/templates deploy/modules/$pluginName/
cp modules/$pluginName/payload.js deploy/modules/$pluginName/
cp modules/$pluginName/config.yml deploy/modules/$pluginName/
cp -r modules/$pluginName/qml deploy/modules/$pluginName/
