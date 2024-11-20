#!/bin/bash

# Define the bin directory and an array of target platforms and architectures
BIN_DIR="bin"
PLATFORMS=("linux" "linux" "windows" "windows" "darwin" "darwin")
ARCHS=("amd64" "arm64" "amd64" "arm64" "amd64" "arm64")
EXTENSIONS=("" "" ".exe" ".exe" "" "")

# Create bin directory if it doesn't exist
mkdir -p $BIN_DIR

# Function to build for a specific platform and architecture
build() {
    local os=$1
    local arch=$2
    local ext=$3
    local output="$BIN_DIR/wol-$os-$arch$ext"
    
    if [ -f ${output} ]; then
        rm -f ${output}
    fi

    echo "Building $os/$arch binary file..."
    CGO_ENABLED=0 GOOS=$os GOARCH=$arch go build -o=$output main.go

    if [ $? -ne 0 ]; then
        echo "Failed to build for $os/$arch"
    else
        echo "Successfully build for $os/$arch"
    fi
}

# Loop over all platforms and architectures
for i in ${!PLATFORMS[@]}; do
    build ${PLATFORMS[$i]} ${ARCHS[$i]} ${EXTENSIONS[$i]} &
done

# Wait for all background processes to finish
wait

echo "Build process completed."