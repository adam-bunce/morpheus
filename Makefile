.PHONY: build run clean test

name = morpheous

ANTLR = java -jar antlr-4.13.0-complete.jar -no-visitor -no-listener -Dlanguage=Go

build:
	make clean
	mkdir -p ./generated
	$(ANTLR) -o ./generated/ ./morpheus.g4
	go build .

run:
	make build
	echo
	./morpheus

clean:
	rm -rf ./generated

test:
	make build
	go test ./tests -v -failfast