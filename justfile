
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

debug-logs:
    @ if test -f /tmp/nix-search.log; then less +G /tmp/nix-search.log; else echo ".log file not found" && exit 1; fi
