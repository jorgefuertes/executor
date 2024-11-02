#!/bin/bash

set -e

# colors
RED='\033[0;31m'
GREEN='\033[0;32m'
NC='\033[0m' # No Color

# os detection
OS="$(uname -s | tr '[:upper:]' '[:lower:]')"
case "$OS" in
    "linux"|"darwin"|"windows")
        ;;
    *)
        echo -e "${RED}Unsupported operating system: $OS${NC}"
        exit 1
        ;;
esac

# arch detection
ARCH="$(uname -m)"
case "$ARCH" in
    x86_64|amd64)
        ARCH="amd64"
        ;;
    i386|i686)
        ARCH="386"
        ;;
    armv7*|armv8*)
        ARCH="arm"
        ;;
    aarch64|arm64)
        ARCH="arm64"
        ;;
    *)
        echo -e "${RED}Unsupported architecture: $ARCH${NC}"
        exit 1
        ;;
esac

# latest
echo "Reading latest version..."
LATEST_VERSION=$(curl -s https://api.github.com/repos/jorgefuertes/executor/releases/latest | grep "tag_name" | cut -d '"' -f 4)

if [ -z "$LATEST_VERSION" ]; then
    echo -e "${RED}Cannot get the latest version${NC}"
    exit 1
fi

# filename
FILENAME="executor_${OS}-${ARCH}_${LATEST_VERSION//./_}"
if [ "$OS" = "windows" ]; then
    FILENAME="${FILENAME}.zip"
else
    FILENAME="${FILENAME}.tar.gz"
fi

DOWNLOAD_URL="https://github.com/jorgefuertes/executor/releases/download/${LATEST_VERSION}/${FILENAME}"

# tmp
TMP_DIR=$(mktemp -d)
cd "$TMP_DIR"

# download
echo -n "Downloading $FILENAME..."
curl -Lso "$FILENAME" "$DOWNLOAD_URL" &> /dev/null
if [[ $? -ne 0 ]]
then
    echo -e "${RED}FAILED${NC}"
    exit 1
else
    echo -e "${GREEN}OK${NC}"
fi

# decrunch
echo "Extracting..."
if [ "$OS" = "windows" ]; then
    unzip "$FILENAME"
else
    tar xzf "$FILENAME"
fi

# install
echo "Installing to /usr/local/bin, may I need sudo password..."
sudo mv executor /usr/local/bin/
echo "Giving execution permission, may I need sudo password..."
sudo chmod +x /usr/local/bin/executor

# clean
echo "Cleaning"
cd - > /dev/null
rm -rf "$TMP_DIR"

echo -e "${GREEN}Installation complete!${NC}"
echo "Installed as /usr/local/bin/executor"
