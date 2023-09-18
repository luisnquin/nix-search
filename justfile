
RED_COLOR := '\e[0;31m'
LCYAN_COLOR := '\e[1;36m'
NO_COLOR := '\e[0m'

LOGS_FILE := '/tmp/nix-search.log'

# Build the program.
build:
    @mkdir -p build/
    @go build -o ./build/nix-search .

# Build and then run the program.
run: build
    @./build/nix-search

# Install the program in your $GOPATH/bin.
install: build
    @if [ ! -z $GOPATH ]; then echo -e "{{ RED_COLOR }}\$GOPATH is undefined{{ NO_COLOR }}" && exit 1; fi
    @mkdir -p $GOPATH/bin
    @mv ./build/nix-search $GOPATH/bin

# Clean "useless" project files.
clean:
    @rm -rf ./result
    @rm -rf ./nix-search
    @rm -rf ./build

# For debugging. Remove the logs file.
clean-logs:
    @rm -f {{ LOGS_FILE }}

# For debugging. Use less command to keep track of the program logs.
show-logs:
    @if [ ! -f "{{ LOGS_FILE }}" ]; then echo -e "{{ LCYAN_COLOR }}Waiting for {{ LOGS_FILE }} to be created{{ NO_COLOR }}"; while [ ! -f "{{ LOGS_FILE }}" ]; do sleep 1; done; fi; less +FG "{{ LOGS_FILE }}"
