
build:
    @mkdir -p build/
    @go build -o ./build/nix-search .

run: build
    @./build/nix-search
