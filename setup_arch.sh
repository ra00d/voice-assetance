#!/bin/bash

# Colors
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
NC='\033[0m' # No Color

echo -e "${YELLOW}Checking system requirements for Kind Echo on Arch Linux...${NC}\n"

# 1. System Packages
PACKAGES=("portaudio" "onnxruntime" "ollama" "go" "gcc" "wget")
MISSING_PACKAGES=()

echo "Checking pacman packages..."
for pkg in "${PACKAGES[@]}"; do
    if pacman -Qi $pkg &> /dev/null; then
        echo -e " - $pkg: ${GREEN}Found${NC}"
    else
        echo -e " - $pkg: ${RED}Missing${NC}"
        MISSING_PACKAGES+=($pkg)
    fi
done

# Check for Piper (usually AUR or manual)
if command -v piper &> /dev/null; then
    echo -e " - piper: ${GREEN}Found${NC}"
else
    echo -e " - piper: ${RED}Missing${NC} (Install via AUR 'piper-tts-bin' or download binary)"
    MISSING_PACKAGES+=("piper-tts-bin (AUR)")
fi

echo ""

# 2. Report Missing
if [ ${#MISSING_PACKAGES[@]} -ne 0 ]; then
    echo -e "${YELLOW}The following packages need to be installed:${NC}"
    for pkg in "${MISSING_PACKAGES[@]}"; do
        echo " - $pkg"
    done
    echo ""
    echo "To install standard packages:"
    echo "  sudo pacman -S ${PACKAGES[*]}"
    echo "To install piper (if missing):"
    echo "  yay -S piper-tts-bin  # or your preferred AUR helper"
else
    echo -e "${GREEN}All system packages appear to be installed.${NC}"
fi

echo ""
read -p "Do you want to proceed with downloading models to ./models? [y/N] " -n 1 -r
echo ""
if [[ ! $REPLY =~ ^[Yy]$ ]]; then
    echo "Aborting."
    exit 1
fi

# 3. Download Models
mkdir -p models
cd models

echo -e "\n${YELLOW}Downloading Models...${NC}"

# Silero VAD
if [ ! -f silero_vad.onnx ]; then
    echo "Downloading Silero VAD..."
    wget -q --show-progress https://github.com/snakers4/silero-vad/raw/master/files/silero_vad.onnx
else
    echo -e "Silero VAD: ${GREEN}Present${NC}"
fi

# Whisper
if [ ! -f ggml-base.en.bin ]; then
    echo "Downloading Whisper Base Model..."
    wget -q --show-progress https://huggingface.co/ggerganov/whisper.cpp/resolve/main/ggml-base.en.bin
else
    echo -e "Whisper Model: ${GREEN}Present${NC}"
fi

# Piper Voice
if [ ! -f en_US-lessac-medium.onnx ]; then
    echo "Downloading Piper Voice (Lessac Medium)..."
    wget -q --show-progress https://huggingface.co/rhasspy/piper-voices/resolve/main/en/en_US/lessac/medium/en_US-lessac-medium.onnx
    wget -q --show-progress https://huggingface.co/rhasspy/piper-voices/resolve/main/en/en_US/lessac/medium/en_US-lessac-medium.onnx.json
else
    echo -e "Piper Voice: ${GREEN}Present${NC}"
fi

echo -e "\n${GREEN}Setup complete!${NC}"
echo "Ensure Ollama is running with: ollama run llama3"
echo "Then run Kind Echo with: go run cmd/kind-echo/main.go"
