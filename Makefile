.PHONY: build

run:
	go run main.go  -desc test -show-env -show-output stdout -show-on-err both echo "!Hola, Mundo!"

build:
	mkdir -p build
	go build -ldflags="-s -w" -o build/executor main.go

clean:
	rm -rf build

install: build
	cp build/executor ~/bin/.
