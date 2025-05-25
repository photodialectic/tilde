#!/bin/bash
# Install vim-plug manually

echo "Installing vim-plug..."

# Create autoload directory if it doesn't exist
mkdir -p ~/.vim/autoload

# Download vim-plug
if command -v curl >/dev/null 2>&1; then
    curl -fLo ~/.vim/autoload/plug.vim --create-dirs \
        https://raw.githubusercontent.com/junegunn/vim-plug/master/plug.vim
elif command -v wget >/dev/null 2>&1; then
    wget -O ~/.vim/autoload/plug.vim \
        https://raw.githubusercontent.com/junegunn/vim-plug/master/plug.vim
else
    echo "Error: Neither curl nor wget found. Please install one of them."
    exit 1
fi

if [ $? -eq 0 ]; then
    echo "✅ vim-plug installed successfully"
    echo "Now run: vim +PlugInstall +qall"
else
    echo "❌ Failed to install vim-plug"
    exit 1
fi
