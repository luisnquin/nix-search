
build:
    @mkdir -p build/
    @go build -o ./build/nix-search .

clean:
    @rm -rf ./result
    @rm -rf ./nix-search
    @rm -rf ./build

run: build
    @./build/nix-search
