#!/bin/bash

pluginName="customSpeak"

echo "* Building plugin: $pluginName"
go build -buildmode=plugin -ldflags "-pluginpath=$pluginName -s" -o deploy/modules/$pluginName/$pluginName.so modules/$pluginName/*.go

cp modules/$pluginName/icon.png deploy/modules/$pluginName/
cp modules/$pluginName/README.md deploy/modules/$pluginName/
cp modules/$pluginName/on.png deploy/modules/$pluginName/
cp modules/$pluginName/off.png deploy/modules/$pluginName/
