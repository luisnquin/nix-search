
build:
    @mkdir -p build/
    @go build -o ./build/nix-search .

install: build
    mkdir -p $GOPATH/bin
    mv ./build/nix-search $GOPATH/bin

clean:
    @rm -rf ./result
    @rm -rf ./nix-search
    @rm -rf ./build

run: build
    @./build/nix-search
