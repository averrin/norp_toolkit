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
# export QT_DIR=/home/alexeynabrodov/Qt5.9.1/
# export QT_VERSION=5.9.1
export QT_WEBKIT=true
qtdeploy $1 build desktop .
cp -r deploy/modules deploy/linux
echo -e ""
