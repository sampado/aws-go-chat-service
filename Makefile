.PHONY: build clean deploy gomodgen

build: gomodgen
	export GO111MODULE=on
	env GOOS=linux go build -ldflags="-s -w" -o bin/connect connect/main.go
	env GOOS=linux go build -ldflags="-s -w" -o bin/disconnect disconnect/main.go
	env GOOS=linux go build -ldflags="-s -w" -o bin/message message/main.go

clean:
	rm -rf ./bin ./vendor go.sum

deploy: clean build
	sls deploy --verbose

gomodgen:
	chmod u+x gomod.sh
	./gomod.sh
