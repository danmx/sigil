#!/bin/bash

set -euo pipefail
IFS=$'\n\t'

if [[ $(uname -s) = *Darwin* ]]; then
  curl "https://s3.amazonaws.com/session-manager-downloads/plugin/latest/mac/sessionmanager-bundle.zip" -o "sessionmanager-bundle.zip"
  unzip sessionmanager-bundle.zip
  sudo ./sessionmanager-bundle/install -i /usr/local/sessionmanagerplugin -b /usr/local/bin/session-manager-plugin
  ./sessionmanager-bundle/install -h
else
  curl "https://s3.amazonaws.com/session-manager-downloads/plugin/latest/ubuntu_64bit/session-manager-plugin.deb" -o "session-manager-plugin.deb"
  sudo dpkg -i session-manager-plugin.deb
  session-manager-plugin
fi
