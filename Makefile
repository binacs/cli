BUILD_TAGS?=
BUILD_FLAGS = -ldflags "-X github.com/binacs/cli/version.GitCommit=`git rev-parse HEAD`"

default: clean build

clean:
	rm -rf bin

build:
	go build $(BUILD_FLAGS) -tags '$(BUILD_TAGS)' -o bin/clid ./cmd/clid
	go build $(BUILD_FLAGS) -tags '$(BUILD_TAGS)' -o bin/cli ./cmd/cli