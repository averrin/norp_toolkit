#!/bin/bash
echo -e ""
echo "Building NoRP Toolkit"
echo -e ""

mkdir -p deploy &> /dev/null

echo -e "Building modules…"
mkdir deploy/modules &> /dev/null
find ./modules/ -iname install.sh -exec bash {} \;
echo -e ""

echo -e "Building main app…"
export QT_DIR=/opt/Qt5.8.0
qtdeploy $1 build desktop .
cp -r deploy/modules deploy/linux
echo -e ""
