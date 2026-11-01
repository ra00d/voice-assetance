#!/bin/bash

# Colors
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
NC='\033[0m'

echo -e "${YELLOW}Setting up Whisper.cpp dependencies...${NC}"

# Create third_party directory
mkdir -p third_party
cd third_party

# Clone Whisper.cpp if not exists
if [ ! -d "whisper.cpp" ]; then
    echo "Cloning whisper.cpp..."
    git clone --depth 1 https://github.com/ggerganov/whisper.cpp.git
else
    echo "whisper.cpp already cloned."
fi

cd whisper.cpp

# Build using CMake
echo "Building whisper.cpp with CMake..."
cmake -B build -DWHISPER_BUILD_TESTS=OFF -DWHISPER_BUILD_EXAMPLES=OFF -DBUILD_SHARED_LIBS=OFF
cmake --build build --config Release

# Consolidate libraries and headers
echo "Consolidating libraries and headers..."
cd .. # Go back to third_party
mkdir -p lib include

# Copy libraries
find whisper.cpp/build -name "libwhisper.a" -exec cp {} lib/ \;
find whisper.cpp/build -name "libggml.a" -exec cp {} lib/ \;
find whisper.cpp/build -name "libggml-cpu.a" -exec cp {} lib/ \; 2>/dev/null || true
find whisper.cpp/build -name "libggml-base.a" -exec cp {} lib/ \; 2>/dev/null || true

# Copy headers
cp whisper.cpp/include/whisper.h include/
cp whisper.cpp/ggml/include/ggml.h include/
cp whisper.cpp/ggml/include/ggml-*.h include/ 2>/dev/null || true

# Get absolute paths
ROOT_DIR=$(pwd)
INCLUDE_PATH="$ROOT_DIR/include"
LIBRARY_PATH="$ROOT_DIR/lib"

echo -e "\n${GREEN}Whisper.cpp built and installed to third_party/lib and third_party/include.${NC}"
echo "To run the application, you need to set CGO environment variables."
echo "Run the following command to start Kind Echo:"
echo ""
echo "export C_INCLUDE_PATH=$INCLUDE_PATH"
echo "export LIBRARY_PATH=$LIBRARY_PATH"
echo "go run cmd/kind-echo/main.go"
echo ""
