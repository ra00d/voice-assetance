#!/bin/bash

# ONNX Runtime version
VERSION="1.22.0"
OS="linux"
ARCH="x64"
FILENAME="onnxruntime-${OS}-${ARCH}-${VERSION}"
URL="https://github.com/microsoft/onnxruntime/releases/download/v${VERSION}/${FILENAME}.tgz"

echo "Downloading ONNX Runtime v${VERSION}..."
wget -q --show-progress "$URL"

echo "Extracting..."
tar -xzf "${FILENAME}.tgz"

echo "Setting up directory..."
mkdir -p third_party/onnxruntime
cp -r "${FILENAME}/lib" third_party/onnxruntime/
cp -r "${FILENAME}/include" third_party/onnxruntime/

# Cleanup
rm -rf "${FILENAME}"
rm "${FILENAME}.tgz"

echo "ONNX Runtime setup complete."
echo "Library location: $(pwd)/third_party/onnxruntime/lib/libonnxruntime.so.${VERSION}"
